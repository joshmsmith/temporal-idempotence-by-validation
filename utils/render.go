package utils

import (
	"log"

	"html/template"
	"net/http"
)

// Render with header and footer templates included
func Render(w http.ResponseWriter, filename string, data interface{}) {

	tmpl, _ := template.ParseFiles("templates/Header.html")
	tmpl.Execute(w, nil)

	tmpl, err := template.ParseFiles(filename)
	if err != nil {
		log.Println("Render: Failed to ParseFiles for template:", filename, err)
		http.Error(w, "<br><br><h3>Sorry, something went wrong parsing template</h3>", http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, data); err != nil {
		log.Print(err)
		http.Error(w, "<br><br><h3>Sorry, something went wrong executing template</h3>", http.StatusInternalServerError)
	}

	tmpl, _ = template.ParseFiles("templates/Footer.html")
	tmpl.Execute(w, nil)
}

// RawRender just the file passed
func RawRender(w http.ResponseWriter, filename string, data interface{}) {

	tmpl, err := template.ParseFiles(filename)
	if err != nil {
		log.Println("Render: Failed to ParseFiles for template:", filename, err)
		http.Error(w, "Sorry, something went wrong parsing template", http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, data); err != nil {
		log.Print(err)
		http.Error(w, "Sorry, something went wrong executing template", http.StatusInternalServerError)
	}
}
