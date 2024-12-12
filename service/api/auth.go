package api

import (
	"errors"
	"net/http"
	"wasa.project/service/api/reqcontext"
)

func checkAuth(w http.ResponseWriter, userId int, ctx reqcontext.RequestContext) error {
	if userId != ctx.UserId {
		Forbidden(w, nil, ctx, "The user is not authorized")
		return errors.New("The user is not authorized")
	}
	return nil
}
