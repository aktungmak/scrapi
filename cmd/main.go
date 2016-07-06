package main

import (
    "github.com/aktungmak/scrapi"
    "net/url"
)


func main() {
u, _ := url.Parse("goob")
scrapi.Process(u)
}
