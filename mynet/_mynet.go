package mynet

import (
	"../myerror"
	"fmt"
	"net"
)

type MyNet struct {
	MyAddr net.TCPAddr
	MyConn *net.TCPConn
	C      net.TCPAddr
}

func NewMyNet(self_addr string, self_port int, client_addr string, client_port int) (*MyNet, error) {
	var n MyNet
	n.MyAddr.IP = net.ParseIP(self_addr)
	n.MyAddr.Port = self_port

	n.C.IP = net.ParseIP(client_addr)
	n.C.Port = client_port

	return &n, nil
}

func (n *MyNet) ListenAndBind() (string, error) {
	if n == nil {
		return "", fmt.Errorf(myerror.ERROR_ARG_NIL)
	}

	//listen
	var err error
	n.MyConn, err = net.ListenUDP("udp", &n.MyAddr)
	if err != nil {
		return "", err

	}

	return n.MyConn.LocalAddr().String(), nil
}
