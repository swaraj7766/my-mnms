// revive:disable-line:package-comments
package mnms

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/jwtauth/v5"
	"github.com/influxdata/go-syslog/v3"
	"github.com/qeof/q"
)

// simplest and stupid http api

var httpRunning bool

// PostWithToken encapsulates the http POST request with the token
func PostWithToken(url, token string, body io.Reader) (resp *http.Response, err error) {
	bearer := "Bearer " + token
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		q.Q(err)
		return nil, err
	}
	req.Header.Add("Authorization", bearer)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		q.Q(err.Error())
		return nil, err
	}
	return resp, nil
}

// GetWithToken encapsulates the http GET request with the token
func GetWithToken(url, token string) (resp *http.Response, err error) {
	bearer := "Bearer " + token
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		q.Q(err)
		return nil, err
	}
	req.Header.Add("Authorization", bearer)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		q.Q(err)
		return nil, err
	}
	return resp, nil
}

// BuildRouter builds http router
func BuildRouter() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
	}))
	r.Use(middleware.SetHeader("Content-Type", "application/json"))
	r.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) { _, _ = w.Write([]byte("mnms says hello")) })

	varifySuperUser := func(next http.Handler) http.Handler {
		return JWTAuthenticatorRole(MNMSSuperUserRole, next)
	}
	varifyAdmin := func(next http.Handler) http.Handler {
		return JWTAuthenticatorRole(MNMSAdminRole, next)
	}

	r.Route("/api/v1", func(r chi.Router) {
		r.HandleFunc("/login", HandleLogin)
		r.HandleFunc("/ws", WsEndpoint)
		r.HandleFunc("/register", HandleRegister)

		// admin permission
		r.Group(func(r chi.Router) {
			r.Use(jwtauth.Verifier(jwtTokenAuth))
			r.Use(varifyAdmin)
			r.Post("/users", HandleAddUser)
			r.Put("/users", HandleUpdateUser)
		})

		// superuser permission
		r.Group(func(r chi.Router) {
			r.Use(jwtauth.Verifier(jwtTokenAuth))
			r.Use(varifySuperUser)

			r.Post("/commands", HandleCommands)
			r.Post("/devices", HandleDevices)
			r.Post("/topology", HandleTopology)
			r.Get("/syslogs", HandleLocalSyslogs)
			r.Post("/logs", HandleLogs)
		})
		// user permission
		r.Group(func(r chi.Router) {
			r.Use(jwtauth.Verifier(jwtTokenAuth))
			r.Use(jwtauth.Authenticator)

			r.Get("/commands", HandleCommands)
			r.Get("/devices", HandleDevices)
			r.Get("/topology", HandleTopology)
			r.Get("/logs", HandleLogs)
			r.Get("/syslogs", HandleLocalSyslogs)
			r.Get("/users", HandleUsers)
		})
	})
	return r
}

// HTTPMain starts http api service, skipTLS = true if you dont want to serve https
func HTTPMain() {
	var wg sync.WaitGroup
	if httpRunning { // hack for the tests
		return
	}
	httpRunning = true

	// Start a go routine
	wg.Add(1)
	go func() {
		defer wg.Done()
		WebSocketStartWriteMessage()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		r := BuildRouter()
		httpAddr := fmt.Sprintf(":%d", QC.Port)
		q.Q(httpAddr)

		// certManager := autocert.Manager{
		// 	Prompt: autocert.AcceptTOS,
		// 	Cache:  autocert.DirCache("certs"),
		// }

		// server := &http.Server{
		// 	Addr:    httpAddr,
		// 	Handler: r,
		// 	TLSConfig: &tls.Config{
		// 		GetCertificate: certManager.GetCertificate,
		// 	},
		// }
		err := http.ListenAndServe(httpAddr, r)
		if err != nil {
			q.Q("error: cannot run http server", httpAddr, err)
		}

		// server.ListenAndServeTLS("", "")
	}()

	wg.Wait()
}

// RespondWithError write error to the response
func RespondWithError(w http.ResponseWriter, err error) {
	q.Q(err)
	// response code = internal server error
	w.WriteHeader(http.StatusInternalServerError)
	errorInfo := make(map[string]string)
	errorInfo["error"] = fmt.Sprintf("%v", err)

	jsonBytes, err := json.Marshal(errorInfo)
	if err != nil {
		q.Q(err)
		return
	}
	_, _ = w.Write(jsonBytes)
}

func marshalCmdInfo(cmds []byte) (map[string]CmdInfo, error) {
	cmddata := make(map[string]CmdInfo)
	err := json.Unmarshal(cmds, &cmddata)
	if err != nil {
		return nil, err
	}
	return cmddata, nil
}

