package main

// todo: add some extra fields to the dump file like
// datetime of dump, and some notes. these can be
// displayed on a little start page which links to
// the serviceroot

import (
	"github.com/aktungmak/scrapi"
	"log"
	"net/url"
	"os"
)

func main() {
	if len(os.Args) < 4 {
		log.Fatal("usage: scrapi <verb> <filename> <target>")
		os.Exit(2)
	}
	var err error
	var user = "sysadmin"
	var pass = "sysadmin123"
	var insecure = true
	switch os.Args[1] {
	case "rec":
		u, err := url.Parse(os.Args[3])
		if err != nil {
			break
		}
		client := scrapi.MakeClient(u.Host, user, pass, insecure)
		err = scrapi.Rec(u, os.Args[2], client)
	case "play":
		err = scrapi.Play(os.Args[2], os.Args[3])
	default:
		log.Fatalf("unknown verb '%s'. Try rec or play", os.Args[1])
	}
	log.Printf("err: %v", err)
}
