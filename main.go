package main

import (
	"strconv"
	"time"

	"github.com/ealfarozi/absen-beacon/common"
	"tinygo.org/x/bluetooth"
)

var adapter = bluetooth.DefaultAdapter
var LOCAL_NAME string
var HASHED string

func main() {
	LOCAL_NAME = common.GetEnv("LOCAL_NAME")
	HASHED = strconv.FormatUint(common.SetHash(common.GetEnv("BEACON_ID")), 10)
	for {
		run()
		time.Sleep(1 * time.Minute)
	}
}

func run() {
	must("enable BLE stack", adapter.Enable())
	adv := adapter.DefaultAdvertisement()

	must("config adv", adv.Configure(bluetooth.AdvertisementOptions{
		LocalName: LOCAL_NAME + HASHED,
	}))
	must("start adv", adv.Start())

	println("advertising...")
	address, _ := adapter.Address()
	for i := 1; i < 60; i++ {
		println(LOCAL_NAME, "/", address.MAC.String())
		time.Sleep(time.Second)
	}
	must("stop adv", adv.Stop())
}

func must(action string, err error) {
	if err != nil {
		panic("failed to " + action + ": " + err.Error())
	}
}
