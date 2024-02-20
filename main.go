package main

import (
	"time"

	"github.com/ealfarozi/absen-beacon/common"
	"tinygo.org/x/bluetooth"
)

var adapter = bluetooth.DefaultAdapter

func main() {
	common.GetVars()

	if common.IS_STATIC == "0" {
		for {
			run()
		}
	} else {
		runStatic()
	}
}

func run() {
	must("enable BLE stack", adapter.Enable())
	common.GetHash()
	adv := adapter.DefaultAdvertisement()

	must("config adv", adv.Configure(bluetooth.AdvertisementOptions{
		LocalName: common.LOCAL_NAME + common.HASHED,
	}))
	must("start adv", adv.Start())

	println("start advertising...")

	address, _ := adapter.Address()
	for i := 1; i < common.REFRESH_INTERVAL; i++ {
		println(common.LOCAL_NAME+common.HASHED, "/", address.MAC.String())

		time.Sleep(time.Second)
	}
	must("stop adv", adv.Stop())
	println("stop advertising...")
}

func runStatic() {
	must("enable BLE stack", adapter.Enable())
	common.GetHash()
	adv := adapter.DefaultAdvertisement()

	must("config adv", adv.Configure(bluetooth.AdvertisementOptions{
		LocalName: common.LOCAL_NAME + common.HASHED,
	}))
	must("start adv", adv.Start())

	println("start advertising...")

	address, _ := adapter.Address()
	for {
		println(common.LOCAL_NAME+common.HASHED, "/", address.MAC.String())
		time.Sleep(time.Second)
	}
}

func must(action string, err error) {
	if err != nil {
		panic("failed to " + action + ": " + err.Error())
	}
}
