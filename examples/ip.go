/*
  IP is a simple example, contacting a provider of public IP resolution, and
  unmarshaling the result into the IP struct.
*/
package main

import (
	"fmt"

	"github.com/jamescun/rest"
)

type IP struct {
	AS       string `json:"as"`
	City     string `json:"city"`
	Code     string `json:"countryCode"`
	Country  string `json:"country"`
	ISP      string `json:"ISP"`
	Org      string `json:"org"`
	Query    string `json:"query"`
	Timezone string `json:"timezone"`
}

var (
	IPAPI, _ = rest.New("http://ip-api.com", rest.EncoderJSON{})

	GetIP = IPAPI.Request("GET", "/json")
)

func main() {
	var ip IP
	err := GetIP(&ip)
	if err != nil {
		fmt.Println("error:", err)
	}

	fmt.Printf("ip: %+v\n", ip)
}
