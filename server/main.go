package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var videoPath = "./videos/video.mp4"

const PORT = "8000"

type VideoInfo struct {
	VideoURL string `json:"video_url"`
}

func main() {

	http.HandleFunc("/videos", createSignedURL)
	http.HandleFunc("/videos/video.mp4", streamVideo)

	fmt.Println("Running on http://localhost:" + PORT)
	if err := http.ListenAndServe(":"+PORT, nil); err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}
}

func createSignedURL(w http.ResponseWriter, r *http.Request) {
	videoName := r.URL.Query().Get("video_name")
	expiresAt := time.Now().Add(15 * time.Minute).Unix()

	signature, err := HashAndSalt(fmt.Sprintf("%d%s", expiresAt, videoName), 6)
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	videoInfo := &VideoInfo{
		VideoURL: fmt.Sprintf("http://localhost:8000/videos/video.mp4?signature=%s&expires_at=%d", signature, expiresAt),
	}

	info, err := json.Marshal(videoInfo)
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(info)
	w.WriteHeader(http.StatusOK)

}

func streamVideo(w http.ResponseWriter, r *http.Request) {
	data, err := os.ReadFile(videoPath)
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusNotFound)
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
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusBadRequest)
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

func HashAndSalt(s string, cost int) (string, error) {
	var err error
	var hashed []byte
	if cost <= bcrypt.MinCost {
		hashed, err = bcrypt.GenerateFromPassword([]byte(s), bcrypt.DefaultCost)
		if err != nil {
			return "", errors.New("couldn't hash")
		}

		return string(hashed), nil
	}

	hashed, err = bcrypt.GenerateFromPassword([]byte(s), cost)
	if err != nil {
		return "", errors.New("couldn't hash")
	}

	return string(hashed), nil

}
