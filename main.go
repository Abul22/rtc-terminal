package main

import (
	"net"
	"net/url"
	"os"
	"fmt"
	"log"
	"flag"
	"os/signal"
	"time"
	"path/filepath"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
    ini "github.com/pierrec/go-ini"
)

type Config struct {
	Uuid     string  `ini:"uuid,identify"`
}

type Session struct {
	Type string   `json:"type"`
	Sdp  string   `json:"sdp"`	
	Error string  `json:"error"`
}

func connect(closeSig <-chan struct{}, query string) *websocket.Conn {
	var u url.URL
	var ws *websocket.Conn
	var err error
	
	u = url.URL{Scheme: "wss", Host: "sqs.io", Path: "/signal", RawQuery: query}
	ws, _, err = websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Println(err)
	}
	
	go func(){
		for{
			select {
			case <-closeSig:
				err := ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
				if err != nil {
					log.Println("write close:", err)
				}
				return
			}
		}
	}()
	return ws
}

func main() {
	log.SetFlags(0)
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	
	var executablePath string
	var conf Config
		
	ex, err := os.Executable()
    	check(err)
	executablePath = filepath.Dir(ex)
	
	uuidKey := flag.String("uuid", "", "uuid remote device")	
	getkey := flag.Bool("getkey", false, "display uuid remote device")
	
	flag.Parse()
	
	loadConf := func(){
		file, err := os.Open(executablePath + "/config.ini") 
		if err == nil {
			err = ini.Decode(file, &conf)
			check(err)
		}
		if os.IsNotExist(err) && *uuidKey == "" {
			log.Println("File config.ini not found, using option -uuid= <UUID key remote device>")
			os.Exit(0)
		} 
		defer file.Close()	
	
	}
	saveConf := func(){
		file, err := os.Create(executablePath + "/config.ini")
		check(err)
		defer file.Close()
		err = ini.Encode(file, &conf)
		check(err)
	} 
		
	loadConf()
			
	if *getkey { 
		fmt.Println("uuid:", conf.Uuid)
	}
	
	if *uuidKey != "" {
		conf.Uuid = *uuidKey
		_, err = uuid.Parse(conf.Uuid)
		if err != nil {
			fmt.Println(err)
			os.Exit(0)
		}
		saveConf()
	}
	
	if conf.Uuid == "" && *uuidKey == ""{
		log.Println("UUID not set, using option -uuid = <UUID key remote device>")
		os.Exit(0)
	} 
	
	done := make(chan struct{})
	closeSig := make(chan struct{})
	
	go func(){ for {
		select {
		case <-done:
			 return
		case <-interrupt:
			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			fmt.Println("interrupt close")
			close(closeSig)
			select {
				case <-done:
				case <-time.After(time.Second):
			}
			os.Exit(0)
		}
	}
	}()
	
	port := 2222
	listen, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil {
		log.Println(err)
	}
	
	go func(){
		ssh, err := listen.Accept()
		if err != nil {
			log.Println(err)
		}
		query := "localUser=" + uuid.New().String() + "&remoteUser=" + conf.Uuid
		ws := connect(closeSig, query)
		defer ws.Close()
		err = hub(ws, ssh)
		if err != nil {
			fmt.Println(err)
			os.Exit(0)
		}
		select {}
		
	}()
	Terminal(port)
}

func check(e error) {
    if e != nil {
        log.Println(e)
    }
}
