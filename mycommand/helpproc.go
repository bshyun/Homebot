//help command processing

package mycommand

import (
	mer "../../homebot/myerror"
	_ "../mydb"
	_ "../svcutil"
	"fmt"
	"strings"
	_ "time"
)

const (
	HelpCmd = "/help"
	ShowCmd = "/show"
)

const (
	_USAGE = "*** homebot USAGE ***\n/show devlist : device list\ndev_name devinfo : get device info"
)

type HelpProc struct {
	commandproc
	message string
	//clist   string
}

func NewHelpProc(cp *commandproc, msg string) (*HelpProc, error) {
	if len(msg) == 0 {
		return nil, fmt.Errorf(mer.ERROR_ARG_NIL)
	}

	return &HelpProc{*cp, msg}, nil
}

func (hp *HelpProc) Process(ID int) (*string, error) {
	spl := strings.Split(hp.message, " ")

	reply := new(string)
	cmd := strings.ToLower(spl[0])
	switch cmd {
	case HelpCmd:
		*reply = _USAGE
	case ShowCmd:
		if len(spl) > 1 {
			switch spl[1] {
			case "devlist":
				devlist, errd := hp.db.GetMyDev(-1)
				if errd != nil {
					*reply = _USAGE
				} else {
					r := strings.Split(devlist, " ")
					for _, va := range r {
						*reply += va + "\n"
					}
					//*reply = fmt.Sprintf("%s", devlist)
				}
			default:
				*reply = _USAGE

			}
		} else {
			*reply = fmt.Sprintf("Not found command or registered dev.\n\n%s", _USAGE)
		}
	default:
		*reply = _USAGE
	}

	return reply, nil
}
