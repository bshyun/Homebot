//command processing

package mycommand

import (
	mer "../../homebot/myerror"
	"../mydb"
	"../svcutil"
	"fmt"
	"strings"
)

const (
	TYPE_DEFAULT = iota
	TYPE_TEXT    = iota
	TYPE_VOICE
)

const (
	DEV_COMMAND    = "cmdlist"
	PUBLIC_COMMAND = "devlist"
)

type T_INST struct {
	Value interface{}
}

type Message interface{}
type commandproc struct {
	Command Message
	CmdType int
	cfg     *svcutil.Config
	db      *mydb.MyDB
}

type Processor interface {
	Process(ID int) (*string, error)
}

func NewCommand(cmd interface{}, cfg *svcutil.Config, db *mydb.MyDB) (*commandproc, error) {
	if cmd == nil || cfg == nil {
		return nil, fmt.Errorf(mer.ERROR_ARG_NIL)
	}

	switch cmd.(type) {
	case string:
		//fmt.Println("TEXT")
		return &commandproc{cmd, TYPE_TEXT, cfg, db}, nil
	case []byte:
		//fmt.Println("VOICE")
		return &commandproc{cmd, TYPE_VOICE, cfg, db}, nil
	}

	return nil, fmt.Errorf(mer.ERROR_INVALID_TYPE)
}

func (cp *commandproc) GetProc(userdevlist string) (inst *T_INST) {
	if cp.CmdType != TYPE_TEXT {
		return nil
	}
	if len(cp.Command.(string)) == 0 {
		return nil
	}

	inst = new(T_INST)
	cmd := cp.Command.(string)
	spl := strings.Split(cmd, " ")

	udls := strings.Split(userdevlist, " ")
	var udl string
	for _, va := range udls {
		udl += va + ":"
	}
	command := strings.ToLower(spl[0])
	command += ":"
	//if strings.Contains(userdevlist, command) { //user device
	if len(spl) > 1 && strings.Contains(udl, command) { //user device
		up, _ := NewUserProc(cp, cmd)
		if up == nil {
			return nil
		}

		inst.Value = up
		return
	} else { //general command
		switch command {
		case RegisterDevCmd:
			fallthrough
		case SignupCmd:
			rp, _ := NewRegisterProc(cp, cmd)
			if rp == nil {
				return nil
			}
			inst.Value = rp
			return
		case HelloCmd:
			ep, _ := NewEchoProc(cp, cmd)
			if ep == nil {
				return nil
			}

			inst.Value = ep
			return
		case ShowCmd:
			fallthrough
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
	}

	return nil
}
