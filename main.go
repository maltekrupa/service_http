package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"sort"
	"strings"
)

// https://stackoverflow.com/a/37382208
func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}

// https://arjanschaaf.github.io/request-headers-webserver-in-go/
func handler(w http.ResponseWriter, r *http.Request) {
	log.Println("Got request")
	ip := GetOutboundIP()
	fmt.Fprintln(w, "Outbound IP: ", ip)

	var keys []string
	for k := range r.Header {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	fmt.Fprintln(w, "Request Headers: ")
	for _, k := range keys {
		key := strings.TrimSpace(k)
		line := fmt.Sprintf("%q: %q", key, r.Header[k])
		fmt.Fprintln(w, line)
	}
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
