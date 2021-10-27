## go-vp9

**_This project has been replaced with [py_vp9](https://github.com/jchamberlain0/py_vp9)!_**

Go script to encode high quality video files from lossless recordings for use in web applications.

Directories, filenames, and desired constant rate factor are set in settings.json.

Besides these, the other arguments to ffmpeg are hard-coded to ensure either best quality or best file size:

- libvpx-vp9 encoding. From the [ffmpeg docs](https://trac.ffmpeg.org/wiki/Encode/VP9#twopass):
  >libvpx-vp9 can save about 20–50% bitrate compared to libx264 (the default H.264 encoder), while retaining the same visual quality.
- Variable bitrate switch `-b:v` is set to 0 because the high quality use case prefers that image quality is not variable.
- Constant rate factor switch `-crf` is set by the user to achieve the desired consistent image quality across the length of the video. Good defaults are between 8 and 30, but more research is needed.
- Two-pass to further optimize IQ
- _(coming soon)_: `-deadline` and `-cpu-used` arguments to further optimize IQ

#### Features:
- _File_: Transcode a single input video to vp9
	- Default crf: 30
- _Batch_: Transcode a single input video to multiple quality option thresholds based on settings (Batch mode)
	- Default crf values: 50, 30, 12, 8 (it may be preferable to use crf 4 or lower for exceptional IQ)
- _(coming soon) Folder_: Transcode all items inside a single folder, with or without batch mode.
