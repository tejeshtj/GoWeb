package main

/*We import 4 important libraries
1. “net/http” to access the core go http functionality
2. “fmt” for formatting our text
3. “html/template” a library that allows us to interact with our html file.
4. "time" - a library for working with date and time.*/
import (
	"fmt"
	"html/template"
	"net/http"

	userHandler "./operations"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/sessions"
)

//Create a struct that holds information to be displayed in our HTML file
type Welcome struct {
	Name string
	Time string
}
type ContactDetails struct {
	Email    string
	Password string
}

var (
	// key must be 16, 24 or 32 bytes long (AES-128, AES-192 or AES-256)
	key   = []byte("super-secret-key")
	store = sessions.NewCookieStore(key)
)

//Go application entrypoint
func main() {

	//welcome := Welcome{"Anonymous", time.Now().Format(time.Stamp)}

	//templates := template.Must(template.ParseFiles("templates/login.html", "templates/signup.html"))

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		templates := template.Must(template.ParseFiles("templates/login.html"))
		fmt.Println(r.FormValue("email"), r.FormValue("password"))
		/*if err := templates.ExecuteTemplate(w, "login.html", nil); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}*/
		if r.Method != http.MethodPost {
			templates.Execute(w, nil)
			return
		}
		details := ContactDetails{
			Email:    r.FormValue("email"),
			Password: r.FormValue("password"),
		}

		userID, role := userHandler.QueryData(details.Email, details.Password)
		fmt.Println(userID)
		if userID > 0 {
			templates.Execute(w, struct{ Success bool }{true})

			session, err := store.Get(r, "session-name")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			session.Values["user-id"] = userID
			session.Values["role"] = role
			// Save it before we write to the response/return from the handler.
			err = session.Save(r, w)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			fmt.Println("sesion value", session.Values["user-id"])
		}

	})
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		templates := template.Must(template.ParseFiles("templates/login.html"))
		if r.Method != http.MethodPost {
			templates.Execute(w, nil)
			return
		}
		details := ContactDetails{
			Email:    r.FormValue("email"),
			Password: r.FormValue("password"),
		}
		userID, role := userHandler.QueryData(details.Email, details.Password)

		if userID > 0 {
			templates.Execute(w, struct{ Success bool }{true})
			session, err := store.Get(r, "session-name")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			session.Values["user-id"] = userID
			session.Values["role"] = role
			// Save it before we write to the response/return from the handler.
			err = session.Save(r, w)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			fmt.Println("sesion value", session.Values["user-id"])
		} else {
			templates.Execute(w, nil)
		}

	})
	http.HandleFunc("/signup", func(w http.ResponseWriter, r *http.Request) {
		templates := template.Must(template.ParseFiles("templates/signup.html"))

		if r.Method != http.MethodPost {
			templates.Execute(w, nil)
			return
		}
		details := ContactDetails{
			Email:    r.FormValue("email"),
			Password: r.FormValue("password"),
		}
		userID := userHandler.InsertData(details.Email, details.Password)
		fmt.Println(userID)
		if userID > 0 {
			templates.Execute(w, struct{ Success bool }{true})
		} else {
			templates.Execute(w, struct{ Success bool }{false})
		}

	})

	fmt.Println("Listening")
	fmt.Println(http.ListenAndServe(":8080", nil))
	//userHandler.InsertData()
}
