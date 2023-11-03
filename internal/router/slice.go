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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Access the "pla" field value
	material := c.Request.FormValue("material")

	// Access the uploaded model file
	file, err := c.FormFile("model") // "model" is the name of the file input field in the form
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Ensure that the model is an STL, 3MF, or OBJ file (Checking MIME type is useless, so we check the extension instead, not really secure though :\ )
	nameParts := strings.Split(file.Filename, ".")
	if len(nameParts) < 2 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file extension"})
		return
	}

	extension := nameParts[len(nameParts)-1]
	extension = strings.ToLower(extension)

	if extension != "stl" && extension != "3mf" && extension != "obj" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file extension"})
		return
	}

	// Select the correct config
	var confPath string
	switch material {
	case "pla":
		confPath = "./presets/pla.ini"
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid material"})
	}

	// Save the model to the temp directory
	filenameNoExt := genRandomFilename()
	filename := filenameNoExt + "." + extension
	err = saveFile(file, "./temp/model/"+filename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error :(": err.Error()})
		return
	}

	// Make the gcode dir and reserve the filename
	dir := filepath.Dir("./temp/gcode/" + filenameNoExt + ".gcode")
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		fmt.Println("Error on os.MkdirAll")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Slice the model
	cmd := exec.Command("prusa-slicer", "--export-gcode", "--load", confPath, "./temp/model/"+filename, "--output", "./temp/gcode/"+filenameNoExt+".gcode", "--info")

	output, err := cmd.CombinedOutput()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			break
		}
	}

	grams := volume * 1.24 * 1000 // 1.24 g/cm^3 is the density of PLA
	fmt.Printf("Estimated weight: %f grams\n", grams)

	// Delete the model from the temp directory
	err = removeFile("./temp/model/" + filename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	/*// Grab the GCODE file
	gcodeFileRef, err := os.Open("./temp/gcode/" + filenameNoExt + ".gcode")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer gcodeFileRef.Close()

	gcodeFile, err := io.ReadAll(gcodeFileRef)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}*/

	/*gcodeFile, err := os.ReadFile("./temp/gcode/" + filenameNoExt + ".gcode")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Respond with a success message
	c.JSON(http.StatusOK, gin.H{
		"file":   gcodeFile,
		"weight": "-1",
	})*/
	c.File("./temp/gcode/" + filenameNoExt + ".gcode")
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
