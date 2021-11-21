package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"text/template"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

type Todo struct {
	Id   int
	Task string
}

type Secret struct {
	User   string `json: "User"`
	Pass   string `json: "Pass"`
	DbName string `json: "DbName"`
}

var Db *sql.DB
var todo Todo

var tmpl = template.Must(template.ParseGlob("template/*"))

func ReadJson() (secret Secret) {
	jsonFile, err := os.Open("secret.json")
	if err != nil {
		panic(err)
	}
	defer jsonFile.Close()
	jsonData, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		panic(err)
	}
	json.Unmarshal(jsonData, &secret)
	return
}

func init() {
	secret := ReadJson()
	Db, err := gorm.Open("postgres", "user="+secret.User+" dbname="+secret.DbName+" password="+secret.Pass+" sslmode=disable")
	if err != nil {
		panic(err)
	}
	Db.AutoMigrate(&Todo{})
	result := Db.Find(&todo)
	fmt.Println(result)
}

func main() {
	server := http.Server{
		Addr:    "localhost:8080",
		Handler: nil,
	}
	http.HandleFunc("/", Index)
	server.ListenAndServe()
}

func Index(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "base", "")
}
