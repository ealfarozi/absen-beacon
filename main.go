package main

import (
	"time"

	"github.com/ealfarozi/absen-beacon/common"
	"tinygo.org/x/bluetooth"
)

var adapter = bluetooth.DefaultAdapter

func main() {
	for {
		run()
		time.Sleep(time.Duration(common.REFRESH_INTERVAL) * time.Second)
	}

}

func run() {
	must("enable BLE stack", adapter.Enable())

	adv := adapter.DefaultAdvertisement()

	must("config adv", adv.Configure(bluetooth.AdvertisementOptions{
		LocalName: common.LOCAL_NAME + common.HASHED,
	}))
	must("start adv", adv.Start())

	println("advertising...")

	address, _ := adapter.Address()
	for i := 1; i < common.REFRESH_INTERVAL; i++ {
		println(common.LOCAL_NAME+common.HASHED, "/", address.MAC.String())

		time.Sleep(time.Second)
	}
	must("stop adv", adv.Stop())
}

func must(action string, err error) {
	if err != nil {
		panic("failed to " + action + ": " + err.Error())
	}
}
