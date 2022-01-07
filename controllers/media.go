package controllers

import(
	"bytes"
	"fmt"
	"log"
	"net/http"
	"youtube-downloader-api/scraper"
	"github.com/gin-gonic/gin"
)

// @BasePath /api/v1

// YoutubeDownloaderApi godoc
// @Summary Download Audio 
// @Description Download Audio from Youtube Url
// @Tags Media
// @Accept json
// @Produce json
// @Param m query string  false  "Youtube url"
// @Success 200 {string} Success
// @Router /audio [get]
func GetAudio(c *gin.Context){
	audio := c.Query("m")

	contentLength, url, title := scraper.ParseMedia(audio, scraper.Audio)

	if contentLength == 0 {
		c.String(http.StatusBadRequest, "Error, Either url is wrong or media is copyrighted")
		return
	}

	if contentLength > 104857600 {
		c.String(http.StatusBadRequest, "Error, Media size is larger than 100 MB")
		return
	}

	bodyResponse := scraper.BodyResponse(contentLength, url, title)

	if bodyResponse == nil{
		c.String(http.StatusBadRequest, "Error Getting Body Response")
		return
	}
	log.Println("Title is : " + title)

	reader := bytes.NewReader(bodyResponse)

	contentType := "audio/mp4"
	
	extraHeaders := map[string]string{
		"Content-Disposition": fmt.Sprintf(`attachment; filename="%v".mp4`, title),
	}

   c.DataFromReader(http.StatusOK, contentLength, contentType, reader, extraHeaders)
}


// YoutubeDownloaderApi godoc
// @Summary Download Video 
// @Description Download Video from Youtube Url
// @Tags Media
// @Accept json
// @Produce json
// @Param m query string  false  "Youtube url"
// @Success 200 {string} Success
// @Router /video [get]
func GetVideo(c *gin.Context){
	video := c.Query("m")

	contentLength, url, title := scraper.ParseMedia(video, scraper.Video)

	if contentLength == 0 {
		c.String(http.StatusBadRequest, "Error, Either url is wrong or media is copyrighted")
		return
	}

	if contentLength > 104857600 {
		c.String(http.StatusBadRequest, "Error, Media size is larger than 100 MB")
		return
	}

	bodyResponse := scraper.BodyResponse(contentLength, url, title)

	if bodyResponse == nil{
		c.String(http.StatusBadRequest, "Error Getting Body Response")
		return
	}
	log.Println("Title is : " + title)

	reader := bytes.NewReader(bodyResponse)
	contentType := "video/mp4"
	
	extraHeaders := map[string]string{
		"Content-Disposition": fmt.Sprintf(`attachment; filename="%v".mp4`, title),
	}

   c.DataFromReader(http.StatusOK, contentLength, contentType, reader, extraHeaders)
	
}