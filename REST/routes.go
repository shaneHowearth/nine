package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	acgrpc "github.com/shanehowearth/nine/createarticles/integration/grpc/client/v1"
	acgrpcProto "github.com/shanehowearth/nine/createarticles/integration/grpc/proto/v1"
	aidgrpc "github.com/shanehowearth/nine/readarticles/integration/grpc/client/v1"
	aidgrpcProto "github.com/shanehowearth/nine/readarticles/integration/grpc/proto/v1"
)

var ac = acgrpc.CreateClient{Address: "creator:5100"}            //TODO get url:port from env
var aid = aidgrpc.ReadArticleClient{Address: "readarticle:5200"} //TODO get url:port from env

func articleRoutes(router *chi.Mux) {
	router.Route("/articles", func(r chi.Router) {
		r.Post("/", CreateArticles)
		r.Route("/{articleID}", func(r2 chi.Router) {
			r2.Get("/", GetArticlesByID)
		})
	})
}

// GetArticlesByID -
func GetArticlesByID(w http.ResponseWriter, req *http.Request) {

	id := chi.URLParam(req, "articleID")
	// validate (id can only be int32 for now)
	if _, err := strconv.Atoi(id); err != nil {
		log.Printf("An invalid article id was supplied, ID: %s Error: %v", id, err)
		respondWithError(w, http.StatusInternalServerError, "Supplied Article ID is an incorrect format")
	}

	article, err := aid.GetArticle(&aidgrpcProto.ArticleRequest{Id: id})
	if err != nil {
		// log the error
		log.Printf("An error occurred with GetArticlesByID, Error: %v", err)
		// We don't want the user to know about the inner workings of the application
		// But we do want to be able to uniquely identify the error
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
	ack := ac.CreateArticle(a)
	respondWithJSON(w, http.StatusOK, ack)

}
