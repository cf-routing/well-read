package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"time"
)

func waitingHandler() http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("received request")

			t1 := time.Now()

			_, err := io.Copy(ioutil.Discard, r.Body)
			if err != nil {
				log.Println("failed to discard the body")
				w.WriteHeader(http.StatusTeapot)
				return
			}

			t2 := time.Now()

			fmt.Printf("received body in %s\n", t2.Sub(t1))

			dump, err := httputil.DumpRequest(r, false)
			if err != nil {
				log.Println("failed to dump the request")
				w.WriteHeader(http.StatusTeapot)
				return
			}

			w.Header().Set("Content-Type", "text/html")
			w.WriteHeader(http.StatusOK)
			fmt.Printf("dumping request: %s", dump)
		})
}

func main() {
	http.HandleFunc("/", waitingHandler())
	log.Fatal(http.ListenAndServe(":8080", nil))
}
