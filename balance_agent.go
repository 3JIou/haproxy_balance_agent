package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"strconv"
	"regexp"
)

// Configuration data
type Configuration struct {
	Host string
	Port string
	Type string
}

var oldStatus int

func main() {
	config := Configuration{}
	config.SetConfig("config.json")
	listen, err := net.Listen(config.Type, config.Host+":"+config.Port)
	if err != nil {
		log.Fatal(err)
	}
	defer listen.Close()
	log.Println("Listening on " + config.Host + ":" + config.Port)
	for {
		connect, err := listen.Accept()
		if err != nil {
			log.Println("Error accespting: ", err)
		}
		go StatusRequest(connect)
	}
}

// SetConfig getting configuration data from file
func (config *Configuration) SetConfig(ConfigFileName string) {
	if _, err := os.Stat(ConfigFileName); err == nil {
		file, _ := os.Open(ConfigFileName)
		decoder := json.NewDecoder(file)
		err := decoder.Decode(&config)
		if err != nil {
			fmt.Println("error:", err)
		}
	} else {
		log.Fatal("Configuration file not found!")
	}
}

// StatusRequest return answer based on postfix queue len
func StatusRequest(conn net.Conn) {
	buffer := make([]byte, 1)
	_, err := conn.Read(buffer)
	newStatus := GetStatus()
	if err != nil {
		log.Fatal(err)
	}
	println(oldStatus)
	println(newStatus)
    switch {
    case oldStatus > newStatus:
        conn.Write([]byte("UP"))
    case oldStatus < newStatus:
        conn.Write([]byte("DOWN"))
	default:
		conn.Write([]byte("OK"))
    }
	conn.Close()
	oldStatus = newStatus
}

// GetStatus return result execution command
func GetStatus() int {
	// find /var/spool/postfix/{deferred,active,maildrop}/ -type f | wc -l
	// mailq
	out, err := exec.Command("cat", "test").Output()
	if err != nil {
		println(err)
	}
	result, _ := strconv.Atoi(regexp.MustCompile(`^[-+]?\d+`).FindString(string(out)))
	return result
}