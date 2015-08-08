package main

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"image"
	"image/draw"
	_ "image/jpeg"
	"image/png"
)

func handleGetMarkerImage(res http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	email := users.Users[vars["user"]].EMail

	gravatarURL := fmt.Sprintf("http://www.gravatar.com/avatar/%x?s=50", md5.Sum([]byte(email)))
	gravatarRaw, err := http.Get(gravatarURL)
	if err != nil {
		http.Error(res, "Internal server error", http.StatusInternalServerError)
		fmt.Printf("Unable to fetch Gravatar: %s\n", err)
		return
	}
	defer gravatarRaw.Body.Close()
	gravatar, _, err := image.Decode(gravatarRaw.Body)
	if err != nil {
		http.Error(res, "Internal server error", http.StatusInternalServerError)
		fmt.Printf("Unable to decode Gravatar: %s\n", err)
		return
	}

	markerRaw, _ := Asset("assets/markerimage.png")
	marker, _, err := image.Decode(bytes.NewBuffer(markerRaw))
	if err != nil {
		http.Error(res, "Internal server error", http.StatusInternalServerError)
		fmt.Printf("Unable to decode marker: %s\n", err)
		return
	}

	dstRect := image.Rect(13, 9, 63, 59)
	dstImage := image.NewRGBA(marker.Bounds())

	draw.Draw(dstImage, marker.Bounds(), marker, image.Pt(0, 0), draw.Src)
	draw.Draw(dstImage, dstRect, gravatar, image.Pt(0, 0), draw.Src)

	res.Header().Set("Content-Type", "image/png")
	res.Header().Set("Cache-Control", "public, max-age=3600")
	png.Encode(res, dstImage)
}
