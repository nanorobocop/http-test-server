package main

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func getHTTPCode(s string) (codeInt int, codeStr string, err error) {
	codeInt, err = strconv.Atoi(s)
	if err != nil {
		return 0, "", errors.New("There's no related http code")
	}
	codeStr = http.StatusText(codeInt)
	if codeStr == "" {
		return 0, "", errors.New("There's no related http code")
	}
	return codeInt, codeStr, nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		urlParts := strings.Split(r.URL.Path, "/")
		if urlParts[1] == "echo" {
			fmt.Fprintf(w, "%s", strings.Join(urlParts[2:], "/"))
			return
		} else if r.URL.Path == "/ip" {
			remoteAddr, _, _ := net.SplitHostPort(r.RemoteAddr)
			fmt.Fprintf(w, "%s", remoteAddr)
			return
		} else if urlParts[1] == "sleep" && len(urlParts) == 3 {
			seconds, err := strconv.Atoi(urlParts[2])
			if err != nil {
				http.Error(w, "Cannot sleep specified time", 500)
				return
			}
			time.Sleep(time.Duration(seconds) * time.Second)
			fmt.Fprintf(w, "Slept for %v second(s)", seconds)
			return
		} else if codeInt, codeStr, err := getHTTPCode(urlParts[1]); err == nil {
			codeMessage := strconv.Itoa(codeInt) + " " + codeStr
			http.Error(w, codeMessage, codeInt)
			return
		} else {
			fmt.Fprintf(w, "Hello there %s!", r.URL.Path[1:])
			return
		}
	case "POST":
		buf := new(bytes.Buffer)
		buf.ReadFrom(r.Body)
		s := buf.String()
		fmt.Fprintf(w, "%v", s)
	default:
		http.Error(w, "Not implemented", 501)
	}
}

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8090", nil))
}
