package mynet

import (
	"../myerror"
	"fmt"
	"net"
)

//MyAddr is local UDP address
//C is remote UDP address
//MyConn is UDPConn.
type MyNet struct {
	MyAddr net.UDPAddr
	MyConn *net.UDPConn
	C      net.UDPAddr
}

//MyNet struct construtor.
//self_addr, self_port is local UDP address.
//client_addr, client_port is remote UDP address.
func NewMyNet(self_addr string, self_port int, client_addr string, client_port int) *MyNet {
	var n MyNet
	n.MyAddr.IP = net.ParseIP(self_addr)
	n.MyAddr.Port = self_port

	n.C.IP = net.ParseIP(client_addr)
	n.C.Port = client_port

	return &n
}

//bind available local UDP port and make UDPConn.
//returns binded local address.
func (n *MyNet) MakeUDPConnAndBindLocalPort() (string, error) {
	if n == nil {
		return "", fmt.Errorf(myerror.ERROR_ARG_NIL)
	}

	var err error
	n.MyConn, err = net.ListenUDP("udp", &n.MyAddr)
	if err != nil {
		return "", err
	}

	return n.MyConn.LocalAddr().String(), nil
}
