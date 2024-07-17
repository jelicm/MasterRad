package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"projekat/service"

	"github.com/gorilla/mux"
	//"projekat/model"
)

type AppHandler struct {
	appservice *service.ApplicationService
	nsservice  *service.NamespaceService
}

func NewAppHandler(appservice *service.ApplicationService, nsservice *service.NamespaceService) *AppHandler {
	return &AppHandler{
		appservice: appservice,
		nsservice:  nsservice,
	}
}

func (handler *AppHandler) RunApp(w http.ResponseWriter, r *http.Request) {
	var app AppDTO
	err := json.NewDecoder(r.Body).Decode(&app)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	rez, err := handler.appservice.RunApplication(app.ApplicationId, app.ParentNamespaceId, app.SizeKB)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(rez.ApplicationId))
}

func (handler *AppHandler) RunDataDiscovery(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	nsID, ok := vars["nsId"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	items := handler.nsservice.RunDataDiscovery(nsID)
	fmt.Println(nsID)
	jsonResponse(items, w)

}
