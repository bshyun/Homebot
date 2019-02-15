//user registered dev processing

package mycommand

import (
	mer "../../homebot/myerror"
	"../mynet"
	_ "../svcutil"
	"fmt"
	_ "net"
	"strconv"
	"strings"
	"time"
)

type UserProc struct {
	commandproc
	message string
	//clist   string
}

func NewUserProc(cp *commandproc, msg string) (*UserProc, error) {
	//func NewUserProc(cp *commandproc, msg string, cmdlist string) (*UserProc, error) {
	if len(msg) == 0 {
		return nil, fmt.Errorf(mer.ERROR_ARG_NIL)
	}

	return &UserProc{*cp, msg}, nil
	//return &UserProc{*cp, msg, cmdlist}, nil
}

func (up *UserProc) Process(ID int) (*string, error) {
	//get client ip addr
	//make userid
	userid := fmt.Sprintf("%s_%d", "telegram", ID)

	spl := strings.SplitN(up.message, " ", 2)

	client_addr, err := up.db.GetMyIP(userid, spl[0])
	if err != nil {
		return nil, fmt.Errorf("Not found dev IP")
	}

	//public command
	switch strings.ToLower(spl[1]) {
	case "info":
		s := fmt.Sprint("IP:", client_addr)
		return &s, nil
	}

	//set UDP
	rport, _ := strconv.Atoi(up.cfg.Client_config.Port)
	n := mynet.NewMyNet("", 0, client_addr, rport)

	//bind avail port
	_, err1 := n.MakeUDPConnAndBindLocalPort()
	if err1 != nil {
		return nil, fmt.Errorf("Failed to set network")
	}
	defer n.MyConn.Close()

	slen, err1 := n.MyConn.WriteToUDP([]byte(spl[1]), &n.C)
	if err1 != nil {
		fmt.Println(err1, slen)
	}
	fmt.Println("Request to", n.C, ":", spl[1])

	p := make([]byte, 1000)
	for {
		n.MyConn.SetReadDeadline(time.Now().Add(time.Second * 30))

		plen, _, err2 := n.MyConn.ReadFromUDP(p)
		if err2 != nil {
			fmt.Println("ReadFromUDP ERROR!! :", err2)

			r := fmt.Sprintf("No response: %s:%d", client_addr, rport)
			p = []byte(r)
			break
		}
		if plen > 0 {
			fmt.Printf("recved from %s > %d:%s", spl[0], plen, string(p))
			break
		}
	}

	reply := new(string)
	for _, v := range spl {
		*reply += v + " "
	}
	*reply += "> " + string(p)
	//*reply += "> " + string(p[:l-1])

	return reply, nil
}
