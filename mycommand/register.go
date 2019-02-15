//register user and dev

package mycommand

import (
	mer "../../homebot/myerror"
	_ "../svcutil"
	"fmt"
	"strings"
)

const (
	SignupCmd      = "/signup"
	RegisterDevCmd = "/register"
)

type RegisterProc struct {
	commandproc
	message string
	//clist   string
}

func NewRegisterProc(cp *commandproc, msg string) (*RegisterProc, error) {
	if len(msg) == 0 {
		return nil, fmt.Errorf(mer.ERROR_ARG_NIL)
	}

	return &RegisterProc{*cp, msg}, nil
}

func (rp *RegisterProc) Process(ID int) (*string, error) {
	reply := new(string)
	//*reply = fmt.Sprint("Register process: ", rp.message)

	spl := strings.Split(rp.message, " ")
	//cmdlen := len(spl)

	command := strings.ToLower(spl[0])
	switch command {
	case SignupCmd:
	case RegisterDevCmd:
	}

	return reply, nil
}
