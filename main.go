package main

import "github.com/adriamb/gopad/server"
import "github.com/CrowdSurge/banner"

func main() {
	banner.Print("gopad")
	server.ExecuteCmd()
}
