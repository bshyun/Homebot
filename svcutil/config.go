package svcutil

import (
	"encoding/xml"
	"fmt"
	"os"
)

type DSN_config struct {
	Driver string `xml:"dsn>driver"`
	DSN    string `xml:"dsn>datasourcename"`
	User   string `xml:"dsn>user"`
	Passwd string `xml:"dsn>passwd"`
	Net    string `xml:"dsn>net"`
	IP     string `xml:"dsn>ip"`
	Port   string `xml:"dsn>port"`
	DB     string `xml:"dsn>database"`
}

type Admin_config struct {
	Id       string `xml:"admin>id"`
	Username string `xml:"admin>username"`
}

type Client_config struct {
	Port  string `xml:"client>port"`
	Limit string `xml:"client>limit"`
}

type InitServer_config struct {
	Port  string `xml:"initserver>port"`
	Limit string `xml:"initserver>limit"`
}

type Config struct {
	XMLName xml.Name `xml:"homebot"`
	Basedir string   `xml:"basedir"`
	Selfip  string   `xml:"selfip"`
	Botkey  string   `xml:"botkey"`
	Admin_config
	Client_config
	InitServer_config
	DSN_config
}

func OpenConfig(ConfigFileName string) (*Config, error) {
	fi, err := os.Open(ConfigFileName)
	if err != nil {
		return nil, err
	}
	defer fi.Close()

	FI, _ := os.Stat(ConfigFileName)
	FileLength := FI.Size()

	configdata := make([]byte, FileLength)
	cnt, err := fi.Read(configdata)
	if cnt == 0 || err != nil {
		return nil, fmt.Errorf("Config file has no data: %s", ConfigFileName)
	}

	c := Config{}
	//c := Config{Selfip: "127.0.0.1", Client_config{Port: "12120"}}

	err = xml.Unmarshal([]byte(configdata), &c)
	if err != nil {
		return nil, err
	}

	if len(c.Selfip) == 0 {
		return nil, fmt.Errorf("%s: No selfip", ConfigFileName)
	} else if len(c.Basedir) == 0 {
		return nil, fmt.Errorf("%s: No basedir", ConfigFileName)
	} else if len(c.Botkey) == 0 {
		return nil, fmt.Errorf("%s: No Botkey", ConfigFileName)
	} else if len(c.Admin_config.Id) == 0 {
		return nil, fmt.Errorf("%s: No admin configuration", ConfigFileName)
	} else if len(c.Client_config.Port) == 0 {
		return nil, fmt.Errorf("%s: No client configuration", ConfigFileName)
	} else if len(c.InitServer_config.Port) == 0 {
		return nil, fmt.Errorf("%s: No initserver configuration", ConfigFileName)
	}

	return &c, nil
}
