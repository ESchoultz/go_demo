package main

import (
	"html/template"
	"net/http"

	uuid "github.com/nu7hatch/gouuid"
)

type user struct {
	UserName string
	Password string
	First    string
	Last     string
}

var tpl *template.Template
var dbUsers = map[string]user{}      // user ID, user
var dbSessions = map[string]string{} // session ID, user ID
var u user

func init() {
	tpl = template.Must(template.ParseGlob("templates/*"))
	dbUsers["ethan@mail.com"] = user{"ethan@mail.com", "pass", "Ethan", "Schoultz"}
	dbUsers["mason@mail.com"] = user{"mason@mail.com", "pass", "Mason", "Hill"}
	dbUsers["cooper@mail.com"] = user{"cooper@mail.com", "pass", "Cooper", "Vasiliou"}
}

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/login", login)
	http.HandleFunc("/register", register)
	http.HandleFunc("/bar", bar)
	http.HandleFunc("/logout", logout)
	//http.Handle("/favicon.ico", http.NotFoundHandler())
	http.ListenAndServe(":8080", nil)
}

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

func bar(w http.ResponseWriter, req *http.Request) {

	// get cookie
	c, err := req.Cookie("session")
	if err != nil {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}
	un, ok := dbSessions[c.Value]
	if !ok {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}
	u := dbUsers[un]
	tpl.ExecuteTemplate(w, "bar.gohtml", u)
}

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

type msg struct {
	message string
}

func register(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		un := req.FormValue("username")
		p := req.FormValue("password")
		f := req.FormValue("firstname")
		l := req.FormValue("lastname")
		u = user{un, p, f, l}
		dbUsers[un] = u
		var msg msg
		msg.message = "Please Enter Valid Credentials!"
		tpl.ExecuteTemplate(w, "register.gohtml", msg)
		return
	}
	tpl.ExecuteTemplate(w, "register.gohtml", nil)
}

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
