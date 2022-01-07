package scraper

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"sort"
	"time"
)

var client *http.Client

func init() {
	t := &http.Transport{
		Dial: (&net.Dialer{
			Timeout:   60 * time.Second,
			KeepAlive: 30 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 60 * time.Second,
	}
	client = &http.Client{Transport: t}
}

func BodyResponse(contentLength int64, url string, title string) (bodyResponse []byte) {
	var channels []<-chan []byte
	var streamChannels []Stream
	var ranges []string

	batchSize := 500000
	byteRange := 0
	param := byteRangeBuilder(0, int64(batchSize))
	div := contentLength / int64(batchSize)

	for i := 0; int64(i) < div; i++ {
		ranges = append(ranges, param)
		byteRange += batchSize
		param = byteRangeBuilder(byteRange+1, int64(byteRange+batchSize))
	}
	ranges = append(ranges, byteRangeBuilder(byteRange+1, contentLength))

	log.Printf("Stream Batch size is : %v \n", len(ranges))

	for _, r := range ranges {
		channels = append(channels, streamData(r, url))
	}

	order, result := fanIn(channels)

	for i := 0; i < len(ranges); i++ {
		streamChannels = append(streamChannels, Stream{Order: <-order, Body: <-result})
	}

	sort.Slice(streamChannels, func(i, j int) bool {
		return streamChannels[i].Order < streamChannels[j].Order
	})

	for _, v := range streamChannels {
		bodyResponse = append(bodyResponse, v.Body...)
	}
	return bodyResponse
}

func streamData(bytesRange string, url string) <-chan []byte {
	chBody := make(chan []byte)
	go func() {
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			fmt.Printf("Erro Http New Request: \n" + err.Error())
			close(chBody)
			return
		}
		req.Header.Set("range", bytesRange)

		resp, err := client.Do(req)

		if err != nil {
			fmt.Printf("Error client request: \n" + err.Error())
			close(chBody)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusPartialContent {
			fmt.Printf("Error status code: %v \n", resp.StatusCode)
			close(chBody)
			return
		}
		bodyResponse, err := ioutil.ReadAll(resp.Body)
		chBody <- bodyResponse
	}()
	return chBody
}

func fanIn(channels []<-chan []byte) (<-chan int, <-chan []byte) {
	chOrder := make(chan int)
	chBody := make(chan []byte)

	for i, v := range channels {
		go func(ch <-chan []byte, order int) {
			for {
				chOrder <- order
				chBody <- <-ch
			}
		}(v, i)
	}
	return chOrder, chBody
}

func byteRangeBuilder(r1 int, r2 int64) string {
	return fmt.Sprintf("bytes=%d-%d", r1, r2)
}
