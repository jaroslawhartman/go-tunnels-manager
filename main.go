package main

import (
	"jhartman.pl/go-tunnels-ui/tunnelsmgr"
)

func main() {
	myTunnelmgr := &tunnelsmgr.Tunnelmgr{}

	myTunnelmgr.Run()
}
