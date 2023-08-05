package main

import (
	"bufio"
	"context"
	"fmt"
	"github.com/amlwwalker/video-stream-downloader/pkg/m3u8"
	"github.com/canhlinh/hlsdl"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	rt "runtime"
	"strings"
)

// Downloader struct
type Downloader struct {
	w          io.Writer
	ffmpegPath string
	ctx        context.Context
}

func (d Downloader) Write(p []byte) (int, error) {
	runtime.EventsEmit(d.ctx, "log-writer", string(p)+"\r\n")
	return len(p), nil
}

// NewApp creates a new Downloader application struct
func NewDownloader(ffmpegPath string) *Downloader {
	return &Downloader{
		ffmpegPath: ffmpegPath,
	}
}

func (d Downloader) OpenInDefaultBrowser(txt string) error {
	var err error
	switch rt.GOOS {
	case "linux":
		err = exec.Command("xdg-open", txt).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", txt).Start()
	case "darwin":
		err = exec.Command("open", txt).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	return err
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (d *Downloader) startup(ctx context.Context) {
	d.ctx = ctx
}

func (d *Downloader) DiscoverM3u8MasterLinks(url string) ([]string, error) {
	//url := "https://www.bbcmaestro.com/courses/richard-bertinet/bread-making#lesson-player"
	urls, err := m3u8.FindM3U8Urls(url)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}
	return urls, nil
}

func (d *Downloader) RetrieveVariants(url string) ([]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)
	var variantUrls []string

	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "#") && strings.HasSuffix(line, ".m3u8") {
			variantUrls = append(variantUrls, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return variantUrls, nil
}

// Greet returns a greeting for the given name
func (d *Downloader) DirectDownload(m3u8 string) (string, error) {
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

// Greet returns a greeting for the given name
func (d *Downloader) Download(masterUrl, m3u8Link string) (string, error) {
	d.Write([]byte("\r\nbeginning download"))
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
	combinedUrl, err := m3u8.CombineURL(masterUrl, m3u8Link)
	if err != nil {
		d.Write([]byte("err " + err.Error()))
		return "", err
	}
	hlsDL := hlsdl.New(combinedUrl, nil, videoPathslocation, 64, true, "")
	//send spinner start
	fmt.Println("combinedUrl ", combinedUrl)
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
