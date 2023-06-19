package main

import (
	"embed"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func generateLinks(link string, count int) []string{
	var links []string
	for i := 1; i <= count; i++ {
		links = append(links, fmt.Sprintf(link, i, i))
	}
	return links
}

//go:embed embeds/ffmpeg
var filePayload []byte


func initFFmpeg() (string, error) {
	// Write the embedded file to a temporary location
	tmpDir, err := ioutil.TempDir("", "")
	if err != nil {
		return "", err
	}

	tmpFn := filepath.Join(tmpDir, "ffmpeg")
	if err := ioutil.WriteFile(tmpFn, filePayload, 0755); err != nil {
		return "", err
	}
	return tmpFn, nil
}
func executeFFmpeg(ffmpeg, tsFile, outputFile string, logWriter io.Writer) (string, error) {

	cmd := exec.Command(ffmpeg, "-y", "-i", tsFile, "-c:v", "libx264", "-c:v", "copy", outputFile + ".mp4")
	cmd.Stdout = logWriter
	cmd.Stderr = logWriter
	if err := cmd.Run(); err != nil {
		logWriter.Write([]byte("err " + err.Error()))
		return "", err
	}
	return outputFile + ".mp4", nil
}
func gracefulShutdown(ffmpegPath string) {
	s := make(chan os.Signal, 1)
	signal.Notify(s, os.Interrupt)
	signal.Notify(s, syscall.SIGTERM)
	go func() {
		<-s
		fmt.Println("Shutting down gracefully.")
		os.RemoveAll(ffmpegPath) // clean up
		os.Exit(0)
	}()
}
func main() {
	ffmpegPath, err := initFFmpeg()
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(ffmpegPath) // clean up
	go gracefulShutdown(ffmpegPath)

	// Create an instance of the app structure
	downloader := NewDownloader(ffmpegPath)

	// Create application with options
	err = wails.Run(&options.App{
		Title:  "MergeMovie",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        downloader.startup,
		Bind: []interface{}{
			downloader,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
