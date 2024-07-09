package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"projekat/service"
	//"projekat/model"
)

type AppHandler struct {
	service *service.ApplicationService
}

func NewAppHandler(service *service.ApplicationService) *AppHandler {
	return &AppHandler{
		service: service,
	}
}

func (handler *AppHandler) RunApp(w http.ResponseWriter, r *http.Request) {
	var app AppDTO
	err := json.NewDecoder(r.Body).Decode(&app)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	rez, err := handler.service.RunApplication(app.ApplicationId, app.ParentNamespaceId, app.SizeKB)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(rez.ApplicationId))
}
