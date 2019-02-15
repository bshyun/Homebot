package main

import (
	"../mydb"
	"../svcutil"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"time"
)

const DEFAULT_CONFIG = string("../cfg/homebot.xml")

func main() {
	fmt.Println("[Session Server]")
	fmt.Println(time.Now(), "\n")

	//read confing
	var ConfigFileName string
	args := os.Args
	if len(args) < 2 {
		fmt.Printf("No configuration file. Default config: %s\n", DEFAULT_CONFIG)
		ConfigFileName = DEFAULT_CONFIG
	} else {
		ConfigFileName = args[1]
	}

	cfg, cfgerr := svcutil.OpenConfig(ConfigFileName)
	if cfgerr != nil {
		log.Fatal(cfgerr)
	}

	//setup db
	db, errdb := mydb.NewMyDB(&cfg.DSN_config)
	if errdb != nil {
		log.Fatal(errdb)
	}
	defer db.Close()

	//setup listener
	lport, _ := strconv.Atoi(cfg.InitServer_config.Port)
	listen, errl := net.ListenTCP("tcp", &net.TCPAddr{net.IP(net.IPv4zero), lport, ""})
	if errl != nil {
		log.Fatal(errl)
	}
	defer listen.Close()
	fmt.Println("Service Address:", listen.Addr(), "\n")

	for {
		//accept connections
		conn, erra := listen.AcceptTCP()
		if erra != nil {
			fmt.Println(erra)
			continue
		}

		//launch with new connection
		go func(tc *net.TCPConn) {
			defer tc.Close()

			errrd := tc.SetReadDeadline(time.Now().Add(time.Second * 5))
			if errrd != nil {
				fmt.Println(errrd)
			}

			raddr := tc.RemoteAddr()
			fmt.Println("Connection accepted! Remote addr:", raddr)

			//get mod s/n
			p := make([]byte, 1000)
			rlen, errr := tc.Read(p)
			if errr != nil {
				fmt.Println(errr)
			}

			sn := p[:rlen-2]
			//search s/n on db.
			if !db.FindMe(string(sn)) {
				fmt.Println("Not registered mod s/n :", string(sn))
			} else {
				fmt.Println("Valid mod s/n :", string(sn))
				erru := db.UpdateModIP(raddr, string(sn))
				if erru != nil {
					fmt.Println("Update mod ip failed", erru)
				}
				fmt.Printf("Updated mod ip: %s|%s\n\n", raddr, sn)
			}

			rip, rport, _ := net.SplitHostPort(raddr.String())
			smsg := fmt.Sprintf("%s:%s", rip, rport)
			_, errs := tc.Write([]byte(smsg))
			if errs != nil {
				fmt.Println(errs)
			}

			//update addr.
			//but the mod is had who?
		}(conn)
	}

}
