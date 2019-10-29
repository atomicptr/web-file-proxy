package proxy

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"net/http"
	"time"

	"golang.org/x/crypto/sha3"

	"github.com/fernet/fernet-go"
)

const hashSalt = "wfp_"

// Create a SHA3-256 hash
func hash(data []byte) string {
	data = append([]byte(hashSalt), data...)

	hasher := sha3.New256()
	hasher.Write(data)
	return hex.EncodeToString(hasher.Sum(nil))
}

// Determine the User IP from a given http.Request
func getUserIp(r *http.Request) string {
	ip := r.Header.Get("X-Real-Ip")
	if ip != "" {
		return ip
	}

	ip = r.Header.Get("X-Forwarded-For")
	if ip != "" {
		return ip
	}

	return r.RemoteAddr
}

// Generate a new random "proof" to encrypt with fernet
func (p *Proxy) generateNewAuthProof() error {
	p.authProof = make([]byte, 128)
	_, err := rand.Read(p.authProof)
	return err
}

// create a user key from the "secret hash" and the user ip
func (p *Proxy) createRequestKey(r *http.Request) string {
	userIp := getUserIp(r)
	return hash([]byte(userIp + p.SecretHash))
}

// Check if the current user is authenticated
func (p *Proxy) isUserAuthenticated(r *http.Request) bool {
	cookie, err := r.Cookie("token")
	if err != nil {
		return false
	}

	key, err := fernet.DecodeKeys(p.createRequestKey(r))
	if err != nil {
		return false
	}

	proof := fernet.VerifyAndDecrypt([]byte(cookie.Value), 30*time.Minute, key)

	if len(proof) != len(p.authProof) {
		return false
	}

	for i := range proof {
		if proof[i] != p.authProof[i] {
			return false
		}
	}

	return true
}

// Authenticate the current user
func (p *Proxy) authenticate(password string, w http.ResponseWriter, r *http.Request) error {
	key, err := fernet.DecodeKey(p.createRequestKey(r))
	if err != nil {
		return err
	}

	token, err := fernet.EncryptAndSign(p.authProof, key)
	if err != nil {
		return err
	}

	tokenCookie := http.Cookie{
		Name:  "token",
		Value: string(token),
	}

	http.SetCookie(w, &tokenCookie)

	log.Println("User with IP:", getUserIp(r), "sucessfully authenticated.")

	return nil
}
