package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
)

type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join(
			"templates", t.filename)))
	})
	t.templ.Execute(w, nil)
	if r.Method == "POST" {
		src, hdr, err := r.FormFile("myFile")
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		defer src.Close()

		dst, err := os.Create(filepath.Join(tmpDir, hdr.Filename))
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		defer dst.Close()

		buf := make([]byte, 250*1024*1024)
		if _, err := io.CopyBuffer(dst, src, buf); err != nil {
			log.Fatal(err)
		}
	}
}

var tmpDir = os.TempDir()

func main() {
	fmt.Println("TEMP DIR:", tmpDir)
	room := newRoom()
	http.Handle("/", &templateHandler{filename: "chat.html"})
	http.Handle("/room", room)
	go room.run()
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
