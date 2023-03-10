package main

import (
	"bytes"
	"context"
	"errors"
	"flag"
	"io"
	"log"
	"net/http"
	"time"
)

func main() {
	// log with time and line
	log.SetFlags(log.Lshortfile | log.Ldate | log.Lmicroseconds)
	log.Println("v2.2")
	log.Println("this is biz server")

	sidecarMode := flag.Bool("sidecar", false, "sidecar mode")
	flag.Parse()

	if *sidecarMode {
		// try to connect 9000 port with 200ms  timeout
		log.Println("try to connect sidecar ...")
		var sidecarReady = false
		for i := 0; i < 2; i++ {
			err := checkSidecar()
			if err != nil {
				log.Printf("sidecar not ready [%d] err:%v \n", i, err)
			} else {
				sidecarReady = true
				break
			}
			time.Sleep(time.Millisecond * 500)
		}
		if !sidecarReady {
			log.Fatal("sidecar not ready, exit ...")
		}
		log.Println("sidecar ready")

	}

	mux := http.NewServeMux()
	mux.HandleFunc("/hello", Hello)
	log.Println("register route : " + "/echo")
	mux.HandleFunc("/readyz", Ping)
	log.Println("register route : " + "/readyz")
	mux.HandleFunc("/healthz", Ping)
	log.Println("register route : " + "/healthz")

	log.Print("start serving ... ")
	err := http.ListenAndServe(":80", mux)
	if err != nil {
		panic(err)
	}

	log.Print("end serving ... ")
}

func checkSidecar() (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*200)
	defer cancel()
	request, err := http.NewRequestWithContext(ctx,
		http.MethodPost,
		"http://localhost:9000/echo",
		bytes.NewReader([]byte("this is biz")))
	if err != nil {
		return
	}
	defer request.Body.Close()
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		err = errors.New("sidecar not ready")
		return
	}

	body, _ := io.ReadAll(response.Body)

	log.Println("response from sidecar : " + string(body))
	return
}

// Ping returns true automatically when checked.
var Ping = func(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

var Hello = func(w http.ResponseWriter, r *http.Request) {
	// write response
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("hello world"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	return
}
