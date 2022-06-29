package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
)

var (
	serverFile = "ports.txt"
	log        = logrus.New()
)

func echo(c net.Conn, timezone string) {
	// Тут надо реализовать бесконечный цикл который будет писать в конект время в нужной часовой зоне каждую секунду
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		log.Error(err)
		return
	}
	for {
		t := time.Now().In(loc)
		fmt.Fprintln(c, t.Format("Monday, 02-Jan-06 15:04:05 MST"))
		time.Sleep(time.Second)
	}
}

func main() {
	log.Out = os.Stdout
	address, timezone := generateServer()
	closeHandle(address)
	server(address, timezone)

	fmt.Println(address, timezone)
}

func server(address, timezone string) {
	// нужно вынести в настройки порт и часовой пояс
	listener, err := net.Listen("tcp", address)
	log.Infof("Address: %v Timezone: %v", listener.Addr().String(), timezone)
	defer listener.Close()
	if err != nil {
		log.Error(err)
		return
	}
	for {
		conn, err := listener.Accept()

		if err != nil {
			log.Error(err)
			return
		}
		go echo(conn, timezone)
	}
}

func generateServer() (address, timezone string) {
	serverPtr := flag.String("s", "localhost:8080", "server address")
	timezonePtr := flag.String("t", "Europe/Moscow", "servers timezone")
	flag.Parse()
	timezone = *timezonePtr
	address = *serverPtr

	file, err := os.OpenFile(serverFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Error(err)
		return
	}
	defer file.Close()
	if _, err := file.WriteString(timezone + " " + address + "\n"); err != nil {
		log.Error(err)
		return
	}

	return
}

func closeHandle(address string) {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		deleteServer(address)
		log.Info("Shutting down server...")
		os.Exit(0)
	}()
}

func deleteServer(address string) {
	str := fmt.Sprintf("sed -i '/%v/d' ports.txt", address)
	cmd := exec.Command("bash", "-c", str)
	stdout, err := cmd.Output()
	if err != nil {
		log.Error(err)
		return
	}
	fmt.Println(string(stdout))
}
