package main

import (
	"fmt"

	ping "github.com/pashanskiy/Minecraft-Go-Scan-Servers/components/ping"
)

func main() {
	PingServer := ping.PingServer{}

	if !checkErr(PingServer.NewPing("kek", "53777")) {
		return
	}
	if !checkErr(PingServer.GetConnect()) {
		return
	}
	if !checkErr(PingServer.RequestInfoAndUnmarshall()) {
		return
	}
	fmt.Println("Success:", PingServer.ServerData)


	}

	func checkErr(err error) bool {
		if err != nil {
			fmt.Println("ERROR:", err)
			return false
		}
		return true
	}