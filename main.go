package main

import "github.com/adriamb/gopad/cmd"
import "github.com/crowdsurge/banner"

func main() {
	banner.Print("gopad")
	cmd.Execute()
}
