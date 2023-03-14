package main

import (
	"encoding/json"
	"github.com/MegaMindInKZ/task-techno.git/db"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func GetRedirectsAdminHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")

	pageNumber, _ := strconv.Atoi(request.URL.Query().Get("page"))
	if pageNumber <= 0 {
		pageNumber = 1
	}

	links, err := db.GetLinksWithPagination(pageNumber)

	if err != nil && err.Error() == "Invalid page number" {
		SendErrorMessage(writer, err.Error())
		return
	}

	if err != nil {
		SendInternalServerErrorMessage(writer)
		return
	}
	SendMessage(writer, 200, links)
}

func GetRedirectAdminHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")

	id, err := strconv.Atoi(mux.Vars(request)["id"])
	if err != nil {
		SendErrorMessage(writer, "ID Format is invalid")
	}

	link, err := db.GetLinkByID(id)
	if err != nil {
		SendErrorMessage(writer, "Invalid id parameter")
	}
	SendMessage(writer, 200, link)
}

func PostRedirectAdminHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")

	var link db.Link
	if err := json.NewDecoder(request.Body).Decode(&link); err != nil {
		SendErrorMessage(writer, "Bad Request")
		return
	}

	if err := link.Create(); err != nil {
		SendInternalServerErrorMessage(writer)
	}

	writer.WriteHeader(http.StatusOK)

}

func PatchRedirectAdminHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")

	var link db.Link
	var err error
	link.ID, err = strconv.Atoi(mux.Vars(request)["id"])

	if err != nil {
		SendErrorMessage(writer, "ID Format is invalid")
	}

	if err := json.NewDecoder(request.Body).Decode(&link); err != nil {
		SendErrorMessage(writer, "Bad Request")
		return
	}

	err = link.Update()

	if err != nil {
		SendInternalServerErrorMessage(writer)
	}
	SendMessage(writer, 201, link)
}

func DeleteRedirectAdminHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")

	var link db.Link
	var err error
	link.ID, err = strconv.Atoi(mux.Vars(request)["id"])

	if err != nil {
		SendErrorMessage(writer, "ID Format is invalid")
	}

	err = link.Delete()

	if err != nil {
		SendInternalServerErrorMessage(writer)
	}

	SendMessage(writer, 204, nil)
}

func GetRedirectHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")

	link := request.URL.Query().Get("link")
	activeLink, isActive := db.GetActiveLinkByLink(link)
	if isActive {
		SendMessage(writer, 200, nil)
	} else {
		data := struct {
			Active_link string `json:"active_link"`
		}{
			Active_link: activeLink,
		}
		SendMessage(writer, 301, data)
	}
}
