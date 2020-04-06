// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"database/sql"
	"fmt"
	"mux-master"
	_ "mysql-master"
	"net/http"
	"strconv"
	"text/template"
	"time"
)

var templates = template.Must(template.ParseFiles("templates/edit.html", "templates/view.html"))

type basedatos struct {
	Campo   string
	Tipo    string
	Null    string
	Key     string
	Defecto string
	Extra   string
}

type datosbase struct {
	Strubase []basedatos
}

var usuario, clave, host, puerto, basedato, tabla string

type data struct {
	Title string
	Body  []byte
}

type pagedata struct {
	Datasql  []entry
	Userdata []entry
}

type entry struct {
	Number                               int
	ID, Codigo, Nombre, Valor, Actufecha string
}

func obtenerBaseDeDatos() (db *sql.DB, e error) {
	usuario = "pepe"
	clave = "P0CMCmnbBsnjWfIW"
	host = "192.168.0.201"
	puerto = "3306"
	basedato = "pruebas"
	tabla = "monedas"
	db, err := sql.Open("mysql", usuario+":"+clave+"@tcp("+host+":"+puerto+")/"+basedato)

	// Debe tener la forma usuario:contraseña@host/nombreBaseDeDatos
	//db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@%s/%s", usuario, clave, host, basedato))
	if err != nil {
		return nil, err
	}
	return db, nil

}

