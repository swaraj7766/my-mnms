package simulator

import (
	"github.com/qeof/q"
)

// TelnetServerRun run telnet server
func (a *AtopGwdClient) TelnetServerRun() {
	q.Q(a.ModelInfo.IPAddress, "telnet server Run")
	go func() {
		_ = a.telnetserver.Run()
	}()
}

// TelnetServerShutdown  shutdonw telnet server
func (a *AtopGwdClient) TelnetServerShutdown() {
	q.Q(a.ModelInfo.IPAddress, "telnet server Shutdown")
	a.telnetserver.Shutdown()
}
