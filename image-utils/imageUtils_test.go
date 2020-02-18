package imageutil_test

import (
	"bytes"
	"os"
	"testing"

	imageutil "github.com/jjkoh95/image-micro-web-framework/image-utils"
)

const resizeDst = "resize.jpg"

func TestMain(m *testing.M) {
	code := m.Run()
	// tear down
	tearDownFile()
	os.Exit(code)
}

func TestIsImage(t *testing.T) {
	r, err := os.Open("test/babe.jpg")
	if err != nil {
		t.Error("Expect to open test/babe.jpg without error")
	}
	defer r.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(r)
	checkIsImage := imageutil.IsImage(buf.Bytes())
	if !checkIsImage {
		t.Error("Expect test/babe.jpg to return true")
	}

	r, err = os.Open("test/resume.pdf")
	if err != nil {
		t.Error("Expect to open test/resume.pdf without error")
	}

	defer r.Close()
	buf = new(bytes.Buffer)
	buf.ReadFrom(r)
	checkIsImage = imageutil.IsImage(buf.Bytes())
	if checkIsImage {
		t.Error("Expect test/resume.pdf to return false")
	}
}

func TestResizeImage(t *testing.T) {
	file, err := os.Open("test/babe.jpg")
	if err != nil {
		t.Error("Expect to read image without error")
	}

	resizedImgBytes, err := imageutil.ResizeImage(file, imageutil.Widths[0])
	if err != nil {
		t.Error("Expect to resize image without error")
	}

	resizedImgReader := bytes.NewReader(resizedImgBytes)
	w, _, err := imageutil.GetImageSize(resizedImgReader)
	if err != nil {
		t.Error("Expect to get size without error")
	}
	if w != imageutil.Widths[0] {
		t.Error("Expect to resize to preselected size")
	}
}

func TestGetImageSize(t *testing.T) {
	file, err := os.Open("test/babe.jpg")
	if err != nil {
		t.Error("Expect to read image without error")
	}
	w, h, err := imageutil.GetImageSize(file)
	if err != nil {
		t.Error("Expect to get size without error")
	}
	if w != 1080 && h != 1350 {
		t.Error("Expect to get correct width and height of the image")
	}
}

func tearDownFile() {
	os.Remove(resizeDst)
}
