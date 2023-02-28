package simulator

func invitePacket() []byte {

	packet := make([]byte, 300)
	packet[0] = 2
	packet[1] = 1
	packet[2] = 6
	packet[4] = 0x92
	packet[5] = 0xDA
	return packet

}

func configPacket() []byte {

	packet := make([]byte, 300)
	packet[0] = 0
	packet[1] = 1
	packet[2] = 6
	packet[4] = 0x92
	packet[5] = 0xDA
	return packet

}
func rebootPacket() []byte {
	packet := make([]byte, 300)
	packet[0] = 5
	packet[1] = 1
	packet[2] = 6
	packet[4] = 0x92
	packet[5] = 0xDA
	return packet
}

type ModelInfo struct {
	Model      string `json:"model"`
	MACAddress string `json:"macAddress"`
	IPAddress  string `json:"iPAddress"`
	Netmask    string `json:"netmask"`
	Gateway    string `json:"gateway"`
	Hostname   string `json:"hostname"`
	Kernel     string `json:"kernel"`
	Ap         string `json:"ap"`
	IsDHCP     bool   `json:"isDHCP"`
}
