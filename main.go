package main

import (
	"github.com/atomicptr/web-file-proxy/proxy"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func getEnv(envVar, defaultValue string) string {
	env := os.Getenv(envVar)
	if env != "" {
		return env
	}

	return defaultValue
}

func main() {
	secretHash := getEnv("SECRET_HASH", "")
	if secretHash == "" || len(secretHash) != 64 {
		log.Fatal("SECRET_HASH is unset, aborting...\n\n\tCreate your own using: sha3_256(\"wfp_\" + your_password)")
	}

	p := proxy.Proxy{
		DatabaseDriver: getEnv("DATABASE_DRIVER", "sqlite3"),
		DatabaseUrl:    getEnv("DATABASE_URL", "./proxy.db"),
		Addr:           getEnv("SERVICE_ADDR", ":8081"),
		SecretHash:     secretHash,
	}
	p.Run()
}
