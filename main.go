package main

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	guuid "github.com/google/uuid"
	fileutil "github.com/jjkoh95/image-micro-web-framework/file-utils"
	imageutil "github.com/jjkoh95/image-micro-web-framework/image-utils"
)

const imageDir = "images"

func main() {
	handleRequest()
}

func handleRequest() {
	// serving static files
	fs := http.FileServer(http.Dir("images"))
	http.Handle("/images/", http.StripPrefix("/images/", fs))

	// endpoints
	http.HandleFunc("/upload-image", uploadImage)
	http.HandleFunc("/upload-zip", uploadZip)
	http.HandleFunc("/generate-thumbnails", generateThumbnails)

	log.Println("Listening to port :3000")
	http.ListenAndServe(":3000", nil)
}

// uploadImage endpoint
func uploadImage(w http.ResponseWriter, r *http.Request) {
	// only allow POST method
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid Method", http.StatusMethodNotAllowed)
	}

	r.ParseMultipartForm(32 << 20) // limit max input length

	var buf bytes.Buffer

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "File error", http.StatusBadRequest)
		return
	}
	defer file.Close()
	io.Copy(&buf, file)

	imageBytes := buf.Bytes()

	isImage := imageutil.IsImage(imageBytes)
	if !isImage {
		http.Error(w, "Invalid image", http.StatusBadRequest)
		return
	}

	filenameSplits := strings.Split(header.Filename, ".")
	extension := filenameSplits[len(filenameSplits)-1]

	id := guuid.New()
	saveImageDir := fmt.Sprintf("%s/%s.%s", imageDir, id.String(), extension)
	err = fileutil.SaveToStorage(imageBytes, saveImageDir)
	if err != nil {
		http.Error(w, "Unable to save image", http.StatusBadGateway)
		return
	}

	responseString := fmt.Sprintf("%s\n", saveImageDir)
	w.Write([]byte(responseString))
}

func uploadZip(w http.ResponseWriter, r *http.Request) {
	// only allow POST method
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid Method", http.StatusMethodNotAllowed)
		return
	}

	r.ParseMultipartForm(32 << 20) // limit max input length

	var buf bytes.Buffer

	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "File error", http.StatusBadRequest)
		return
	}
	defer file.Close()
	io.Copy(&buf, file)

	zipReader, err := zip.NewReader(file, int64(len(buf.Bytes())))
	if err != nil {
		http.Error(w, "Unable to unzip", http.StatusBadRequest)
	}

	var responseStrings []string

	ch := make(chan string)

	for _, zf := range zipReader.File {
		go func(zf *zip.File, ch chan string) {
			src, err := zf.Open()
			if err != nil {
				ch <- ""
				return
			}
			defer src.Close()

			var tempBuf bytes.Buffer
			io.Copy(&tempBuf, src)

			tempBufBytes := tempBuf.Bytes()
			isZfImage := imageutil.IsImage(tempBufBytes)
			if !isZfImage {
				ch <- ""
				return
			}

			filenameSplits := strings.Split(zf.Name, ".")
			extension := filenameSplits[len(filenameSplits)-1]

			id := guuid.New()
			zfImgPath := fmt.Sprintf("%s/%s.%s", imageDir, id.String(), extension)

			err = fileutil.SaveToStorage(tempBufBytes, zfImgPath)
			if err != nil {
				ch <- ""
				return
			}

			ch <- zfImgPath
		}(zf, ch)
	}

	for {
		path := <-ch
		responseStrings = append(responseStrings, path)
		if len(responseStrings) == len(zipReader.File) {
			break
		}
	}

	response, _ := json.Marshal(responseStrings)

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

// generateThumbnailsPayload is the expected request body
type generateThumbnailsPayload struct {
	ImagePath string `json:"imagePath"`
	WidthSize int    `json:"widthSize"`
}

// generateThumbnails endpoint
func generateThumbnails(w http.ResponseWriter, r *http.Request) {
	// only allow POST method
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid Method", http.StatusMethodNotAllowed)
		return
	}

	// decode payload
	var payload generateThumbnailsPayload
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}

	// check if image exists
	hasFile := fileutil.HasFile(payload.ImagePath)
	if !hasFile {
		http.Error(w, "File not found", http.StatusNotFound)
	}

	imgReader, _ := os.Open(payload.ImagePath)

	resizedImgBytes, err := imageutil.ResizeImage(imgReader, payload.WidthSize)
	if err != nil {
		http.Error(w, "Error resizing image", http.StatusInternalServerError)
	}

	pathSplits := strings.Split(payload.ImagePath, "/")
	resizedFilename := fmt.Sprintf("%s/%dx%s", pathSplits[0], payload.WidthSize, pathSplits[1])

	err = fileutil.SaveToStorage(resizedImgBytes, resizedFilename)
	if err != nil {
		http.Error(w, "Unable to save image", http.StatusInternalServerError)
	}

	responseString := fmt.Sprintf("%s\n", resizedFilename)
	w.Write([]byte(responseString))
}
