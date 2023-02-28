package devicetype

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"
)

type Simulator_type int

const split = "_"

var (
	stringMapSimulator = map[string]Simulator_type{
		"EH7506":  EH7506,
		"EH7508":  EH7508,
		"EH7512":  EH7512,
		"EH7520":  EH7520,
		"EHG750x": EHG750x,
	}
	simulatorMapSring = map[Simulator_type]string{
		EH7506:  "EH7506",
		EH7508:  "EH7508",
		EH7512:  "EH7512",
		EH7520:  "EH7520",
		EHG750x: "EHG750x",
	}
	ArraySimulator = []string{EH7506.String(), EH7508.String(), EH7512.String(), EH7520.String(), EHG750x.String()}
)

const (
	EH7506 Simulator_type = iota
	EH7508
	EH7512
	EH7520
	EHG750x
	Unknow
)

func (s Simulator_type) Port() int {
	switch s {
	case EH7506:
		return 6
	case EH7508:
		return 8
	case EH7512:
		return 12
	case EH7520, EHG750x:
		return 20
	default:
		return 20
	}
}

func (s Simulator_type) String() string {
	switch s {
	case EH7506:
		return "EH7506"
	case EH7508:
		return "EH7508"
	case EH7512:
		return "EH7512"
	case EH7520:
		return "EH7520"
	case EHG750x:
		return "EHG750x"
	}
	return ""
}

func ParseString(str string) (Simulator_type, bool) {
	c, ok := stringMapSimulator[str]
	return c, ok
}

func (s Simulator_type) IsValid() bool {
	_, ok := simulatorMapSring[s]
	return ok
}

// GetRandomModelAp  get Random model name and ap info
func GetRandomModelAp() (string, string) {
	array := []Simulator_type{EH7506, EH7508, EH7512, EH7520, EHG750x}
	rand.Seed(time.Now().Unix())
	v := array[rand.Intn(len(array))]
	return fmt.Sprintf("Simu%v%v", split, v), fmt.Sprintf("Simu%v_atop device", v)

}

// GetModelAp  get model  name and ap info by value
func GetModelAp(name Simulator_type) (string, string) {
	if name.IsValid() {
		return fmt.Sprintf("Simu%v%v", split, name), fmt.Sprintf("Simu%v_atop device", name)
	} else {
		log.Fatalf("Simulator_type:%v error", name)
		return "", "'"
	}
}

// ParsingType parsing device name from model name
func ParsingType(str string) (Simulator_type, error) {
	for k, v := range simulatorMapSring {
		b := strings.Contains(str, v)
		if b {
			return k, nil
		}
	}

	/*
		str1 := strings.Split(str, split)
		if len(str1) != 2 {
			return Unknow, errors.New("format error")
		}
		r, ok := ParseString(str1[1])
		if ok {
			return r, nil
		}*/
	return Unknow, errors.New("format error")
}
