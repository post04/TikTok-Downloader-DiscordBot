package tiktok

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"strings"
	"time"

	goquery "github.com/PuerkitoBio/goquery"
)

// I forked https://github.com/pawanpaudel93/go-tiktok-downloader for what I needed

// Profile - Tiktok Profile
type Profile struct {
	URL        string
	filePath   string
	data       PageProps
	httpClient *http.Client
	Proxy      string
	BaseDIR    string
}

// Video - Tiktok Video
type Video struct {
	URL        string
	filePath   string
	data       VideoData
	Proxy      string
	httpClient *http.Client
	BaseDIR    string
}

var headers map[string]string
var cookies []*http.Cookie

func generateRandomNumber() string {
	max := 1999999999999999999
	min := 1000000000000000000
	return strconv.Itoa(min + rand.Intn(max-min))
}

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
	headers = map[string]string{
		"User-Agent":      "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.66 Safari/537.36",
		"Accept":          "*/*",
		"Connection":      "keep-alive",
		"Pragma":          "no-cache",
		"Cache-Control":   "no-cache",
		"Sec-Fetch-Site":  "same-site",
		"Sec-Fetch-Mode":  "no-cors",
		"Sec-Fetch-Dest":  "video",
		"Referer":         "https://www.tiktok.com/",
		"Accept-Language": "en-US,en;q=0.9,bs;q=0.8,sr;q=0.7,hr;q=0.6",
		"sec-gpc":         "1",
		"DNT":             "1",
		"Range":           "bytes=0-",
	}
	webID := generateRandomNumber()
	expiration := time.Now().Add(365 * 24 * time.Hour)
	cookies = []*http.Cookie{
		{Name: "tt_webid", Value: webID, Path: "/", Expires: expiration, Domain: ".tiktok.com"},
		{Name: "tt_webid_v2", Value: webID, Path: "/", Expires: expiration, Domain: ".tiktok.com"},
	}
}

func (video *Video) setProxy() {
	if !(strings.Contains(video.Proxy, "http://") || strings.Contains(video.Proxy, "https://")) {
		video.Proxy = "http://" + string(video.Proxy)
	}
}

func (video *Video) setClient(jar *cookiejar.Jar) {
	video.httpClient = &http.Client{
		Jar: jar,
	}
	if video.Proxy != "" {
		if proxyURL, err := url.Parse(video.Proxy); err == nil {
			transport := &http.Transport{Proxy: http.ProxyURL(proxyURL)}
			video.httpClient = &http.Client{
				Jar:       jar,
				Transport: transport,
			}
		} else {
			fmt.Println(err)
			fmt.Println("Not Using Proxy")
		}
	}
}

func (profile *Profile) setProxy() {
	if !(strings.Contains(profile.Proxy, "http://") || strings.Contains(profile.Proxy, "https://")) {
		profile.Proxy = "http://" + string(profile.Proxy)
	}
}

func (profile *Profile) setClient(jar *cookiejar.Jar) {
	profile.httpClient = &http.Client{
		Jar: jar,
	}
	if profile.Proxy != "" {
		if proxyURL, err := url.Parse(profile.Proxy); err == nil {
			transport := &http.Transport{Proxy: http.ProxyURL(proxyURL)}
			profile.httpClient = &http.Client{
				Jar:       jar,
				Transport: transport,
			}
		} else {
			fmt.Println(err)
			fmt.Println("Not Using Proxy")
		}
	}
}

// Download - Download Tiktok video
func (video *Video) Download() (io.Reader, error) {
	jar, _ := cookiejar.New(nil)
	URL := video.data.Video.URL
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		return nil, err
	}
	parsedURL, _ := url.Parse(URL)

	for k, v := range headers {
		req.Header.Set(k, v)
	}
	jar.SetCookies(parsedURL, cookies)
	video.setClient(jar)
	resp, err := video.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	return resp.Body, err
}

// FetchInfo - Get Tiktok video Information.
func (video *Video) FetchInfo() error {
	jar, _ := cookiejar.New(nil)
	videoURL := video.URL
	if video.Proxy != "" {
		video.setProxy()
	}

	req, err := http.NewRequest("GET", videoURL, nil)
	if err != nil {
		return err
	}
	parsedURL, _ := url.Parse(videoURL)
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	jar.SetCookies(parsedURL, cookies)
	video.setClient(jar)
	resp, err := video.httpClient.Do(req)
	if err != nil {
		return err
	}
	VideoData := VideoData{}
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return err
	}
	doc.Find("#__NEXT_DATA__").Each(func(i int, s *goquery.Selection) {
		err = json.Unmarshal([]byte(s.Text()), &VideoData)
		video.data = VideoData
	})
	return err
}
