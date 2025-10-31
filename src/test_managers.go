package main

import (
	"fmt"
	"vpn-client/src/managers"
)

func main() {
	cm := managers.NewConnectionManager()
	fmt.Printf("Connection manager created with status: %v\n", cm.GetStatus())
	
	sm := managers.NewSubscriptionManager(nil)
	fmt.Printf("Subscription manager created with %d subscriptions\n", len(sm.GetAllSubscriptions()))
}