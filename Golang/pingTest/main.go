package main

import (
	"log"
	"os"
	"os/exec"
	"strconv"

	"strings"
	"time"
)

func main() {
	target := os.Args[1]
	t := os.Args[2]
	timeDelay, err := strconv.Atoi(t)
	if err != nil {
		log.Println("ERR : ", err)
	}
	log.Println("Start")
	for {
		if !pingCheck(target) {
			log.Println("Stop")
			break
		}
		time.Sleep(time.Second * time.Duration(timeDelay))
	}
}

func pingCheck(target string) bool {
	out, err := exec.Command("ping", target, "-c 3").Output()
	if err != nil {
		log.Println("ERR : ", err)
	}
	log.Println(string(out))
	if strings.Contains(string(out), "Destination Host Unreachable") || strings.Contains(string(out), "errors") {
		log.Println("Down")
		return false
	}
	if strings.Contains(string(out), "0% packet loss") {
		n := strings.Index(string(out), "0% packet loss")
		if string(out[n-1]) == "0" && string(out[n-2]) == "1" {
			log.Println("Down")
			return false
		}
	} else {
		log.Println("Down")
		return false
	}

	log.Println("Alive")
	return true
}
