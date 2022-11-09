package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func loggingHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("%s - %s - %s - %s\n",
			time.Now().Format(time.Stamp),
			r.RemoteAddr,
			r.RequestURI,
			r.Referer())

		h.ServeHTTP(w, r)
	})
}

func checkFilepath(p string) (string, error) {
	fp, err := filepath.Abs(p)
	if err != nil {
		return "", err
	}
	_, err = os.Stat(fp)
	if err != nil {
		return "", err
	}
	return fp, nil
}

func checkPortNumber(p int) error {
	if p < 1 || p > 65535 {
		return errors.New("invalid port number")
	}
	return nil
}

func listNetworkInterfacesAndIP() {
	ifaces, err := net.Interfaces()
	if err != nil {
		log.Println(err)
		return
	}

	for _, iface := range ifaces {
		fmt.Printf("\n%s\n", iface.Name)
		addrs, err := iface.Addrs()
		if err != nil {
			fmt.Println("Error:", err.Error())
			continue
		}
		for _, addr := range addrs {
			fmt.Printf("  %s\n", addr.String())
		}
	}
}

func main() {
	var ip, dir string
	var port int
	var listInterfaces bool

	flag.StringVar(&ip, "ip", "127.0.0.1", "network address to listen on. \"0.0.0.0\" is all addresses")
	flag.IntVar(&port, "p", 9001, "network port to listen on")
	flag.StringVar(&dir, "dir", ".", "path to the directory to serve")
	flag.BoolVar(&listInterfaces, "l", false, "list network interfaces and their ip addresses")
	flag.Parse()

	if listInterfaces {
		listNetworkInterfacesAndIP()
		return
	}

	fp, err := checkFilepath(dir)
	if err != nil {
		log.Fatalln(err)
	}

	err = checkPortNumber(port)
	if err != nil {
		log.Fatalf("%s: %d", err, port)
	}

	addr := fmt.Sprintf("%v:%v", ip, port)

	fmt.Println("http server starting...")
	fmt.Println("address:", addr)
	fmt.Println("directory:", fp)

	fs := http.FileServer(http.Dir(fp))
	http.Handle("/", loggingHandler(fs))
	log.Fatal(http.ListenAndServe(addr, nil))
}
