// revive:disable-line:package-comments
package mnms

import (
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
	"github.com/icza/backscanner"
	"github.com/influxdata/go-syslog/v3"
	"github.com/pquerna/otp/totp"
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

func CheckStaticFilesFolder() (string, error) {
	// check static files folder
	nmswd, err := CheckMNMSFolder()
	if err != nil {
		q.Q(err)
		nmswd = "."
	}
	fileDir := path.Join(nmswd, "files")
	// check if folder exist
	if _, err := os.Stat(fileDir); os.IsNotExist(err) {
		err := os.Mkdir(fileDir, 0755)
		if err != nil {
			q.Q(err)
			return "", err
		}
	}
	return fileDir, nil
}

// FileServer conveniently sets up a http.FileServer handler to serve
// static files from a http.FileSystem.
func FileServer(r chi.Router, path string, root http.FileSystem) {
	if len(path) == 0 {
		q.Q("FileServer cannot be used with an empty path")

		return
	}

	if strings.ContainsAny(path, "{}*") {
		q.Q("FileServer does not permit any URL parameters.")
		return
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		ctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(ctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
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
	r.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("mnms says hello"))
		if err != nil {
			q.Q(err)
		}
	})

	varifySuperUser := func(next http.Handler) http.Handler {
		return JWTAuthenticatorRole(MNMSSuperUserRole, next)
	}
	varifyAdmin := func(next http.Handler) http.Handler {
		return JWTAuthenticatorRole(MNMSAdminRole, next)
	}

	// static file directory
	fileDir, err := CheckStaticFilesFolder()
	if err != nil {
		q.Q(err)
	}

	r.Route("/api/v1", func(r chi.Router) {
		r.HandleFunc("/login", HandleLogin)
		r.Post("/2fa/validate", HandleValidate2FA)
		r.HandleFunc("/ws", WsEndpoint)
		r.HandleFunc("/register", HandleRegister)

		// admin permission
		r.Group(func(r chi.Router) {
			r.Use(jwtauth.Verifier(jwtTokenAuth))
			r.Use(varifyAdmin)
			r.Post("/users", HandleAddUser)
			r.Put("/users", HandleUpdateUser)
			r.Delete("/users", HandleDeleteUser)
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
			r.Get("/users", HandleUsers)
			r.HandleFunc("/2fa/secret", Handle2FA)

			FileServer(r, "/files", http.Dir(fileDir))
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
	_, err = w.Write(jsonBytes)
	if err != nil {
		q.Q(err)
	}
}

func unmarshalCmdInfo(cmds []byte) (map[string]CmdInfo, error) {
	cmddata := make(map[string]CmdInfo)
	err := json.Unmarshal(cmds, &cmddata)
	if err != nil {
		return nil, err
	}
	return cmddata, nil
}

// HandleCommands accepts the commands and returns command info and history
//
// POST /api/v1/commands
//
//	   Example parameter: (map[string]CmdInfo)
//		{ "beep 01-22-33-44-55-66 10.1.1.1" : {}, "devices publish", : {"all": true} }
//
// GET /api/v1/commands?id=client1
//
//	retrieve commands intended for service client1
//
// GET /api/v1/commands?cmd=beep 01-22-33-44-55-66 10.1.1.1
//
//	retrieve command status of a particular command
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
		cmddata, err := unmarshalCmdInfo(bodyText)
		if err != nil {
			q.Q(err)
			RespondWithError(w, err)
			return
		}
		for k, v := range cmddata {
			found := false
			ws := strings.Split(k, " ")
			if len(ws) < 2 {
				err = fmt.Errorf("error: invalid short command")
				RespondWithError(w, err)
				return
			}
			acmd := ws[0]
			for _, c := range ValidCommands {
				if c == acmd {
					found = true
				}
			}
			if !found {
				//err := fmt.Errorf("error: invalid command name %v", acmd)
				//RespondWithError(w, err)
				v.Result = "error: invalid command"
				cmddata[k] = v
			}
		}
		retrieveRootCmd(cmddata)
		UpdateCmds(&cmddata)
		_, err = w.Write(bodyText)
		if err != nil {
			q.Q(err)
		}
		return
	}
	id := r.URL.Query().Get("id")
	cmd := r.URL.Query().Get("cmd")
	q.Q(id, cmd)
	// GET
	cmddata := make(map[string]CmdInfo)
	if cmd != "" {
		q.Q("get cmd info", cmd)
		// clients wants info on a specific cmd
		if cmd == "all" {
			q.Q("client wants to get all commands")
			QC.CmdMutex.Lock()
			cmddata = QC.CmdData
			QC.CmdMutex.Unlock()
		} else {
			q.Q("clients wants", cmd)
			QC.CmdMutex.Lock()
			found, ok := QC.CmdData[cmd]
			QC.CmdMutex.Unlock()
			if ok {
				cmddata[cmd] = found
			}
			q.Q(found)
		}
		jsonBytes, err := json.Marshal(cmddata)
		if err != nil {
			RespondWithError(w, err)
			return
		}
		_, err = w.Write(jsonBytes)
		if err != nil {
			q.Q(err)
		}
		return
	}
	q.Q("get all non status or pending cmds")
	QC.CmdMutex.Lock()
	for k, v := range QC.CmdData {
		if v.Status == "" || strings.HasPrefix(v.Status, "pending:") {
			if id != "" && v.Client != "" && id != v.Client {
				continue
			}
			if v.Name == QC.Name && QC.IsRoot {
				continue
			}
			cmddata[k] = v
		}
	}
	QC.CmdMutex.Unlock()
	q.Q("sending to client", cmddata)
	jsonBytes, err := json.Marshal(cmddata)
	if err != nil {
		RespondWithError(w, err)
		return
	}
	_, err = w.Write(jsonBytes)
	if err != nil {
		q.Q(err)
	}
}

// HandleDevices accepts devices information and stores them and can return device info
//
// POST /api/v1/devices
//
//		Example parameter: (map[string]DevInfo)
//	         {"00-60-E9-2D-91-3E": {"mac":"00-60-E9-2D-91-3E","modelname":...},
//	          "00-60-E9-1F-A6-02":{"mac":"00-60-E9-1F-A6-02","modelname":...}}
//
// GET /api/v1/devices
//
//	retrieve devices information
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
			// InsertDev(v)
			InsertDev(v)
		}
		_, err = w.Write(body)
		if err != nil {
			q.Q(err)
		}
		return
	}

	devid := r.URL.Query().Get("dev")
	q.Q(devid)
	if len(devid) > 0 {
		dev, err := FindDev(devid)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			_, err = w.Write([]byte("error: " + err.Error()))
			if err != nil {
				q.Q(err)
			}
			return
		}
		jsonBytes, err := json.Marshal(dev)
		if err != nil {
			q.Q(err)
			w.WriteHeader(http.StatusInternalServerError)
			_, err = w.Write([]byte("error: " + err.Error()))
			if err != nil {
				q.Q(err)
			}
			return
		}
		_, err = w.Write(jsonBytes)
		if err != nil {
			q.Q(err)
		}
		return
	}
	specialDev := DevInfo{Mac: specialMac, Timestamp: lastTimestamp}
	QC.DevMutex.Lock()
	QC.DevData[specialMac] = specialDev
	QC.DevMutex.Unlock()
	jsonBytes, err := json.Marshal(QC.DevData)
	if err != nil {
		RespondWithError(w, err)
		return
	}

	_, err = w.Write(jsonBytes)
	if err != nil {
		q.Q(err)
	}
}

