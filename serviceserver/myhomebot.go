// service main

//ng Copyright 2015-2016 mrd0ll4r and contributors. All rights reserved.
// Use of this source code is governed by the MIT license, which can be found in
// the LICENSE file.
//56913263

package main

import (
	"fmt"
	"log"
	_ "net"
	"os"
	_ "runtime"
	"time"

	"../boilerplate"
	MC "../mycommand"
	"../mydb"
	mer "../myerror"
	_ "../mynet"
	"../svcutil"
	"github.com/mrd0ll4r/tbotapi"
	_ "strings"
)

const (
	ME             = int(56913263)
	DEFAULT_CONFIG = "../cfg/homebot.xml"
)

func main() {
	//read configuration
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

	/*
		//set network
		n, _ := mynet.NewMyNet("127.0.0.1")
		n.ListenAndRequest()
	*/

	db, errdb := mydb.NewMyDB("", &cfg.DSN_config)
	if errdb != nil {
		fmt.Print(errdb)
		return
	}
	defer db.Close()

	apiToken := cfg.Botkey

	updateFunc := func(update tbotapi.Update, api *tbotapi.TelegramBotAPI) {
		fmt.Println("[Message recved at", time.Now(), "]")

		switch update.Type() {
		case tbotapi.MessageUpdate:
			msg := update.Message
			typ := msg.Type()
			//fmt.Printf("msg type : %T\n", msg.Text)

			//errMsg := new(string)
			//fmt.Sprintf(*errMsg, "Ignoring non-text message")
			//errMsg := string("Sending Error: Ignoring non-text message")
			errMsg := string("No support message type")

			if typ != tbotapi.TextMessage {
				// Ignore non-text messages for now.
				fmt.Println("Ignoring non-text message : ", typ)
				//_, err := api.NewOutgoingMessage(tbotapi.NewChatRecipient(ME), errMsg).Send()
				_, err := api.NewOutgoingMessage(tbotapi.NewChatRecipient(msg.Chat.ID), errMsg).Send()
				if err != nil {
					fmt.Printf("Error: %s[%s]\n", err, errMsg)
				}
				return
			}
			fmt.Printf("<-%d, From:\t%d, Text: %s \n", msg.ID, msg.Chat.ID, *msg.Text)
			//fmt.Printf("<-%d, From:\t%s, Text: %s \n", msg.ID, msg.Chat, *msg.Text)

			go func(tapi *tbotapi.TelegramBotAPI, tmsg *tbotapi.Message) {
				//ToDo: find user on db. if not, signup the user.
				//alloc cmd and deter the msg type
				//if errid := db.SetUserID(msg.Chat.ID); errid != nil {
				if errid := db.SetUserID(tmsg.Chat.ID); errid != nil {
					fmt.Print(errid)
					return
				}
				//cp, err := MC.NewCommand(*msg.Text, cfg, db)
				cp, err := MC.NewCommand(*tmsg.Text, cfg, db)
				if err != nil {
					fmt.Print(err)
					return
				}

				//change voice msg to string
				if cp.CmdType == MC.TYPE_VOICE {
					cp.V2S()
					fmt.Printf("command: %s", cp.Command)
				}

				//va, errcl := db.GetMyDev(msg.Chat.ID)
				va, errcl := db.GetMyDev(tmsg.Chat.ID)
				//va, errcl := db.GetMyCommand(msg.Chat.ID)
				if errcl != nil {
					fmt.Println(errcl)
					return
				}

				ti := cp.GetProc(va)
				if ti == nil {
					fmt.Println(mer.ERROR_ARG_NIL)
					return
				}

				//msg.Text, err = ti.Value.(MC.Processor).Process(msg.Chat.ID)
				tmsg.Text, err = ti.Value.(MC.Processor).Process(tmsg.Chat.ID)
				if err != nil {
					fmt.Println(err)
					return
				}
				fmt.Println("outgoing msg: ", *tmsg.Text)

				// Now simply echo that back.
				outMsg, err := tapi.NewOutgoingMessage(tbotapi.NewRecipientFromChat(tmsg.Chat), *tmsg.Text).Send()
				//outMsg, err := api.NewOutgoingMessage(tbotapi.NewRecipientFromChat(msg.Chat), *msg.Text).Send()

				if err != nil {
					fmt.Printf("Error sending: %s\n", err)
					return
				}
				fmt.Printf("->%d, To:\t%d, Text: %s\n", outMsg.Message.ID, outMsg.Message.Chat.ID, *outMsg.Message.Text)
				//fmt.Printf("->%d, To:\t%s, Text: %s\n", outMsg.Message.ID, outMsg.Message.Chat, *outMsg.Message.Text)

				fmt.Println("Message replied at", time.Now())
			}(api, msg)

		case tbotapi.InlineQueryUpdate:
			fmt.Println("Ignoring received inline query: ", update.InlineQuery.Query)
		case tbotapi.ChosenInlineResultUpdate:
			fmt.Println("Ignoring chosen inline query result (ID): ", update.ChosenInlineResult.ID)
		default:
			fmt.Printf("Ignoring unknown Update type.")
		}
	}

	// Run the bot, this will block.
	boilerplate.RunBot(apiToken, updateFunc, "[homebot", "Serves things info]\n")
}
