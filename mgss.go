package main

import (
	"fmt"
	ping "github.com/pashanskiy/Minecraft-Go-Scan-Servers/internal/ping"
)

func main() {
	PingServer := ping.PingServer{}

	if !checkErr(PingServer.NewPing("kek", "53777")) {
		return
	}
	if !checkErr(PingServer.GetConnect()) {
		return
	}
	ServerData, err := PingServer.RequestInfoAndUnmarshall()

	if checkErr(err) {
		fmt.Println("Success:", ServerData)
		fmt.Println("Ping:", ServerData.Ping)
	}

}

func checkErr(err error) bool {
	if err != nil {
		fmt.Println("ERROR:", err)
		return false
	}
	return true
}