// HandleTopology handles the topology
//
// POST /api/v1/topology
//
//	Example parameter: map[string]Topology
//
// GET /api/v1/topology
//
//	returns topology
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
		// q.Q(string(body))
		_, err = w.Write(body)
		if err != nil {
			q.Q(err)
		}
		return
	}
	jsonBytes, err := json.Marshal(QC.TopologyData)
	if err != nil {
		RespondWithError(w, err)
		return
	}

	_, err = w.Write(jsonBytes)
	if err != nil {
		q.Q(err)
	}
}

// HandleRegister accept client information and returns cluster client information
//
// POST /api/v1/register
//
//	    Example parameter: ClientInfo{}
//	             { "client1": { "Name": "client1", "NumDevices": 0, ... },
//			"client2": { "Name": "client2", "NumDevices": 0, "NumCmds": 0, ...}}
//
// GET /api/v1/register
//
//	returns cluster client information
func HandleRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			RespondWithError(w, err)
			return
		}
		ci := ClientInfo{}
		err = json.Unmarshal(body, &ci)
		if err != nil {
			RespondWithError(w, err)
			return
		}
		q.Q(ci)
		name := ci.Name
		if name == "" {
			RespondWithError(w, fmt.Errorf("name required"))
			return
		}
		QC.ClientMutex.Lock()
		QC.Clients[name] = ci
		QC.ClientMutex.Unlock()
		_, err = w.Write(body)
		if err != nil {
			q.Q(err)
		}
		return
	}
	jsonBytes, err := json.Marshal(QC.Clients)
	if err != nil {
		RespondWithError(w, err)
		return
	}
	_, err = w.Write(jsonBytes)
	if err != nil {
		q.Q(err)
	}
}

