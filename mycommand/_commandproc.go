//command processing

package mycommand

import (
	mer "../../homebot/myerror"
	"../svcutil"
	"fmt"
	"strings"
)

const (
	TYPE_DEFAULT = iota
	TYPE_TEXT    = iota
	TYPE_VOICE
)

type T_INST struct {
	Value interface{}
}

type Message interface{}
type commandproc struct {
	Command Message
	CmdType int
	cfg     *svcutil.Config
}

type Processor interface {
	Process(ID int) (*string, error)
}

func NewCommand(cmd interface{}, cfg *svcutil.Config) (*commandproc, error) {
	if cmd == nil || cfg == nil {
		return nil, fmt.Errorf(mer.ERROR_ARG_NIL)
	}

	switch cmd.(type) {
	case string:
		//fmt.Println("TEXT")
		return &commandproc{cmd, TYPE_TEXT, cfg}, nil
	case []byte:
		//fmt.Println("VOICE")
		return &commandproc{cmd, TYPE_VOICE, cfg}, nil
	}

	return nil, fmt.Errorf(mer.ERROR_INVALID_TYPE)
}

func (cp *commandproc) GetProc() (inst *T_INST) {
	if cp.CmdType != TYPE_TEXT {
		return nil
	}
	if len(cp.Command.(string)) == 0 {
		return nil
	}

	inst = new(T_INST)
	cmd := cp.Command.(string)
	spl := strings.Split(cmd, " ")
	//fmt.Println(cmd, spl)
	switch spl[0] {
	case HelloCmd:
		ep, _ := NewEchoProc(cp, cmd)
		if ep == nil {
			return nil
		}

		inst.Value = ep
		return
	case CurTemp:
		fallthrough
	case TodayTemp:
		fallthrough
	case YesterdayTemp:
		fallthrough
	case ThisweekTemp:
		tp, _ := NewTemperProc(cp, cmd)
		if tp == nil {
			return nil
		}

		inst.Value = tp
		return
	case HelpCmd:
		fallthrough
	default:
		hp, _ := NewHelpProc(cp, cmd)
		if hp == nil {
			return nil
		}

		inst.Value = hp
		return
	}
	return nil
}
