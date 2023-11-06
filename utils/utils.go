package utils

import (
	"fmt"
	"math/rand"
	"strings"
)

// getRandomIpAddress generates the ip address in IPv4 format
func GetRandomIpAddress() string {
	min := 1
	max := 255
	octets := [4]int{
		rand.Intn(max-min) + min,
		rand.Intn(max-min) + min,
		rand.Intn(max-min) + min,
		rand.Intn(max-min) + min,
	}

	return strings.Trim(strings.Replace(fmt.Sprint(octets), " ", ".", -1), "[]")
}
