package main

import "github.com/adriamb/runes/server"
import "github.com/CrowdSurge/banner"

func main() {
	banner.Print("runes")
	server.ExecuteCmd()
}
