package main

import (
	//"context"
	//"fmt"

	"net/http"

	"encoding/json"

	"github.com/gorilla/mux"
)

var task []Message

func GetHandler(w http.ResponseWriter, r *http.Request) {

	find := DB.Find(&task)
	if find.Error != nil {
		http.Error(w, "Ошибка при получении ", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

func PostHandler(w http.ResponseWriter, r *http.Request) {
	var task Message

	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		http.Error(w, "Ошибка при отправлении", http.StatusBadRequest)
		return
	}

	create := DB.Create(&task)
	if create.Error != nil {
		http.Error(w, "Ошибка при сохранении задачи", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)

}

func PatchHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var reqBody Message

	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		http.Error(w, "Ошибка при отправлении", http.StatusBadRequest)
		return
	}

	var task Message
	first := DB.First(&task, id)
	if first.Error != nil {
		http.Error(w, "Задача не найдена", http.StatusInternalServerError)
		return
	}

	// Обновляем только те поля, которые были переданы
	if reqBody.Task != "" {
		task.Task = reqBody.Task
	}
	task.IsDone = reqBody.IsDone

	save := DB.Save(&task)
	if save.Error != nil {
		http.Error(w, "Ошибка при обновлении задачи", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

// DeleteHandler - обработчик для удаления задачи по ID
func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	delete := DB.Delete(&Message{}, id)
	if delete.Error != nil {
		http.Error(w, "Ошибка при удалении задачи", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent) // Успешное удаление, без содержимого
}

func main() {

	InitDB()

	DB.AutoMigrate(&Message{})

	router := mux.NewRouter()
	// наше приложение будет слушать запросы на localhost:8080/api/hello

	router.HandleFunc("/api/task", PostHandler).Methods("POST") // Обработчик для POST

	router.HandleFunc("/api/hello", GetHandler).Methods("GET")

	router.HandleFunc("/api/task/{id}", PatchHandler).Methods("PATCH") // Обработчик для PATCH

	router.HandleFunc("/api/task/{id}", DeleteHandler).Methods("DELETE") // Обработчик для DELETE

	http.ListenAndServe(":8080", router)

}
