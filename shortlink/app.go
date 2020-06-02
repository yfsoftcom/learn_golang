package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type App struct {
	Router *mux.Router
}

type shortenReq struct {
	URL     string `json:"url"`
	Expired int64  `json:"exp"`
}

type shortlinkResp struct {
	ShortLink string `json:"short_link"`
	CreateAt  string `json:"create_at"`
}

func (app *App) Initialize() {
	app.Router = mux.NewRouter()
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	app.Router.HandleFunc("/api/shorten", app.shortenHandler).Methods("POST")
	app.Router.HandleFunc("/api/info", app.infoHandler).Methods("GET")
	app.Router.HandleFunc("/{shortlink}", app.shortlinkHandler).Methods("GET")
}

func writeError(w http.ResponseWriter, err error) {
	switch e := err.(type) {
	case Error:
		// user defined error
		writeJson(w, e.Status(), e.Error())
	default:
		writeJson(w, http.StatusInternalServerError, e.Error())
	}
}

func writeJson(w http.ResponseWriter, code int, payload interface{}) {
	data, _ := json.Marshal(payload)
	w.Write(data)
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")
}
func (app *App) shortenHandler(w http.ResponseWriter, r *http.Request) {
	var req shortenReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("error: %+v\n", err)
		writeError(w, &StatusError{http.StatusBadRequest, fmt.Errorf("bad json")})
		return
	}
	defer r.Body.Close()
	log.Printf(`data: %s, %d`, req.URL, req.Expired)
	return
}

func (app *App) infoHandler(w http.ResponseWriter, r *http.Request) {
	vals := r.URL.Query()
	shortlink := vals.Get("shortlink")
	log.Printf("shortlink: %s", shortlink)

	return
}

func (app *App) shortlinkHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shortlink := vars["shortlink"]
	log.Printf("shortlink: %s", shortlink)
	return
}

func (app *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, app.Router))
}
