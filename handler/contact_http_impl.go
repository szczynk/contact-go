package handler

import (
	"contact-go/model"
	"contact-go/usecase"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type contactHTTPHandler struct {
	ContactUC usecase.ContactUsecase
}

func NewContactHTTPHandler(contactUC usecase.ContactUsecase) ContactHTTPHandler {
	return &contactHTTPHandler{
		ContactUC: contactUC,
	}
}

func (handler *contactHTTPHandler) List(w http.ResponseWriter, r *http.Request) {
	contacts, err := handler.ContactUC.List()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(&contacts)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (handler *contactHTTPHandler) Add(w http.ResponseWriter, r *http.Request) {
	// err := r.ParseForm()
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusBadRequest)
	// 	return
	// }

	// name := r.PostForm.Get("name")
	// if name == "" {
	// 	http.Error(w, "name yang dimasukkan tidak valid", http.StatusBadRequest)
	// 	return
	// }

	// noTelp := r.PostForm.Get("no_telp")
	// if noTelp == "" {
	// 	http.Error(w, "no_telp yang dimasukkan tidak valid", http.StatusBadRequest)
	// 	return
	// }

	var contactRequest model.ContactRequest
	err := json.NewDecoder(r.Body).Decode(&contactRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	contact, err := handler.ContactUC.Add(&contactRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// msg := fmt.Sprintf(`{ "message":"Berhasil add contact with id %d" }`, contact.ID)
	// w.WriteHeader(http.StatusOK)
	// w.Write([]byte(msg))

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(&contact)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (handler *contactHTTPHandler) Detail(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/contacts/")
	if idStr == "" {
		http.Error(w, "Invalid ID.", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID.", http.StatusBadRequest)
		return
	}

	contact, err := handler.ContactUC.Detail(int64(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(&contact)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (handler *contactHTTPHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/contacts/")
	if idStr == "" {
		http.Error(w, "Invalid ID.", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID.", http.StatusBadRequest)
		return
	}

	var contactRequest model.ContactRequest
	err = json.NewDecoder(r.Body).Decode(&contactRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	contact, err := handler.ContactUC.Update(int64(id), &contactRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(&contact)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (handler *contactHTTPHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/contacts/")
	if idStr == "" {
		http.Error(w, "Invalid ID.", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID.", http.StatusBadRequest)
		return
	}

	err = handler.ContactUC.Delete(int64(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	msg := fmt.Sprintf(`{ "message":"Berhasil delete contact with id %d" }`, id)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(msg))
}
