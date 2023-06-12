package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func prompt() string {
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	return strings.TrimSpace(input)
}

func resolvePath(path string) string {
	absPath, err := filepath.Abs(path)
	if err != nil {
		log.Fatal(err)
	}
	return absPath
}

func ASCII() {
	lines := []string{
		"   ____                          _             _      _                                ",
		"  / ___|___  ___ _ __ ___   ___ | | ___   __ _(_) ___( )___                            ",
		" | |   / _ \\/ __| '_ ` _ \\ / _ \\| |/ _ \\ / _` | |/ __|// __|                           ",
		" | |__| (_) \\__ \\ | | | | | (_) | | (_) | (_| | | (__  \\__ \\                           ",
		"  \\____\\___/|___/_| |_| |_|\\___/|_|\\___/ \\__, |_ |\\___| |___/               _           ",
		" | ____|/ _|/ _| ___  _ __| |_| | ___  __|___/  | ____|_ __   ___ ___   __| | ___ _ __ ",
		" |  _| | |_| |_ / _ \\| '__| __| |/ _ \\/ __/ __| |  _| | '_ \\ / __/ _ \\ / _` |/ _ \\ '__|",
		" | |___|  _|  _| (_) | |  | |_| |  __/\\__ \\__ \\ | |___| | | | (_| (_) | (_| |  __/ |   ",
		" |_____|_| |_|  \\___/|_|   \\__|_|\\___||___/___/ |_____|_| |_|\\___\\___/ \\__,_|\\___|_|   ",
		"                                                                                       ",
	}

	for _, line := range lines {
		fmt.Println(line)
	}
}

