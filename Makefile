clientgo:
	go run client/main.go
	
servergo:
	 go run server/main.go

servergo2:
	 go run server/main.go -s=localhost:8081 -t=Europe/Berlin

