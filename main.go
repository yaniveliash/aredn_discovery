package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Nodes struct {
	Nodes []Node `json:"hosts"`
}

type Node struct {
	Name string `json:"name"`
	Ip   string `json:"ip"`
}

type Services struct {
	Services []Service `json:"services"`
}

type Service struct {
	Name     string `json:"name"`
	Protocol string `json:"protocol"`
	Link     string `json:"link"`
}

func getUrl(url string) []byte {

	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	respByte := buf.Bytes()

	return respByte
}

// getServices list all the visible services in the network
func getService() {
	var services Services

	respByte := getUrl("http://localnode:8080/cgi-bin/sysinfo.json?services=1")

	err := json.Unmarshal(respByte, &services)
	if err != nil {
		panic(err)
	}

	fmt.Println("Number of services in the network: ", len(services.Services))

	for i := 0; i < len(services.Services); i++ {
		fmt.Println(services.Services[i].Name)
		fmt.Println(services.Services[i].Protocol)
		fmt.Println(services.Services[i].Link)
		fmt.Println("-------------")
	}
}

// getNodes list all the visible nodes in the network
func getNodes() {
	var nodes Nodes

	respByte := getUrl("http://localnode:8080/cgi-bin/sysinfo.json?hosts=1")
	err := json.Unmarshal(respByte, &nodes)
	if err != nil {
		panic(err)
	}

	fmt.Println("Number of hosts in the network: ", len(nodes.Nodes))

	for i := 0; i < len(nodes.Nodes); i++ {
		fmt.Println(nodes.Nodes[i].Name)
		fmt.Println(nodes.Nodes[i].Ip)
		fmt.Println("-------------")
	}
}

func main() {

	getNodes()
	getService()
}
