package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/ealfarozi/absen-beacon/common"
	"tinygo.org/x/bluetooth"
)

var adapter = bluetooth.DefaultAdapter

func main() {
	common.GetVars()

	/*
		if common.IS_STATIC == "0" {
			for {
				run()
			}
		} else {
			runStatic()
		}
	*/
	for {
		run()
	}
}

func run() {
	must("enable BLE stack", adapter.Enable())
	common.GetHash()

	//musti bikin retry
	hitBeaconAPI(common.HASHED)

	adv := adapter.DefaultAdvertisement()

	must("config adv", adv.Configure(bluetooth.AdvertisementOptions{
		LocalName: common.LOCAL_NAME + "|" + common.UUID + "|" + common.HASHED,
	}))
	must("start adv", adv.Start())
	println("start advertising...")

	address, _ := adapter.Address()
	for i := 1; i < common.REFRESH_INTERVAL; i++ {
		println(common.LOCAL_NAME+"|"+common.UUID+"|"+common.HASHED, "/", address.MAC.String())

		time.Sleep(time.Second)
	}
	must("stop adv", adv.Stop())
	println("stop advertising...")
}

func runStatic() {
	must("enable BLE stack", adapter.Enable())
	common.GetHash()

	//musti bikin retry
	hitBeaconAPI(common.HASHED)

	adv := adapter.DefaultAdvertisement()

	must("config adv", adv.Configure(bluetooth.AdvertisementOptions{
		LocalName: common.LOCAL_NAME + "|" + common.UUID + "|" + common.HASHED,
	}))
	must("start adv", adv.Start())

	println("start advertising...")

	address, _ := adapter.Address()
	for {
		println(common.LOCAL_NAME+"|"+common.UUID+"|"+common.HASHED, "/", address.MAC.String())
		time.Sleep(time.Second)
	}
}

func must(action string, err error) {
	if err != nil {
		panic("failed to " + action + ": " + err.Error())
	}
}

func hitBeaconAPI(data string) bool {
	url := common.BASE_URL + "/api-iot/v1/localname-beacon"
	strToken := "Basic " + common.TOKEN

	b := common.BeaconRequest{}
	b.BeaconID = common.UUID
	loc, _ := time.LoadLocation("Asia/Jakarta")
	startTime := time.Now().In(loc)
	endTime := startTime.Add(time.Second * time.Duration(common.REFRESH_INTERVAL))
	b.StartTime = startTime.Format("2006-01-02 15:04:05")
	b.EndTime = endTime.Format("2006-01-02 15:04:05")
	b.ExpTimeMin = common.REFRESH_INTERVAL / 60
	b.Data = data

	br, _ := json.Marshal(b)

	/*
		_, resp, body, err := common.HitAPI(url, br, "POST", strToken, time.Duration(120))
		fmt.Println("[Hit BeaconRequest]:", resp)
		if err != nil {
			fmt.Println("[Error Hit BeaconRequest]:", body)
			return false
		}

		if resp.StatusCode != 200 {
			fmt.Println("[Error Hit BeaconRequest]:", body)
			fmt.Println(body)
			return false
		}
	*/

	//maxRetries := 10
	//for i := 0; i < maxRetries; i++ {
	for {
		_, resp, body, err := common.HitAPI(url, br, "POST", strToken, time.Duration(120))
		fmt.Println("[Hit BeaconRequest]:", resp)
		if err != nil {
			fmt.Println("[Error Hit BeaconRequest]:", body)
		} else if resp.StatusCode == 200 {
			return true
		} else {
			fmt.Println("[Error Hit BeaconRequest]:", body)
		}

		//if i < maxRetries-1 {
		fmt.Println("Retrying...")
		time.Sleep(2 * time.Second)
		//}
	}

	return true
}