// HandleCommands handles the commands
func HandleCommands(w http.ResponseWriter, r *http.Request) {
	// enableCors(&w)
	if r.Method == "POST" {
		bodyText, err := ioutil.ReadAll(r.Body)
		if err != nil {
			q.Q(err)
			RespondWithError(w, err)
			return
		}
		defer r.Body.Close()
		cmddata, err := marshalCmdInfo(bodyText)
		if err != nil {
			RespondWithError(w, err)
			return
		}

		if QC.IsRoot {
			rootcmds := make(map[string]CmdInfo)
			for k, v := range cmddata {
				if isRootCommand(k) {
					if v.Status == "" {
						q.Q(v, "root command, execute right away")
						if v.Timestamp == "" {
							v.Timestamp = time.Now().Format(time.RFC3339)
						}
						if v.Command == "" {
							v.Command = k
						}
						v.Name = QC.Name
						rootcmd := RunCmd(&v)
						QC.CmdMutex.Lock()
						QC.CmdData[k] = *rootcmd
						QC.CmdMutex.Unlock()
						rootcmds[k] = *rootcmd
					}
				}
			}

			if len(rootcmds) > 0 {
				resData := cmddata
				for k, v := range rootcmds {
					resData[k] = v
				}
				bodyText, err = json.Marshal(resData)
				if err != nil {
					q.Q(err)
					RespondWithError(w, err)
					return
				}
			}
		}

		UpdateCommands(&cmddata)

		_, _ = w.Write(bodyText)
		return
	}
	id := r.URL.Query().Get("id")
	cmd := r.URL.Query().Get("cmd")
	q.Q(id, cmd)
	// GET

	cmddata := make(map[string]CmdInfo)

	if cmd != "" {
		if cmd == "all" {
			q.Q("all")
			QC.CmdMutex.Lock()
			cmddata = QC.CmdData
			QC.CmdMutex.Unlock()
		} else {
			QC.CmdMutex.Lock()
			found, ok := QC.CmdData[cmd]
			QC.CmdMutex.Unlock()
			if ok {
				q.Q(found)
				cmddata[cmd] = found
			}
		}
		jsonBytes, err := json.Marshal(cmddata)
		if err != nil {
			RespondWithError(w, err)
			return
		}

		_, _ = w.Write(jsonBytes)
		return
	}

	for k, v := range QC.CmdData {
		if v.Status == "" || strings.HasPrefix(v.Status, "pending:") {
			cmddata[k] = v
		}
	}

	if id != "" {
		QC.ClientMutex.Lock()
		client, ok := QC.Clients[id]
		QC.ClientMutex.Unlock()

		if ok && client != "" {
			ci := CmdInfo{
				Timestamp: time.Now().Format(time.RFC3339),
				Command:   client,
			}
			cmddata[client] = ci
			QC.ClientMutex.Lock()
			QC.Clients[id] = ""
			QC.ClientMutex.Unlock()
		}
	}

	jsonBytes, err := json.Marshal(cmddata)
	if err != nil {
		RespondWithError(w, err)
		return
	}

	_, _ = w.Write(jsonBytes)
}

// HandleDevices handles the devices
func HandleDevices(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			RespondWithError(w, err)
			return
		}
		devinfo := make(map[string]DevInfo)
		err = json.Unmarshal(body, &devinfo)
		if err != nil {
			RespondWithError(w, err)
			return
		}
		for _, v := range devinfo {
			//InsertDev(v)
			InsertDev(v)
		}
		_, _ = w.Write(body)
		return
	}

	devid := r.URL.Query().Get("dev")
	q.Q(devid)
	if len(devid) > 0 {
		dev, err := FindDev(devid)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte("error: " + err.Error()))
			return
		}
		jsonBytes, err := json.Marshal(dev)
		if err != nil {
			q.Q(err)
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte("error: " + err.Error()))
			return
		}
		_, _ = w.Write(jsonBytes)
		return
	}
	specialDev := DevInfo{Mac: specialMac, UnixTime: lastUnixTime}
	QC.DevMutex.Lock()
	QC.DevData[specialMac] = specialDev
	QC.DevMutex.Unlock()
	jsonBytes, err := json.Marshal(QC.DevData)
	if err != nil {
		RespondWithError(w, err)
		return
	}

	_, _ = w.Write(jsonBytes)
}

// HandleTopology handles the topology
func HandleTopology(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			RespondWithError(w, err)
			return
		}
		topoinfo := make(map[string]Topology)
		err = json.Unmarshal(body, &topoinfo)
		if err != nil {
			RespondWithError(w, err)
			return
		}
		for k, v := range topoinfo {
			InsertTopology(k, v)
		}
		//q.Q(string(body))
		_, _ = w.Write(body)
		return
	}
	jsonBytes, err := json.Marshal(QC.TopologyData)
	if err != nil {
		RespondWithError(w, err)
		return
	}

	_, _ = w.Write(jsonBytes)
}

// HandleRegister handles the register
func HandleRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			RespondWithError(w, err)
			return
		}
		name := string(body)
		c, ok := QC.Clients[name]
		if ok {
			q.Q("already registered", c, r.RemoteAddr)
		}
		QC.ClientMutex.Lock()
		QC.Clients[name] = ""
		QC.ClientMutex.Unlock()

		_, _ = w.Write(body)
		return
	}
}

