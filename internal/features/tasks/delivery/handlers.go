package delivery

import (
	"chek-list/internal/features/tasks/service"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

type HandlerStruct struct {
	service *service.TaskService
}

func NewHandler(s *service.TaskService) *HandlerStruct {
	return &HandlerStruct{service: s}
}

func (s *HandlerStruct) CreateHTTP(w http.ResponseWriter, r *http.Request)  {
	var newTask TaskRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&newTask); err != nil {
		ErrorHTTP(w, err, http.StatusBadRequest)
		return
	}
	task, err := s.service.Create(r.Context(), newTask.Title, newTask.Description)
	if err != nil {
		ErrorHTTP(w, err, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(task); err != nil {
		ErrorHTTP(w, err, http.StatusInternalServerError)
		return
	}
}

func (s *HandlerStruct) GetAllHTTP(w http.ResponseWriter, r *http.Request)  {
	tasksList, err := s.service.GetAll(r.Context())
	if err != nil {
		ErrorHTTP(w, err, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(tasksList); err != nil {
		ErrorHTTP(w, err, http.StatusInternalServerError)
		return
	}
}

func (s *HandlerStruct) DeleteHTTP(w http.ResponseWriter, r *http.Request)  {
	id := r.PathValue("id")
	if id == "" {
		ErrorHTTP(w, errors.New("id is required"), http.StatusBadRequest)
		return
	}
	idInt, err := strconv.Atoi(id)
	if err != nil {
		ErrorHTTP(w, err, http.StatusBadRequest)
		return
	}
	if err := s.service.Delete(r.Context(), idInt); err != nil {
		ErrorHTTP(w, err, http.StatusInternalServerError)
		return
	}

}

func (s *HandlerStruct) DoneHTTP(w http.ResponseWriter, r *http.Request)  {
	id := r.PathValue("id")
	if id == "" {
		ErrorHTTP(w, errors.New("id is required"), http.StatusBadRequest)
		return
	}
	idInt, err := strconv.Atoi(id)
	if err != nil {
		ErrorHTTP(w, err, http.StatusBadRequest)
		return
	}
	task, err := s.service.Done(r.Context(), idInt)
	if err != nil {
		ErrorHTTP(w, err, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(task); err != nil {
		ErrorHTTP(w, err, http.StatusInternalServerError)
		return
	}
}

func ErrorHTTP(w http.ResponseWriter, err error, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	response := map[string]string{"error": err.Error()}
	if err_ := json.NewEncoder(w).Encode(response); err_ != nil {
		panic(err_)
	}
}
