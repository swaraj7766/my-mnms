package mnms

import (
	"fmt"
	"log"

	"github.com/kardianos/service"
	"github.com/qeof/q"
)

var Usage = "run | start | stop | restart | install |  uninstall  | status "
var servicName = "mnmsctl"

const displayName = "mnmsctl software"
const description = "atop mnmsctl software application."
const DaemonFlag = "daemon"

func filterArgs(args []string) []string {
	found := false
	foundindex := []int{}
	for i, v := range args {
		if v == "-"+DaemonFlag {
			next := i + 1
			if len((args)) > next {
				found = true
				foundindex = append(foundindex, i, next)
			}
		}
	}
	if found {
		array := []string{}
		for i, v := range args {
			if i == foundindex[0] || i == foundindex[1] {
				continue
			}
			array = append(array, v)
		}
		return array
	} else {
		return args
	}
}

func NewDaemon(name string, args []string) (*Daemon, error) {
	args = filterArgs(args)
	if len(name) != 0 {
		servicName = name
	}
	options := make(service.KeyValue)
	options["OnFailure"] = "restart"
	svcConfig := &service.Config{
		Name:         servicName,
		DisplayName:  fmt.Sprintf("%v %v", servicName, displayName),
		Description:  description,
		Dependencies: []string{},
		Option:       options,
		Arguments:    args[1:],
	}
	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		return nil, err
	}
	d := &Daemon{srv: s, prog: prg}
	return d, nil
}

type Daemon struct {
	srv  service.Service
	prog *program
}

func (d *Daemon) RunMode(mode string) error {
	switch mode {
	case "run":
		if err := d.srv.Run(); err != nil {
			return err
		}
	case "start":
		if err := d.srv.Start(); err != nil {
			return err
		}

	case "stop":
		if err := d.srv.Stop(); err != nil {
			return err
		}

	case "restart":
		if err := d.srv.Restart(); err != nil {
			log.Fatal(err)
		}

	case "install":
		if err := d.srv.Install(); err != nil {
			return err
		}
		q.Q("complete install")
		if err := d.srv.Start(); err != nil {
			return err
		}
		q.Q("complete start")
	case "uninstall":
		d.srv.Stop()
		if err := d.srv.Uninstall(); err != nil {
			return err
		}
	case "status":
		s, err := d.srv.Status()
		if err != nil {
			return err
		}
		switch s {
		case service.StatusUnknown:
			log.Print("Unknown")
		case service.StatusRunning:
			log.Print("Running")
		case service.StatusStopped:
			log.Print("Stopped")
		}

	default:
		if err := d.srv.Run(); err != nil {
			return err
		}
	}

	return nil
}

//RegisterRunEvent register run service
func (d *Daemon) RegisterRunEvent(run RunEvent) {
	d.prog.registerRunEvent(run)
}

//RegisterStopEvent register stop event like defer
func (d *Daemon) RegisterStopEvent(stop StopEvent) {
	d.prog.registerStopEvent(stop)
}

type program struct {
	runevent  RunEvent
	stopevent StopEvent
}

type RunEvent func()
type StopEvent func()

func (p *program) registerRunEvent(run RunEvent) {
	p.runevent = run
}
func (p *program) registerStopEvent(stop StopEvent) {
	p.stopevent = stop
}
func (p *program) Start(s service.Service) error {
	// Start should not block. Do the actual work async.
	go p.run()
	return nil
}

func (p *program) run() error {
	if p.runevent != nil {
		p.runevent()
	}
	return nil

}
func (p *program) Stop(s service.Service) error {
	// Stop should not block. Return with a few seconds.
	if p.stopevent != nil {
		p.stopevent()
	}
	return nil
}
