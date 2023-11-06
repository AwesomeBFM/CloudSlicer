package router

import (
	"fmt"
	"github.com/AwesomeBFM/CloudSlicer/internal/slicer"
	"github.com/gin-gonic/gin"
	"io"
	"math/rand"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func SliceFile(c *gin.Context) {
	err := c.Request.ParseMultipartForm(10 << 20) // 10 MB limit (adjust as needed)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse multipart form data"})
		return
	}

	// Access the uploaded model file
	model, err := c.FormFile("model") // "model" is the name of the file input field in the form
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Please include a model file"})
		return
	}

	// Ensure that the model is an STL, 3MF, or OBJ file (Checking MIME type is useless,
	// 	so we check the extension instead, not really secure though :\ )
	modelExt := extractExtension(model.Filename)
	modelExt = strings.ToLower(modelExt)
	if modelExt != "stl" && modelExt != "3mf" && modelExt != "obj" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid model file extension"})
		return
	}

	// Access the uploaded config file
	config, err := c.FormFile("config")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Please include a config file"})
		return
	}

	// Ensure that the config is an INI file
	configExt := extractExtension(config.Filename)
	configExt = strings.ToLower(configExt)
	if configExt != "ini" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid config file extension"})
		return
	}

	// Save the model to the temp directory
	saved := genRandomFilename()
	savedModel := saved + "." + modelExt
	err = saveFile(model, "./temp/model/"+savedModel)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not save model file"})
		return
	}

	// Save the config to the temp directory
	savedConfig := saved + "." + configExt
	err = saveFile(config, "./temp/config/"+savedConfig)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not save config file"})
		return
	}

	gcode, err := slicer.SliceFile(savedModel, savedConfig)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Could not slice model"})
		return
	}

	// Return the GCode file
	c.File(gcode)

	// Remove the model file
	err = removeFile("./temp/model/" + savedModel)
	if err != nil {
		// Report to the internal error handler
	}

	// Remove the config file
	err = removeFile("./temp/config/" + savedConfig)
	if err != nil {
		// Report to the internal error handler
	}

	// Remove the GCode file
	err = removeFile(gcode)
	if err != nil {
		// Report to the internal error handler
	}
}

func extractExtension(filename string) string {
	for i := len(filename) - 1; i >= 0; i-- {
		if filename[i] == '.' {
			return filename[i+1:]
		}
	}
	return ""
}

func saveFile(file *multipart.FileHeader, path string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		fmt.Println("Error on os.MkdirAll")
		return err
	}

	fileDest, err := os.Create(path) // temp/model/ + filename
	if err != nil {
		fmt.Println("Error on os.Create")
		return err
	}
	defer fileDest.Close()

	fileSrc, err := file.Open()
	if err != nil {
		fmt.Println("Error on *multipart.FileHeader.Open")
		return err
	}
	defer fileSrc.Close()

	_, err = io.Copy(fileDest, fileSrc)
	if err != nil {
		fmt.Println("Error on io.Copy")
		return err
	}

	return nil
}

func removeFile(path string) error {
	err := os.Remove(path)
	if err != nil {
		return err
	}
	return nil
}

func genRandomFilename() string {
	const letters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	result := make([]byte, 16)
	for i := range result {
		result[i] = letters[rand.Intn(len(letters))]
	}
	return string(result)
}
