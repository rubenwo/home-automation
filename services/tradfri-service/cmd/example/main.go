package main

import (
	"encoding/json"
	"fmt"
	"github.com/eriklupander/tradfri-go/model"
	coap2 "github.com/moroen/gocoap/v3"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	//fmt.Println("Hello World")
	//
	//client := tradfri.NewTradfriClient()
	//devices, err := client.ListDevices()
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//for _, device := range devices {
	//	fmt.Printf("%+v\n", device)
	//}

	log.Println("Observe called")

	param := coap2.ObserveParams{}

	endpoints := `["15001/65536"]`

	var uris []string

	err := json.Unmarshal([]byte(endpoints), &uris)
	if err != nil {
		panic(err.Error())
	}

	param.URI = uris

	msg := make(chan []byte)
	sign := make(chan bool)
	errSign := make(chan error)

	// state := 0
	termChan := make(chan os.Signal)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)

	err = coap2.Observe(param, msg, sign, errSign)
	if err != nil {
		log.Println(err.Error())
	} else {
		for {
			select {
			case message, isOpen := <-msg:
				if isOpen {
					fmt.Println(string(message))
					device := &model.Device{}
					err = json.Unmarshal(message, &device)
					if err != nil {
						log.Fatal(err)
					}
					fmt.Printf("Parsed: %+v\n", device)
				} else {
					return
				}
			case err = <-errSign:
				log.Fatal(err)
			case <-termChan:
				fmt.Println("Shutdown requested")
				sign <- true
			}
		}
	}

}
