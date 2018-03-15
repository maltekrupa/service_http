package main

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"net/http"
	"os/exec"
	"sort"
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

// https://stackoverflow.com/a/34331660
func GetFqdn() string {
	cmd := exec.Command("/bin/hostname", "-f")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Panic(err)
	}
	fqdn := out.String()
	return fqdn[:len(fqdn)-1]
}

// https://arjanschaaf.github.io/request-headers-webserver-in-go/
func handler(w http.ResponseWriter, r *http.Request) {
	log.Println("Got request")
	oip := GetOutboundIP()
	fmt.Fprintln(w, "Outbound IP:", oip)
	fqdn := GetFqdn()
	fmt.Fprintln(w, "FQDN:", fqdn)
	fmt.Fprintln(w, "")

	var keys []string
	for k := range r.Header {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	fmt.Fprintln(w, "Request Headers: ")
	for _, k := range keys {
		line := fmt.Sprintf("%s: %s", k, r.Header[k])
		fmt.Fprintln(w, line)
	}
}

func handlerHealth(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "jo")
}

func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/health", handlerHealth)
	http.ListenAndServe(":8080", nil)
}
