package utils

import (
	"fmt"
	"github.com/google/uuid"
	exifremove "github.com/scottleedavis/go-exif-remove"
	"gopkg.in/gographics/imagick.v3/imagick"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func GetFileContentType(out multipart.File) (string, error) {
	buffer := make([]byte, 512)
	_, err := out.Read(buffer)
	if err != nil {
		return "", err
	}
	contentType := http.DetectContentType(buffer)
	return contentType, nil
}

func ValidImageContentType(contentType string) bool {
	switch contentType {
	case "image/bmp", "image/jpeg", "image/gif", "image/png", "image/webp", "image/heif", "image/heic":
		return true
	}
	return false
}

type CreatedImage struct {
	ID            uuid.UUID
	Filename      string
	Width, Height uint
	Original      bool
	Size          int64
}

func CreateImageFile(fh *multipart.FileHeader, directory string) ([]CreatedImage, error) {
	// Fix directory path.
	if strings.HasSuffix(directory, "/") {
		directory = directory[:len(directory)-1]
	}
	// Initialize imagick.
	var result []CreatedImage
	imagick.Initialize()
	defer imagick.Terminate()
	mw := imagick.NewMagickWand()
	defer mw.Destroy()

	// Open form file.
	multipartFile, err := fh.Open()
	if err != nil {
		return result, err
	}

	// Create temporary folder.
	os.Mkdir("temp-images", 0755)

	// Create temporary file.
	out, pathError := ioutil.TempFile("temp-images", "upload-*.png")
	if pathError != nil {
		return result, pathError
	}
	_, copyError := io.Copy(out, multipartFile)
	if copyError != nil {
		return result, copyError
	}
	multipartFile.Close()
	filename1 := out.Name()
	out.Close()
	defer os.Remove(filename1)

	// Remove exif data from image.
	f1, _ := os.ReadFile(filename1)
	noExifBytes, err := exifremove.Remove(f1)
	if err == nil {
		ioutil.WriteFile(filename1, noExifBytes, 0644)
	}

	// Read image from temporary file.
	err = mw.ReadImage(filename1)
	if err != nil {
		return result, err
	}
	mw.SetIteratorIndex(0) // This being the page offset
	err = mw.SetImageFormat("jpg")
	if err != nil {
		return result, err
	}

	// Original Image Compression
	imageId := uuid.New()
	originalWidth := mw.GetImageWidth()
	originalHeight := mw.GetImageHeight()
	filename := fmt.Sprintf("%s-%dx%d.jpg", imageId.String(), originalWidth, originalHeight)
	err = mw.WriteImage(fmt.Sprintf("%s/%s", directory, filename))
	if err != nil {
		return result, err
	}
	originalFile, err := os.Open(fmt.Sprintf("%s/%s", directory, filename))
	if err != nil {
		return result, err
	}
	defer originalFile.Close()
	originalFileStat, _ := originalFile.Stat()
	result = append(result, CreatedImage{
		ID:       imageId,
		Filename: filename,
		Width:    originalWidth,
		Height:   originalHeight,
		Original: true,
		Size:     originalFileStat.Size(),
	})

	// Create thumbnail.
	targetWidth := uint(512)
	targetHeight := uint(512)

	// Check width and height.
	if originalWidth > targetWidth && originalHeight > targetHeight {
		if originalWidth >= originalHeight {
			rate := float64(originalWidth) / float64(targetWidth)
			targetHeight = uint(float64(originalHeight) / rate)
		} else {
			rate := float64(originalHeight) / float64(targetHeight)
			targetWidth = uint(float64(originalWidth) / rate)
		}
	}
	err = mw.ThumbnailImage(targetWidth, targetHeight)
	if err != nil {
		os.Remove(result[0].Filename)
		return result, err
	}
	thumbnailWidth := mw.GetImageWidth()
	thumbnailHeight := mw.GetImageHeight()
	filename = fmt.Sprintf("%s-%dx%d.jpg", imageId.String(), thumbnailWidth, thumbnailHeight)
	err = mw.WriteImage(fmt.Sprintf("%s/%s", directory, filename))
	if err != nil {
		os.Remove(result[0].Filename)
		return result, err
	}
	thumbnailFile, err := os.Open(fmt.Sprintf("%s/%s", directory, filename))
	if err != nil {
		os.Remove(result[0].Filename)
		return result, err
	}
	defer thumbnailFile.Close()
	thumbnailFileStat, _ := thumbnailFile.Stat()
	result = append(result, CreatedImage{
		ID:       imageId,
		Filename: filename,
		Width:    thumbnailWidth,
		Height:   thumbnailHeight,
		Original: false,
		Size:     thumbnailFileStat.Size(),
	})
	return result, nil
}

func CreateProfileImageFile(fh *multipart.FileHeader, directory string) ([]CreatedImage, error) {
	// Fix directory path.
	if strings.HasSuffix(directory, "/") {
		directory = directory[:len(directory)-1]
	}
	// Initialize imagick.
	var result []CreatedImage
	imagick.Initialize()
	defer imagick.Terminate()
	mw := imagick.NewMagickWand()
	defer mw.Destroy()

	// Open form file.
	multipartFile, err := fh.Open()
	if err != nil {
		return result, err
	}

	// Create temporary folder.
	os.Mkdir("temp-images", 0755)

	// Create temporary file.
	extension := filepath.Ext(fh.Filename)
	out, pathError := ioutil.TempFile("temp-images", "upload-*"+extension)
	if pathError != nil {
		return result, pathError
	}
	_, copyError := io.Copy(out, multipartFile)
	if copyError != nil {
		return result, copyError
	}
	multipartFile.Close()
	filename1 := out.Name()
	out.Close()
	defer os.Remove(filename1)

	// Remove exif data from image.
	f1, _ := os.ReadFile(filename1)
	noExifBytes, err := exifremove.Remove(f1)
	if err == nil {
		ioutil.WriteFile(filename1, noExifBytes, 0644)
	}

	// Read image from temporary file.
	err = mw.ReadImage(filename1)
	if err != nil {
		return result, err
	}
	mw.SetIteratorIndex(0) // This being the page offset
	if strings.ToLower(mw.GetImageFormat()) == "gif" {
		mw = mw.CoalesceImages()
	}
	err = mw.SetImageFormat("jpg")
	if err != nil {
		return result, err
	}

	// Original Image Compression
	imageId := uuid.New()
	originalWidth := mw.GetImageWidth()
	originalHeight := mw.GetImageHeight()
	filename := fmt.Sprintf("%s-%dx%d.jpg", imageId.String(), originalWidth, originalHeight)
	err = mw.WriteImage(fmt.Sprintf("%s/%s", directory, filename))
	if err != nil {
		return result, err
	}

	// Create 400x400~ image.
	targetWidth := uint(400)
	targetHeight := uint(400)

	// Check width and height.
	if originalWidth > targetWidth && originalHeight > targetHeight {
		if originalWidth >= originalHeight {
			rate := float64(originalWidth) / float64(targetWidth)
			targetHeight = uint(float64(originalHeight) / rate)
		} else {
			rate := float64(originalHeight) / float64(targetHeight)
			targetWidth = uint(float64(originalWidth) / rate)
		}
	}
	err = mw.ThumbnailImage(targetWidth, targetHeight)
	if err != nil {
		os.Remove(result[0].Filename)
		return result, err
	}
	thumbnailWidth := mw.GetImageWidth()
	thumbnailHeight := mw.GetImageHeight()
	filename = fmt.Sprintf("%s-%dx%d.jpg", imageId.String(), thumbnailWidth, thumbnailHeight)
	err = mw.WriteImage(fmt.Sprintf("%s/%s", directory, filename))
	if err != nil {
		os.Remove(result[0].Filename)
		return result, err
	}
	thumbnailFile, err := os.Open(fmt.Sprintf("%s/%s", directory, filename))
	if err != nil {
		os.Remove(result[0].Filename)
		return result, err
	}
	defer thumbnailFile.Close()
	thumbnailFileStat, _ := thumbnailFile.Stat()
	result = append(result, CreatedImage{
		ID:       imageId,
		Filename: filename,
		Width:    thumbnailWidth,
		Height:   thumbnailHeight,
		Original: false,
		Size:     thumbnailFileStat.Size(),
	})

	// Create 200x200~ image.
	targetWidth2 := uint(200)
	targetHeight2 := uint(200)

	// Check width and height.
	if originalWidth > targetWidth2 && originalHeight > targetHeight2 {
		if originalWidth >= originalHeight {
			rate := float64(originalWidth) / float64(targetWidth2)
			targetHeight2 = uint(float64(originalHeight) / rate)
		} else {
			rate := float64(originalHeight) / float64(targetHeight2)
			targetWidth2 = uint(float64(originalWidth) / rate)
		}
	}
	err = mw.ThumbnailImage(targetWidth2, targetHeight2)
	if err != nil {
		os.Remove(result[0].Filename)
		return result, err
	}
	thumbnailWidth2 := mw.GetImageWidth()
	thumbnailHeight2 := mw.GetImageHeight()
	filename = fmt.Sprintf("%s-%dx%d.jpg", imageId.String(), thumbnailWidth2, thumbnailHeight2)
	err = mw.WriteImage(fmt.Sprintf("%s/%s", directory, filename))
	if err != nil {
		os.Remove(result[0].Filename)
		return result, err
	}
	thumbnailFile2, err := os.Open(fmt.Sprintf("%s/%s", directory, filename))
	if err != nil {
		os.Remove(result[0].Filename)
		return result, err
	}
	defer thumbnailFile2.Close()
	thumbnailFileStat2, _ := thumbnailFile2.Stat()
	result = append(result, CreatedImage{
		ID:       imageId,
		Filename: filename,
		Width:    thumbnailWidth2,
		Height:   thumbnailHeight2,
		Original: false,
		Size:     thumbnailFileStat2.Size(),
	})
	return result, nil
}
