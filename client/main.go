package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
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

	for _, port := range ports() {
		go func(server string) {
			write := client(server)
			for {
				send := <-write
				fanninTime <- send
			}
		}(port)
	}

	for {
		read := <-fanninTime
		hr, min, sec := read.Clock()
		fmt.Printf("Time: %v:%v:%v\n", hr, min, sec)
		fmt.Printf("Timezone: %v\n", read.Location())
	}

}

func client(address string) <-chan time.Time {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		log.Fatal(err)
	}

	fannoutTime := make(chan time.Time)
	go func() {
		for {
			message, err := bufio.NewReader(conn).ReadString('\n')
			if err != nil {
				log.Println(err)
				conn.Close()
				return
			}
			message = strings.Trim(message, "\n")
			parsedTime, err := time.Parse(timeFormat, message)
			fannoutTime <- parsedTime
		}

	}()

	return fannoutTime

}

func ports() []string {
	var ports []string
	file, err := os.Open("ports.txt")
	if err != nil {
		log.Println(err)
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		split := strings.Split(scanner.Text(), " ")
		ports = append(ports, split[1])
	}

	return ports
}
