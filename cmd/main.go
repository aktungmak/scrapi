package main

import (
	"flag"
	"github.com/aktungmak/scrapi"
	"log"
	"net/url"
	"os"
)

var mode string   // rec to capture state, play to serve a capture file
var file string   // location to store the results or play back from
var target string // the host to record or the host:port to serve on
var user string   // http Basic Auth username
var pass string   // http Basic Auth password
var insecure bool // accept self signed/bad certificates?
var note string   // some notes to attach to this capture
var err error

func init() {
	flag.StringVar(&mode, "m", "rec", "rec to capture state, play to serve a capture file")
	flag.StringVar(&file, "f", "dump.json", "location to store the results or play back from")
	flag.StringVar(&target, "t", "", "the service root to record or the host:port to serve on")
	flag.StringVar(&user, "u", "localhost\\sysadmin", "HTTP basic auth username")
	flag.StringVar(&pass, "p", "Sysadmin123@", "HTTP basic auth password")
	flag.BoolVar(&insecure, "k", true, "accept bad certificates (e.g. self signed)")
	flag.StringVar(&note, "n", "no notes", "some notes to attach to this capture")
}

func main() {
	flag.Parse()
	if target == "" {
		log.Print("the target is required! check -h for help")
		os.Exit(2)
	}
	switch mode {
	case "rec":
		u, err := url.Parse(target)
		if err != nil {
			break
		}
		client := scrapi.MakeClient(u.Host, user, pass, insecure)
		log.Printf("capturing %s into file %s with credentials %s:%s", target, file, user, pass)
		if insecure {
			log.Printf("ignoring any bad certificates")
		}
		err = scrapi.Rec(u, file, client, note)
	case "play":
		err = scrapi.Play(file, target)
	default:
		log.Fatalf("unknown mode '%s'. Try rec or play", mode)
	}
	if err != nil {
		log.Printf("error: %v", err)
	} else {
		log.Printf("complete!")
	}
}
