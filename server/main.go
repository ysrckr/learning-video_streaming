package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
)

var videoPath = "./videos/video.mp4"

const PORT = "8000"

func main() {

	http.HandleFunc("/videos", streamVideo)

	fmt.Println("Running on http://localhost:" + PORT)
	if err := http.ListenAndServe(":"+PORT, nil); err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}
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
