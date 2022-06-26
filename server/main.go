package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var (
	wg         sync.WaitGroup
	serverFile = "ports.txt"
)

func echo(c net.Conn, timezone string) {
	// Тут надо реализовать бесконечный цикл который будет писать в конект время в нужной часовой зоне каждую секунду
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		panic(err)
	}
	for {
		t := time.Now().In(loc)
		fmt.Fprintln(c, t.Format("Monday, 02-Jan-06 15:04:05 MST"))
		time.Sleep(time.Second)
	}
}

func main() {
	address, timezone := generateServer()
	closeHandle(address)
	server(address, timezone)

	fmt.Println(address, timezone)
}

func server(address, timezone string) {
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
		log.Println(err)
	}
	defer file.Close()
	if _, err := file.WriteString(timezone + " " + address + "\n"); err != nil {
		log.Println(err)
	}

	return
}

func closeHandle(address string) {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		deleteServer(address)
		fmt.Println("Shutting down server...")
		os.Exit(0)
	}()
}

func deleteServer(address string) {
	str := fmt.Sprintf("sed -i '/%v/d' ports.txt", address)
	cmd := exec.Command("bash", "-c", str)
	stdout, err := cmd.Output()
	if err != nil {
		log.Println(err)
	}
	fmt.Println(string(stdout))
}
