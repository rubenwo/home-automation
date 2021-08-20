package main

import (
	"github.com/rubenwo/home-automation/services/led-strip-service/internal/service"
	"log"
	"net/http"
)

func main() {
	//client, err := ledstrip.NewClient("192.168.2.135", 10)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//client.SetOnLedStripOnlineCallback(func(id string) {
	//	time.Sleep(time.Second * 1)
	//	fmt.Println(id)
	//
	//	info, err := client.InformationById("123-456-789-000")
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	fmt.Println(info)
	//	time.Sleep(time.Second * 1)
	//
	//	if err := client.SetSolidColor(ledstrip.Color{R: 255, G: 255, B: 255}); err != nil {
	//		log.Fatal(err)
	//	}
	//	time.Sleep(time.Second * 1)
	//
	//	if err := client.SetAnimationColorCycle(); err != nil {
	//		log.Fatal(err)
	//	}
	//})
	//
	//select {}

	cfg, err := service.LoadConfigFromPath("./config.json")
	if err != nil {
		log.Fatal(err)
	}
	router, err := service.New(cfg)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("led-strip-service is online!")
	if err := http.ListenAndServe(":80", router); err != nil {
		log.Fatal(err)
	}
	log.Println("led-strip-service is offline!")
}
