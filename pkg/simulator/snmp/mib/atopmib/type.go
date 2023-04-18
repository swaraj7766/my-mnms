package atopmib

import (
	"mnms/pkg/simulator/devicetype"
)

func OidType(mode devicetype.Simulator_type) uint {
	switch mode {
	case devicetype.EH7506:
		return 14
	case devicetype.EH7508:
		return 23
	case devicetype.EH7512:
		return 15
	case devicetype.EH7520, devicetype.EHG750x:
		return 21
	}
	return 21
}
