package controllers

////////////////////////////////////////////////////////////////////////////////////////////////
/// Imports
////////////////////////////////////////////////////////////////////////////////////////////////

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"net/http"
)


////////////////////////////////////////////////////////////////////////////////////////////////
/// Index Controller
////////////////////////////////////////////////////////////////////////////////////////////////

type IndexController struct {
	*mux.Router
}

func NewIndexController(r *mux.Router) *IndexController {
	router := mux.NewRouter()

	ic := IndexController{}

	router.HandleFunc("/", ic.indexHandler)

	handlerChain := alice.New().Then(router)

	r.Handle("/", http.Handler(handlerChain))

	return &ic
}

func (c *IndexController) indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html")
	fmt.Fprint(w, "Hello World.")
}
