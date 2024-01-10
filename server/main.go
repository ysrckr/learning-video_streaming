package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/kodefluence/aurelia"
	"github.com/rs/cors"
)

var videoPath = "./videos/video.mp4"

const PORT = "8000"

const SECRET = "super secret"

type VideoInfo struct {
	VideoURL string `json:"video_url"`
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/videos", createSignedURL)
	mux.HandleFunc("/videos/video.mp4", streamVideo)

	handler := cors.Default().Handler(mux)
	fmt.Println("Running on http://localhost:" + PORT)
	if err := http.ListenAndServe(":"+PORT, handler); err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}
}

func createSignedURL(w http.ResponseWriter, r *http.Request) {
	videoName := r.URL.Query().Get("video_name")
	expiresAt := time.Now().Add(15 * time.Minute).Unix()

	signature := aurelia.Hash(SECRET, fmt.Sprintf("%d%s", expiresAt, videoName))

	videoInfo := &VideoInfo{
		VideoURL: fmt.Sprintf("http://localhost:8000/videos/video.mp4?signature=%s&expires_at=%d", signature, expiresAt),
	}

	info, err := json.Marshal(videoInfo)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(info)

}

func streamVideo(w http.ResponseWriter, r *http.Request) {
	signature := r.URL.Query().Get("signature")
	expiresAt := r.URL.Query().Get("expires_at")

	if signature == "" || expiresAt == "" {
		message := "signature and expires_at cannot be empty"

		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(message))

		return
	}

	expiresAtUnix, err := strconv.Atoi(expiresAt)
	if err != nil {
		message := "invalid expires date"
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(message))
		return
	}

	if !aurelia.Authenticate(SECRET, fmt.Sprintf("%d%s", expiresAtUnix, "video.mp4"), signature) {
		message := "unauthorized"
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(message))
		return
	}

	if time.Now().After(time.Unix(int64(expiresAtUnix), 0)) {
		message := "video not found"
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(message))
		return
	}

	data, err := os.ReadFile(videoPath)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))
	}

	getNumber := regexp.MustCompile(`\d`)
	videoRange := r.Header.Get("Range")
	startString := ""
	for _, match := range getNumber.FindAllString(videoRange, -1) {
		startString += match
	}
	fileSize := len(data)
	chunkSize := 1024 * 1024
	start, err := strconv.Atoi(startString)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	}
	end := min((start + chunkSize), (fileSize - 1))

	contentLength := end - start

	w.Header().Set("Content-Type", "video/mp4")
	w.Header().Set("Content-Length", strconv.Itoa(contentLength))
	w.Header().Set("Content-Range", "bytes "+strconv.Itoa(start)+"-"+strconv.Itoa(end)+"/"+strconv.Itoa(fileSize))
	w.Header().Set("Accept-Ranges", "bytes")

	w.WriteHeader(http.StatusPartialContent)
	w.Write(data[start:end])
}