func main() {
	ASCII()
	fmt.Printf("Please choose an input file: ")
	inputFile := prompt()
	resolvedPath := resolvePath(inputFile)
	inputDirectory := filepath.Dir(resolvedPath)
	inputFileName := filepath.Base(inputFile)
	inputFileNameNoExt := strings.TrimSuffix(inputFileName, filepath.Ext(inputFileName))

	fmt.Printf("Please choose the output file's location [%s]: ", inputDirectory)
	outputFolder := prompt()
	if outputFolder == "" {
		outputFolder = inputDirectory
	}

	fmt.Printf("Please choose a name for the output file [%s]: ", inputFileNameNoExt)
	outputName := prompt()
	if outputName == "" {
		outputName = inputFileNameNoExt
	}
	outputFile := filepath.Join(outputFolder, outputName)

	fmt.Printf("Please choose a file format for the output file [mp4]: ")
	outputFormat := prompt()
	if outputFormat == "" {
		outputFormat = "mp4"
	}

	fmt.Println("    Video Codecs")
	fmt.Println("1. H.264 (libx264)")
	fmt.Println("2. H.265 (libx265)")
	fmt.Println("3. VP9 (libvpx-vp9)")
	fmt.Println("4. AV1 (libaom-av1)")
	fmt.Printf("Please select a video codec [1]: ")
	videoCodec := prompt()
	if videoCodec == "" {
		videoCodec = "1"
	}

	var codec string
	switch videoCodec {
	case "1":
		fmt.Println("You selected H.264 (libx264). This codec is widely used and has good compatibility with most devices.")
		codec = "libx264"
	case "2":
		fmt.Println("You selected H.265 (libx265). This codec is more efficient than H.264 but has less compatibility with devices.")
		codec = "libx265"
	case "3":
		fmt.Println("You selected VP9 (libvpx-vp9). This codec provides high compression efficiency.")
		codec = "libvpx-vp9"
	case "4":
		fmt.Println("You selected AV1 (libaom-av1). This codec offers excellent compression but may require higher computational resources.")
		codec = "libaom-av1"
	default:
		fmt.Printf("Invalid video codec selected. Using the raw input as the codec: %s\n", videoCodec)
		codec = videoCodec
	}

	fmt.Println("    Audio Codecs")
	fmt.Println("1. AAC (aac)")
	fmt.Println("2. MP3 (libmp3lame)")
	fmt.Println("3. Opus (libopus)")
	fmt.Println("4. FLAC (flac)")
	fmt.Printf("Please select an audio codec [1]: ")
	audioCodec := prompt()
	if audioCodec == "" {
		audioCodec = "1"
	}

	var audio string
	switch audioCodec {
	case "1":
		fmt.Println("You selected AAC (aac). This codec is widely used and has good compatibility with most devices.")
		audio = "aac"
	case "2":
		fmt.Println("You selected MP3 (libmp3lame). This codec is widely used but has less compatibility with some devices.")
		audio = "libmp3lame"
	case "3":
		fmt.Println("You selected Opus (libopus). This codec provides high audio quality at low bitrates.")
		audio = "libopus"
	case "4":
		fmt.Println("You selected FLAC (flac). This codec provides lossless audio compression.")
		audio = "flac"
	default:
		fmt.Printf("Invalid audio codec selected. Using the raw input as the codec: %s\n", audioCodec)
		audio = audioCodec
	}

	cmd := exec.Command("ffprobe", "-v", "error", "-select_streams", "v:0", "-show_entries", "stream=width,height", "-of", "csv=s=x:p=0", inputFile)
	output, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	inputResolution := strings.TrimSpace(string(output))
	outputResolution := inputResolution

	fmt.Printf("The input file's resolution is: %s\n", inputResolution)
	fmt.Printf("Do you want to scale the video? (Y/[N]): ")
	scaleVideo := prompt()
	if len(scaleVideo) > 0 {
		scaleVideo = strings.ToLower(scaleVideo[:1])
	} else {
		scaleVideo = "n"
	}
	if scaleVideo == "y" {
		fmt.Println("    Scaling Options")
		fmt.Println("1. Scale to a specific resolution")
		fmt.Println("2. Scale by a percentage of the input resolution")
		fmt.Printf("Please select a scaling option: ")
		scalingOption := prompt()

		switch scalingOption {
		case "1":
			fmt.Printf("Please enter the desired resolution [%s]: ", inputResolution)
			outputResolution = prompt()
			if outputResolution == "" {
				outputResolution = inputResolution
			}
		case "2":
			fmt.Printf("Please enter the desired scaling percentage number [100]: ")
			scalingPercentage := prompt()
			if scalingPercentage == "" {
				scalingPercentage = "100"
			}
			outputResolution = fmt.Sprintf("iw*%s/100:ih*%s/100", scalingPercentage, scalingPercentage)
		default:
			fmt.Println("Invalid scaling option selected. Scaling will be skipped.")
			scaleVideo = "n"
		}
	}

	fmt.Println("    Hardware Acceleration APIs")
	fmt.Println("1. VDPAU")
	fmt.Println("2. VAAPI")
	fmt.Println("3. QSV")
	fmt.Println("4. NVDEC")
	fmt.Println("5. CUVID")
	fmt.Println("6. CUDA")
	fmt.Println("7. DXVA2")
	fmt.Println("8. Auto")
	fmt.Println("9. Off")
	fmt.Printf("Please select a Hardware Acceleration API [Auto]: ")
	hwaccel := prompt()
	if hwaccel == "" {
		hwaccel = "8"
	}

	var accelerationOption string
	switch hwaccel {
	case "1":
		fmt.Println("You selected VDPAU.")
		accelerationOption = "vdpau"
	case "2":
		fmt.Println("You selected VAAPI.")
		accelerationOption = "vaapi"
	case "3":
		fmt.Println("You selected QSV.")
		accelerationOption = "qsv"
	case "4":
		fmt.Println("You selected NVDEC.")
		accelerationOption = "nvdec"
	case "5":
		fmt.Println("You selected CUVID.")
		accelerationOption = "cuvid"
	case "6":
		fmt.Println("You selected CUDA.")
		accelerationOption = "cuda"
	case "8":
		fmt.Println("You selected Auto.")
		accelerationOption = "auto"
	case "9":
		fmt.Println("You selected Off.")
		accelerationOption = "off"
	default:
		fmt.Println("Invalid hardware acceleration API selected. Using Auto [default].")
		accelerationOption = ""
	}

	fmt.Printf("Do you want to include the resolution and video&audio codecs in the name of the output file? (Y/[N]): ")
	includeDetails := prompt()
	includeDetails = strings.ToLower(includeDetails[:1])
	if len(includeDetails) > 0 {
		includeDetails = strings.ToLower(scaleVideo[:1])
	} else {
		includeDetails = "n"
	}
	if includeDetails == "y" {
		outputFile = fmt.Sprintf("%s_%s_%s_%s", outputFile, outputResolution, codec, audio)
	}
	outputFile = fmt.Sprintf("%s.%s", outputFile, outputFormat)

	cmd = exec.Command("ffmpeg", "-hwaccel", accelerationOption, "-i", inputFile)
	if scaleVideo == "y" {
		cmd.Args = append(cmd.Args, "-vf", fmt.Sprintf("scale=%s", outputResolution))
	}
	cmd.Args = append(cmd.Args, "-c:v", codec, "-c:a", audio, outputFile)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Println("Executing FFmpeg command...")
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Encoding completed successfully!")
}