// HandleLogin handles login requests
//
// POST  /api/v1/login
//
//		Example parameter:
//		         {"user":"user1@example.com","password":"Pas$word1"}
//
//		Response:
//	     need 2fa : {"sessionID": "sessionID", "user":"user1"}
//		   {"token": "AAA...", "user": "user1", "role": "admin"}
func HandleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		if !QC.IsRoot {
			w.WriteHeader(http.StatusUnauthorized)
			_, err := w.Write([]byte("error: only root can issue tokens"))
			if err != nil {
				q.Q(err)
			}
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
			_, err = w.Write([]byte(err.Error()))
			if err != nil {
				q.Q(err)
			}
			return
		}
		user, err := GetUserConfig(body.User)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			_, err = w.Write([]byte(err.Error()))
			if err != nil {
				q.Q(err)
			}
			return
		}
		res := make(map[string]interface{})
		// check 2fa
		if user.Enable2FA {
			sessionID := createLoginSession(*user)
			res["sessionID"] = sessionID
			res["user"] = user.Name
			res["email"] = user.Email
			err = json.NewEncoder(w).Encode(res)
			if err != nil {
				RespondWithError(w, err)
			}
			return
		}
		res["token"] = token
		res["user"] = body.User
		res["role"] = user.Role
		q.Q(res)
		err = json.NewEncoder(w).Encode(res)
		if err != nil {
			RespondWithError(w, err)
		}
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

// HandleDeleteUser handles delete user requests
//
//	 DELETE /api/v1/users
//
//	Example parameter: { "name": "abc"}
func HandleDeleteUser(w http.ResponseWriter, r *http.Request) {
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
	err = DeleteUserConfig(body.Name)
	if err != nil {
		RespondWithError(w, err)
		return
	}
	_, err = w.Write([]byte("ok"))
	if err != nil {
		q.Q(err)
	}
}

// HandleUpdateUser handles update user requests
//
// PUT /api/v1/users
//
//	Example parameter: { "name": "abc", "email": "abc@def.com", "password": "password1" , "role": "admin"}
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
	_, err = w.Write([]byte("ok"))

	if err != nil {
		RespondWithError(w, err)
	}
	return

}

// HandleAddUser handles add user requests
//
// POST /api/v1/users
//
//	Example parameter: { "name": "abc", "email": "abc@def.com", "password": "Pas$Word1" , "role": "admin"}
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
	_, err = w.Write([]byte("ok"))
	if err != nil {
		q.Q(err)
	}

}

