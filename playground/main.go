package main

import (
	"log"

	u "illumio.com/iplist/utils"
)

func main() {
	log.Printf("Random IP address -> %s", u.GetRandomIpAddress())
}
