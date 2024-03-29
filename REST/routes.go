package main

import (
	"encoding/json"
	"log"
	"net/http"
	"regexp"
	"strconv"

	"github.com/go-chi/chi"
	acgrpc "github.com/shanehowearth/nine/createarticles/integration/grpc/client/v1"
	acgrpcProto "github.com/shanehowearth/nine/createarticles/integration/grpc/proto/v1"
	ragrpc "github.com/shanehowearth/nine/readarticles/integration/grpc/client/v1"
	ragrpcProto "github.com/shanehowearth/nine/readarticles/integration/grpc/proto/v1"
)

// Create articles instance
var ac = acgrpc.CreateClient{Address: "creator:5100"} //TODO get url:port from env

// Read articles instance
var ra = ragrpc.ReadArticleClient{Address: "readarticle:5200"} //TODO get url:port from env

func articleRoutes(router *chi.Mux) {
	// Articles related routes
	router.Route("/articles", func(r chi.Router) {
		r.Post("/", CreateArticles)
		r.Route("/{articleID}", func(r2 chi.Router) {
			r2.Get("/", GetArticlesByID)
		})
	})

	// Tags related routes
	router.Route("/tags", func(r chi.Router) {
		r.Route("/{tagname}/{date}", func(r2 chi.Router) {
			r2.Get("/", GetArticlesByTag)
		})
	})
}

// GetArticlesByTag -
func GetArticlesByTag(w http.ResponseWriter, req *http.Request) {
	tagname := chi.URLParam(req, "tagname")
	// Disallow tagnames that contain characters not in the set [0-9a-zA-Z]
	if match, _ := regexp.MatchString("^[0-9a-zA-Z ]+$", tagname); !match {
		respondWithError(w, http.StatusBadRequest, "Bad tagname supplied.")
		return
	}
	date := chi.URLParam(req, "date")
	// validations
	var badDate bool
	// date must be an int and exactly 8 chars
	_, err := strconv.Atoi(date)
	if len(date) != 8 || err != nil {
		badDate = true
	}
	// ensure date is possible
	if v, _ := strconv.Atoi(date[4:6]); v > 12 {
		//bad month
		badDate = true
	}
	if v, _ := strconv.Atoi(date[6:]); v > 31 {
		//bad day
		badDate = true
	}
	// TODO check combinations are possible (eg. leap years, correct days for given month, etc)
	if badDate {
		log.Printf("An invalid date was supplied, date: %s", date)
		respondWithError(w, http.StatusInternalServerError, "Supplied Date is an incorrect format, it must be YYYYMMDD")
	}

	taginf, err := ra.GetTagInfo(&ragrpcProto.ArticleRequest{Date: date, Tag: tagname})
	if err != nil {

		log.Printf("Taginfo for tag: %s, date: %s, error: %v", tagname, date, err)
		respondWithError(w, http.StatusInternalServerError, "Unable to fetch tag information")
	}
	respondWithJSON(w, http.StatusOK, taginf)
}

// GetArticlesByID -
func GetArticlesByID(w http.ResponseWriter, req *http.Request) {

	id := chi.URLParam(req, "articleID")
	// validate (id can only be int32 for now)
	if _, err := strconv.Atoi(id); err != nil {
		log.Printf("An invalid article id was supplied, ID: %s Error: %v", id, err)
		respondWithError(w, http.StatusInternalServerError, "Supplied Article ID is an incorrect format")
	}

	article, err := ra.GetArticle(&ragrpcProto.ArticleRequest{Id: id})
	if err != nil {
		// TODO Check if the Read server is alive.
		log.Printf("An error occurred with GetArticlesByID, Error: %v", err)
		// We don't want the user to know about the inner workings of the application
		respondWithError(w, http.StatusInternalServerError, "An unexpected error has occurred, the issue has been reported to our engineers and will be looked into")
	}
	respondWithJSON(w, http.StatusOK, article)
}

// CreateArticles -
func CreateArticles(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var a *acgrpcProto.Article
	err := decoder.Decode(&a)
	if err != nil {
		// log the error
		log.Printf("An error occurred with CreateArticles, Error: %v", err)
		// We don't want the user to know about the inner workings of the application
		// But we do want to be able to uniquely identify the error
		respondWithError(w, http.StatusInternalServerError, "An unexpected error has occurred, the issue has been reported to our engineers and will be looked into")
	}
	// TODO validate Date
	if match, _ := regexp.MatchString("^[0-9]{4}-[0-9]{2}-[0-9]{2}$", a.GetDate()); !match {
		respondWithError(w, http.StatusBadRequest, "Bad date format supplied. Date must be YYYY-MM-DD")
		return
	}
	ack := ac.CreateArticle(a)
	if ack.GetErrormessage() != "" {

		// TODO Check if the Create server is alive.
		log.Printf("An error occurred with CreateArticles, Error: %v", err)
		respondWithError(w, http.StatusInternalServerError, ack.GetErrormessage())
		return
	}
	respondWithJSON(w, http.StatusOK, ack)

}
