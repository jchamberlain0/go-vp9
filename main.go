package main

/*
ffmpeg wrapper for libvpx-vp9 video encoding
Features:
	- File: Transcode a single input video to vp9
		- Default crf: 30
	- Batch: Transcode a single input video to multiple quality option thresholds based on settings (Batch mode)
		- Default crf values: 50, 30, 12, 8 (it may be preferable to use crf 4 or even lower for exceptional IQ)
	- Folder: Transcode all items inside a single folder, with or without batch mode (coming soon?)
*/

import (
	"os/exec"
	"log"
	"go-libvpx/settings"
)

// encode a video with high image quality to compression ratio.
// libvpx-vp9 can save 20-50% bitrate with the same visual quality as H.264.
// these commands use two passes without an upper bound for bitrate, using the
// constant rate factor switch (-crf) instead of using an average or target bitrate.
// https://trac.ffmpeg.org/wiki/Encode/VP9#twopass
func encodeVP9(crf string) bool {

	// s is for sugar.
	// do some string building to keep the command arg list clean
	// TODO: make the distinction between q and crf more clear
	s := settings.Settings
	i := s.FileDir+s.InputFilename+s.InputExtension
	o := s.OutFolderDir+s.InputFilename+crf+s.OutputExtension

	// check crf against 
	if crf == "" || &crf == nil {
		log.Printf("Invalid quality level passed. Using default from settings: %v",s.CRFDefault)
		crf = s.CRFDefault
	}

	if s.Debug {
		log.Printf("s: %#+v\n", s)
		log.Printf("i: %#+v\n", i)
		log.Printf("o: %#+v\n", o)
		log.Printf("crf: %#+v\n", crf)
	}
	
	// Build the ffmpeg commands
	firstPass  := exec.Command("ffmpeg", "-y", "-i", i, "-c:v", "libvpx-vp9", "-b:v", "0", "-crf",  crf, "-pass", "1", "-an", "-f", "webm", "NUL")
	secondPass := exec.Command("ffmpeg",       "-i", i, "-c:v", "libvpx-vp9", "-b:v", "0", "-crf",  crf, "-pass", "2", "-c:a", "libopus", o)

	// First pass
	log.Printf("Running first pass.")
	firstLog, err1 := firstPass.CombinedOutput()
	if err1 != nil {
		log.Printf(string(firstLog))
		log.Printf("First Pass stopped with error: %v",err1)
		log.Printf("Encoding failed for crf %v. Exiting",crf)
		return false
	}
	log.Printf(string(firstLog))
	log.Printf("Encoding succeeded. Beginning second pass crf %v",crf)
	
	// Second pass
	secondLog, err2 := secondPass.CombinedOutput()
	if err2 != nil {
		log.Printf(string(secondLog))
		log.Printf("Second Pass stopped with error: %v",err1)
		log.Printf("Encoding failed on second pass for crf %v. Exiting",crf)
		return false
	}
		log.Printf(string(secondLog))
		log.Printf("Encoding succeeded. Exiting")
	
		return true
}

func main() {
	settings.LoadSettings()
	if settings.Settings.Debug {	
		log.Printf("settings: %v",settings.Settings)
	}

	if settings.Settings.Batch {
		// batch processing: multiple quality options of the same input file.
		// update the values in settings.CRF to specify constant rate factor.
		log.Printf("Batch mode: %v", settings.Settings.CRF)
    for i, v := range settings.Settings.CRF {
			log.Printf("Batch mode: Starting encode for quality: %v (%v of %v)",v,i+1,len(settings.Settings.CRF))
			// encode a video with the current quality selection
			if encodeVP9(v) {
				log.Printf("Batch mode: crf %v succeeded. (%v of %v)",v,i+1,len(settings.Settings.CRF))
			} else {
				log.Printf("Batch mode: ")
			}
		}
	} else {
		//encode a single video with default quality from settings
		encodeVP9(settings.Settings.CRFDefault)
	}

	log.Printf("Successfully processed %v movies: %v",len(settings.Settings.CRF),settings.Settings.CRF)


}