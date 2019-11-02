package proxy

import (
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (p *Proxy) initHttp() error {
	r := mux.NewRouter()
	r.HandleFunc("/", p.handleIndex)
	r.HandleFunc("/new", p.handleNew)
	r.HandleFunc("/actions-auth", p.handleAuthAction)
	r.HandleFunc("/actions-add", p.handleAddLinkAction)
	r.HandleFunc("/actions-delete", p.handleDeleteLinkAction)
	r.HandleFunc("/p/{ident}", p.handleProxy)

	http.Handle("/", r)
	return http.ListenAndServe(p.Addr, nil)
}

func (p *Proxy) redirectToIndex(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Hostname()
	http.Redirect(w, r, url, 302)
}

func (p *Proxy) handleIndex(w http.ResponseWriter, r *http.Request) {
	if !p.isUserAuthenticated(r) {
		err := p.renderLoginPage(w)
		if err != nil {
			log.Println(err)
		}
		return
	}

	links, err := p.linkRepository.FindAll()
	if err != nil {
		log.Println(err)
		p.handleError(500, "Could not reach database", w)
		return
	}

	w.WriteHeader(200)
	w.Header().Add("Content-Type", "text/html")

	err = p.renderLinksPage(links, w)
	if err != nil {
		log.Println(err)
	}
}

func (p *Proxy) handleNew(w http.ResponseWriter, r *http.Request) {
	if !p.isUserAuthenticated(r) {
		p.redirectToIndex(w, r)
		return
	}

	w.WriteHeader(200)
	w.Header().Add("Content-Type", "text/html")

	err := p.renderNewLinkPage(w)
	if err != nil {
		log.Println(err)
	}
}

func (p *Proxy) handleAuthAction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		p.redirectToIndex(w, r)
		return
	}

	password := r.FormValue("password")
	passwordHash := hash([]byte(password))

	if passwordHash == p.SecretHash {
		err := p.authenticate(password, w, r)
		if err != nil {
			log.Println(err)
		}
	}

	p.redirectToIndex(w, r)
}

func (p *Proxy) handleAddLinkAction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost || !p.isUserAuthenticated(r) {
		p.redirectToIndex(w, r)
		return
	}

	identifier := r.FormValue("identifier")
	url := r.FormValue("url")
	contentType := r.FormValue("content-type")

	err := p.linkRepository.InsertNew(identifier, url, contentType)
	if err != nil {
		log.Println(err)
		p.handleError(500, "Could not insert record", w)
		return
	}

	p.redirectToIndex(w, r)
}

func (p *Proxy) handleDeleteLinkAction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost || !p.isUserAuthenticated(r) {
		p.redirectToIndex(w, r)
		return
	}

	uid, err := strconv.Atoi(r.FormValue("uid"))
	if err != nil {
		log.Println(err)
		p.handleError(500, "Could not determine record", w)
		return
	}

	err = p.linkRepository.DeleteByUid(uid)
	if err != nil {
		log.Println(err)
		p.handleError(500, "Could not delete record", w)
		return
	}

	p.redirectToIndex(w, r)
}

func (p *Proxy) handleProxy(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ident := vars["ident"]

	link, err := p.linkRepository.FindByIdentifier(ident)
	if err != nil {
		log.Println(err)
		p.handleError(404, "Item not found", w)
		return
	}

	res, err := http.Get(link.Url)
	if err != nil {
		log.Println(err)
		p.handleError(500, "Target resource is unreachable", w)
		return
	}
	defer res.Body.Close()

	// if resource had a "Content-Type" set, use that
	if contentType := r.Header.Get("Content-Type"); contentType != "" {
		w.Header().Set("Content-Type", contentType)
	}

	// if user specified one, use it instead
	if link.ContentType != "" {
		w.Header().Set("Content-Type", link.ContentType)
	}

	w.WriteHeader(200)

	_, err = io.Copy(w, res.Body)
	if err != nil {
		log.Println(err)
		p.handleError(500, "Could not write resource", w)
		return
	}
}

func (p *Proxy) handleError(statusCode int, message string, w http.ResponseWriter) {
	w.WriteHeader(statusCode)
	w.Header().Add("Content-Type", "text/html")

	err := p.renderErrorPage(message, w)
	if err != nil {
		log.Println(err)
	}
}
