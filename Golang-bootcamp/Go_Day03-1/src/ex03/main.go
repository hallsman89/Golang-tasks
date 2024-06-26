package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"main/db"
	"math"
	"net/http"
	"strconv"
)

type Store interface {
	GetPlaces(limit int, offset int) ([]db.Place, int, error)
	GetPlacesRecommend(lat, lon float64) ([]db.Place, error)
}

var (
	base Store
	err  error
)

func init() {
	base, err = db.NewElasticsearch()
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	http.HandleFunc("/", home)
	http.HandleFunc("/api/places", apiPlaces)
	http.HandleFunc("/api/recommend", apiRecommend)
	log.Fatal(http.ListenAndServe(":8888", nil))
}

type Data struct {
	Places []db.Place
	Total  int
	Page   int
	Last   int
}

func home(w http.ResponseWriter, r *http.Request) {
	var res Data
	if res.Page, err = strconv.Atoi(r.URL.Query().Get("page")); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	limit := 10
	offset := (res.Page - 1) * limit

	res.Places, res.Total, err = base.GetPlaces(limit, offset)
	res.Last = int(math.Ceil(float64(res.Total) / float64(limit)))
	tmpl, err := template.New("index.html").Funcs(
		template.FuncMap{
			"sum": sum,
			"sub": sub,
		},
	).ParseFiles("template/index.html")

	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, res)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func sum(x, y int) int {
	return x + y
}

func sub(x, y int) int {
	return x - y
}

func apiPlaces(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var res Data
	pageStr := r.URL.Query().Get("page")
	if res.Page, err = strconv.Atoi(pageStr); err != nil {
		http.Error(w, fmt.Sprintf("\"error\" : Invalid 'page' value: '%s'", pageStr), http.StatusBadRequest)
		return
	}
	limit := 10
	offset := (res.Page - 1) * limit

	if res.Places, res.Total, err = base.GetPlaces(limit, offset); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	res.Last = int(math.Ceil(float64(res.Total) / float64(limit)))
	response, err := json.MarshalIndent(res, "", "    ")
	if err != nil {
		log.Println(err)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(response)
}

func apiRecommend(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	lat, err := strconv.ParseFloat(r.URL.Query().Get("lat"), 64)
	if err != nil {
		http.Error(w, "Invalid 'lat' number", http.StatusBadRequest)
	}
	lon, err := strconv.ParseFloat(r.URL.Query().Get("lon"), 64)
	if err != nil {
		http.Error(w, "Invalid 'lon' number", http.StatusBadRequest)
	}
	place, err := base.GetPlacesRecommend(lat, lon)
	if err != nil {
		log.Println(err)
		return
	}
	res := struct {
		Name   string     `json:"name"`
		Places []db.Place `json:"places"`
	}{
		Name:   "Recommendation",
		Places: place,
	}
	response, err := json.MarshalIndent(res, "", "    ")
	if err != nil {
		log.Println(err)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(response)
}
