package handler

import (
	"contact-go/helper/apperrors"
	"contact-go/helper/response"
	"contact-go/model"
	"contact-go/usecase"
	"encoding/json"
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
		panic(err)
	}

	if err := response.NewJsonResponse(w, http.StatusOK, "OK", contacts); err != nil {
		panic(err)
	}
}

func (handler *contactHTTPHandler) Add(w http.ResponseWriter, r *http.Request) {
	var contactRequest model.ContactRequest
	err := json.NewDecoder(r.Body).Decode(&contactRequest)
	if err != nil {
		panic(err)
	}

	if contactRequest.Name == "" {
		_ = response.NewJsonResponse(w, http.StatusBadRequest, apperrors.ErrContactNameNotValid, nil)
		return
	}

	if contactRequest.NoTelp == "" {
		_ = response.NewJsonResponse(w, http.StatusBadRequest, apperrors.ErrContactNoTelpNotValid, nil)
		return
	}

	contact, err := handler.ContactUC.Add(&contactRequest)
	if err != nil {
		code, message := apperrors.HandleAppError(err)
		_ = response.NewJsonResponse(w, code, message, nil)
		return
	}

	if err := response.NewJsonResponse(w, http.StatusCreated, "Created", contact); err != nil {
		panic(err)
	}
}

func (handler *contactHTTPHandler) Detail(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/contacts/")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		_ = response.NewJsonResponse(w, http.StatusBadRequest, apperrors.ErrContactIdNotValid, nil)
		return
	}

	if id <= 0 {
		_ = response.NewJsonResponse(w, http.StatusBadRequest, apperrors.ErrContactIdNotValid, nil)
		return
	}

	contact, err := handler.ContactUC.Detail(int64(id))
	if err != nil {
		code, message := apperrors.HandleAppError(err)
		_ = response.NewJsonResponse(w, code, message, nil)
		return
	}

	if err := response.NewJsonResponse(w, http.StatusOK, "OK", contact); err != nil {
		panic(err)
	}
}

func (handler *contactHTTPHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/contacts/")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		_ = response.NewJsonResponse(w, http.StatusBadRequest, apperrors.ErrContactIdNotValid, nil)
		return
	}

	if id <= 0 {
		_ = response.NewJsonResponse(w, http.StatusBadRequest, apperrors.ErrContactIdNotValid, nil)
		return
	}

	var contactRequest model.ContactRequest
	err = json.NewDecoder(r.Body).Decode(&contactRequest)
	if err != nil {
		_ = response.NewJsonResponse(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	if contactRequest.Name == "" {
		_ = response.NewJsonResponse(w, http.StatusBadRequest, apperrors.ErrContactNameNotValid, nil)
		return
	}

	if contactRequest.NoTelp == "" {
		_ = response.NewJsonResponse(w, http.StatusBadRequest, apperrors.ErrContactNoTelpNotValid, nil)
		return
	}

	contact, err := handler.ContactUC.Update(int64(id), &contactRequest)
	if err != nil {
		code, message := apperrors.HandleAppError(err)
		_ = response.NewJsonResponse(w, code, message, nil)
		return
	}

	if err := response.NewJsonResponse(w, http.StatusOK, "OK", contact); err != nil {
		panic(err)
	}
}

func (handler *contactHTTPHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/contacts/")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		_ = response.NewJsonResponse(w, http.StatusBadRequest, apperrors.ErrContactIdNotValid, nil)
		return
	}

	if id <= 0 {
		_ = response.NewJsonResponse(w, http.StatusBadRequest, apperrors.ErrContactIdNotValid, nil)
		return
	}

	err = handler.ContactUC.Delete(int64(id))
	if err != nil {
		code, message := apperrors.HandleAppError(err)
		_ = response.NewJsonResponse(w, code, message, nil)
		return
	}

	if err := response.NewJsonResponse(w, http.StatusOK, "OK", nil); err != nil {
		panic(err)
	}
}
