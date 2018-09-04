package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"os/exec"
	"strings"
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

var netstatStates = []string{
	"LISTEN",
	"ESTABLISHED",
	"SYN_SENT",
	"SYN_RECV",
	"LAST_ACK",
	"CLOSE_WAIT",
	"TIME_WAIT",
	"CLOSED",
	"CLOSING",
	"FIN_WAIT1",
	"FIN_WAIT2",
}

func collectNetstat() {
	nt, err := exec.Command("netstat", "-at").CombinedOutput()
	if err != nil {
		fmt.Printf("failed to get netstat: %s\n", err)
		return
	}

	counts := make([]string, len(netstatStates))

	for i, state := range netstatStates {
		counts[i] = fmt.Sprintf("%s=%d", state, strings.Count(string(nt), state))
	}

	fmt.Printf("netstat stats: %s\n", strings.Join(counts, " "))
}

func main() {
	go func() {
		defer func() {
			if e := recover(); e != nil {
				fmt.Println("Recovered in netstat: %s", e)
			}
		}()

		for {
			collectNetstat()
			time.Sleep(1 * time.Second)
		}
	}()

	http.HandleFunc("/", waitingHandler())
	log.Fatal(http.ListenAndServe(":8080", nil))
}
