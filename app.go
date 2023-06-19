package main

import (
	"context"
	"github.com/canhlinh/hlsdl"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Downloader struct
type Downloader struct {
	w     io.Writer
	ffmpegPath string
	ctx context.Context
}
func (d Downloader) Write(p []byte) (int, error) {
	runtime.EventsEmit(d.ctx, "log-writer", string(p) + "\r\n")
	return len(p), nil
}


// NewApp creates a new Downloader application struct
func NewDownloader(ffmpegPath string) *Downloader {
	return &Downloader{
		ffmpegPath: ffmpegPath,
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (d *Downloader) startup(ctx context.Context) {
	d.ctx = ctx
}

// Greet returns a greeting for the given name
func (d *Downloader) Download(m3u8 string) (string, error) {
	d.Write([]byte("beginning download"))
	dir, err := os.UserHomeDir()
	if err != nil {
		d.Write([]byte("err " + err.Error()))
		return "", err
	}
	outputFile, err := runtime.SaveFileDialog(d.ctx, runtime.SaveDialogOptions{
		DefaultDirectory:           dir,
		DefaultFilename:            "",
		Title:                      "Save as mp4 to...",
		Filters:                    nil,
		ShowHiddenFiles:            false,
		CanCreateDirectories:       true,
		TreatPackagesAsDirectories: false,
	})
	if err != nil {
		return "", err
	}
	outputFile = strings.TrimSuffix(outputFile, filepath.Ext(outputFile))
	videoPathslocation := filepath.Join(os.TempDir(), "ts-download")
	os.RemoveAll(videoPathslocation)
	os.Mkdir(videoPathslocation, 0750)
	hlsDL := hlsdl.New(m3u8, nil, videoPathslocation, 64, true, "")
	//send spinner start
	filepath, err := hlsDL.Download()
	if err != nil {
		d.Write([]byte("error downloading " + filepath + " - " + err.Error()))
		return "", err
	}
	d.Write([]byte("download finished"))
	//send spinner end
	d.Write([]byte("converting - " + filepath))

	newPath, err := executeFFmpeg(d.ffmpegPath, filepath, outputFile, d)
	if err != nil {
		d.Write([]byte("err " + err.Error()))
		return "", err
	}
	if err := os.Remove(filepath); err != nil {
		d.Write([]byte("err " + err.Error()))
		return "", err
	}
	return newPath, nil
}
