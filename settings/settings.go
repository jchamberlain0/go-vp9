package settings

import (
	"log"
	"io/ioutil"
	"encoding/json"
)

// Settings file path
var settingsFile = "./settings.json"

// Settings variable export
var Settings SettingsStruct

// SettingsStruct - structure for global settings read from file
type SettingsStruct struct {
	InputFilename    string    // name of movie file
	InputExtension   string    // extension of input file, probably ".avi" if recording with OBS
	OutputExtension  string    // extension of output file, probably ".webm" if encoding to vp9
	Debug            bool      // debug logging flag
	Batch            bool      // flag for using the quality option array instead of default CRF
	Environment      string    // maybe there's a use for this
	Mode             string    // "file" for 1 file, or "folder" to encode all matching in the folder (not supported yet)
	FileDir          string    // input directory for "file" mode
	OutFileDir       string    // output directory for "file" mode
	FolderDir        string    // input directory for "folder" mode
	OutFolderDir     string    // output directory for "folder"
	CRFDefault       string    // default quality option used when batch mode is off
	CRF              []string  // quality option array for batch mode
}

// LoadSettings - try to read from file and load settings
func LoadSettings() bool {
	log.Printf("LoadSettings()")
	bytes, fErr := ioutil.ReadFile(settingsFile)
	if fErr != nil {
		log.Printf("Failed to read settings file.")
		return false		
	}

	if err := json.Unmarshal([]byte(bytes), &Settings); err != nil {
		log.Printf("Error unmarshaling settings file")
		return false
	}

	return true
}