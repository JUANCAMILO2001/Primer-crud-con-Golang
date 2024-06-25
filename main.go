package main

import (
	"database/sql"

	//"log"
	"net/http"
	"text/template"

	_ "github.com/go-sql-driver/mysql"
)

func conexionDB() (conexion *sql.DB) {
	Driver := "mysql"
	Usuario := "root"
	Contrasenia := ""
	Nombre := "prueba_goland"
	conexion, err := sql.Open(Driver, Usuario+":"+Contrasenia+"@tcp(127.0.0.1)/"+Nombre)
	if err != nil {
		panic(err.Error())
	}
	return conexion
}

var plantillas = template.Must(template.ParseGlob("plantillas/*"))

func main() {
	http.HandleFunc("/", Inicio)
	http.HandleFunc("/crear", Crear)
	http.HandleFunc("/insertar", Insertar)
	http.HandleFunc("/borrar", Borrar)
	http.HandleFunc("/edit", Edit)
	http.HandleFunc("/update", Update)
	//fmt.Println("Servidor corriendo")
	http.ListenAndServe(":8080", nil)
}

type Employees struct {
	Id    int
	Name  string
	Email string
}

func Inicio(w http.ResponseWriter, r *http.Request) {
	conexionEsrablecida := conexionDB()
	list, err := conexionEsrablecida.Query("SELECT * FROM employees")
	if err != nil {
		panic(err.Error())
	}
	employees := Employees{}
	arrayEmployees := []Employees{}
	for list.Next() {
		var id int
		var name, email string
		err = list.Scan(&id, &name, &email)
		if err != nil {
			panic(err.Error())
		}
		employees.Id = id
		employees.Name = name
		employees.Email = email

		arrayEmployees = append(arrayEmployees, employees)
	}
	//fmt.Println(arrayEmployees)
	plantillas.ExecuteTemplate(w, "inicio", arrayEmployees)
}

func Crear(w http.ResponseWriter, r *http.Request) {
	plantillas.ExecuteTemplate(w, "crear", nil)
}

func Insertar(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		name := r.FormValue("name")
		email := r.FormValue("email")
		conexionEsrablecida := conexionDB()
		insert, err := conexionEsrablecida.Prepare("INSERT INTO employees(name,email) VALUES(?,?)")
		if err != nil {
			panic(err.Error())
		}
		insert.Exec(name, email)
		http.Redirect(w, r, "/", 301)
	}
}

func Borrar(w http.ResponseWriter, r *http.Request) {
	idEmployee := r.URL.Query().Get("id")
	conexionEsrablecida := conexionDB()
	delete, err := conexionEsrablecida.Prepare("DELETE FROM employees WHERE id =?")
	if err != nil {
		panic(err.Error())
	}
	delete.Exec(idEmployee)
	http.Redirect(w, r, "/", 301)
}

func Edit(w http.ResponseWriter, r *http.Request) {
	idEmployee := r.URL.Query().Get("id")
	//fmt.Println(idEmployee)
	conexionEsrablecida := conexionDB()
	list, err := conexionEsrablecida.Query("SELECT * FROM employees WHERE id = ?", idEmployee)
	employee := Employees{}
	for list.Next() {
		var id int
		var name, email string
		err = list.Scan(&id, &name, &email)
		if err != nil {
			panic(err.Error())
		}
		employee.Id = id
		employee.Name = name
		employee.Email = email
	}
	plantillas.ExecuteTemplate(w, "edit", employee)
}

func Update(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		id := r.FormValue("id")
		name := r.FormValue("name")
		email := r.FormValue("email")
		conexionEsrablecida := conexionDB()
		update, err := conexionEsrablecida.Prepare(" UPDATE employees SET name = ? , email = ? WHERE id = ?")
		if err != nil {
			panic(err.Error())
		}
		update.Exec(name, email, id)
		http.Redirect(w, r, "/", 301)
	}
}
