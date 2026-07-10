package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

const maxBytes = 2000

func main() {
	r := gin.Default()
	r.POST("/upload", func(context *gin.Context) {
		context.Request.Body = http.MaxBytesReader(context.Writer, context.Request.Body, maxBytes)
		f, err := context.MultipartForm()
		if err != nil {
			var err *http.MaxBytesError
			if errors.As(err, &err) {
				context.Status(http.StatusRequestEntityTooLarge)
				return
			}
			context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		files := f.File["files"]
		for _, file := range files {
			fmt.Println("FILE: ", file.Filename)
			err = context.SaveUploadedFile(file, filepath.Join("./files/", filepath.Base(file.Filename)))
			if err != nil {
				context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
		}
		context.String(http.StatusOK, "OK")
	})

	err := r.Run("localhost:8084")
	if err != nil {
		log.Fatal(err)
	}
}
