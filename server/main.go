package main

import (
	"fmt"
	"github.com/srjchsv/clientNservers/configs"
	"log"
	"net"
	"sync"
	"time"
)

var (
	wg sync.WaitGroup
)

func echo(c net.Conn) {
	// Тут надо реализовать бесконечный цикл который будет писать в конект время в нужной часовой зоне каждую секунду
	loc, err := time.LoadLocation(configs.Timezones())
	if err != nil {
		panic(err)
	}
	for {
		t := time.Now().In(loc)
		c.Write([]byte(t.Format("Monday, 02-Jan-06 15:04:05 MST")))
		time.Sleep(time.Second)
	}
}

func main() {
	for _, address := range configs.Servers {
		wg.Add(1)
		go server(address)
	}
	wg.Wait()
}

func server(address string) {
	// нужно вынести в настройки порт и часовой пояс
	listener, err := net.Listen("tcp", address)
	fmt.Println(listener.Addr().String())
	defer listener.Close()
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := listener.Accept()

		if err != nil {
			log.Print(err)
			continue
		}
		go echo(conn)
	}
}
