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
	timeFormat = "Monday, 02-Jan-06 15:04:05 MST"
	wg         sync.WaitGroup
)

func main() {
	// здесь на вход нужно принимать адреса серверов и читать данные в отдельных горутинах
	// там надо будет читать сообщения в беск цикле с в кром с помощью буфио читать до новой строки
	// полученную строку парсить как время и отправлять в канал с типом времени
	// в основной горутине читать время из цикла и печатать, разделяя время и таймзону
	fanninTime := make(chan time.Time)
	for _, server := range configs.Servers {
		wg.Add(1)
		go func(server string) {
			for {
				select {
				case read := <-client(server, &wg):
					parsedTime, err := time.Parse(timeFormat, read)
					if err != nil {
						log.Println(err)
					}
					fanninTime <- parsedTime
				}
			}
		}(server)
	}

	go func() {
		for {
			select {
			case read := <-fanninTime:
				hr, min, sec := read.Clock()
				fmt.Printf("Time: %v:%v:%v\n", hr, min, sec)
				fmt.Printf("Timezone: %v\n", read.Location())
			}
		}
	}()
	wg.Wait()
}

func client(address string, wg *sync.WaitGroup) <-chan string {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		log.Fatal(err)
	}
	fannoutTime := make(chan string)
	go func() {
		buffer := make([]byte, 1024)
		n, err := conn.Read(buffer)
		message := string(buffer[:n])

		//fmt.Printf("Server address: %v\n", conn.RemoteAddr().String())

		if n > 0 {
			fannoutTime <- message
		}
		if err != nil {
			log.Println(err)
			conn.Close()
			wg.Done()
			return
		}
	}()

	return fannoutTime
}
