package aniModule

import(
	"database/sql"
	"github.com/gorilla/sessions"
	"log"
	"net/http"
)

type user struct {
	Id int
	Username string
	Email string
	Password string
	Status string
}

var store = sessions.NewCookieStore([]byte("i don't give a fuck"))


func Login(w http.ResponseWriter, r *http.Request){
	if r.Method == http.MethodGet{
		http.ServeFile(w,r, "./temp/login.html")
	} else if r.Method == http.MethodPost {
		email := r.PostFormValue("email")
		password := r.PostFormValue("pass")

		db, err := sql.Open("sqlite3", "./anime.db")
		if err != nil {
			panic(err)
		}
		defer db.Close()
		rows, err := db.Query("select * from users where email = $1 and password = $2", email, password)
		if err != nil {
			panic(err)
		}
		defer rows.Close()
		user1 := []user{}

		for rows.Next(){
			u := user{}
			err := rows.Scan(&u.Id, &u.Username, &u.Email, &u.Password, &u.Status)
			if err != nil{
				log.Println(err)
				continue
			}
			user1 = append(user1, u)
		}

		if user1[0].Email == email && user1[0].Password == password{
			session, err := store.Get(r, "session-mdma")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// Set some session values.
			session.Values["userID"] = user1[0].Id

			// Save it before we write to the response/return from the handler.
			err = session.Save(r, w)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}
}

func Reg(w http.ResponseWriter, r *http.Request){
	if r.Method == http.MethodGet{
		http.ServeFile(w,r, "./temp/reg.html")
	} else if r.Method == http.MethodPost {
		username := r.PostFormValue("username")
		email := r.PostFormValue("email")
		password1 := r.PostFormValue("pass1")
		password2 := r.PostFormValue("pass2")

		if password1 == password2{
			db, err := sql.Open("sqlite3", "./anime.db")
			if err != nil {
				panic(err)
			}
			defer db.Close()
			result, err := db.Exec("insert into users (username, email, password, status) values ($1, $2, $3, $4)", username, email, password1, "user")
			if err != nil{
				panic(err)
			}

			log.Print(result.LastInsertId())

			http.Redirect(w,r, "/", 200)
		}
	}
}
