package iplookup

import (
	"encoding/json"
	"net/http"
)

type IpifyResponse struct {
	IP string `json:"ip"`
}

func Ipify() (ip string, err error) {
	var resp *http.Response

	resp, err = http.Get("https://api.ipify.org?format=json")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	ipify := &IpifyResponse{}
	err = decoder.Decode(ipify)
	ip = ipify.IP

	return
}
