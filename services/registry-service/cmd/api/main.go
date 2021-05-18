package main

import (
	"github.com/rubenwo/home-automation/services/registry-service/pkg/registry"
	"log"
	"net/http"
)

//var js = `
//console.log("hello world");
//resp = HttpGet("https://homeautomation.rubenwoldhuis.nl/api/v1/devices");
//console.log(resp);
//
//data = {
//	username:"",
//	password:""
//}
//
//resp = HttpPost("https://homeautomation.rubenwoldhuis.nl/auth/login", data);
//
//loginResp = JSON.parse(resp);
//
//console.log(loginResp.authorization_token);
//
//headers = {
//	Authorization:"Bearer " + loginResp.authorization_token
//}
//
//resp = HttpGet("https://homeautomation.rubenwoldhuis.nl/api/v1/devices", headers);
//console.log(resp);
//
//`


func main() {
	//vm := otto.New()
	//
	//_ = vm.Set("HttpGet", script.HttpGet)
	//_ = vm.Set("HttpPost", script.HttpPost)
	//_ = vm.Set("HttpDelete", script.HttpDelete)
	//_ = vm.Set("HttpPut", script.HttpPut)
	//
	//res, err := vm.Run(js)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//log.Println(res)


	cfg, err := registry.LoadConfigFromPath("./config.json")
	if err != nil {
		log.Fatal(err)
	}
	router, err := registry.New(cfg)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("registry-service is online!")
	if err := http.ListenAndServe(":80", router); err != nil {
		log.Fatal(err)
	}
}
