package slicer

import (
	"errors"
	"os/exec"
	"strings"
)

func SliceFile(model string, config string) (string, error) {

	gcode := stripExtension(model) + ".gcode"

	cmd := exec.Command(
		"prusa-slicer",
		"--export-gcode",
		"--load",
		"./temp/config/"+config,
		"./temp/model/"+model,
		"--output", "./temp/gcode/"+gcode,
		"--info")

	_, err := cmd.CombinedOutput()
	if err != nil {
		return "", errors.New("could not slice the model")
	}
	//outputStr := string(output)
	// We may use the output in the future but not right now

	return "./temp/gcode/" + gcode, nil
}

func stripExtension(filename string) string {
	// Find the last dot in the filename
	dotIndex := strings.LastIndex(filename, ".")

	// If there's no dot, return the original filename
	if dotIndex == -1 {
		return filename
	}

	// Extract the substring before the last dot
	strippedFilename := filename[:dotIndex]

	return strippedFilename
}
