package main

import (
	"fmt"

	"net/http"

	"encoding/json"

	"github.com/gorilla/mux"
)

var task []Message

func GetHandler(w http.ResponseWriter, r *http.Request) {

	result := DB.Find(&task)
	if result.Error != nil {
		http.Error(w, "Ошибка при получении ", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

func PostHandler(w http.ResponseWriter, r *http.Request) {
	var reqBody Message

	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		http.Error(w, "Ошибка при отправлении", http.StatusBadRequest)
		return
	}

	task := Message{
		Task:   reqBody.Task,
		IsDone: reqBody.IsDone,
	}

	result := DB.Create(&task)
	if result.Error != nil {
		http.Error(w, "Ошибка при сохранении задачи", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Сообщение добавлено: %s\n", reqBody.Task)

}

func main() {

	InitDB()

	DB.AutoMigrate(&Message{})

	router := mux.NewRouter()
	// наше приложение будет слушать запросы на localhost:8080/api/hello

	router.HandleFunc("/api/task", PostHandler).Methods("POST") // Обработчик для POST

	router.HandleFunc("/api/hello", GetHandler).Methods("GET")

	http.ListenAndServe(":8080", router)

}
