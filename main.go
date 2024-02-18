package main

import (
	"time"

	"crypto/rand"
	"encoding/hex"
	"fmt"

	"github.com/cloudflare/circl/hpke"
	"tinygo.org/x/bluetooth"
)

var adapter = bluetooth.DefaultAdapter

func main() {
	for {
		run()
		time.Sleep(1 * time.Minute)
	}
}

func run() {
	must("enable BLE stack", adapter.Enable())
	adv := adapter.DefaultAdvertisement()
	ciph := encryption("RSCM:01")

	must("config adv", adv.Configure(bluetooth.AdvertisementOptions{
		LocalName: "Go " + ciph,
	}))
	must("start adv", adv.Start())

	println("advertising...")
	address, _ := adapter.Address()
	for i := 1; i < 60; i++ {
		println("Go", ciph, "/", address.MAC.String())
		time.Sleep(time.Second)
	}
	must("stop adv", adv.Stop())
}

func must(action string, err error) {
	if err != nil {
		panic("failed to " + action + ": " + err.Error())
	}
}

func encryption(localName string) string {

	kemID := int(hpke.KEM_P256_HKDF_SHA256)
	kdfID := int(hpke.KDF_HKDF_SHA256)
	aeadID := int(hpke.AEAD_AES128GCM)
	msg := localName

	suite := hpke.NewSuite(hpke.KEM(kemID), hpke.KDF(kdfID), hpke.AEAD(aeadID))

	info := []byte("some_info_key")

	Bob_pub, Bob_private, _ := hpke.KEM(kemID).Scheme().GenerateKeyPair()

	Bob, _ := suite.NewReceiver(Bob_private, info)

	Alice, _ := suite.NewSender(Bob_pub, info)

	enc, sealer, _ := Alice.Setup(rand.Reader)

	Alice_msg := []byte(msg)
	aad := []byte("some additional data")
	ct, _ := sealer.Seal(Alice_msg, aad)

	opener, _ := Bob.Setup(enc)

	Bob_msg, _ := opener.Open(ct, aad)

	// fmt.Printf("Public key type:\t%s\n", Bob_pub.Scheme().Name())
	// fmt.Printf(" Params\t%s\n", suite.String())
	// fmt.Printf("Key exchange parameters:\n")
	// fmt.Printf(" Ciphersize:\t%d\n", hpke.KEM(kemID).Scheme().CiphertextSize())
	// fmt.Printf(" EncapsulationSeedSize:\t%d\n", hpke.KEM(kemID).Scheme().EncapsulationSeedSize())
	// fmt.Printf(" PrivateKeySize:\t%d\n", hpke.KEM(kemID).Scheme().PrivateKeySize())
	// fmt.Printf(" PublicKeySize:\t%d\n", hpke.KEM(kemID).Scheme().PublicKeySize())
	// fmt.Printf(" SeedSize:\t%d\n", hpke.KEM(kemID).Scheme().SeedSize())
	// fmt.Printf(" SharedKeySize:\t%d\n", hpke.KEM(kemID).Scheme().SharedKeySize())

	// fmt.Printf("Cipher parameters:\n")
	// fmt.Printf(" Key Length:\t%d\n", hpke.AEAD(aeadID).KeySize())

	// fmt.Printf("Key derivation function:\n")
	// fmt.Printf(" Extract size:\t%d\n", hpke.KDF(kdfID).ExtractSize())

	fmt.Printf("\nMessage:\t%s\n", Alice_msg)
	ciph := hex.EncodeToString(ct)
	fmt.Printf("Cipher:\t%x\n", ciph)
	fmt.Printf("Decipher:\t%s\n", Bob_msg)

	return ciph

}
