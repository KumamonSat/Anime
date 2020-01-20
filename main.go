package main

import (
	"./aniModule"
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"html/template"
	"log"
	"net/http"
	"strings"
)

//maybe i need it
//select *FROM anime WHERE  id >= 38 and id <= 47 AND genre like '%2019%'


type anime19 struct{
	Id int
	Title string
	Img string
	About string
	Genre string
	Video string
	TitleSM string
}

type animeGenre struct{
	Id int
	Title string
	Img string
	About string
	Genre string
	Video string
	TitleSM string
}

type watch struct{
	Id int
	Title string
	Img string
	About string
	Genre string
	Video string
	TitleSM string
}

type genre struct {
	Id int
	Genre string
}

type searchRes struct {
	Id int
	Title string
	Img string
	About string
	Genre string
	Video string
	TitleSM string
}

func main(){

	r := mux.NewRouter()
	r.HandleFunc("/",index)
	r.HandleFunc("/{id:[0-9]+}",anime)
	r.HandleFunc("/genre/{id:[0-9]+}", genr)
	r.HandleFunc("/search", search)
	r.HandleFunc("/ebalo", ebalo)
	r.HandleFunc("/reg", aniModule.Reg)
	r.HandleFunc("/login", aniModule.Login)
	r.HandleFunc("/admin", admin)

	r.PathPrefix("/static/").Handler(http.FileServer(http.Dir("./")))


	log.Print("Server started")
	err := http.ListenAndServe("127.0.0.1:5000", r)
	if err != nil {
		log.Print(err)
	}
}

func admin(w http.ResponseWriter, r *http.Request){
	http.ServeFile(w,r, "temp/admin.html")
}

func ebalo(w http.ResponseWriter, r *http.Request){
	http.ServeFile(w,r, "temp/ebalo.html")
}

func index(w http.ResponseWriter, r *http.Request){
	db, err := sql.Open("sqlite3", "anime.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	rows, err := db.Query("SELECT *FROM anime WHERE genre like '%2019%' LIMIT 10 OFFSET 10")
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	anime := []anime19{}

	for rows.Next(){
		w := anime19{}
		err := rows.Scan(&w.Id, &w.Title, &w.Img, &w.About, &w.Genre, &w.Video, &w.TitleSM)
		if err != nil{
			fmt.Println(err)
			continue
		}
		anime = append(anime, w)
	}

	for i := 0; i < len(anime); i++ {
		short := anime[i].About
		runes := []rune(short)
		// ... Convert back into a string from rune slice.
		short = string(runes[0:400])
		short = short + "..."
		anime[i].About = short
	}

	//for _, a := range animes19{
	//	fmt.Println(a.Title)
	//}


	tmpl, _ := template.ParseFiles("temp/index.html")
	_ = tmpl.Execute(w, anime)
}

func anime(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	id := vars["id"]

	db, err := sql.Open("sqlite3", "anime.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	rows, err := db.Query("select * from anime where id = $1", id)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	watchs := []watch{}

	for rows.Next(){
		w := watch{}
		err := rows.Scan(&w.Id, &w.Title, &w.Img, &w.About, &w.Genre, &w.Video, &w.TitleSM)
		if err != nil{
			fmt.Println(err)
			continue
		}
		watchs = append(watchs, w)
	}

	tmpl, _ := template.ParseFiles("temp/anime.html")
	_ = tmpl.Execute(w, watchs)
}

func genr(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	id := vars["id"]

	db, err := sql.Open("sqlite3", "anime.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	rows, err := db.Query("select * from genres where id = $1", id)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	genres := []genre{}

	for rows.Next(){
		g := genre{}
		err := rows.Scan(&g.Id, &g.Genre)
		if err != nil{
			fmt.Println(err)
			continue
		}
		genres = append(genres, g)
	}
	fmt.Print(genres[0].Genre)
	lol := "%"+genres[0].Genre+"%"
	rows, err = db.Query("select * from anime where genre like $1", lol)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	aGenre := []animeGenre{}

	for rows.Next(){
		ag := animeGenre{}
		err := rows.Scan(&ag.Id, &ag.Title, &ag.Img, &ag.About, &ag.Genre, &ag.Video, &ag.TitleSM)
		if err != nil{
			fmt.Println(err)
			continue
		}
		aGenre = append(aGenre, ag)
	}

	for i := 0; i < len(aGenre); i++ {
		short := aGenre[i].About
		runes := []rune(short)
		// ... Convert back into a string from rune slice.
		short = string(runes[0:400])
		short = short + "..."
		aGenre[i].About = short
	}

	//fmt.Print(aGenre[0].Title)
	tmpl, _ := template.ParseFiles("temp/genre.html")
	_ = tmpl.Execute(w, aGenre)
}

func search(w http.ResponseWriter, r *http.Request){
	anime := r.PostFormValue("search")
	anime1 := strings.ToLower(anime)
	anime2 := "%"+anime1+"%"
	log.Print(anime)

	db, err := sql.Open("sqlite3", "anime.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	row, err := db.Query("select * from anime where titleSM like $1", anime2)
	if err != nil {
		panic(err)
	}
	defer row.Close()
	searches := []searchRes{}

	for row.Next(){
		s := searchRes{}
		err := row.Scan(&s.Id, &s.Title, &s.Img, &s.About, &s.Genre, &s.Video, &s.TitleSM)
		if err != nil{
			fmt.Println(err)
			continue
		}
		searches = append(searches, s)
	}

	tmpl, _ := template.ParseFiles("temp/search.html")
	_ = tmpl.Execute(w, searches)
}