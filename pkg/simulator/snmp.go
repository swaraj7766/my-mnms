package simulator

import (
	"github.com/qeof/q"
)

// SnmpRun run snmp
func (a *AtopGwdClient) SnmpRun(ip string) error {
	q.Q("snmp:", a.ModelInfo.IPAddress, " Run")
	go func() {
		err := a.snmp.Run(ip)
		if err != nil {
			q.Q(err)
		}
	}()
	return nil
}

// SnmpShutdown  shutdonw snmp
func (a *AtopGwdClient) SnmpShutdown() {
	q.Q("snmp:", a.ModelInfo.IPAddress, " Shutdown")
	a.snmp.Shutdown()
}

// updataSnmp update snmp data
func (a *AtopGwdClient) upSnmpData() {
	v := a.snmp.GetData()
	v.SetSystem(a.ModelInfo.Hostname)
	v.SetIP(a.ModelInfo.IPAddress)
	v.SetMask(a.ModelInfo.Netmask)
	v.SetGateWay(a.ModelInfo.Gateway)
	v.SetMac(a.ModelInfo.MACAddress)
}

// bindValue bind value with snmp
func (a *AtopGwdClient) bindValue() {
	a.snmp.GetData().BindSystem(&a.ModelInfo.Hostname)
	a.snmp.GetData().BindAp(&a.ModelInfo.Ap)
	a.snmp.GetData().BindKernel(&a.ModelInfo.Kernel)
	a.snmp.GetData().BindIP(&a.ModelInfo.IPAddress)
	a.snmp.GetData().BindMask(&a.ModelInfo.Netmask)
	a.snmp.GetData().BindGateWay(&a.ModelInfo.Gateway)
	a.snmp.GetData().BindMac(&a.ModelInfo.MACAddress)

	a.snmp.GetData().BindUser(&a.user)
	a.snmp.GetData().BindPwd(&a.pwd)
}
