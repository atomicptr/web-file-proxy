package proxy

import (
	"log"

	"github.com/atomicptr/web-file-proxy/link"
)

type Proxy struct {
	DatabaseDriver string
	DatabaseUrl    string
	Addr           string
	SecretHash     string

	linkRepository *link.Repository
	authProof      []byte
}

func (p *Proxy) Run() {
	err := p.generateNewAuthProof()
	if err != nil {
		log.Fatal(err)
	}

	err = p.initDatabase()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Web File Proxy running at", p.Addr)

	err = p.initHttp()
	if err != nil {
		log.Fatal(err)
	}
}