// HandleValidate2FA handles 2FA validation requests
// POST /api/v1/2fa/validate
// request :{"sessionID":"id", "code":"123456"}
// Validate 2fa code
// example response: {"valid": true, "user": "user1", "token": "token", "role": "admin""}

func HandleValidate2FA(w http.ResponseWriter, r *http.Request) {
	var data map[string]string
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		RespondWithError(w, err)
		return
	}
	defer r.Body.Close()
	sessionID := data["sessionID"]
	code := data["code"]
	user, err := getLoginSession(sessionID)
	if err != nil {
		RespondWithError(w, err)
		return
	}

	if code == "" {
		RespondWithError(w, fmt.Errorf("code is empty"))
		return
	}

	if !user.Enable2FA {
		RespondWithError(w, fmt.Errorf("2fa not enabled"))
		return
	}
	secret := user.Secret

	valid := totp.Validate(code, secret)
	q.Q(code, secret, valid)
	var token string
	if !valid {

		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte("invalid code"))
		return

	}
	token, err = generateJWT(user.Name, user.Password)
	if err != nil {
		RespondWithError(w, err)
		return
	}

	res := make(map[string]interface{})
	res["user"] = user.Name
	res["token"] = token
	res["role"] = user.Role
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		RespondWithError(w, err)
	}
	return
}

// Handle2FA handles 2FA requests
// GET /api/v1/2fa/secret?user=user1
// get current user's 2fa secret
// response : {"user":"user1", "secret":"secret"}
//
// POST /api/v1/2fa/secret
// request body {"user":"user1"}
// response : {"user":"user1", "secret":"secret"}
// Generate 2fa secret
//
// PUT /api/v1/2fa/secret
// request body {"user":"user1"}
// Update 2fa secret
//
// DELETE /api/v1/2fa/secret
// request body {"user":"user1"}
// Disable user's 2fa
func Handle2FA(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" {
		userID := r.URL.Query().Get("user")
		if userID == "" {
			RespondWithError(w, fmt.Errorf("user is empty"))
			return
		}
		user, err := GetUserConfig(userID)
		if err != nil {
			RespondWithError(w, err)
			return
		}
		if !user.Enable2FA {
			RespondWithError(w, fmt.Errorf("2fa not enabled"))
			return
		}
		secret := user.Secret
		res := make(map[string]interface{})
		res["user"] = userID
		res["secret"] = secret
		res["account"] = user.Email
		res["issuer"] = IssuerOf2FA
		res["enable2fa"] = user.Enable2FA
		err = json.NewEncoder(w).Encode(res)
		if err != nil {
			RespondWithError(w, err)
		}
		return
	}

	if r.Method == "POST" {
		var data map[string]string
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			RespondWithError(w, err)
			return
		}
		defer r.Body.Close()
		userID := data["user"]

		user, err := GetUserConfig(userID)
		if err != nil {
			RespondWithError(w, err)
			return
		}
		if user.Enable2FA {
			RespondWithError(w, fmt.Errorf("2fa already enabled, use PUT to update"))
			return
		}
		if len(user.Email) == 0 {
			RespondWithError(w, fmt.Errorf("email is empty"))
			return
		}
		// generate 2fa secret
		secret, err := generate2FASecret(user.Email)
		if err != nil {
			RespondWithError(w, err)
			return
		}
		// save 2fa secret
		user.Secret = secret
		user.Enable2FA = true
		err = MergeUserConfig(*user)
		if err != nil {
			RespondWithError(w, err)
			return
		}
		// write response
		res := make(map[string]interface{})
		res["secret"] = secret
		res["account"] = user.Email
		res["issuer"] = IssuerOf2FA
		res["user"] = userID
		err = json.NewEncoder(w).Encode(res)
		if err != nil {
			RespondWithError(w, err)
		}
		return
	}

	if r.Method == "PUT" {
		var data map[string]string
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			RespondWithError(w, err)
			return
		}
		defer r.Body.Close()
		userID := data["user"]
		user, err := GetUserConfig(userID)
		if err != nil {
			RespondWithError(w, err)
			return
		}
		if !user.Enable2FA {
			RespondWithError(w, fmt.Errorf("2fa not enabled"))
			return
		}
		if len(user.Email) == 0 {
			RespondWithError(w, fmt.Errorf("email is empty"))
			return
		}
		// generate 2fa secret
		secret, err := generate2FASecret(user.Email)
		if err != nil {
			RespondWithError(w, err)
			return
		}
		// save 2fa secret
		user.Secret = secret
		err = MergeUserConfig(*user)
		if err != nil {
			RespondWithError(w, err)
			return
		}
		// write response
		res := make(map[string]interface{})
		res["secret"] = secret
		res["account"] = user.Email
		res["issuer"] = IssuerOf2FA
		res["user"] = userID
		err = json.NewEncoder(w).Encode(res)
		if err != nil {
			RespondWithError(w, err)
		}
		return
	}
	if r.Method == "DELETE" {
		// disable 2fa and delete secret
		var data map[string]string
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			RespondWithError(w, err)
			return
		}
		defer r.Body.Close()
		userID := data["user"]
		user, err := GetUserConfig(userID)
		if err != nil {
			RespondWithError(w, err)
			return
		}
		if !user.Enable2FA {
			RespondWithError(w, fmt.Errorf("2fa not enabled"))
			return
		}
		user.Enable2FA = false
		user.Secret = ""
		err = MergeUserConfig(*user)
		if err != nil {
			RespondWithError(w, err)
			return
		}
		// write response
		_, err = w.Write([]byte("ok"))
		if err != nil {
			RespondWithError(w, err)
		}
		return
	}
}

