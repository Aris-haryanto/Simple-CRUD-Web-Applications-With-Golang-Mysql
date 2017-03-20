// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"text/template"
	"net/http"
	"database/sql"
	"fmt"
    _ "github.com/go-sql-driver"
    "github.com/gorilla/mux"
)
var templates = template.Must(template.ParseFiles("edit.html", "view.html"))

type Page struct {
	Title string
	Body  []byte
}
type Pagedata struct {
	// Title string
	Datasql []Entry
    Userdata []Entry
}

type Entry struct {
    Number int
    ID, Name, Email, Address string
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Pagedata) {
	// p.Datasql = template.HTML()

    err := templates.ExecuteTemplate(w, tmpl+".html", p)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

func editHandler(w http.ResponseWriter, r *http.Request) {
    
    vars := mux.Vars(r)
    fmt.Println(vars["id"]);
    db, err := sql.Open("mysql", "root@tcp(localhost:3306)/cc-dolcegusto") //localhost:3306
    if err != nil {
        panic(err.Error())  // Just for example purpose. You should use proper error handling instead of panic
    }
    defer db.Close()

    //SELECT DATA BY ID
    users, err := db.Query("SELECT * FROM users WHERE id="+vars["id"])
    if err != nil {
        panic(err.Error()) // proper error handling instead of panic in your app
    }
    defer users.Close()

    results := []Entry{}
    setData := Entry{}

    for users.Next() {
        var id, name, email, address string
        users.Scan(&id, &name, &email, &address)
        setData.ID = id
        setData.Name = name
        setData.Email = email
        setData.Address = address
        results = append(results, setData)
    }

    fmt.Println(results)


    //SELECT ALL DATA
    rows, err := db.Query("SELECT * FROM users")
    if err != nil {
        panic(err.Error()) // proper error handling instead of panic in your app
    }
    defer rows.Close()

    // Fetch rows

    resultsRow := []Entry{}
    setDataRow := Entry{}
    i := 0
    for rows.Next() {
        var id, name, email, address string
        rows.Scan(&id, &name, &email, &address)
        i += 1
        setDataRow.Number = i
        setDataRow.ID = id
        setDataRow.Name = name
        setDataRow.Email = email
        setDataRow.Address = address
        resultsRow = append(resultsRow, setDataRow)
    }

    if err = rows.Err(); err != nil {
        panic(err.Error()) // proper error handling instead of panic in your app
    }
    p := &Pagedata{Datasql: resultsRow, Userdata: results}
    w.Header().Set("Content-Type", "text/html")
    renderTemplate(w, "edit", p)
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	
	db, err := sql.Open("mysql", "root@tcp(localhost:3306)/cc-dolcegusto") //localhost:3306
    if err != nil {
        panic(err.Error())  // Just for example purpose. You should use proper error handling instead of panic
    }
    defer db.Close()

	rows, err := db.Query("SELECT * FROM users")
    if err != nil {
        panic(err.Error()) // proper error handling instead of panic in your app
    }

    // text := ""
    resultsRow := []Entry{}
    setDataRow := Entry{}
    i := 0
    for rows.Next() {
        var id, name, email, address string
        rows.Scan(&id, &name, &email, &address)
        i += 1
        setDataRow.Number = i
        setDataRow.ID = id
        setDataRow.Name = name
        setDataRow.Email = email
        setDataRow.Address = address
        resultsRow = append(resultsRow, setDataRow)
    }

    if err = rows.Err(); err != nil {
        panic(err.Error()) // proper error handling instead of panic in your app
    }
	p := &Pagedata{Datasql: resultsRow}
	w.Header().Set("Content-Type", "text/html")
	renderTemplate(w, "view", p)
}
func updateHandler(w http.ResponseWriter, r *http.Request){
    id := r.FormValue("id")
    name := r.FormValue("name")
    email := r.FormValue("email")
    address := r.FormValue("address")

    db, err := sql.Open("mysql", "root@tcp(localhost:3306)/cc-dolcegusto") //localhost:3306
    if err != nil {
        panic(err.Error())  // Just for example purpose. You should use proper error handling instead of panic
    }
    defer db.Close()

        // Execute the query
    stmtIns, err := db.Query("UPDATE users SET name='"+name+"', email='"+email+"', address='"+address+"' WHERE id="+id)
    
    if err != nil {
        panic(err.Error()) // proper error handling instead of panic in your app
    }
    defer stmtIns.Close()

    http.Redirect(w, r, "/view/", http.StatusFound)
}
func saveHandler(w http.ResponseWriter, r *http.Request) {

    name := r.FormValue("name")
    email := r.FormValue("email")
    address := r.FormValue("address")

    db, err := sql.Open("mysql", "root@tcp(localhost:3306)/cc-dolcegusto") //localhost:3306
    if err != nil {
        panic(err.Error())  // Just for example purpose. You should use proper error handling instead of panic
    }
    defer db.Close()

        // Execute the query
    stmtIns, err := db.Query("INSERT INTO users (name, email, address) VALUES ('"+name+"','"+email+"','"+address+"')")
    
    if err != nil {
        panic(err.Error()) // proper error handling instead of panic in your app
    }
    defer stmtIns.Close()

    http.Redirect(w, r, "/view/", http.StatusFound)
}
func deleteHandler(w http.ResponseWriter, r *http.Request){

    vars := mux.Vars(r)
    fmt.Println(vars["id"]);

    db, err := sql.Open("mysql", "root@tcp(localhost:3306)/cc-dolcegusto") //localhost:3306
    if err != nil {
        panic(err.Error())  // Just for example purpose. You should use proper error handling instead of panic
    }
    defer db.Close()

        // Execute the query
    stmtIns, err := db.Query("DELETE FROM users WHERE id="+vars["id"])
    
    if err != nil {
        panic(err.Error()) // proper error handling instead of panic in your app
    }
    defer stmtIns.Close()

    http.Redirect(w, r, "/view/", http.StatusFound)
}

func main() {
    r := mux.NewRouter()
    r.HandleFunc("/", viewHandler)
	r.HandleFunc("/view/", viewHandler)
	r.HandleFunc("/edit/{id:[0-9]+}", editHandler)
	r.HandleFunc("/save/", saveHandler)
    r.HandleFunc("/update/", updateHandler)
    r.HandleFunc("/delete/{id:[0-9]+}", deleteHandler)
    r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
    http.ListenAndServe(":8080", r)
}