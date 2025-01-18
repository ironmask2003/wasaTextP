package api

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strconv"

	"wasa.project/service/api/imageFunctions"

	"github.com/julienschmidt/httprouter"
	"wasa.project/service/api/reqcontext"
)

func (rt *_router) setMyPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Get the user id from the request
	userId, err := strconv.Atoi(ps.ByName("user"))
	if err != nil {
		BadRequest(w, err, ctx, "Bad request")
		return
	}

	if checkAuth(w, userId, ctx) != nil {
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
		InternalServerError(w, err, "Error reading the image file", ctx)
		return
	}

	response := base64.StdEncoding.EncodeToString(data)

	// Check if the file is a jpeg
	fileType := http.DetectContentType(data)
	if fileType != "image/jpeg" {
		http.Error(w, "Bad Request wrong file type", http.StatusBadRequest)
		return
	}
	defer func() { err = file.Close() }()

	// Create the file
	path := imageFunctions.SetDefaultPhoto(userId) // Take the path of the image of the user profile
	err = os.WriteFile(path, data, 0644)           // Write the new image in the path selected
	if err != nil {
		InternalServerError(w, err, "Error setting the profile photo", ctx)
		return
	}

	// Crop the image
	err = imageFunctions.SaveAndCrop(path, 250, 250)
	if err != nil {
		InternalServerError(w, err, "Error cropping the image", ctx)
	}

	type Response struct {
		Photo string `json:"photo"`
	}

	var res Response
	res.Photo = response

	// Resposne
	w.WriteHeader(http.StatusOK)
	w.Header().Set("content-type", "plain/text")
	if err := json.NewEncoder(w).Encode(res); err != nil {
		InternalServerError(w, err, "Error encoding the response", ctx)
		return
	}
}
