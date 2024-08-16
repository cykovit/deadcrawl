package main

import (
	"bufio"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/publicsuffix"
)

type LinkChecker struct {
	baseURL  *url.URL
	client   *http.Client
	deadLink bool // flag to track if any dead links are found
}

func NewLinkChecker(baseURL *url.URL) *LinkChecker {
	return &LinkChecker{
		baseURL:  baseURL,
		client:   &http.Client{},
		deadLink: false,
	}
}

// recursively go through the HTML document tree and search for <a> tags
func (lc *LinkChecker) CheckLinks(node *html.Node) {
	if node.Type == html.ElementNode && node.Data == "a" {
		for _, attr := range node.Attr {
			if attr.Key == "href" {
				link, err := lc.parseLink(attr.Val)
				if err != nil {
					fmt.Println("Error parsing link:", err)
					continue
				}
				// check if the link is dead or alive
				if err := lc.checkLink(link); err != nil {
					fmt.Printf("Error checking link %s: %s\n", link, err)
				}
			}
		}
	}

	// recursively check all other links inside the document
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		lc.CheckLinks(child)
	}
}

// if the link is relative, resolve it against the base URL
func (lc *LinkChecker) parseLink(href string) (*url.URL, error) {
	link, err := url.Parse(href)
	if err != nil {
		return nil, err
	}
	// if the link doesn't have a scheme (ex:http), assume it's relative and resolve it anyway
	if link.Scheme == "" {
		return lc.baseURL.ResolveReference(link), nil
	}
	return link, nil
}

// send an HTTP request to the given URL to check if it can be reached
func (lc *LinkChecker) checkLink(link *url.URL) error {
	// ignore URLs that don't have a valid domain (ex:mailto)
	if _, err := publicsuffix.EffectiveTLDPlusOne(link.Hostname()); err != nil {
		return nil
	}

	// create a new GET request
	req, err := http.NewRequest("GET", link.String(), nil)
	if err != nil {
		return err
	}

	// set headers to mimic a real web browser and avoid 403 errors
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Referer", lc.baseURL.String())

	resp, err := lc.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// if the response status is 400 or higher, the link is considered dead
	if resp.StatusCode >= 400 {
		fmt.Printf("Dead link found: %s (%d)\n", link, resp.StatusCode)
		lc.deadLink = true
	}
	return nil
}

// retrieve the HTML document from the given URL
func fetchHTML(url string) (*html.Node, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// parse the HTML document into a node tree structure
	doc, err := html.Parse(resp.Body)
	if err != nil {
		return nil, err
	}
	return doc, nil
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	// prompt user for the target URL + start, error and failure messages
	fmt.Print("Enter the URL of the website to check for dead links: ")
	targetURL, _ := reader.ReadString('\n')
	targetURL = strings.TrimSpace(targetURL)
	baseURL, err := url.Parse(targetURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing base URL: %v\n", err)
		os.Exit(1)
	}

	checker := NewLinkChecker(baseURL)

	fmt.Println("Checking for dead links. This may take some time...")

	doc, err := fetchHTML(targetURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching HTML: %v\n", err)
		os.Exit(1)
	}

	checker.CheckLinks(doc)

	if !checker.deadLink {
		fmt.Println("No dead links found.")
	}
}
