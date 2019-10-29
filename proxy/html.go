package proxy

import (
	"github.com/atomicptr/web-file-proxy/link"
	"html/template"
	"io"
	"io/ioutil"
	"os"
)

func readFile(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}

	data, err := ioutil.ReadAll(f)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func renderTemplate(path string) (*template.Template, error) {
	templateString, err := readFile(path)
	if err != nil {
		return nil, err
	}

	tpl, err := template.New(path).Parse(templateString)
	if err != nil {
		return nil, err
	}

	return tpl, nil
}

func (p *Proxy) renderLinksPage(links []*link.Link, w io.Writer) error {
	tpl, err := renderTemplate("./templates/links.html")
	if err != nil {
		return err
	}

	return tpl.Execute(w, links)
}

func (p *Proxy) renderNewLinkPage(w io.Writer) error {
	tpl, err := renderTemplate("./templates/add-new-link.html")
	if err != nil {
		return err
	}

	return tpl.Execute(w, nil)
}

func (p *Proxy) renderLoginPage(w io.Writer) error {
	tpl, err := renderTemplate("./templates/login.html")
	if err != nil {
		return err
	}

	return tpl.Execute(w, nil)
}

func (p *Proxy) renderErrorPage(message string, w io.Writer) error {
	tpl, err := renderTemplate("./templates/error.html")
	if err != nil {
		return err
	}

	return tpl.Execute(w, message)
}
