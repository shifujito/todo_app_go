package main

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"strconv"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

type Todo struct {
	Id   uint `gorm:"primary_key"`
	Task string
}

type Secret struct {
	User   string `json: "User"`
	Pass   string `json: "Pass"`
	DbName string `json: "DbName"`
}

var db *gorm.DB
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
	var err error
	secret := ReadJson()
	db, err := gorm.Open("postgres", "user="+secret.User+" dbname="+secret.DbName+" password="+secret.Pass+" sslmode=disable")
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&Todo{})
}

func ConnectDB() (db *gorm.DB) {
	secret := ReadJson()
	db, err := gorm.Open("postgres", "user="+secret.User+" dbname="+secret.DbName+" password="+secret.Pass+" sslmode=disable")
	if err != nil {
		panic("failed to connect database")
	}
	return
}

func main() {
	server := http.Server{
		Addr:    ":" + os.Getenv("PORT"),
		Handler: nil,
	}
	http.HandleFunc("/index", Index)
	http.HandleFunc("/create", Create)
	http.HandleFunc("/delete", Delete)
	http.HandleFunc("/edit", Update)
	server.ListenAndServe()
}

func Index(w http.ResponseWriter, r *http.Request) {
	db := ConnectDB()
	query := []Todo{}
	db.Find(&query)
	sort.Slice(query, func(i, j int) bool {
		return query[i].Id < query[j].Id
	})
	tmpl.ExecuteTemplate(w, "base", query)
	defer db.Close()
}

func Create(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tmpl.ExecuteTemplate(w, "new", "")
	} else if r.Method == "POST" {
		db := ConnectDB()
		task := template.HTML(r.FormValue("task"))
		newTodo := Todo{Task: string(task)}
		db.Create(&newTodo)
		http.Redirect(w, r, "/index", 301)
	}
}

func Delete(w http.ResponseWriter, r *http.Request) {
	deleteId, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		panic("user id is not intger but string")
	}
	db := ConnectDB()
	db.Delete(&Todo{}, deleteId)
	http.Redirect(w, r, "/index", 301)
	defer db.Close()
}

func Update(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		UpContent := Todo{}
		updateId, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil {
			panic(err)
		}
		db := ConnectDB()
		db.First(&UpContent, updateId)
		tmpl.ExecuteTemplate(w, "edit", UpContent)
		defer db.Close()
	} else if r.Method == "POST" {
		task := template.HTML(r.FormValue("task"))
		id := r.FormValue("id")
		intId, _ := strconv.Atoi(id)
		newTodo := Todo{Id: uint(intId)}
		db := ConnectDB()
		db.Model(&newTodo).Update("task", task)
		http.Redirect(w, r, "/index", 301)
	}
}
