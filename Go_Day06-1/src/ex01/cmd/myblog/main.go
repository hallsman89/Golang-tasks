package main

import (
	"day06/ex01/internal/handlers"
	"log"
	"net/http"
)

const (
	PORT         = "8888"
	articlesHTML = "/home/hallsman/school21/Go_Day06-1/src/ex01/ui/html/articles.html"
	adminHTML    = "/home/hallsman/school21/Go_Day06-1/src/ex01/ui/html/admin.html"
)

func main() {
	log.Println("Starting the server...")

	img := http.FileServer(http.Dir("../../ui/img"))
	http.Handle("/img/", http.StripPrefix("/img/", img))

	md := http.FileServer(http.Dir("../../ui/md/"))
	http.Handle("/md/", http.StripPrefix("/md/", md))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.ServeFile(w, r, articlesHTML)
			return
		}
		handlers.HandleDefault(w, r)
	})

	http.HandleFunc("/admin", handlers.HandleAdmin)

	log.Println("Server started at PORT:", PORT+"...")
	log.Fatalln(http.ListenAndServe("localhost:"+PORT, nil))
}
