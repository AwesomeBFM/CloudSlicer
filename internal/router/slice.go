package router

import (
	"bufio"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"math/rand"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

func SliceFile(c *gin.Context) {
	err := c.Request.ParseMultipartForm(10 << 20) // 10 MB limit (adjust as needed)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse multipart form data"})
		return
	}

	// Access the inputted material
	material := c.Request.FormValue("material")

	// Access the uploaded model file
	file, err := c.FormFile("model") // "model" is the name of the file input field in the form
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Please include a model file"})
		return
	}

	// Ensure that the model is an STL, 3MF, or OBJ file (Checking MIME type is useless,
	// 	so we check the extension instead, not really secure though :\ )
	extension := extractExtension(file.Filename)
	extension = strings.ToLower(extension)
	if extension != "stl" && extension != "3mf" && extension != "obj" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file extension"})
		return
	}

	// Select the correct config & density
	var confPath string
	var density float64
	switch material {
	case "pla":
		density = 1.24
		confPath = "./presets/pla.ini"
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid material"})
	}

	// Save the model to the temp directory
	filenameNoExt := genRandomFilename()
	filename := filenameNoExt + "." + extension
	err = saveFile(file, "./temp/model/"+filename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not save model file"})
		return
	}

	// Slice the model
	cmd := exec.Command(
		"prusa-slicer",
		"--export-gcode",
		"--load",
		confPath,
		"./temp/model/"+filename, "--output",
		"./temp/gcode/"+filenameNoExt+".gcode",
		"--info")

	output, err := cmd.CombinedOutput()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not slice the model"})
		return
	}
	outputStr := string(output)

	// Get the estimated weight of filament required
	scanner := bufio.NewScanner(strings.NewReader(outputStr))
	pattern := "volume = "

	var volume float64
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, pattern) {
			volumeStr := strings.Split(line, pattern)[1]
			volumeStr = strings.Split(volumeStr, " ")[0]
			volume, err = strconv.ParseFloat(volumeStr, 64)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not calculate weight"})
				return
			}
			break
		}
	}

	grams := volume * density * 1000
	fmt.Printf("Estimated weight: %f grams\n", grams)

	// Return the gcode file
	c.File("./temp/gcode/" + filenameNoExt + ".gcode")

	// Delete the model from the temp directory
	err = removeFile("./temp/model/" + filename)
	if err != nil {
		// Report to the internal error handler
	}

	// Remove the GCode file
	err = removeFile("./temp/gcode/" + filenameNoExt + ".gcode")
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
