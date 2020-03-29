package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	logLevel "github.com/llimllib/loglevel"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func main() {
	logLevel.SetPriorityString("info")
	logLevel.SetPrefix("crawler")

	logLevel.Debug(os.Args)

	if len(os.Args) < 2 {
		logLevel.Fatalln("Missing Url arg")
	}

	recurDownloader(os.Args[1], 0)
}

func recurDownloader(url string, depth int) {
	page, err := downloader(url)
	if err != nil {
		logLevel.Error(err)
		return
	}
	links := LinkReader(page, depth)

	for _, link := range links {
		fmt.Println(link)
		if depth+1 < MaxDepth {
			recurDownloader(link.url, depth+1)
		}
	}
}

func downloader(url string) (resp *http.Response, err error) {
	logLevel.Debugf("Downloading %s", url)
	resp, err = http.Get(url)
	if err != nil {
		logLevel.Debugf("Error: %s", err)
		return
	}

	if resp.StatusCode > 299 {
		err = HTTPError{fmt.Sprintf("Error (%d): %s", resp.StatusCode, url)}
		logLevel.Debug(err)
		return
	}
	return

}

//LinkReader take link and select only A field in html
func LinkReader(resp *http.Response, depth int) []Link {
	page := html.NewTokenizer(resp.Body)
	links := []Link{}

	var start *html.Token
	var text string

	for {
		_ = page.Next()
		token := page.Token()
		if token.Type == html.ErrorToken {
			break
		}

		if start != nil && token.Type == html.TextToken {
			text = fmt.Sprintf("%s%s", text, token.Data)
		}

		if token.DataAtom == atom.A {
			switch token.Type {
			case html.StartTagToken:
				if len(token.Attr) > 0 {
					start = &token
				}
			case html.EndTagToken:
				if start == nil {
					logLevel.Warnf("Link End found without Start: %s", text)
					continue
				}
				link := NewLink(*start, text, depth)
				if link.Valid() {
					links = append(links, link)
					logLevel.Debugf("Link Found %v", link)
				}

				start = nil
				text = ""
			}
		}
	}

	logLevel.Debug(links)
	return links
}

//Link define a  link struct
type Link struct {
	url   string
	text  string
	depth int
}

//HTTPError define an error
type HTTPError struct {
	original string
}

/*Ex 2*/
func linkReader(resp *http.Response, depth int) []Link {
	page := html.NewTokenizer(resp.Body)
	links := []Link{}

	var start *html.Token
	var text string

	for {
		_ = page.Next()
		token := page.Token()

		if token.Type == html.ErrorToken {
			break
		}

		if start != nil && token.Type == html.TextToken {
			text = fmt.Sprintf("%s%s", text, token.Data)
		}

		if token.DataAtom == atom.A {
			switch token.Type {
			case html.StartTagToken:
				if len(token.Attr) > 0 {
					start = &token
				}
			case html.EndTagToken:
				if start == nil {
					logLevel.Warnf("link End without start: %s", text)
					continue
				}
				link := NewLink(*start, text, depth)
				if link.Valid() {
					links = append(links, link)
					logLevel.Debugf("link found %v", link)
				}
				start = nil
				text = ""
			}
		}

	}

	logLevel.Debug(links)
	return links

}

//NewLink generate new link
func NewLink(tag html.Token, text string, depth int) Link {
	link := Link{text: strings.TrimSpace(text), depth: depth}

	for i := range tag.Attr {
		if tag.Attr[i].Key == "href" {
			link.url = strings.TrimSpace(tag.Attr[i].Val)
		}
	}
	return link
}

func (t Link) String() string {
	spacer := strings.Repeat("\t", t.depth)
	return fmt.Sprintf("%s%s (%d) - %s", spacer, t.text, t.depth, t.url)
}

//MaxDepth will be the max level fro crawler
const MaxDepth = 2

//Valid validate link
func (t Link) Valid() bool {
	if t.depth >= MaxDepth {
		return false
	}

	if len(t.text) == 0 {
		return false
	}
	if len(t.url) == 0 || strings.Contains(strings.ToLower(t.url), "javascript") {
		return false
	}

	return true
}

func (t HTTPError) Error() string {
	return t.original
}

func initializeServer() {
	port := flag.String("p", "8000", "por")
	dir := flag.String("d", ".", "dir")

	flag.Parse()

	http.Handle("/", http.FileServer(http.Dir(*dir)))
	log.Printf("Serving on %s on Http port: %s\n", *dir, *port)
	log.Fatal(http.ListenAndServe(":"+*port, nil))
}
