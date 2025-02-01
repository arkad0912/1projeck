package main

import (
	"fmt"
	"net/http"

	"encoding/json"

	"github.com/gorilla/mux"
)

var task string = "по умолчанию"

type requestBody struct {
	Task string `json:"task"` // Поле Message для получения данных из JSON
}

func HelloHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintf(w, "hello, %s", task)

}

func TaskHandler(w http.ResponseWriter, r *http.Request) {
	var reqBody requestBody

	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	task = reqBody.Task

	fmt.Fprintf(w, "Сообщение обновлено: %s\n", task)

}

func main() {
	router := mux.NewRouter()
	// наше приложение будет слушать запросы на localhost:8080/api/hello

	router.HandleFunc("/api/task", TaskHandler).Methods("POST") // Обработчик для POST

	router.HandleFunc("/api/hello", HelloHandler).Methods("GET")

	http.ListenAndServe(":8080", router)

}