// HandleUsers handles users requests
//
// GET  /api/v1/users
//
//	Example return: { "name": "abc", "email": "abc@def.com", "password": "pass1" , "role": "admin"}
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
}

// HandleLogs accepts log messages for storage in memory and returns them
//
// POST /api/v1/logs
//
//	Example paramter: (map[string]Log)
//
// GET /api/v1/logs
//
//	returns logs from memory
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
		_, err = w.Write(body)
		if err != nil {
			q.Q(err)
		}
		return
	}
	jsonBytes, err := json.Marshal(QC.Logs)
	if err != nil {
		RespondWithError(w, err)
		return
	}

	_, err = w.Write(jsonBytes)
	if err != nil {
		q.Q(err)
	}
}

// HandleLocalLogs handles syslogs requests
//
// GET /api/v1/syslogs
//
// Example parameter: { "start": "2023/02/21 22:06:00", "end": "2023/02/23 22:08:00", "number": 3}
//
// returns local syslogs from files
func HandleLocalSyslogs(w http.ResponseWriter, r *http.Request) {
	start := r.URL.Query().Get("start")
	end := r.URL.Query().Get("end")
	n := r.URL.Query().Get("number")

	number, _ := strconv.Atoi(n)
	if number <= 0 {
		number = 0
	}

	f, err := os.Open(QC.SyslogLocalPath)
	if err != nil {
		_, err = w.Write([]byte(""))
		if err != nil {
			q.Q(err)
		}
		return
	}
	defer func() {
		_ = f.Close()
	}()
	fi, err := f.Stat()
	if err != nil {
		_, err = w.Write([]byte(""))
		if err != nil {
			q.Q(err)
		}
		return
	}
	scanner := backscanner.New(f, int(fi.Size()))
	//fileScanner := bufio.NewScanner(f)
	//fileScanner.Split(bufio.ScanLines)
	logs := []syslog.Base{}

	for {
		line, _, err := scanner.LineBytes()
		if err != nil {
			if err == io.EOF {
				q.Q(QC.SyslogLocalPath, "found to EOF")
			} else {
				q.Q("err:", err)
			}
			break
		}
		b, t, err := parsingDataofSyslog(string(line))
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
		_, err = w.Write([]byte(""))
		if err != nil {
			q.Q(err)
		}
		return
	}
	_, err = w.Write(jsonBytes)
	if err != nil {
		q.Q(err)
	}
}
