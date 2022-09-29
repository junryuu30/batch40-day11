package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"personal-web/connection"
	"strconv"
	"strings"
	"text/template"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

func main() {

	route := mux.NewRouter()

	connection.DatabaseConnect()

	//root for public
	route.PathPrefix("/public").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))

	route.HandleFunc("/", home).Methods("GET")
	route.HandleFunc("/contact", contact).Methods("GET")
	route.HandleFunc("/formAddProject", formAddProject).Methods("GET")
	route.HandleFunc("/projectDetail/{id}", projectDetail).Methods("GET")
	route.HandleFunc("/addProject", addProject).Methods("POST")
	route.HandleFunc("/delete-project/{id}", deleteProject).Methods("GET")
	// route.HandleFunc("/edit-project/{index}", editProject).Methods("GET")

	route.HandleFunc("/form-register", formRegister).Methods("GET")
	route.HandleFunc("/register", register).Methods("POST")

	route.HandleFunc("/form-login", formLogin).Methods("GET")
	route.HandleFunc("/login", login).Methods("POST")

	route.HandleFunc("/logout", logout).Methods("GET")

	fmt.Println("server running in port 8080")
	http.ListenAndServe("localhost:8080", route)
}

// func helloWorld(w http.ResponseWriter, r *http.Request) {
// 	w.Write([]byte("Hello World jihan hallo woy ayo pasti bisa"))
// }

type SessionData struct {
	IsLogin   bool
	UserName  string
	FlashData string
}

var Data = SessionData{}

type Project struct {
	ID           int
	Title        string
	Description  string
	Technologies string
	NodeJs       string
	Python       string
	ReactJs      string
	Golang       string
	// StartDate    string
	// EndDate      string
	// Duration     string
}

type User struct {
	ID       int
	Name     string
	Email    string
	Password string
}

// var dataProject = []Project{
// 	{
// 		Title:        "Hallo Title",
// 		Description:  "Ini deskripsinya",
// 		Technologies: "node-js",
// 		NodeJs:       "node-js",
// 		ReactJs:      "react",
// 		Golang:       "golang",
// 	},
// }

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	var tmpl, err = template.ParseFiles("views/index.html")

	if err != nil {
		w.Write([]byte("message:" + err.Error()))
	}

	var store = sessions.NewCookieStore([]byte("SESSION_KEY"))
	session, _ := store.Get(r, "SESSION_KEY")

	if session.Values["IsLogin"] != true {
		Data.IsLogin = false
	} else {
		Data.IsLogin = session.Values["IsLogin"].(bool)
		Data.UserName = session.Values["Name"].(string)
	}

	fm := session.Flashes("message")

	var flashes []string
	if len(fm) > 0 {
		session.Save(r, w)
		for _, f1 := range fm {
			flashes = append(flashes, f1.(string))
		}
	}

	Data.FlashData = strings.Join(flashes, "")

	data, _ := connection.Conn.Query(context.Background(), "SELECT id, title, description FROM table_project2 ORDER BY id DESC")
	fmt.Println(data)

	var result []Project
	for data.Next() {
		var each = Project{}

		var err = data.Scan(&each.ID, &each.Title, &each.Description)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		result = append(result, each)
	}

	resData := map[string]interface{}{
		"DataSession": Data,
		"Project":     result,
		"Pesan":       "Anda Berhasil Log In, Selamat Datang",
	}

	fmt.Println(result)
	w.WriteHeader(http.StatusOK)

	tmpl.Execute(w, resData)

}

func contact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	var tmpl, err = template.ParseFiles("views/contact.html")

	if err != nil {
		w.Write([]byte("message:" + err.Error()))
	}

	tmpl.Execute(w, nil)

}

func formAddProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	var tmpl, err = template.ParseFiles("views/addProject.html")

	if err != nil {
		w.Write([]byte("message:" + err.Error()))
	}

	tmpl.Execute(w, nil)

	// http.Redirect(w, r, "/", http.StatusMovedPermanently)

}

func addProject(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	title := r.PostForm.Get("inputName")
	description := r.PostForm.Get("description")
	// start_date := r.PostForm.Get("startDate")
	// end_date := r.PostForm.Get("endDate")

	// nodeJs := r.PostForm.Get("nodeJs")
	// python := r.PostForm.Get("python")
	// reactJs := r.PostForm.Get("react")
	// golang := r.PostForm.Get("golang")

	// layout := "2006-01-02"
	// dateStart, _ := time.Parse(layout, start_date)
	// dateEnd, _ := time.Parse(layout, end_date)

	// //duration = (dateEnd - dateStart) in go:
	// hours := dateEnd.Sub(dateStart).Hours()
	// daysInHours := (hours / 24)
	// monthInDay := (daysInHours / 30)
	// fmt.Println(daysInHours)
	// var duration string
	// var month, _ float64 = math.Modf(monthInDay)

	// if monthInDay > 0 {
	// 	duration = strconv.FormatFloat(month, 'f', 0, 64) + " Bulan"
	// } else if daysInHours <= 31 {
	// 	duration = strconv.FormatFloat(daysInHours, 'f', 0, 64) + " Hari"
	// }
	// fmt.Println("ini durasinya: ", duration)

	_, err = connection.Conn.Exec(context.Background(), "INSERT INTO table_project2(title, description) VALUES ($1, $2, $3, $4)", title, description)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message: " + err.Error()))
		return
	}

	http.Redirect(w, r, "/", http.StatusMovedPermanently)

}

