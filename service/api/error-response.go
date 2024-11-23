package api

import (
	"net/http"
	"wasa.project/service/api/reqcontext"
)

func InternalServerError(w http.ResponseWriter, err error, ctx reqcontext.RequestContext) {
	ctx.Logger.WithError(err).Error("can't encode the response")
	w.WriteHeader(http.StatusInternalServerError)
}

func BadRequest(w http.ResponseWriter, err error, ctx reqcontext.RequestContext, message string) {
	http.Error(w, "message: "+err.Error(), http.StatusBadRequest)
}
