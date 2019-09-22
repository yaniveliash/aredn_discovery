package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
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

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func getUrl(url string) []byte {

	resp, err := http.Get(url)
	check(err)

	defer resp.Body.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	respByte := buf.Bytes()

	return respByte
}

// getServices list all the visible services in the network
func getService() string {
	var services Services

	respByte := getUrl("http://localnode:8080/cgi-bin/sysinfo.json?services=1")

	err := json.Unmarshal(respByte, &services)
	check(err)
	result := "<h1>AREDN VISIBLE SERVICES</h1>"
	result = result + "<table><tr>"
	for i := 0; i < len(services.Services); i++ {
		result = result + "<th><a href=\""
		result = result + (services.Services[i].Link)
		result = result + "/\">"
		result = result + (services.Services[i].Name)
		result = result + "</a></th>"
		// New table line for each 5 records
		if (i+1)%5 == 0 {
			result = result + "</tr><tr>"
		}
	}

	result = result + "</tr></table><br>"
	result = result + ("Number of services in the network: ")
	result = result + strconv.Itoa(len(services.Services))

	return result
}

// getNodes list all the visible nodes in the network
func getNodes() string {
	var nodes Nodes

	respByte := getUrl("http://localnode:8080/cgi-bin/sysinfo.json?hosts=1")
	err := json.Unmarshal(respByte, &nodes)
	check(err)

	//Creating table
	result := "<h1>AREDN VISIBLE NODES</h1>"
	result = result + "<table><tr>"
	for i := 0; i < len(nodes.Nodes); i++ {
		result = result + "<th><a href=\"http://"
		result = result + (nodes.Nodes[i].Name)
		result = result + "/\">"
		result = result + (nodes.Nodes[i].Name)
		result = result + "</a></th>"
		// New table line for each 5 records
		if (i+1)%5 == 0 {
			result = result + "</tr><tr>"
		}
	}

	result = result + "</tr></table><br>"
	result = result + "Number of hosts in the network: "
	result = result + strconv.Itoa(len(nodes.Nodes))
	result = result + "<br><br>"

	return result
}

func main() {
	// Check or Create static directory
	if _, err := os.Stat("static/"); os.IsNotExist(err) {
		os.Mkdir("static/", 0755)
		os.Mkdir("static/style", 0755)
	}
	style, err := os.Create("static/style/style.css")
	check(err)
	defer style.Close()

	ws := bufio.NewWriter(style)

	// Generating CSS Style file
	ws.WriteString(".content { \n")
	ws.WriteString("	max-width: 500px; \n")
	ws.WriteString("	margin: auto; \n}\n\n")
	ws.WriteString("table { \n")
	ws.WriteString("	border-collapse: collapse; \n}\n\n")
	ws.WriteString("table, th, td { \n")
	ws.WriteString("	border: 1px solid black; \n}\n\n")
	ws.WriteString("th, td { \n")
	ws.WriteString("	padding: 15px; \n")
	ws.WriteString("	text-align: left; \n}\n")

	ws.Flush()

	html, err := os.Create("static/index.html")
	check(err)
	defer html.Close()

	w := bufio.NewWriter(html)

	w.WriteString("<link rel=\"stylesheet\" type=\"text/css\" href=\"style/style.css\">")
	w.WriteString("<body><div class=\"content\">")

	w.WriteString(getNodes())
	w.WriteString(getService())

	w.WriteString("</div></body>")

	w.Flush()

	http.Handle("/", http.FileServer(http.Dir("./static")))

	log.Fatal(http.ListenAndServe(":80", nil))
}
