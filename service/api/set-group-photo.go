package api

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strconv"

	"wasa.project/service/api/imageFunctions"

	"github.com/julienschmidt/httprouter"
	"wasa.project/service/api/reqcontext"
)

func (rt *_router) setGroupPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	userId, err := strconv.Atoi(ps.ByName("user"))
	if err != nil {
		BadRequest(w, err, ctx, "Bad Request")
		return
	}

	// Check if the user is authorized
	if userId != ctx.UserId {
		Forbidden(w, err, ctx, "Forbidden")
		return
	}

	// Take the group id from the endpoint
	groupId, err := strconv.Atoi(ps.ByName("group"))
	if err != nil {
		BadRequest(w, err, ctx, "Bad Request")
		return
	}

	// Check if the user is a member of the group
	isMember, err := rt.db.CheckMember(userId, groupId)
	if !isMember || err != nil {
		BadRequest(w, err, ctx, "The user is not a member of the group")
		return
	}

	// Check if the size of the image is less than 5MB
	err = r.ParseMultipartForm(5 << 20)
	if err != nil {
		BadRequest(w, err, ctx, "Image too big")
		return
	}

	// Access the file from the request
	file, _, err := r.FormFile("image")
	if err != nil {
		BadRequest(w, err, ctx, "Bad request")
		return
	}
	// Read the file
	data, err := io.ReadAll(file) // In data we have the image file taked in the request
	if err != nil {
		InternalServerError(w, err, ctx)
		return
	}
	// Check if the file is a jpeg
	fileType := http.DetectContentType(data)
	if fileType != "image/jpeg" {
		http.Error(w, "Bad Request wrong file type", http.StatusBadRequest)
		return
	}
	defer func() { err = file.Close() }()

	// Create the file
	path := imageFunctions.SetDefaultPhotoGroup(groupId) // Take the path of the image of the user profile
	err = os.WriteFile(path, data, 0644)                 // Write the new image in the path selected
	if err != nil {
		InternalServerError(w, err, ctx)
		return
	}

	// Crop the image
	err = imageFunctions.SaveAndCrop(path, 250, 250)
	if err != nil {
		InternalServerError(w, err, ctx)
	}

	// Resposne
	w.WriteHeader(http.StatusOK)
	w.Header().Set("content-type", "plain/text")
	if err := json.NewEncoder(w).Encode("Photo changed"); err != nil {
		InternalServerError(w, err, ctx)
		return
	}

}
