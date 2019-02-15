//*** homebot ***
//themometer device simulator.

package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"time"
)

const (
	SERIAL_NUMBER       = string("26f7262e-b4df-42dd-b95a-d21df9\r\n")
	SERIAL_NUMBER2      = string("fd1a7dac-4cf6-4f53-bb85-16c3ec665528")
	INIT_SERVER_ADDR    = string("192.168.39.1:12126")
	SERVICE_SERVER_ADDR = string("192.168.50.173:12125")
)

func main() {
	args := os.Args

	for idx, arg := range args {
		fmt.Println(idx, arg)
	}

	var self_addr net.UDPAddr

	/*
		////////////////////////////////////////////////////////////
		//setup mod ip to init server
		fmt.Println("==================================\nRegister device to homebot")
		var self_addr net.UDPAddr
		var sess_addr net.TCPAddr

		self_addr.IP = net.ParseIP("")
		self_addr.Port = 0

		aip, aport, _ := net.SplitHostPort(INIT_SERVER_ADDR)
		sess_addr.IP = net.ParseIP(aip)
		sess_addr.Port, _ = strconv.Atoi(aport)

		//init_conn, erri := net.ListenUDP("udp", &self_addr)
		init_conn, erri := net.DialTCP("tcp", nil, &sess_addr)
		//init_conn, erri := net.DialUDP("udp", &self_addr, &sess_addr)
		if erri != nil {
			fmt.Println("dial tcp", erri)
			return
		}
		defer init_conn.Close()

		wlen, errw := init_conn.Write([]byte(SERIAL_NUMBER))
		if errw != nil {
			fmt.Println("send error: ", errw)
		}
		fmt.Println("sess_addr: ", sess_addr.String())
		fmt.Println("self_addr: ", self_addr.String())
		fmt.Println("send len: ", wlen)

		p := make([]byte, 1000)
		plen, err := init_conn.Read(p)
		if err != nil {
			log.Fatalln(err)
		}
		if plen > 0 {
			fmt.Println("recv: ", string(p)) //[:plen]))
			fmt.Println("Registered!")
		}
	*/

	/////////////////////////////////////////////////////
	//var self_addr net.UDPAddr
	fmt.Println("=================================\nWaiting to request from homebot...")
	_, aport1, _ := net.SplitHostPort(SERVICE_SERVER_ADDR)
	//aip1, aport1, _ := net.SplitHostPort(SERVICE_SERVER_ADDR)
	self_addr.IP = net.ParseIP("")
	self_addr.Port, _ = strconv.Atoi(aport1)

	self_conn, err := net.ListenUDP("udp", &self_addr)
	if err != nil {
		log.Fatal("udp listen error!", err)
	}
	fmt.Println("UDP Listen:", self_conn.LocalAddr().String(), "\n")
	defer self_conn.Close()

	for {
		self_conn.SetReadDeadline(time.Now().Add(time.Second * 10))

		p := make([]byte, 1000)
		plen, ua, err := self_conn.ReadFromUDP(p)
		if err != nil {
			oerr, _ := err.(*net.OpError)
			if !oerr.Timeout() {
				fmt.Println(err)
			} else {
				continue
			}
		}
		if plen > 0 {
			//fmt.Println("udp len: ", plen)
			//fmt.Println("udp recv from: ", ua)
			fmt.Println("Command: ", string(p)) //[:plen-2]))

			reply := fmt.Sprintf("this is simulator reply!!\n")
			_, errw := self_conn.WriteToUDP([]byte(reply), ua)
			if errw == nil {
				fmt.Println("Reply msg:", reply)
			}
		}
	}

	fmt.Println("exit!")
}
