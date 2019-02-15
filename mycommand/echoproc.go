//temperature command processing

package mycommand

import (
	mer "../../homebot/myerror"
	_ "../svcutil"
	"fmt"
)

const (
	HelloCmd = "/hello"
)

type EchoProc struct {
	commandproc
	message string
	//clist   string
}

func NewEchoProc(cp *commandproc, msg string) (*EchoProc, error) {
	//func NewEchoProc(cp *commandproc, msg string, cmdlist string) (*EchoProc, error) {
	if len(msg) == 0 {
		return nil, fmt.Errorf(mer.ERROR_ARG_NIL)
	}

	return &EchoProc{*cp, msg}, nil
	//return &EchoProc{*cp, msg, cmdlist}, nil
}

func (ep *EchoProc) Process(ID int) (*string, error) {
	reply := new(string)
	*reply = fmt.Sprint("echo process: ", ep.message)

	return reply, nil
}
