package configs

import (
	"math/rand"
	"strconv"
)

var Servers = GenerateServers(5)
var timezonesArr = []string{"America/Guatemala", "Europe/Moscow", "Antarctica/Macquarie"}
var Timezones = func() string {
	return timezonesArr[rand.Intn(2)]
}

func GenerateServers(quantity int) []string {
	var servers []string
	for i := 0; i < quantity; i++ {
		servers = append(servers, "localhost:808"+strconv.Itoa(i))
	}
	return servers
}
