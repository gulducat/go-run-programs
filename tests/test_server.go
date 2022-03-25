package main

// simple lil test program

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"syscall"
	"time"
)

func main() {
	help := `
usage:
  test-server command [arg]
commands:
  test-server start <port>    | start an http server on port
  test-server check <port>    | confirm http connectivity
  test-server sleep <seconds> | sleep for some seconds
  test-server death <message> | exit with error
`
	if len(os.Args) < 3 {
		fmt.Println(help)
		syscall.Exit(1)
	}

	command := os.Args[1]
	arg := strings.Join(os.Args[2:], " ")

	switch command {
	default:
		fmt.Println(help)

	case "start":
		port := os.Args[2]
		if os.Getenv("PORT") != "" {
			port = os.Getenv("PORT")
		}
		log.Fatal(webStart(port))

	case "check":
		port := os.Args[2]
		if os.Getenv("PORT") != "" {
			port = os.Getenv("PORT")
		}
		resp, err := http.Get("http://127.0.0.1:" + port)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(resp.Status)

	case "sleep":
		i, err := strconv.ParseFloat(arg, 32)
		if err != nil {
			log.Fatal(err)
		}
		millis := time.Duration(i * 1000)
		time.Sleep(time.Millisecond * millis)

	case "death":
		log.Fatal(arg + "?  SQUEAK.") // TODO: arg here makes it unclear whether the output we check in tests is from here or from go-run-programs
	}
}

func webStart(port string) error {
	// exit early if port is taken
	ln, err := net.Listen("tcp", "127.0.0.1:"+port)
	if err != nil {
		return err
	}
	if err = ln.Close(); err != nil {
		return err
	}

	// pretend to be slow
	time.Sleep(time.Second * 2) // TODO
	log.Println("oh no i'm so slow i take like 3 seconds to start")

	// actually start a server
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("uri:", r.RequestURI)
		// spit the URI back out to the requestor
		if _, err := w.Write([]byte(r.RequestURI)); err != nil {
			log.Println(err)
		}
	})

	return http.ListenAndServe("127.0.0.1:"+port, nil)
}
