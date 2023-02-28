package mnms

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/qeof/q"
	cron "github.com/robfig/cron/v3"
)

type CronInfo struct {
	Job     string       `json:"job"`
	EntryID cron.EntryID `json:"entryid"`
}

func init() {
	q.Q("init crontab")
	QC.CronJobs = make([]CronInfo, 0) // like a crontab but in memory and there is an entry id for each job
	QC.Cron = cron.New()              // scheduler instance
	// load crontab if exists
	_ = LoadCrontab()
}

/* AddCronJob add a cron job in system
 * @param[in] cmd: standard cron format, eg "0 0 0 * * scan gwd"
 * @param[in] id: cron job id, leave it empty in request, will be filled in once added, can be used to remove the job
 */
func AddCronJob(info CronInfo) {
	q.Q(info)

	cmd := info.Job
	ws := strings.Split(cmd, " ")
	if len(ws) < 6 {
		q.Q("error", len(ws))
		return
	}

	cronTime := strings.Join(ws[0:5], " ")
	job := strings.Join(ws[5:], " ")
	q.Q(cronTime, job)
	id, err := QC.Cron.AddFunc(cronTime, func() {
		q.Q("cron job", job)
		cd := make(map[string]CmdInfo)
		ci := &CmdInfo{
			Command:   job,
			Timestamp: time.Now().Format(time.RFC3339),
		}
		cd[job] = *ci
		if isRootCommand(job) {
			ci.Name = QC.Name
			ret := RunCmd(ci)
			QC.CmdMutex.Lock()
			QC.CmdData[job] = *ret
			QC.CmdMutex.Unlock()
		} else {
			UpdateCommands(&cd)
		}
	})
	if err != nil {
		q.Q(err)
		return
	}
	QC.Cron.Start()
	QC.CmdMutex.Lock()
	QC.CronJobs = append(QC.CronJobs, CronInfo{cmd, id})
	QC.CmdMutex.Unlock()
}

// DeleteCronJob delete a cron job in system by EntryID
func DeleteCronJob(id cron.EntryID) {
	q.Q(id)

	QC.Cron.Remove(cron.EntryID(id))
	QC.CmdMutex.Lock()
	defer QC.CmdMutex.Unlock()
	for i, v := range QC.CronJobs {
		if v.EntryID == cron.EntryID(id) {
			QC.CronJobs = append(QC.CronJobs[:i], QC.CronJobs[i+1:]...)
			break
		}
	}
}

// DumpCrontab dump cron jobs to file as crontab
func DumpCrontab() error {
	fn := fmt.Sprintf("%s_crontab", QC.Name)
	QC.DevMutex.Lock()
	defer QC.DevMutex.Unlock()
	q.Q(fn)
	var data []byte
	for _, v := range QC.CronJobs {
		data = append(data, []byte(v.Job)...)
		data = append(data, []byte("\n")...)
	}
	err := os.WriteFile(fn, data, 0644)
	if err != nil {
		q.Q(err)
		return err
	}
	return nil
}

// LoadCrontab load crontab from file
func LoadCrontab() error {
	readFile, err := os.Open(fmt.Sprintf("./%s_crontab", QC.Name))
	if err != nil {
		q.Q(err)
		return err
	}
	fileScanner := bufio.NewScanner(readFile)
	defer readFile.Close()
	fileScanner.Split(bufio.ScanLines)
	for fileScanner.Scan() {
		q.Q(fileScanner.Text())
		job := CronInfo{
			Job: fileScanner.Text(),
		}
		AddCronJob(job)
	}
	return nil
}

// RemoveCrontab delete crontab file
func RemoveCrontab() error {
	fn := fmt.Sprintf("./%s_crontab", QC.Name)
	q.Q(fn)
	err := os.Remove(fn)
	if err != nil {
		q.Q(err)
		return err
	}
	return nil
}

/*
	 CrontabCmd handle crontab command
		crontab add: add a cron job
		crontab delete: delete a cron job
	    crontab dump: dump cron jobs to file as crontab
		crontab load: load crontab from file
		crontab remove: remove crontab file
*/
func CrontabCmd(cmdinfo *CmdInfo) *CmdInfo {
	cmd := cmdinfo.Command

	if !QC.IsRoot {
		q.Q("not root, no crontab")
		cmdinfo.Status = "error: not root, no crontab"
		return cmdinfo
	}

	if strings.HasPrefix(cmd, "config crontab add") {
		// config crontab add 0 0 0 * * all scan gwd
		info := CronInfo{Job: strings.Join(strings.Split(cmd, " ")[3:], " ")}
		AddCronJob(info)
		cmdinfo.Status = "ok"
		return cmdinfo
	}

	if strings.HasPrefix(cmd, "config crontab delete") {
		// crontab delete all
		ws := strings.Split(cmd, " ")
		if len(ws) < 4 {
			cmdinfo.Status = "error: invalid command"
			return cmdinfo
		}
		if ws[3] == "all" {
			QC.Cron.Stop()
			for _, v := range QC.CronJobs {
				QC.Cron.Remove(v.EntryID)
			}
			QC.CmdMutex.Lock()
			QC.CronJobs = make([]CronInfo, 0)
			QC.CmdMutex.Unlock()
			cmdinfo.Status = "ok"
			return cmdinfo
		}

		// crontab delete id
		id := ws[3]
		intid, err := strconv.Atoi(id)
		if err != nil {
			q.Q(err)
			cmdinfo.Status = "error: invalid id"
			return cmdinfo
		}
		DeleteCronJob(cron.EntryID(intid))
		cmdinfo.Status = "ok"
		return cmdinfo
	}

	if strings.HasPrefix(cmd, "config crontab list") {
		// crontab list
		jsonString, err := json.Marshal(QC.CronJobs)
		if err != nil {
			q.Q(err)
			cmdinfo.Status = "error: " + err.Error()
			cmdinfo.Result = "error: " + err.Error()
			return cmdinfo
		}
		cmdinfo.Status = "ok"
		cmdinfo.Result = string(jsonString)
		return cmdinfo
	}

	if cmd == "config crontab dump" {
		err := DumpCrontab()
		if err != nil {
			q.Q(err)
			cmdinfo.Status = fmt.Sprintf("dump crontab failed: %v", err)
			return cmdinfo
		}
		cmdinfo.Status = "ok"
		return cmdinfo
	}

	if cmd == "config crontab load" {
		err := LoadCrontab()
		if err != nil {
			q.Q(err)
			cmdinfo.Status = fmt.Sprintf("load crontab failed: %v", err)
			return cmdinfo
		}
		cmdinfo.Status = "ok"
		return cmdinfo
	}

	if cmd == "config crontab remove" {
		err := RemoveCrontab()
		if err != nil {
			q.Q(err)
			cmdinfo.Status = fmt.Sprintf("remove crontab failed: %v", err)
			return cmdinfo
		}
		cmdinfo.Status = "ok"
		return cmdinfo
	}

	cmdinfo.Status = "error: invalid command"
	return cmdinfo
}
