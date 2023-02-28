package simulator

import (
	"github.com/sirupsen/logrus"
)

// SnmpRun run snmp
func (a *AtopGwdClient) SnmpRun(ip string) {
	logrus.Printf("snmp:%v Run", a.ModelInfo.IPAddress)
	go func() {
		err := a.snmp.Run(ip)
		if err != nil {
			logrus.Fatal(err)
		}
	}()
}

// SnmpRun run snmp
func (a *AtopGwdClient) SetLogger(log *logrus.Logger) {
	a.snmp.SetLogger(log)
}

// SnmpShutdown  shutdonw snmp
func (a *AtopGwdClient) SnmpShutdown() {
	logrus.Printf("snmp:%v Shutdown", a.ModelInfo.IPAddress)
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
