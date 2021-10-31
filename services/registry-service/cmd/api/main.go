package main

import (
	"github.com/rubenwo/home-automation/services/registry-service/pkg/registry"
	"github.com/rubenwo/home-automation/services/registry-service/pkg/registry/routines/plugins"
	"log"
	"net"
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
//console.log(resp);
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
//resp
//`

//enp3s0: flags=4163<UP,BROADCAST,RUNNING,MULTICAST>  mtu 1500
//inet 192.168.178.45  netmask 255.255.255.0  broadcast 192.168.178.255
//inet6 fe80::468a:5bff:fe83:9037  prefixlen 64  scopeid 0x20<link>
//ether 44:8a:5b:83:90:37  txqueuelen 1000  (Ethernet)
//RX packets 72748  bytes 99836762 (99.8 MB)
//RX errors 0  dropped 0  overruns 0  frame 0
//TX packets 18144  bytes 1648484 (1.6 MB)
//TX errors 0  dropped 0 overruns 0  carrier 0  collisions 0
//device interrupt 17
//

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
	//i, err := res.Export()
	//var m map[string]interface{}
	//err = json.Unmarshal([]byte(i.(string)), &m)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Println(m)

	//44:8a:5b:83:90:37

	mac, err := net.ParseMAC("44:8a:5b:83:90:37")
	if err != nil {
		log.Fatal(err)
	}

	err = new(plugins.WoL).Run(&plugins.WoLCfg{
		Addr:         "192.168.178.45:9",
		HardwareAddr: mac,
	})

	log.Fatal(err)

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
