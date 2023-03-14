package db

import (
	"errors"
)

type Link struct {
	ID          int    `json:"id"`
	ActiveLink  string `json:"active_link"`
	HistoryLink string `json:"history_link"`
}

const paginationSize = 30

func GetActiveLinkByLink(link string) (activeLink string, isValid bool) {
	if IsActiveLink(link) {
		return activeLink, true
	}
	return RetrieveActiveLink(link), false
}

func IsActiveLink(link string) bool {
	var activeLink string
	err := DB.QueryRow("SELECT ACTIVE_LINK FROM LINKS WHERE ACTIVE_LINK = $1", link).Scan(&activeLink)
	return err == nil
}

func RetrieveActiveLink(link string) string {
	var activeLink string
	DB.QueryRow("SELECT ACTIVE_LINK FROM LINKS WHERE HISTORY_LINK = $1", link).Scan(&activeLink)
	return activeLink
}

func (link *Link) Create() error {
	st, err := DB.Prepare("INSERT INTO LINKS(ACTIVE_LINK, HISTORY_LINK) VALUES ($1, $2) RETURNING ID")
	defer st.Close()
	if err != nil {
		return err
	}

	err = st.QueryRow(link.ActiveLink, link.HistoryLink).Scan(&link.ID)
	return err
}

func (link *Link) Update() error {
	st, err := DB.Prepare("UPDATE LINKS SET ACTIVE_LINK = $1, HISTORY_LINK = $2 WHERE ID = $3")
	if err != nil {
		return err
	}

	_, err = st.Exec(link.ActiveLink, link.HistoryLink, link.ID)
	return err
}

func (link *Link) Delete() error {
	st, err := DB.Prepare("DELETE FROM LINKS WHERE ID=$1")
	if err != nil {
		return err
	}
	_, err = st.Exec(link.ID)
	return err
}

func GetLinkByID(id int) (link Link, err error) {
	err = DB.QueryRow("SELECT ID, ACTIVE_LINK, HISTORY_LINK FROM LINKS WHERE ID=$1", id).Scan(
		&link.ID, &link.ActiveLink, &link.HistoryLink,
	)
	return
}

func getLinks() (links []Link, err error) {
	rows, err := DB.Query("SELECT ID, ACTIVE_LINK, HISTORY_LINK FROM LINKS")
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var link Link
		err = rows.Scan(&link.ID, &link.ActiveLink, &link.HistoryLink)
		links = append(links, link)
	}
	return
}

func GetLinksWithPagination(pageNumber int) (links []Link, err error) {
	links, err = getLinks()
	if err != nil {
		return nil, errors.New("Internal server error")
	}

	if len(links) < (pageNumber-1)*paginationSize {
		return nil, errors.New("Invalid page number")
	}
	if len(links) < pageNumber*paginationSize {
		return links[(pageNumber-1)*paginationSize:], nil
	}
	return links[(pageNumber-1)*paginationSize : pageNumber*paginationSize], nil
}
