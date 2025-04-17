package web

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func Init() {
	// Регистрируем обработчик для маршрута "/"
	http.HandleFunc("/config", fileWrapper("config.json"))
	http.HandleFunc("/menu", fileWrapper("menu.json"))
	http.HandleFunc("/actions", fileWrapper("actions.json"))
	http.HandleFunc("/subitem", fileWrapper("subitem.json"))
	http.HandleFunc("/subitem2", fileWrapper("subitem2.json"))
	http.HandleFunc("/balance", fileWrapper("balance.json"))

	// Запуск сервера на порту 8080
	fmt.Println("Starting server at port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Server failed to start:", err)
	}
}

func fileWrapper(filename string) http.HandlerFunc {
	var resp = readFile(filename)

	return func(w http.ResponseWriter, r *http.Request) {
		params := r.URL.Query()

		// Read the body of the request
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Println("Failed to read request body")
		}
		// It's good practice to close the body reader
		defer r.Body.Close()

		// Convert the body to a string and print it
		jsonString := string(body)

		log.Printf("TEST WEB API INCOMING REQUEST :: params: %s, post data: %s\n", params, jsonString)

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(resp))
	}
}

func readFile(filename string) string {
	file, err := os.Open(fmt.Sprintf("pkg/web/resources/%s", filename))
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Read the file content into a byte slice
	content, err := io.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}
	return string(content)

}
