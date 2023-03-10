package main

import (
	"flag"
	"io"
	"log"
	"net/http"
	"time"
)

func main() {
	// log with time and line
	log.SetFlags(log.Lshortfile | log.Ldate | log.Lmicroseconds)

	noDelay := flag.Bool("nodelay", false, "disable Nagle's algorithm")
	flag.Parse()
	if !*noDelay {
		log.Println("delay 3 sec ....")
		// 模拟慢启动
		time.Sleep(time.Second * 3)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/echo", Echo)
	log.Println("register route : " + "/echo")
	mux.HandleFunc("/readyz", Ping)
	log.Println("register route : " + "/readyz")
	mux.HandleFunc("/healthz", Ping)
	log.Println("register route : " + "/healthz")

	log.Print("start serving ... ")
	err := http.ListenAndServe(":9000", mux)
	if err != nil {
		panic(err)
	}

	log.Print("end serving ... ")
}

// Ping returns true automatically when checked.
var Ping = func(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

var Echo = func(w http.ResponseWriter, r *http.Request) {
	// get body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()
	log.Println("receive request : " + string(body))
	// get time
	now := time.Now().Format(time.RFC3339)

	// write response
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(now + " " + string(body)))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	return
}
