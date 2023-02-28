package simulator

import (
	"github.com/qeof/q"
)

// FirmwareRun run Firmware upgrade server
func (a *AtopGwdClient) FirmwareRun() {
	q.Q(a.ModelInfo.IPAddress, "firmwarm Run")
	go func() {
		err := a.firmware.Run()
		if err != nil {
			q.Q(err)
		}
	}()
}

// FirmwareShutdown shutdonw Firmware upgrade server
func (a *AtopGwdClient) FirmwareShutdown() {
	q.Q(a.ModelInfo.IPAddress, "firmware Shutdown")
	a.firmware.Shutdown()
}
