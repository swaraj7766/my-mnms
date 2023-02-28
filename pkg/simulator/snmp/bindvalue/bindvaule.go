package bindValue

import (
	"strings"

	"mnms/pkg/simulator/devicetype"
)

func NewBindValue(model *string) *Value {
	return &Value{model: model}
}

type Value struct {
	ker     *string
	ap      *string
	model   *string
	system  *string
	ip      *string
	mac     *string
	mask    *string
	gateway *string
	user    *string
	pwd     *string
}

func (b *Value) Port() int {
	v, err := devicetype.ParsingType(*b.model)
	if err != nil {
		return 12
	}
	return v.Port()
}

// BindUser bind pointer
func (b *Value) BindUser(name *string) {
	b.user = name
}

// BindUser bind pointer
func (b *Value) BindPwd(pwd *string) {
	b.pwd = pwd
}

// BindModel bind pointer
func (b *Value) BindKernel(name *string) {
	b.ker = name
}

func (b *Value) GetKernel() string {
	return *b.ker
}

// BindModel bind pointer
func (b *Value) BindAp(name *string) {
	b.ap = name
}

func (b *Value) GetAp() string {
	return *b.ap
}

// BindSystem bind pointer
func (b *Value) GetModel() string {
	return *b.model
}

// BindSystem bind pointer
func (b *Value) BindSystem(name *string) {
	b.system = name
}

func (b *Value) BindIP(ip *string) {
	b.ip = ip
}

func (b *Value) BindMask(mask *string) {
	b.mask = mask
}

func (b *Value) BindGateWay(gateway *string) {
	b.gateway = gateway
}

func (b *Value) BindMac(mac *string) {
	b.mac = mac
}

func (b *Value) SetSystem(name string) {
	*b.system = name
}

func (b *Value) GetSystem() string {
	return *b.system
}

func (b *Value) SetIP(ip string) {
	*b.ip = ip
}

func (b *Value) GetIP() string {
	return *b.ip
}

func (b *Value) SetMac(mac string) {
	*b.mac = mac
}

func (b *Value) GetMac() string {
	v := *b.mac
	m := strings.Replace(v, "-", "", -1)
	mac := strings.Replace(m, ":", "", -1)
	return mac
}

func (b *Value) SetMask(mask string) {
	*b.mask = mask
}
func (b *Value) GetMask() string {
	return *b.mask
}

func (b *Value) SetGateWay(gateway string) {
	*b.gateway = gateway
}
func (b *Value) GetGateWay() string {
	return *b.gateway
}

func (b *Value) SetUser(user string) {
	*b.user = user
}
func (b *Value) GetUser() string {
	return *b.user
}

func (b *Value) SetPwd(pwd string) {
	*b.pwd = pwd
}
func (b *Value) GetPwd() string {
	return *b.pwd
}

func (b *Value) SetModel(model string) {
	*b.model = model
}
