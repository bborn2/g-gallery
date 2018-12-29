package main

import (
    "fmt"
    "net/http"
    "log"
    "html/template"
    "path"
    "flag"
    "os"
    "os/signal"
    "syscall"
    "strings"
    "context"
    "time"
    // "reflect"
    // "bytes"
    "github.com/GeertJohan/go.rice"
)


var basePath string

func IndexHandler(w http.ResponseWriter, r *http.Request){

    uri := r.URL.Path
 
    folders, err := listDir(path.Join(basePath, uri))

    if err != nil {
        fmt.Fprintf(w, "filepath.Walk() returned %v\n", err)
        return
    }
    
    templateBox, err := rice.FindBox("templates")
    if err != nil {
        log.Fatal(err)
    }

    // get file contents as string
    templateString, err := templateBox.String("index.html")
    if err != nil {
        log.Fatal(err)
    }

    funcMap := template.FuncMap{"makeUrl": makeUrl,  "inc": inc}


    tmplMessage, err := template.New("message").Funcs(funcMap).Parse(templateString)
    // tmplMessage, err := template.New("message").Funcs(funcMap).ParseFiles("index.html")
    if err != nil {
        log.Fatal(err)
    }

    data := indexData{Folders: folders, Breads: strings.Split(uri, "/")}
    tmplMessage.Execute(w, data)

}

func AlbumHandler(w http.ResponseWriter, r *http.Request){

    uri := r.URL.Path

    alb, err := listImage( path.Join(basePath, uri))

    templateBox, err := rice.FindBox("templates")
    if err != nil {
        log.Fatal(err)
    }

    // get file contents as string
    templateString, err := templateBox.String("album.html")
    if err != nil {
        log.Fatal(err)
    }

    tmplMessage, err := template.New("message").Parse(templateString)

    // tmplMessage, err := template.New("message").ParseFiles("album.html")

    
    if err != nil {
        log.Fatal(err)
    }

    tmplMessage.Execute(w, alb)
}


func startHttpServer(port *string) *http.Server{

    mux := http.NewServeMux()

    mux.Handle("/", http.HandlerFunc(IndexHandler))

    fh := http.FileServer(rice.MustFindBox("static").HTTPBox())
    mux.Handle("/static/", http.StripPrefix("/static/", fh))
    // mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))


    mux.Handle("/folder/", http.StripPrefix("/folder/", http.HandlerFunc(IndexHandler)))

    mux.Handle("/file/", http.StripPrefix("/file/", http.FileServer(http.Dir(basePath))))

    mux.Handle("/album/", http.StripPrefix("/album/", http.HandlerFunc(AlbumHandler)))

    srv := &http.Server{Addr: ":" + *port , Handler:mux}

    go func() {
        log.Println("Listening...")

        if err := srv.ListenAndServe(); err != http.ErrServerClosed {
            // cannot panic, because this probably is an intentional close
            log.Printf("Httpserver: ListenAndServe() error: %s", err)
        }

        log.Println("close...")
    }()

    return srv
}

func main() {

    port := flag.String("p", "7000", "http server port")
    dir := flag.String("d", "", "image folder")
 
    flag.Parse()

    if *dir == ""{
        fmt.Println("run as -d /mnt/gallery -p 7000")
        return
    }

    var err error

    basePath, err = getBasePath(dir)

    log.Println("base = " + basePath)

    if err != nil {
        log.Println(err)
        return
    }

    log.Println(basePath)

    stop := make(chan os.Signal, 1)
    // signal.Notify(_stop, os.Interrupt, os.Kill, syscall.SIGUSR1, syscall.SIGUSR2)
    signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

    srv := startHttpServer(port)
    <-stop

    timeout := 5*time.Second

    ctx, cancel := context.WithTimeout(context.Background(), timeout)
    defer cancel()

    log.Printf("\nShutdown with timeout: %s\n", timeout)
    
    log.Println("shutdown...")
    if err := srv.Shutdown(ctx); err != nil {
        panic(err) // failure/timeout shutting down the server gracefully
        log.Println(err )
    }else{
        log.Println("server stop")
    }

    log.Println("exit...")
}

