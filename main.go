package main

import (
	"html/template"
	"net/http"

	uuid "github.com/nu7hatch/gouuid"
)

// Declarations
var tpl *template.Template
var dbUsers = map[string]user{}
var dbSessions = map[string]string{}
var u user

// Declare user structure
type user struct {
	Username string
	Password string
	First    string
	Last     string
}

// Runs before main()
func init() {
	tpl = template.Must(template.ParseGlob("templates/*"))
	dbUsers["ethan@mail.com"] = user{"ethan@mail.com", "pass", "Ethan", "Schoultz"}
}

// Runs when binary is executed
func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/login", login)
	http.HandleFunc("/register", register)
	http.HandleFunc("/inventory", inventory)
	http.HandleFunc("/logout", logout)
	http.ListenAndServe(":8080", nil)
}

// Function to handle index get and post methods
func index(w http.ResponseWriter, req *http.Request) {
	// Request cookie
	c, err := req.Cookie("session")
	// If no cookies exists, create one
	if err != nil {
		sID, _ := uuid.NewV4()
		c = &http.Cookie{
			Name:  "session",
			Value: sID.String(),
		}
		// Set cookie
		http.SetCookie(w, c)
	}

	if un, ok := dbSessions[c.Value]; ok {
		u = dbUsers[un]
	}

	tpl.ExecuteTemplate(w, "index.gohtml", u)
}

// Function to handle /login get and post methods
func login(w http.ResponseWriter, req *http.Request) {

	if req.Method == http.MethodPost {
		un := req.FormValue("username")
		p := req.FormValue("password")
		// does this username exist
		u, ok := dbUsers[un]
		if !ok {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		// does the username/password combo have a match?
		if u.Password != p {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		// create a session
		sID, _ := uuid.NewV4()
		c := &http.Cookie{
			Name:  "session",
			Value: sID.String(),
		}
		http.SetCookie(w, c)
		dbSessions[c.Value] = un
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}
	tpl.ExecuteTemplate(w, "login.gohtml", nil)
}

// Function to handle /register get and post methods
func register(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		var un string = req.FormValue("username")
		var p string = req.FormValue("password")
		var f string = req.FormValue("firstname")
		var l string = req.FormValue("lastname")
		// u = user{un, p, f, l}
		//dbUsers[un] = u
		dbUsers[un] = user{un, p, f, l}
		http.Redirect(w, req, "/index", http.StatusSeeOther)
		return
	}
	tpl.ExecuteTemplate(w, "register.gohtml", nil)
}

// Function to handle /logout get and post methods
func logout(w http.ResponseWriter, req *http.Request) {
	c, err := req.Cookie("session")
	if err != nil {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}
	c = &http.Cookie{
		Name:   "session",
		Value:  "0",
		MaxAge: -1,
	}
	http.SetCookie(w, c)

	tpl.ExecuteTemplate(w, "index.gohtml", nil)
}

// Function to handle /books get and post methods
func inventory(w http.ResponseWriter, req *http.Request) {

}
