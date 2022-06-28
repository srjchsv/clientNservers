package main

import (
	"bufio"
	"flag"
	"net"
	"os"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

var (
	timeFormat = "Monday, 02-Jan-06 15:04:05 MST"
	log        = logrus.New()
)

func main() {
	// здесь на вход нужно принимать адреса серверов и читать данные в отдельных горутинах
	// там надо будет читать сообщения в беск цикле с в кром с помощью буфио читать до новой строки
	// полученную строку парсить как время и отправлять в канал с типом времени
	// в основной горутине читать время из цикла и печатать, разделяя время и таймзону
	log.Out = os.Stdout
	fanninTime := make(chan time.Time)
	addressCh := make(chan string)

	for _, port := range ports() {
		go func(server string) {
			write := client(server)
			for v := range write {
				fanninTime <- v
				addressCh <- server
			}
		}(port)
	}

	for read := range fanninTime {
		hr, min, sec := read.Clock()
		log.Infof("Server: %v Time: %v:%v:%v Timezone: %v\n", <-addressCh, hr, min, sec, read.Location())

	}

}

func client(address string) <-chan time.Time {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		log.Error(err)
	}

	fannoutTime := make(chan time.Time)
	go func() {
		for {
			message, err := bufio.NewReader(conn).ReadString('\n')
			if err != nil {
				log.Error(err)
				conn.Close()
				return
			}
			message = strings.Trim(message, "\n")
			parsedTime, err := time.Parse(timeFormat, message)
			if err != nil {
				log.Error(err)
			}
			fannoutTime <- parsedTime
		}

	}()

	return fannoutTime

}

func ports() []string {
	var ports []string

	portsPtr := flag.String("p", "ports.txt", "ports addresses")
	flag.Parse()
	portsFile := *portsPtr

	file, err := os.Open(portsFile)
	if err != nil {
		log.Error(err)
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		split := strings.Split(scanner.Text(), " ")
		ports = append(ports, split[1])
	}

	return ports
}
