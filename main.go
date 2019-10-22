package main

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

type Page struct {
	Title string
	Body []byte
	Links map[string]template.URL
}

type EditPage struct {
	Page

}

const (
	JSONFILE string = "links.json"
)

func (p *Page) save() error {
	return ioutil.WriteFile("home.html", p.Body, 0600)
}

func loadPage(title string, links *map[string]template.URL) (*Page, error) {
	filename := title + ".html"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	if links == nil {
		links = new(map[string]template.URL)
	}
	return &Page{Title: title, Body: body, Links: *links}, nil
}

func editLinksHandler(w http.ResponseWriter, r *http.Request) {
	title := "templates/edit-links"
	links := getLinks()
	p, err := loadPage(title, &links)
    if err != nil {
        p = &Page{Title: title}
    }
    renderTemplate(w, "templates/edit-links", p)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	title := "templates/home"
	links := getLinks()
	p, err := loadPage(title, &links)
	if err != nil {
        p = &Page{Title: title}
	}
    renderTemplate(w, "templates/home", p)
}

func saveLinksHandler(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()
	links := make(map[string]template.URL)
	for _, value := range r.Form {
		name := value[0]
		link := value[1]
		links[name] = template.URL(link)
	}
	
	jsonFile, _ := json.MarshalIndent(links, "", " ")
	_ = ioutil.WriteFile(JSONFILE, jsonFile, 0644)
	
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func getLinks() map[string]template.URL {
	jsonFile, err := ioutil.ReadFile(JSONFILE)
	if err != nil {
		panic(err)
	}
	var links map[string]template.URL
	err = json.Unmarshal(jsonFile, &links)
	if err != nil {
		panic(err)
	}
	
	linksWithProtocol := addProtocol(links)
	return linksWithProtocol
}

func addProtocol(links map[string]template.URL) map[string]template.URL {
	linksWithProtocol := make(map[string]template.URL)

	// Set bool in env through docker? If so, may just need `if https {}`
	https := os.Getenv("HTTPS_LINKS")
	var protocol string
	if https == "true" {
		protocol = "https"
	} else {
		protocol = "http"
	}
	for linkName, link := range links {
		protocolLink := string(link)
		if !strings.Contains(string(link), protocol + "://") {
			protocolLink = protocol + "://" + string(link)
		}
		linksWithProtocol[linkName] = template.URL(protocolLink)
	}
	return linksWithProtocol
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
    t, _ := template.ParseFiles(tmpl + ".html")
	t.Execute(w, p)
}

func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/edit-links/", editLinksHandler)
	http.HandleFunc("/save/", saveLinksHandler)
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("static/css"))))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("static/js"))))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
