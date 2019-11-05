package main 

import (
	    "database/sql"
		_ "github.com/go-sql-driver/mysql"
		"log"
        "net/http"
		"text/template"
		//"fmt"
)

type student struct {
	Id  int
	Fname  string
	Lname string
}

func dbConn() (db *sql.DB) {

	dbDriver := "mysql"
	dbUser := "root"
	dbPass := ""
	dbName := "uni"
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		panic(err.Error())
}
	return db
}

var tmpl = template.Must(template.ParseGlob("form/*"))

func Index(w http.ResponseWriter, r *http.Request) {
		db := dbConn()
		selDB, err := db.Query("SELECT * FROM tbl_student ORDER BY Id DESC")
		if err != nil {
			panic(err.Error())
}
		st := student{}
		res := []student{}
		for selDB.Next() {
				var Id int
				var Fname, Lname string
				err = selDB.Scan(&Id, &Fname, &Lname)
				if err != nil {
					panic(err.Error())
}
				st.Id = Id
				st.Fname = Fname
				st.Lname = Lname
				res = append(res, st)
}
			tmpl.ExecuteTemplate(w, "Index", res)
			defer db.Close()
}
				
func Show(w http.ResponseWriter, r *http.Request) {
 	db := dbConn()
	nId := r.URL.Query().Get("id")
	selDB, err := db.Query("SELECT * FROM tbl_student WHERE Id=?", nId)
	if err != nil {
    	panic(err.Error())
}
	st := student{}
	for selDB.Next() {
		var Id int
		var Fname, Lname string
		err = selDB.Scan(&Id, &Fname, &Lname)
		if err != nil {
			panic(err.Error())
}
	    st.Id = Id
		st.Fname = Fname
		st.Lname = Lname
}
	tmpl.ExecuteTemplate(w, "Show", st)
	defer db.Close()
}

func Insert(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
    if r.Method == "POST" {
        Fname := r.FormValue("Fname")
        Lname := r.FormValue("Lname")
        insForm, err := db.Prepare("INSERT INTO tbl_student(Fname, Lname) VALUES(?,?)")
        if err != nil {
            panic(err.Error())
}
		_, err = insForm.Exec(Fname, Lname)
		if err != nil {
			panic(err.Error())
		}
		log.Println("INSERT: Name: " + Fname + " | lastname: " + Lname)
}
    defer db.Close()
    http.Redirect(w, r, "/", 301)
}

func New(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "New", nil)
}

		
func Edit(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	nId := r.URL.Query().Get("id")
	selDB, err := db.Query("SELECT * FROM tbl_student WHERE Id=?", nId)
	if err != nil {
		panic(err.Error())
}

	st := student{}
	for selDB.Next() {
		var Id int
		var Fname, Lname string
		err = selDB.Scan(&Id, &Fname, &Lname)
		if err != nil {
			panic(err.Error())
}
		st.Id = Id
		st.Fname = Fname
		st.Lname = Lname
}
	tmpl.ExecuteTemplate(w, "Edit", st)
	defer db.Close()
}

func Update(w http.ResponseWriter, r *http.Request) {
	 db := dbConn()
    if r.Method == "POST" {
        Fname := r.FormValue("Fname")
        Lname := r.FormValue("Lname")
        Id := r.FormValue("Id")
        insForm, err := db.Prepare("UPDATE tbl_student  SET Fname=?, Lname=? WHERE Id=?")
        if err != nil {
            panic(err.Error())
}
        insForm.Exec(Fname, Lname, Id)
		log.Println("UPDATE: Name: " + Fname + " | Lname: " + Lname)
}
    defer db.Close()
    http.Redirect(w, r, "/", 301)
}

func Delete(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	st := r.URL.Query().Get("id")
	delForm, err := db.Prepare("DELETE FROM tbl_student WHERE Id=?")
	if err != nil {
		panic(err.Error())
}
	delForm.Exec(st)
	log.Println("DELETE")
	defer db.Close()
	http.Redirect(w, r, "/", 301)
}

func main() {
	log.Println("Server started on: http://localhost:8080")
	http.HandleFunc("/", Index)
	http.HandleFunc("/show", Show)
	http.HandleFunc("/new", New)
	http.HandleFunc("/edit", Edit)
	http.HandleFunc("/insert", Insert)
	http.HandleFunc("/update", Update)
    http.HandleFunc("/delete", Delete)
	http.ListenAndServe(":8080", nil)
}
