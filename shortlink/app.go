package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
)

type App struct {
	Router     *mux.Router
	Middleware *Middleware
	Config     *Env
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
	app.Config = getEnv()
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	m := alice.New(app.Middleware.LoggerMiddleware, app.Middleware.RecoverMiddleware)
	app.Router.Handle("/api/shorten", m.ThenFunc(app.shortenHandler)).Methods("POST")
	app.Router.Handle("/api/info", m.ThenFunc(app.infoHandler)).Methods("GET")
	app.Router.Handle("/{shortlink}", m.ThenFunc(app.shortlinkHandler)).Methods("GET")
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
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
	w.WriteHeader(code)
}
func (app *App) shortenHandler(w http.ResponseWriter, r *http.Request) {
	var req shortenReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("error: %+v\n", err)
		writeError(w, &StatusError{http.StatusBadRequest, fmt.Errorf("bad json")})
		return
	}
	defer r.Body.Close()
	url, err := app.Config.S.Shorten(req.URL, req.Expired)
	if err != nil {
		panic(err)
	}
	writeJson(w, 200, url)
	return
}

func (app *App) infoHandler(w http.ResponseWriter, r *http.Request) {
	vals := r.URL.Query()
	shortlink := vals.Get("shortlink")
	detail, err := app.Config.S.Detail(shortlink)
	if err != nil {
		panic(err)
	}
	writeJson(w, 200, detail)
	return
}

func (app *App) shortlinkHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shortlink := vars["shortlink"]
	log.Printf("shortlink: %s", shortlink)
	url, err := app.Config.S.UnShorten(shortlink)
	if err != nil {
		panic(err)
	}
	writeJson(w, 301, url)
	return
}

func (app *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, app.Router))
}