func inistruc() {

	db, err := obtenerBaseDeDatos()
	//db, err := sql.Open("mysql", "root:13001300@tcp(localhost:3306)/gilardi")
	if err != nil {
		fmt.Printf("Error obteniendo base de datos: %v", err)
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic

	}

	// Ahora vemos si tenemos conexión
	err = db.Ping()
	if err != nil {
		fmt.Printf("Error conectando: %v", err)
		return
	}
	// Listo, aquí ya podemos usar a db!
	fmt.Println("Conectado correctamente a Base de datos " + basedato)

	//tabla = "monedas"
	rows, err := db.Query("SHOW COLUMNS FROM " + basedato + "." + tabla)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer db.Close()

	// text := ""
	//campobd := basedatos{}
	camposbd := []basedatos{}
	for rows.Next() {
		c := new(basedatos)
		//var campo, tipo, null, key, defecto, extra string
		rows.Scan(&c.Campo, &c.Tipo, &c.Null, &c.Key, &c.Defecto, &c.Extra)
		//campobd.Campo = campo
		//campobd.Tipo = tipo
		//campobd.Null = null
		//campobd.Key = key
		//campobd.Defecto = defecto
		//campobd.Extra = extra

		camposbd = append(camposbd, *c)
	}
	//pstru := &datosbase{Strubase: camposbd}
	//fmt.Println(pstru)

	for j := 0; j <= len(camposbd)-1; j++ {
		//fmt.Println(camposbd[j].Campo + " " + camposbd[j].Tipo)
	}
	fmt.Println("Son " + strconv.Itoa(len(camposbd)) + " Campos en tabla ")
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *pagedata) {
	//p.Datasql = template.HTML()

	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func viewHandler(w http.ResponseWriter, r *http.Request) {

	db, err := obtenerBaseDeDatos()

	if err != nil {
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM " + tabla)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	// text := ""
	setDataRow := entry{}
	resultsRow := []entry{}

	i := 0
	for rows.Next() {
		var id, codigo, nombre, valor, actufecha string
		rows.Scan(&id, &codigo, &nombre, &valor, &actufecha)
		i++
		setDataRow.Number = i
		setDataRow.ID = id
		setDataRow.Codigo = codigo
		setDataRow.Nombre = nombre
		setDataRow.Valor = valor
		setDataRow.Actufecha = actufecha
		resultsRow = append(resultsRow, setDataRow)
	}

	if err = rows.Err(); err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	p := &pagedata{Datasql: resultsRow}
	w.Header().Set("Content-Type", "text/html")
	renderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	//fmt.Println(vars["id"])
	//listadrivers := sql.Drivers()
	//fmt.Println(listadrivers)
	db, err := obtenerBaseDeDatos()
	if err != nil {
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}
	defer db.Close()

	//SELECT DATA BY ID
	users, err := db.Query("SELECT * FROM " + tabla + " WHERE ID=" + vars["id"])
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer users.Close()
	setData := entry{}
	results := []entry{}

	for users.Next() {
		var id, codigo, nombre, valor, actufecha string
		users.Scan(&id, &codigo, &nombre, &valor, &actufecha)
		setData.ID = id
		setData.Codigo = codigo
		setData.Nombre = nombre
		setData.Valor = valor
		setData.Actufecha = actufecha
		results = append(results, setData)
	}

	//fmt.Println(results)

	//SELECT ALL DATA
	rows, err := db.Query("SELECT * FROM " + tabla)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer rows.Close()

	// Fetch rows
	setDataRow := entry{}
	resultsRow := []entry{}

	i := 0
	for rows.Next() {
		var id, codigo, nombre, valor, actufecha string
		rows.Scan(&id, &codigo, &nombre, &valor, &actufecha)
		i++
		setDataRow.Number = i
		setDataRow.ID = id
		setDataRow.Codigo = codigo
		setDataRow.Nombre = nombre
		setDataRow.Valor = valor
		setDataRow.Actufecha = actufecha
		resultsRow = append(resultsRow, setDataRow)
	}

	if err = rows.Err(); err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	p := &pagedata{Datasql: resultsRow, Userdata: results}
	w.Header().Set("Content-Type", "text/html")
	renderTemplate(w, "edit", p)
}

func updateHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("ID")
	codigo := r.FormValue("Codigo")
	nombre := r.FormValue("Nombre")
	valor := r.FormValue("Valor")
	actufecha := r.FormValue("Actufecha") + " " + r.FormValue("Hora")
	//hora := r.FormValue("Hora")

	db, err := obtenerBaseDeDatos()
	if err != nil {
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}
	defer db.Close()

	// Execute the query
	//fmt.Println(hora)
	stmtIns, err := db.Query("UPDATE " + tabla + " SET codigo='" + codigo + "', nombre='" + nombre + "', valor='" + valor + "', actufecha='" + actufecha + "' WHERE id=" + id)

	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer stmtIns.Close()

	http.Redirect(w, r, "/view/", http.StatusFound)
}
func saveHandler(w http.ResponseWriter, r *http.Request) {

	codigo := r.FormValue("Codigo")
	nombre := r.FormValue("Nombre")
	valor := r.FormValue("Valor")
	//actufecha := r.FormValue("Actufecha")

	db, err := obtenerBaseDeDatos()
	if err != nil {
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}
	defer db.Close()

	// Execute the query
	stmtIns, err := db.Query("INSERT INTO " + tabla + " (codigo, nombre, valor, actufecha) VALUES ('" + codigo + "','" + nombre + "','" + valor + "',now())")

	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer stmtIns.Close()

	http.Redirect(w, r, "/view/", http.StatusFound)
}
func deleteHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	//fmt.Println(vars["id"])

	db, err := obtenerBaseDeDatos()
	if err != nil {
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}
	defer db.Close()

	// Execute the query
	stmtIns, err := db.Query("DELETE FROM " + tabla + " WHERE id=" + vars["id"])

	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer stmtIns.Close()

	http.Redirect(w, r, "/view/", http.StatusFound)
}

func main() {
	r := mux.NewRouter()
	inistruc()
	r.HandleFunc("/", viewHandler)
	r.HandleFunc("/home", viewHandler)
	r.HandleFunc("/view/", viewHandler)
	r.HandleFunc("/edit/{id:[0-9]+}", editHandler)
	r.HandleFunc("/save/", saveHandler)
	r.HandleFunc("/update/", updateHandler)
	r.HandleFunc("/delete/{id:[0-9]+}", deleteHandler)
	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	fmt.Println("Servidor escuchando en puerto 5000")
	t := time.Now()
	//fmt.Println(t)

	//fecha := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d",
	fecha := fmt.Sprintf("%d-%02d-%02d",
		t.Year(), t.Month(), t.Day())
	fmt.Println(fecha)
	fmt.Println("--")
	http.ListenAndServe(":5000", r)

}
