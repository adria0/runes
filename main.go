package main

import "github.com/amassanet/gopad/web"
import "github.com/amassanet/gopad/server"

import "github.com/crowdsurge/banner"

func main() {
	banner.Print("gopad")
	server.NewServer(server.ServerConfiguration{Port: 8080, Prefix: ""})
	web.InitWeb()
	server.StartServer()
}
