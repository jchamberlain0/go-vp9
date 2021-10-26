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
	"go-libvpx/settings"
	"log"
	"os/exec"
)

// encode a video with high image quality to compression ratio.
// libvpx-vp9 can save 20-50% bitrate with the same visual quality as H.264.
// these commands use two passes without an upper bound for bitrate, using the
// constant rate factor switch (-crf) instead of using an average or target bitrate.
// https://trac.ffmpeg.org/wiki/Encode/VP9#twopass
func encodeVP9(crf string) bool {

	// todo check if file already exists like seriously
	// jesus don't encode a video for an hour and then fail to save it to disk

		// s is for sugar.
	// do some string building to keep the command arg list clean
	// TODO: change the name of settings so we don't have the stutter seen below
	s := settings.Settings

	var resolution string
	var horizontalLines string
	// var swsflags string
	
	if s.Scale {
		resolution = "scale=" + s.OutResolution //" -sws_flags " + s.ScaleMode
		// TODO: fix this line so it works on quad-digit resolutions? I'm not scaling up right now, so..
		horizontalLines = s.OutResolution[len(s.OutResolution)-3:len(s.OutResolution)]
		log.Printf(horizontalLines)
		log.Printf(resolution)
		} else {
			// resolution = ""
			// horizontalLines = "noscale"
	}

	// input file meta string
	i := s.FileDir+s.InputFilename+s.InputExtension
	// output file meta string              crf=quality     horizontal lines
	o := s.OutFolderDir+s.InputFilename+"_crf"+crf+"_"+horizontalLines+s.OutputExtension

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
  // firstPass := ""
	// secondPass := ""

	// var firstPass  Cmd
	// var secondPass Cmd

	// Running this logic twice because I don't know how to declare a Cmd type without errors
	firstPass  := exec.Command("ffmpeg", "-y", "-i", i,  "-c:v", "libvpx-vp9", "-b:v", "0", "-crf", crf, "-pass", "1", "-f", "webm", "NUL")
  secondPass := exec.Command("ffmpeg",       "-i", i,  "-c:v", "libvpx-vp9", "-b:v", "0", "-crf", crf, "-pass", "2", "-acodec", "copy", "-c:a", "libopus", o)


	firstPassScaled :=  exec.Command("ffmpeg", "-y", "-i", i, "-vf", resolution, "-sws_flags", s.ScaleMode, "-c:v", "libvpx-vp9", "-b:v", "0", "-crf", crf, "-pass", "1", "-f", "webm", "NUL")
	secondPassScaled := exec.Command("ffmpeg",       "-i", i, "-vf", resolution, "-sws_flags", s.ScaleMode, "-c:v", "libvpx-vp9", "-b:v", "0", "-crf", crf, "-pass", "2", "-c:a", "libopus", o)

	// if s.Scale {
	// 	// add -vf, resolution, -sws_flags, and s.ScaleMode
	// 	firstPass  = exec.Command("ffmpeg", "-y", "-i", i, "-vf", resolution, "-sws_flags", s.ScaleMode, "-c:v", "libvpx-vp9", "-b:v", "5000k", "-crf", crf, "-pass", "1", "-f", "webm", "NUL")
	// 	secondPass = 
	// 	// firstPass  := exec.Command("ffmpeg", "-y", "-i", i, "-vf", resolution, "-sws_flags", s.ScaleMode, "-c:v", "libvpx-vp9", "-b:v", "0", "-crf", crf, "-pass", "1", "-f", "webm", "NUL")
	// 	// firstPass  := exec.Command("ffmpeg", "-y", "-i", i, "-vf", resolution, "-c:v", "libvpx-vp9", "-b:v", "0", "-crf", crf, "-pass", "1", "-an", "-f", "webm", "NUL")
	// 	// secondPass := exec.Command("ffmpeg",       "-i", i, "-vf", resolution, "-c:v", "libvpx-vp9", "-b:v", "0", "-crf", crf, "-pass", "2", "-c:a", "libopus", o)

	// }

	// uncomment to debug
	log.Printf("encodeVP9debug: Running first pass.\n")
	log.Printf("encodeVP9debug: \n%v\n",firstPass.String())
	log.Printf("encodeVP9debug: Running second pass.\n")
	log.Printf("encodeVP9debug: \n%v\n",secondPass.String())

	log.Printf("encodeVP9debug: Running first pass, scaled.\n")
	log.Printf("encodeVP9debug: \n%v\n",firstPassScaled.String())
	log.Printf("encodeVP9debug: Running second pass, scaled.\n")
	log.Printf("encodeVP9debug: \n%v\n",secondPassScaled.String())
	// return false;

	// stdout, errStdOut := firstPass.StdoutPipe()
	// if errStdOut != nil {
	// 		log.Fatal(errStdOut)
	// }

	// if err := firstPass.Start(); err != nil {
  //       log.Fatal(err)
  //   }



	// go func() {
  //       defer stdin.Close()
  //       io.WriteString(stdin, "an old falcon")
  //   }()

	// firstPassReady := fir

	// First pass
	log.Printf("encodeVP9: Running first pass.\n")
	log.Printf("encodeVP9: \n%v\n",firstPassScaled.String())
	firstLog, err1 := firstPassScaled.CombinedOutput()
	if err1 != nil {
		log.Printf(string(firstLog))
		log.Printf("First Pass stopped with error: %v",err1)
		log.Printf("Encoding failed for crf %v. Exiting",crf)
		return false
	}
	log.Printf(string(firstLog))
	log.Printf("Encoding succeeded. Beginning second pass crf %v",crf)
	
	// Second pass
	log.Printf("encodeVP9: Running second pass.\n")
	log.Printf("encodeVP9: \n%v\n",secondPassScaled.String())	
	secondLog, err2 := secondPassScaled.CombinedOutput()
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


func encodeAV1(crf string) bool {

	// todo check if file already exists like seriously
	// jesus don't encode a video for an hour and then fail to save it to disk

		// s is for sugar.
	// do some string building to keep the command arg list clean
	// TODO: change the name of settings so we don't have the stutter seen below
	s := settings.Settings

	var resolution string
	var horizontalLines string
	
	if s.Scale {
		resolution = "-vf scale="+s.OutResolution+":flags="+s.ScaleMode
		// TODO: fix this line so it works on quad-digit resolutions? I'm not scaling up right now, so..
		horizontalLines = s.OutResolution[len(s.OutResolution)-3:len(s.OutResolution)]
		log.Printf("Scaling might have been broken again...")
		log.Printf(horizontalLines)
		log.Printf(resolution)
		} else {
			log.Printf("No scaling mode")
			resolution = ""
			horizontalLines = ""
	}

	// input file meta string
	i := s.FileDir+s.InputFilename+s.InputExtension
	// output file meta string              crf=quality     horizontal lines
	o := s.OutFolderDir+""+s.InputFilename+"_crf"+crf+"_"+horizontalLines+s.OutputExtension

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
  // firstPass := ""
	// secondPass := ""

	firstPass  := exec.Command("ffmpeg", "-y", "-i", i, "-c:v", "libaom-av1", "-b:v", "0", "-crf", crf, "-pass", "1", "-an", "-f", "matroska", "NUL")//, "&&", "^")
	secondPass := exec.Command("ffmpeg",       "-i", i, "-c:v", "libaom-av1", "-b:v", "0", "-crf", crf, "-pass", "2", "-c:a", "libopus", o)
	// firstPass  := exec.Command("ffmpeg", "-y", "-i", i, "-vf", resolution, "-c:v", "libvpx-vp9", "-b:v", "0", "-crf", crf, "-pass", "1", "-an", "-f", "webm", "NUL")
	// secondPass := exec.Command("ffmpeg",       "-i", i, "-vf", resolution, "-c:v", "libvpx-vp9", "-b:v", "0", "-crf", crf, "-pass", "2", "-c:a", "libopus", o)

	// uncomment to debug
	// log.Printf("encodeAV1: Running first pass.\n")
	// log.Printf("encodeAV1: \n%v\n",firstPass.String())
	// log.Printf("encodeAV1: Running second pass.\n")
	// log.Printf("encodeAV1: \n%v\n",secondPass.String())
	// return false;

	
	// return false;

	// First pass
	log.Printf("encodeAV1: Running first pass.\n")
	log.Printf("encodeAV1: \n%v\n",firstPass.String())


	firstLog, err1 := firstPass.CombinedOutput()
	if err1 != nil {
		log.Printf(string(firstLog))
		log.Printf("First Pass stopped with error: %v",err1)
		log.Printf("Encoding failed for crf %v. Exiting",crf)
		return false
	}
	log.Printf(string(firstLog))
	log.Printf("Encoding succeeded. Beginning second pass crf %v",crf)
	
	log.Printf("encodeAV1: Running second pass.\n")
	log.Printf("encodeAV1: \n%v\n",secondPass.String())

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
		// if encodeAV1(settings.Settings.CRFDefault) {
		// 	log.Printf("Successfully encoded video file %v",settings.Settings.CRFDefault)

		// }
		if encodeVP9(settings.Settings.CRFDefault) {
			log.Printf("Successfully encoded video file %v",settings.Settings.CRFDefault)

		}
	}

	if settings.Settings.Batch {
		log.Printf("Successfully processed %v movies: %v",len(settings.Settings.CRF),settings.Settings.CRF)
	} else {
	}


}