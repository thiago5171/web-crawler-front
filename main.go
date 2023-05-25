package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	_ "os"
)

type CrawlerForm struct {
	URL         string `json:"url"`
	Email       string `json:"email"`
	NumberLinks int    `json:"number_links"`
}

func main() {
	http.HandleFunc("/", crawlerHandler)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Erro ao iniciar servidor:", err)
	}
}

func crawlerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Falha ao realizar o parse para form data", http.StatusBadRequest)
			return
		}

		url := r.Form.Get("url")
		email := r.Form.Get("email")
		numberLinks := r.Form.Get("number_links")

		var numLinks int
		fmt.Sscanf(numberLinks, "%d", &numLinks)

		// Create crawler form object
		crawlerForm := CrawlerForm{
			URL:         url,
			Email:       email,
			NumberLinks: numLinks,
		}

		jsonData, err := json.Marshal(crawlerForm)
		if err != nil {
			http.Error(w, "Falha ao converter para JSON", http.StatusInternalServerError)
			return
		}

		apiURL := "http://localhost:8000/api/search-link"
		resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			http.Error(w, "Erro ao realizar a requisição.", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusCreated {
			fmt.Fprintln(w, "Requisição bem sucedida.")
		} else {
			fmt.Fprintln(w, "Falaha ao realizar a requisição.")
		}
	} else {
		crawlerTemplate := `
			<html>
<head>
	<title>Web Crawler</title>
	<style>
		body {
			display: flex;
			justify-content: center;
			align-items: center;
			height: 100vh;
			margin: 0;
			font-family: Arial, sans-serif;
		}

		.container {
			display: flex;
			flex-direction: column;
			align-items: center;
			width: 300px;
			padding: 20px;
			background-color: #f1f1f1;
			border-radius: 5px;
		}

		label {
			display: block;
			margin-bottom: 10px;
		}

		input[type="text"],
		input[type="number"] {
			width: 100%;
			padding: 5px;
			margin-bottom: 10px;
			border: 1px solid #ccc;
			border-radius: 3px;
			box-sizing: border-box;
		}

		input[type="submit"] {
			width: 100%;
			padding: 10px;
			background-color: #4CAF50;
			border: none;
			color: white;
			cursor: pointer;
			border-radius: 3px;
			font-size: 16px;
		}

		input[type="submit"]:hover {
			background-color: #45a049;
		}
	</style>
</head>
<body>
	<div class="container">
		<h1>Crawler Form</h1>
		<form method="POST" action="/">
			<label for="url">URL:</label>
			<input type="text" name="url" id="url"><br>

			<label for="email">Email para receber os links:</label>
			<input type="text" name="email" id="email"><br>

			<label for="number_links">Número de links a serem encontrados:</label>
			<input type="number" name="number_links" id="number_links"><br>

			<input type="submit" value="Enviar">
		</form>
	</div>
</body>
</html>


		`

		tmpl := template.Must(template.New("crawler").Parse(crawlerTemplate))
		err := tmpl.Execute(w, nil)
		if err != nil {
			http.Error(w, "Failed to render crawler form", http.StatusInternalServerError)
			return
		}
	}
}
