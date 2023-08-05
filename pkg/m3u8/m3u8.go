package m3u8

import (
	"bufio"
	"fmt"
	"golang.org/x/net/html"
	"net/http"
	"net/url"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

var m3u8URLRegex = regexp.MustCompile(`https?://[^\s'"]+\.m3u8`)

func retrieveVariants(url string) ([]string, error) {
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
func FindM3U8Urls(url string) ([]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return nil, err
	}

	var urls []string
	var f func(*html.Node)
	f = func(n *html.Node) {
		// Check all attributes of the element node for URLs ending with .m3u8
		if n.Type == html.ElementNode {
			for _, a := range n.Attr {
				if strings.HasSuffix(a.Val, ".m3u8") {
					urls = append(urls, a.Val)
				}
			}
		}
		// Check text nodes for URLs ending with .m3u8
		if n.Type == html.TextNode {
			matches := m3u8URLRegex.FindAllString(n.Data, -1)
			urls = append(urls, matches...)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	return urls, nil
}
func downloadSegments(variantURL string) error {
	resp, err := http.Get(variantURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)

	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "#") && strings.HasSuffix(line, ".ts") {
			segmentURL := line // Construct the full URL if it's relative
			destPath := filepath.Join("output", filepath.Base(segmentURL))
			fmt.Println("Downloading segment:", segmentURL, " to ", destPath)
			//err := downloadFile(segmentURL, destPath)
			if err != nil {
				return err
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}
func CombineURL(masterURL, variantPath string) (string, error) {
	// If the variant is already a full URL, return it as-is
	if strings.HasPrefix(variantPath, "http://") || strings.HasPrefix(variantPath, "https://") {
		return variantPath, nil
	}

	masterURLParsed, err := url.Parse(masterURL)
	if err != nil {
		return "", err
	}

	// If the variant is an absolute path, combine it with the scheme and host of the master URL
	if strings.HasPrefix(variantPath, "/") {
		return masterURLParsed.Scheme + "://" + masterURLParsed.Host + variantPath, nil
	}

	// Otherwise, combine the variant path with the directory portion of the master URL's path
	return masterURLParsed.Scheme + "://" + masterURLParsed.Host + path.Join(path.Dir(masterURLParsed.Path), variantPath), nil
}
