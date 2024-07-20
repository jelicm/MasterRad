package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"projekat/model"
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

func (handler *AppHandler) AddDataItem(w http.ResponseWriter, r *http.Request) {
	var di DataItemDTO
	err := json.NewDecoder(r.Body).Decode(&di)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	dataItem := model.DataSpaceItem{Path: di.Path, Name: di.Name, Permissions: di.Permissions,
		Scheme: di.HasScheme, SizeKB: di.SizeKB, State: model.State(di.State)}

	rez, err := handler.appservice.CreateDataItem(di.AppID, &dataItem, di.Scheme, false)

	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(rez.Path + "/" + rez.Name))
}

func (handler *AppHandler) DeleteApp(w http.ResponseWriter, r *http.Request) {
	req := &struct {
		NamespaceId  string
		AplicationID string
	}{}
	err := readReq(req, r, w)
	if err != nil {
		return
	}

	err = handler.nsservice.DeleteAppDefault(req.NamespaceId, req.AplicationID)
	if err != nil {
		writeErrorResp(err, w)
		return
	}

	writeResp(nil, http.StatusCreated, w)
}

func (handler *AppHandler) CreateSoftlink(w http.ResponseWriter, r *http.Request) {
	var slDTO SoftlinkDTO
	err := json.NewDecoder(r.Body).Decode(&slDTO)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	app1 := model.Application{ApplicationId: slDTO.Application1Id, ParentNamespaceId: slDTO.Namespace1Id}
	app2 := model.Application{ApplicationId: slDTO.Application2Id, ParentNamespaceId: slDTO.Namespace2Id}
	rez, err := handler.appservice.CreateSoftlink(&app1, &app2, slDTO.DataItemPath)

	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(*rez))
}