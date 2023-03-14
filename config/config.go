package config

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/MegaMindInKZ/task-techno.git/db"
	_ "github.com/lib/pq"
	"io/ioutil"
	"log"
	"os"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "200103287sdu"
	dbname   = "link_service"
)

func init() {
	initDB()
	insert()
}

func initDB() {
	psqlInfo := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname,
	)
	var err error
	db.DB, err = sql.Open("postgres", psqlInfo)

	if err != nil {
		log.Fatal(err)
	}

	st, ioErr := ioutil.ReadFile("setup.sql")
	if ioErr != nil {
		fmt.Println("Cannot read setup.sql")
		os.Exit(1)
	}
	if _, err := db.DB.Exec(string(st)); err != nil {
		fmt.Println(err)
	}
}

func insert() {
	jsonFile, err := os.Open("links.json")
	if err != nil {
		log.Fatal(err)
	}
	defer jsonFile.Close()

	var links []db.Link
	byteValue, _ := ioutil.ReadAll(jsonFile)

	err = json.Unmarshal(byteValue, &links)
	if err != nil {
		log.Fatal(err)
	}

	for _, link := range links {
		link.Create()
	}
}

func End() {
	db.DB.Close()
}
