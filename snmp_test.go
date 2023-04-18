package mnms

import (
	"testing"
	"time"
)

func TestSnmpCommunity(t *testing.T) {
	// q.P = "snmp_test"
	// ctx, cancel := context.WithCancel(context.Background())
	// defer cancel()
	// err := runTestSnmpSimulators(ctx)
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// defer func() {
	// 	for _, v := range testSnmpSimulators {
	// 		_ = v.Shutdown()
	// 	}
	// }()
	// time.Sleep(time.Millisecond * 1000)
	t.Skip("skipping test, simulator has issue can't work, this test just for local test")
	go func() {
		GwdMain()
	}()
	err := GwdInvite()
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(time.Millisecond * 1000)
	t.Log(QC.DevData)
	// dev, ok := QC.DevData[targetMac]
	// if !ok {
	// 	t.Fatal("no such device: ", targetMac)
	// }
	// t.Log(dev)

	r, rw, err := GetSNMPCommunity("admin", "default", "192.168.12.245")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(r, rw)
}
