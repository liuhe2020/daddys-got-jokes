package main

// import (
// 	"fmt"
// 	"html/template"
// 	"net/http"
// 	"os"
// 	"strings"
// )

// var (
// 	staticDir = getAbsDirPath() + "/static/"
// 	templates = template.Must(template.ParseFiles(
// 		staticDir + "/index.html",
// 	))
// )

// func getAbsDirPath() string {
// 	pwd, err := os.Getwd()
// 	if err != nil {
// 		fmt.Println(err)
// 		os.Exit(1)
// 	}
// 	return pwd
// }

// func indexHandler(w http.ResponseWriter, r *http.Request) {
// 	err := templates.ExecuteTemplate(w, "index.html", nil)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 	}
// }

// func staticHandler(w http.ResponseWriter, r *http.Request) {
// 	path := r.URL.Path
// 	if strings.HasSuffix(path, "js") {
// 		w.Header().Set("Content-Type", "text/javascript")
// 	}
// 	// make sure you reference the correct absolute path
// 	data, err := os.ReadFile(staticDir + path[1:])
// 	if err != nil {
// 		fmt.Print(err)
// 	}
// 	_, err = w.Write(data)
// 	if err != nil {
// 		fmt.Print(err)
// 	}
// }
