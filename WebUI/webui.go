package webui

import (
	"context"
	"html/template"
	"net/http"
	logger "github.com/Diamon0/rns-babel/Logger"
)

// I figured I'd make a web UI first, it's easier to maintain and contribute to
// The plan is to make a Qt GUI later down the line

var templates *template.Template
const WEBDIR string = "WebUI/"

func init() {
    templates = template.Must(template.ParseFiles(WEBDIR+"templates/index.html"))

    http.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir(WEBDIR+"static"))))
    http.HandleFunc("GET /", Home)
}

// Why the name "signaler"?
// Because it can later be used to signal more things other than shutdown
func StartWeb(addr string, signaler chan int8) error {
    server := http.Server{
        Addr: addr,
    }
    isShutdown := make(chan struct{})

    go func() {
        <-signaler
        if err := server.Shutdown(context.Background()); err != nil {
            logger.DefaultLogger.Println(err)
        }
        logger.DefaultLogger.Println("Server stopped")
        close(isShutdown)
    }()

    logger.DefaultLogger.Println("Now listening on port", "http://"+addr)
    err := server.ListenAndServe()

    <-isShutdown
    return err
}

func Home(w http.ResponseWriter, r *http.Request) {
    if err := templates.ExecuteTemplate(w, "base", nil); err != nil {
        logger.DefaultLogger.Println("Could not execute Home templates:", err)
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }
    return
}
