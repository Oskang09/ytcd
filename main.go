package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/kkdai/youtube/v2"
)

func trimmer(in string) string {
	return strings.Trim(strings.Trim(in, "\r\n"), "\n")
}

func appendIfNotExist(slice []string, i string) []string {
	for _, ele := range slice {
		if ele == i {
			return slice
		}
	}
	return append(slice, i)
}

func quit(reader *bufio.Reader, err string) {
	if err != "" {
		log.Println(err)
	}
	fmt.Print("Press enter to quit")
	reader.ReadString('\n')
	os.Exit(1)
}

func main() {
	client := youtube.Client{}
	wd, _ := os.Getwd()
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Youtube URL : ")
	videoURL, _ := reader.ReadString('\n')
	videoURL = trimmer(videoURL)

	values, err := url.Parse(videoURL)
	if err != nil {
		quit(reader, "not a valid youtube video url ( example: https://www.youtube.com/watch?v=LMVe_WtN5fE )")
	}

	videoId := values.Query().Get("v")
	video, err := client.GetVideo(videoId)
	if err != nil {
		quit(reader, "internal error: "+err.Error())
	}

	fmt.Print("Type ( video, audio ) : ")
	formatType, _ := reader.ReadString('\n')
	formatType = trimmer(formatType)

	possibleFormats := []string{}
	filteredFormats := video.Formats.Type(formatType + "/")
	for _, format := range filteredFormats {
		possibleFormats = appendIfNotExist(possibleFormats, format.Quality)
	}

	fmt.Print("Format ( " + strings.Join(possibleFormats, ", ") + " ): ")
	format, _ := reader.ReadString('\n')
	format = trimmer(format)

	quality := filteredFormats.FindByQuality(format)
	if quality == nil {
		quit(reader, "not a valid video format, received value '"+format+"'")
	}

	suffix := "mp4"
	if formatType == "audio" {
		suffix = "mp3"
	}

	filename := path.Join(wd, video.Title+"."+suffix)
	fmt.Print("File Name ( " + filename + " ) : ")
	inFilename, _ := reader.ReadString('\n')
	if strings.TrimSpace(inFilename) != "" {
		filename = trimmer(inFilename)
		if !strings.HasSuffix(filename, "."+suffix) {
			filename += "." + suffix
		}
	}

	fmt.Println("Downloading : "+filename, ", From: "+videoURL)
	stream, _, err := client.GetStream(video, quality)
	if err != nil {
		quit(reader, "internal error: "+err.Error())
	}

	file, err := os.Create(filename)
	if err != nil {
		quit(reader, "internal error: "+err.Error())
	}
	defer file.Close()

	_, err = io.Copy(file, stream)
	if err != nil {
		quit(reader, "internal error: "+err.Error())
	}

	fmt.Print("Download Complete : " + filename)
	quit(reader, "")
}
