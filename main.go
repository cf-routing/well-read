package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"os/exec"
	"strings"
	"time"
)

func waitingHandler() http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("received request")

			customHeader := r.Header.Get("User-Agent")
			if customHeader == "well-read" {
				fmt.Printf("received request with user-agent=well-read\n")
			}

			t1 := time.Now()

			_, err := io.Copy(ioutil.Discard, r.Body)
			if err != nil {
				log.Println("failed to discard the body")
				w.WriteHeader(http.StatusTeapot)
				return
			}

			t2 := time.Now()

			fmt.Printf("received body in %s\n", t2.Sub(t1))

			if t2.Sub(t1) > time.Second {
				fmt.Printf("received body slower than 1s: %s\n", t2.Sub(t1))
			}

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
	fmt.Printf("args: %#v\n", os.Args)

	if len(os.Args) == 4 && os.Args[1] == "slowpost" {
		err := makeSlowReq(os.Args[2], os.Args[3])
		if err != nil {
			fmt.Println("slow req: %s", err)
			os.Exit(1)
		}
		os.Exit(0)
	}

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

func makeSlowReq(reqURL string, readDelay string) error {
	dur, err := time.ParseDuration(readDelay)
	if err != nil {
		return err
	}

	fmt.Printf("will sleep for %s\n", dur)

	req, err := http.NewRequest("POST", reqURL, nil)
	if err != nil {
		return err
	}

	beginningData := []byte(strings.Repeat("k", 9000))
	remainingData := []byte("=v\n")

	req.ContentLength = int64(len(beginningData) + len(remainingData))
	req.Body = &slowReader{dur, bytes.NewReader(beginningData), remainingData, 0}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("User-Agent", "well-read")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	fmt.Printf("response: %s\n%s\n", resp.Status, respBytes)

	return nil
}

type slowReader struct {
	delay         time.Duration
	reader        io.Reader
	remainingData []byte
	state         int
}

func (r *slowReader) Read(p []byte) (int, error) {
	switch r.state {
	case 0:
		n, err := r.reader.Read(p)
		if err == io.EOF {
			r.state = 1
			return n, nil
		}
		return n, err

	case 1:
		r.state = 2
		fmt.Printf("sleeping...\n")
		time.Sleep(r.delay)
		fmt.Printf("done sleeping...\n")

	case 2:
		r.state = 3
		return copy(p, r.remainingData), nil

	case 3:
		return 0, io.EOF
	}

	return 0, nil
}

func (r *slowReader) Close() (err error) { return nil }
