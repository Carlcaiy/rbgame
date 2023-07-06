.PHONY: game

all: gate hall game login

client:
	go run ./examples/client/main.go

gate:
	go run ./examples/gateway/main.go -p 6666 -s 1

gate1:
	go run ./examples/gateway/main.go -p 6667 -s 2

hall:
	go run ./examples/hall/main.go -p 6676 -s 1

hall1:
	go run ./examples/hall/main.go -p 6677 -s 2

game:
	go run ./examples/game/main.go

login:
	go run ./examples/login/main.go