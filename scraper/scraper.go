package scraper

import (
	"encoding/json"
	"log"
	"strings"
	"time"
	"strconv"
	"github.com/gocolly/colly"
)

type Stream struct{
	Order int 
	Body []byte
}

type Media struct {
	StreamingData struct {
		ExpiresInSeconds string `json:"expiresInSeconds"`
		AdaptiveFormats  []struct {
			Itag          int    `json:"itag"`
			Url           string `json:"url"`
			MimeType      string `json:"mimeType"`
			ContentLength string `json:"contentLength"`
		}
	}
}

const (
	Audio = iota
	Video
)

func ParseMedia(mediaUrl string, mediaType int) (contentLength int64, url string, title string) {
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2228.0 Safari/537.36"),
	)
	c.SetRequestTimeout(300 * time.Second)

	c.OnHTML("meta[name=\"title\"]", func(h *colly.HTMLElement) {
		title = h.Attr("content")

	})

	c.OnHTML("body > script:first-of-type", func(e *colly.HTMLElement) {
		jsonData := e.Text[strings.Index(e.Text, "{") : len(e.Text)-1]
		data := Media{}

		err := json.Unmarshal([]byte(jsonData), &data)
		if err != nil {
			log.Println("Error on marshall: "+ err.Error())
			return
		}
		if len(data.StreamingData.AdaptiveFormats) <= 0 {
			log.Println("No Data")
			return
		}

		for _, v := range data.StreamingData.AdaptiveFormats {
			if (mediaType == Audio && strings.Contains(v.MimeType, "audio/mp4")) ||
				(mediaType == Video && strings.Contains(v.MimeType, "video/mp4")) {
				
				contentLength, err = strconv.ParseInt(v.ContentLength, 10, 64)

				if err != nil {
					log.Println("Error Converting Content Length: " + err.Error())
				}

				url = v.Url

				log.Println("Media Url is: " + v.Url)
				log.Println("Content length is: " + v.ContentLength)
				break
			}
		}
	})

	c.OnError(func(r *colly.Response, e error) {
		log.Println("Error: ", e, r.Request.URL, string(r.Body))
	})

	c.Visit(mediaUrl)

	return contentLength, url, title
}