// HandleLogin handles login requests
// curl -X POST -H 'Accept: application/json'  https://localhost:27182/api/v1/login -d '{"user":"austinchiang@atop.com.tw","password":"admin"}'
// ! beware this API was changed request body and response body
func HandleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		if !QC.IsRoot {
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte("error: only root can issue tokens"))
			return
		}
		type loginBody struct {
			User     string `json:"user"`
			Password string `json:"password"`
		}
		var body loginBody

		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			RespondWithError(w, err)
			return
		}

		token, err := generateJWT(body.User, body.Password)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte(err.Error()))
			return
		}
		user, err := GetUserConfig(body.User)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte(err.Error()))
			return
		}
		res := make(map[string]interface{})
		res["token"] = token
		res["user"] = body.User
		res["role"] = user.Role
		q.Q(res)
		json.NewEncoder(w).Encode(res)
		// w.Write([]byte(token))
		return
	}
}

type usersBody struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

// HandleUpdateUser handles update user requests
func HandleUpdateUser(w http.ResponseWriter, r *http.Request) {
	var body usersBody
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		RespondWithError(w, err)
		return
	}
	defer r.Body.Close()
	if !UserExist(body.Name) {
		RespondWithError(w, fmt.Errorf("user %s not exist", body.Name))
		return
	}
	err = UpdateUserConfig(body.Name, body.Role, body.Password, body.Email)
	if err != nil {
		RespondWithError(w, err)
		return
	}
	w.Write([]byte("ok"))
	return
}

// HandleAddUser handles add user requests
func HandleAddUser(w http.ResponseWriter, r *http.Request) {
	var body usersBody
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		RespondWithError(w, err)
		return
	}
	defer r.Body.Close()

	if UserExist(body.Name) {
		RespondWithError(w, fmt.Errorf("user %s already exist", body.Name))
		return
	}

	err = AddUserConfig(body.Name, body.Role, body.Password, body.Email)
	if err != nil {
		RespondWithError(w, err)
		return
	}
	_, _ = w.Write([]byte("ok"))
	return
}

// HandleUsers handles users requests
// /api/v1/users
func HandleUsers(w http.ResponseWriter, r *http.Request) {
	// get mnms config
	c, err := GetMNMSConfig()
	if err != nil {
		RespondWithError(w, err)
		return
	}

	// hash password
	for k, v := range c.Users {
		v.Password = "#####"
		c.Users[k] = v
	}

	err = json.NewEncoder(w).Encode(c.Users)
	if err != nil {
		RespondWithError(w, err)
		return
	}
	return
}

// HandleLogs handles logs requests
func HandleLogs(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			RespondWithError(w, err)
			return
		}
		logs := make(map[string]Log)
		err = json.Unmarshal(body, &logs)
		if err != nil {
			RespondWithError(w, err)
			return
		}
		for _, v := range logs {
			InsertLogKind(&v)
		}
		_, _ = w.Write(body)
		return
	}
	jsonBytes, err := json.Marshal(QC.Logs)
	if err != nil {
		RespondWithError(w, err)
		return
	}

	_, _ = w.Write(jsonBytes)
}

// HandleLogs handles logs requests
func HandleLocalSyslogs(w http.ResponseWriter, r *http.Request) {

	start := r.URL.Query().Get("start")
	end := r.URL.Query().Get("end")
	n := r.URL.Query().Get("number")

	number, _ := strconv.Atoi(n)
	if number <= 0 {
		number = 0
	}

	readFile, err := os.Open(path.Join(QC.SyslogLocalPath, file))
	if err != nil {
		RespondWithError(w, err)
		return
	}
	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)
	logs := []syslog.Base{}
	defer func() {
		_ = readFile.Close()
	}()

	for fileScanner.Scan() {
		b, t, err := parsingDataofSyslog(fileScanner.Text())
		if err != nil {
			continue
		}

		r, err := compareTime(start, end, t.Format(foramt))

		if err != nil {
			logs = append(logs, b)
		} else {
			if r {
				logs = append(logs, b)
			}
		}

		if number != 0 {
			if len(logs) >= number {
				break
			}
		}

	}

	jsonBytes, err := json.Marshal(&logs)
	if err != nil {
		RespondWithError(w, err)
		return
	}
	_, _ = w.Write(jsonBytes)

}

// HandleNewUser handles new user creation
// body {account, role} example {"account":"abc@abc.com", "role":"admin"}
// example2 {"account":"alan", "role":"admin"}
func HandleNewUser(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		type userdata struct {
			Account string `json:"account"`
			Role    string `json:"role"`
		}
		var user userdata
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			RespondWithError(w, err)
			return
		}

		// TODO: hardcode master is not a good idea
		password, err := GenPassword("mnmsmaster", user.Account)
		if err != nil {
			RespondWithError(w, err)
			return
		}
		jsonByte, err := json.Marshal(map[string]string{
			"account":  user.Account,
			"password": password,
		})
		if err != nil {
			RespondWithError(w, err)
			return
		}

		_, _ = w.Write(jsonByte)
		return
	}
}
