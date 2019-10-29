package main

import (
	"github.com/atomicptr/web-file-proxy/proxy"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	p := proxy.Proxy{
		DatabaseDriver: "sqlite3",
		DatabaseUrl:    "./proxy.db",
		Addr:           ":8081",
		SecretHash:     "f841c5abf4c6de3ca5db764a1a85d3b645944fb2feea27059827bb09f7bddd08",
	}
	p.Run()
}
