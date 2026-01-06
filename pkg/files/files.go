package files

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"time"
)

// change fileName, here fileName is games
func SaveFile(file multipart.File, header *multipart.FileHeader) (string, error) {
	defer file.Close()
	os.MkdirAll("uploads/games", 0755)
	timestamp := time.Now().Unix()
	filename := fmt.Sprintf("uploads/games/%d-%s", timestamp, header.Filename)
	out, err := os.Create(filename)
	if err != nil {
		return "", err
	}
	defer out.Close()
	_, err = io.Copy(out, file)
	if err != nil {
		return "", err
	}
	url := fmt.Sprintf("/uploads/games/%d-%s", timestamp, header.Filename)
	return url, nil
}
func SaveRecord(file multipart.File, header *multipart.FileHeader) (string, error) {
	defer file.Close()
	os.MkdirAll("uploads/record", 0755)
	timestamp := time.Now().Unix()
	filename := fmt.Sprintf("uploads/record/%d-%s", timestamp, header.Filename)
	out, err := os.Create(filename)
	if err != nil {
		return "", err
	}
	defer out.Close()
	_, err = io.Copy(out, file)
	if err != nil {
		return "", err
	}
	url := fmt.Sprintf("/uploads/record/%d-%s", timestamp, header.Filename)
	return url, nil
}
func DeleteFile(path string) error {
	if path == "" {
		return nil
	}
	if path[0] == '/' {
		path = path[1:]
	}
	err := os.Remove(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	return nil
}