func deleteProject(w http.ResponseWriter, r *http.Request) {

	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	// fmt.Println(index)

	_, err := connection.Conn.Exec(context.Background(), "DELETE FROM table_project2 WHERE id=$1", id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message: " + err.Error()))
	}

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

// func editProject(w http.ResponseWriter, r *http.Request) {

// 	index, _ := strconv.Atoi(mux.Vars(r)["index"])
// 	fmt.Println(index)

// }

func projectDetail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	var tmpl, err = template.ParseFiles("views/projectDetail.html")

	if err != nil {
		w.Write([]byte("message:" + err.Error()))
		return
	}

	var ProjectDetail = Project{}

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	err = connection.Conn.QueryRow(context.Background(), "SELECT id, title, description FROM table_project2 WHERE id=$1", id).Scan(&ProjectDetail.ID, &ProjectDetail.Title, &ProjectDetail.Description)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message: " + err.Error()))
	}
	// for i, data := range dataProject {
	// 	if i == index {
	// 		ProjectDetail = Project{
	// 			Title:       data.Title,
	// 			Description: data.Description,
	// 		}
	// 	}
	// }

	data := map[string]interface{}{
		"Project": ProjectDetail,
	}
	// fmt.Println(data)

	// data := map[string]interface{}{
	// 	"Title":   "Hello Title",
	// 	"Content": "Hello Content",
	// 	"Id":      index,
	// }

	tmpl.Execute(w, data)

}

func formRegister(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	var tmpl, err = template.ParseFiles("views/form-register.html")

	if err != nil {
		w.Write([]byte("message: " + err.Error()))
		return
	}

	tmpl.Execute(w, nil)

}

func register(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}
	var name = r.PostForm.Get("input-name")
	var email = r.PostForm.Get("input-email")
	var password = r.PostForm.Get("input-password")

	passwordHash, _ := bcrypt.GenerateFromPassword([]byte(password), 10)

	_, err = connection.Conn.Exec(context.Background(), "INSERT INTO tb_user2(name, email, password) VALUES ($1, $2, $3)", name, email, passwordHash)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	http.Redirect(w, r, "/form-login", http.StatusMovedPermanently)

}

func formLogin(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	var tmpl, err = template.ParseFiles("views/form-login.html")

	if err != nil {
		w.Write([]byte("message: " + err.Error()))
		return
	}

	var store = sessions.NewCookieStore([]byte("SESSION_KEY"))
	session, _ := store.Get(r, "SESSION_KEY")

	fm := session.Flashes("message")

	var flashes []string
	if len(fm) > 0 {
		session.Save(r, w)
		for _, f1 := range fm {
			// meamasukan flash message
			flashes = append(flashes, f1.(string))
		}
	}

	Data.FlashData = strings.Join(flashes, "")

	tmpl.Execute(w, nil)

}

func login(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	var email = r.PostForm.Get("input-email")
	var password = r.PostForm.Get("input-password")

	user := User{}

	err = connection.Conn.QueryRow(context.Background(), "SELECT * FROM tb_user2 WHERE email=$1", email).Scan(&user.ID, &user.Name, &user.Email, &user.Password)

	if err != nil {
		http.Redirect(w, r, "/form-login", http.StatusMovedPermanently)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("message: Email belum terdaftar" + err.Error()))

		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		http.Redirect(w, r, "/form-login", http.StatusMovedPermanently)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("message : Email belum terdaftar " + err.Error()))
		return
	}
	var store = sessions.NewCookieStore([]byte("SESSION_KEY"))
	session, _ := store.Get(r, "SESSION_KEY")

	session.Values["Name"] = user.Name
	session.Values["Email"] = user.Email
	session.Values["ID"] = user.ID
	session.Values["IsLogin"] = true
	session.Options.MaxAge = 10800 // 3 JAM in second

	session.AddFlash("succesfull login", "message")
	session.Save(r, w)

	http.Redirect(w, r, "/", http.StatusMovedPermanently)

}

func logout(w http.ResponseWriter, r *http.Request) {
	var store = sessions.NewCookieStore([]byte("SESSION_KEY"))
	session, _ := store.Get(r, "SESSION_KEY")
	session.Options.MaxAge = -1
	session.Save(r, w)

	http.Redirect(w, r, "/form-login", http.StatusSeeOther)

}
