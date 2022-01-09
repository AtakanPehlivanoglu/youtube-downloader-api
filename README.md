# youtube-downloader-api
This is a simple Non-Copyrighted Youtube Media Downloader project which demonstrates how [go-colly](https://github.com/gocolly/colly) scrapper can be combined with [gin-gonic](https://github.com/gin-gonic/gin) web framework to stream media data using goroutines

Documentation about API can be found in swagger

# Usage with Go
1. Run ```go run main.go``` in the root folder
2. Go to http://localhost:8080/swagger/index.html 
 
# Usage with Docker 
1. Build ```docker build --tag youtube-downloader-api .```
2. Run ```docker run -p 8080:8080  youtube-downloader-api```
3. Go to http://localhost:8080/swagger/index.html  

### Notes
- Only **Non-Copyrighted** Youtube media is allowed to be downloaded otherwise request will be rejected
- `http.Transport` parameters can be modified in order to download larger files more efficiently
