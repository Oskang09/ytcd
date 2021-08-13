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
	"unicode/utf8"

	"github.com/kkdai/youtube/v2"
)

func trimmer(in string) string {
	return strings.Trim(strings.Trim(in, "\r\n"), "\n")
}

func quit(reader *bufio.Reader, err string) {
	if err != "" {
		log.Println(err)
	}
	fmt.Print("Press enter to quit")
	reader.ReadString('\n')
	panic(err)
}

func stringToAscii(s string) string {
	t := make([]byte, utf8.RuneCountInString(s))
	i := 0
	for _, r := range s {
		t[i] = byte(r)
		i++
	}
	return string(t)
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

	possibleFormats := []string{}
	for _, format := range video.Formats {
		if strings.HasPrefix(format.Quality, "hd") {
			possibleFormats = append(possibleFormats, format.Quality)
		}
	}

	fmt.Print("Video Format ( " + strings.Join(possibleFormats, ", ") + " ): ")
	format, _ := reader.ReadString('\n')
	format = trimmer(format)

	quality := video.Formats.FindByQuality(format)
	if quality == nil {
		quit(reader, "not a valid video format, received value '"+format+"'")
	}

	log.Println(video.Title)

	filename := path.Join(wd, video.Title+".mp4")
	fmt.Print("File Name ( " + filename + " ) : ")
	inFilename, _ := reader.ReadString('\n')
	if strings.TrimSpace(inFilename) != "" {
		filename = trimmer(inFilename)
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
