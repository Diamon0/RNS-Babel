package webui

import (
	"html/template"
	"log"
	"net/http"
)

// I figured I'd make a web UI first, it's easier to maintain and contribute to
// The plan is to make a Qt GUI later down the line

var templates *template.Template
const WEBDIR string = "WebUI/"

func init() {
    templates = template.Must(template.ParseFiles(WEBDIR+"templates/index.html"))
}

func StartWeb(addr string) error {
    http.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir(WEBDIR+"static"))))
    http.HandleFunc("GET /", Home)

    log.Println("Now listening on port", "http://"+addr)
    return http.ListenAndServe(addr, nil)
}

func Home(w http.ResponseWriter, r *http.Request) {
    if err := templates.ExecuteTemplate(w, "base", nil); err != nil {
        log.Println("Could not execute Home templates:", err)
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }
    return
}
