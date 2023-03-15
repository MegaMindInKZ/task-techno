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
	host     = "db"
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
		"postgres://%v:%v@%v:%v/%v?sslmode=disable",
		user,
		password,
		host,
		port,
		dbname,
	)
	var err error
	db.DB, err = sql.Open("postgres", psqlInfo)

	if err != nil {
		log.Fatal(err)
		os.Exit(0)
	}

}

func insert() {
	st, ioErr := ioutil.ReadFile("links_table.sql")
	if ioErr != nil {
		fmt.Println("Cannot read links_table.sql")
		os.Exit(1)
	}
	if _, err := db.DB.Exec(string(st)); err != nil {
		fmt.Println(err)
		return
	}

	jsonFile, err := os.Open("./links.json")
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
