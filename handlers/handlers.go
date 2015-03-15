package handlers

////////////////////////////////////////////////////////////////////////////////////////////////
/// Imports
////////////////////////////////////////////////////////////////////////////////////////////////

import(
	"io"
    "net/http"
	"github.com/gorilla/handlers"
)


////////////////////////////////////////////////////////////////////////////////////////////////
/// Handlers
////////////////////////////////////////////////////////////////////////////////////////////////

func LoggingHandler(out io.Writer) (func(h http.Handler) http.Handler) {
	return func(h http.Handler) http.Handler {
		return handlers.LoggingHandler(out, h)
	}
}