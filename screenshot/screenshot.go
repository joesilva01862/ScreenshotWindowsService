/**
 * 
 * Program to take a screenshot and send it to the target server (webdav
 * in most cases). It requires you to create these entries in the Windows registry:
 * 
 *  TARGET_SERVER_TYPE (-d for webdav)
 *  D_UPLOADURL (your webdav url)
 *  D_USERNAME (your webdav username)
 *  D_PASSWORD (your webdav password)
 * 
 * It was originally intended to be called from an accompanying Windows service,
 * but you can invoke it from the command line after creating the Registry
 * enties above.
 * 
 */

package main

import (
	"net/url"
	"flag"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
    "net/http/httputil"
   	"os"
	"time"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"errors"

    "github.com/studio-b12/gowebdav"    
    "golang.org/x/sys/windows/registry"
    "github.com/flopp/go-findfont"
	"github.com/golang/freetype"
    "github.com/golang/freetype/truetype"
    "github.com/kbinani/screenshot"
)

var (
	username string
	password string
	uploadURL string
    client  *gowebdav.Client
	dpi      = flag.Float64("dpi", 72, "screen resolution in Dots Per Inch")
	fontfile = flag.String("fontfile", "luxisr.ttf", "filename of the ttf font")
	hinting  = flag.String("hinting", "none", "none | full")
	size     = flag.Float64("size", 60, "font size in points")
	spacing  = flag.Float64("spacing", 1.5, "line spacing (e.g. 2 means double spaced)")
	wonb     = flag.Bool("whiteonblack", false, "white text on a black background")
)

func postLogin() *http.Cookie {
	//log.Println("here in postLogin()")
	client := http.Client{}

    // define behavior when there is a redirect
    client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
        //log.Println("here at CheckRedirect")
        return http.ErrUseLastResponse
    }

	apiUrl := uploadURL
    resource := "/autoshots/login"
    loginData := url.Values{}
 	loginData.Set("username", username)
 	loginData.Set("password", password)
    u, _ := url.ParseRequestURI(apiUrl)
    u.Path = resource
    urlStr := u.String() 
	
    bodyString := `username=admin_user&password=password` 
    bodyData := bytes.NewReader([]byte(bodyString))

	loginRequest, err := http.NewRequest("POST", urlStr, bodyData) // URL-encoded payload
    loginRequest.Header.Set("User-Agent", "Go-Program")
    loginRequest.Header.Set("Accept", "*/*")
    loginRequest.Header.Set("Content-Type", "application/x-www-form-urlencoded")
    loginRequest.Header.Set("Content-Length", "37")
    resp, err := client.Do(loginRequest)
    if err != nil {
       return nil
	}
    defer resp.Body.Close()
    return resp.Cookies()[0]
}

func printRequest(req *http.Request) {
	output, err := httputil.DumpRequest(req, true)
	if err != nil {
		panic(err)
	}
	log.Println("-----------------BEGIN---------------------------")
	log.Println(string(output))
	log.Println("------------------END----------------------------")
}

func printResponse(resp *http.Response, msg string) {
    bodyBytes, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        log.Fatal(err)
    }
    bodyString := string(bodyBytes)
    log.Println(msg, "\n", bodyString)
}

func printResponse1(resp *http.Response) {
	body := &bytes.Buffer{}
	
    _, err := body.ReadFrom(resp.Body)

	if err != nil {
		log.Fatal(err)
	}

	resp.Body.Close()
	fmt.Println(resp.StatusCode)
	fmt.Println(resp.Header)
	fmt.Println(body)
}

func mustOpen(filePath string) *os.File {
    r, err := os.Open(filePath)
    if err != nil {
        pwd, _ := os.Getwd()
        fmt.Println("PWD: ", pwd)
        panic(err)
    }
    return r
}

/* take a screenshot and place a timestamp on it */
func takeScreenshot(timestamp string) bytes.Buffer {
    // find a font by name and read into memory
    // pick a font from the system
    fontName := "arial.ttf"
    fontAlternative := "FreeSans.ttf"
    fontPath, err := findfont.Find(fontName)
    if err != nil {
       fontPath, err = findfont.Find(fontAlternative)
       if err != nil {
          panic(err)
       }
    }
    
    // load the font with the freetype library
    fontData, err := ioutil.ReadFile(fontPath)
    if err != nil {
      panic(err)
    }
    font, err := truetype.Parse(fontData)
    if err != nil {
      panic(err)
    }
  
	// Initialize the context.
	fg := image.White
	
	rectImage := image.NewRGBA(image.Rect(0, 20, 580, 100))
	
    // Colors are defined by Red, Green, Blue, Alpha uint8 values.
    red := color.NRGBA{255, 0, 0, 100}
    //    green := color.NRGBA{0, 255, 0, 60}
    //    blue := color.NRGBA{0, 0, 255, 60}
    //    yellow := color.NRGBA{255, 255, 0, 50}
    
	var all image.Rectangle = image.Rect(0, 0, 0, 0)
	bounds := screenshot.GetDisplayBounds(0)
	all = bounds.Union(all)
	rgba, err := screenshot.CaptureRect(bounds)
	if err != nil {
		panic(err)
	}
	
	// draw background into the main image
	// func Draw(dst Image, r image.Rectangle, src image.Image, sp image.Point, op Op)
	//draw.Draw(rgba, myImage.Bounds(), myImage, image.ZP, draw.Src)

    // draw a red rectangle atop the green surface
    draw.Draw(rgba, rectImage.Bounds(), &image.Uniform{red}, image.ZP, draw.Over)
    	
	c := freetype.NewContext()
	c.SetDPI(*dpi)
	c.SetFont(font)
	c.SetFontSize(*size)
	c.SetClip(rgba.Bounds())
	c.SetDst(rgba)
	c.SetSrc(fg)
	
	// Draw the text.
    pt := freetype.Pt(0, 20+int(c.PointToFixed(*size)>>6))
	_, err = c.DrawString(timestamp, pt)

    // create a bytes Buffer
    var imageBuf bytes.Buffer

	// Encode that RGBA image with PNG format into an byte buffer.
	err = png.Encode(&imageBuf, rgba)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	
	return imageBuf
}

func takeScreenshotAndSend() error {
    // create the data for the json objects
    hostname, _ := os.Hostname()
    t := time.Now()
    timestamp :=  fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d",
        t.Year(), t.Month(), t.Day(),
        t.Hour(), t.Minute(), t.Second())

    // prepare vars to create the various parts
    var b bytes.Buffer
    var err error
    w := multipart.NewWriter(&b)
    var fw io.Writer
    
    // create PART 1
    if fw, err = w.CreateFormFile("pictureFile", "anything"); err != nil {
        panic(err)
    }
    
    screenshot := takeScreenshot(timestamp)
    fw.Write(screenshot.Bytes())
    
    
    // create PART 2
    if fw, err = w.CreateFormFile("jsonObject", ""); err != nil {
        panic( err)
    }

    for _, line := range getPictureData(timestamp, hostname) {
		if _, err = io.WriteString(fw, line); err != nil {
			panic( err)
		}
    }


    // create PART 3
    if fw, err = w.CreateFormFile("jsonEventObject", ""); err != nil {
        panic( err)
    }
    
    for _, line := range getEventList(timestamp) {
		if _, err = io.WriteString(fw, line); err != nil {
			panic( err)
		}
    }

    // now close multipart writer
    w.Close()

    // now prepare request
    req, err := http.NewRequest("POST", uploadURL, &b)
    if err != nil {
        return err
    }
    
    // Don't forget to set the content type, this will contain the boundary.
    req.Header.Set("Content-Type", w.FormDataContentType())
    
    cookie := postLogin()
    
    if cookie == nil {
        return errors.New("Server unavailable")	
	}
	
	req.AddCookie( cookie )
	
	// send request
	client := &http.Client{}
	_, err = client.Do(req)
    
	if err != nil {
		log.Fatal(err)
		panic(err)
    }
    
    return nil
}

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func takeScreenshotSendToWebdav() error {
    t := time.Now()
    targetName :=  fmt.Sprintf("screenshot_%d%02d%02d_%02d%02d%02d%s",
        t.Year(), t.Month(), t.Day(),
        t.Hour(), t.Minute(), t.Second(), ".jpg")
    timestamp :=  fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d",
        t.Year(), t.Month(), t.Day(),
        t.Hour(), t.Minute(), t.Second())
        
    // take a screenshot with a timestamp on it    
    screenshot := takeScreenshot(timestamp)
    
    // send to a webdav server (nextcloud in my case)
	err := client.Write(targetName, screenshot.Bytes(), 0664)
    check(err)    
    return nil
}

func getEventList(timestamp string) []string {
	events := make([]string, 1)
	events = append(events, "[")
	events = append(events, fmt.Sprintf("{ \"timestamp\": \"%v\", \"event\" : \"%v\" }, ", timestamp, "Fake mouse event") )
	events = append(events, fmt.Sprintf("{ \"timestamp\": \"%v\", \"event\" : \"%v\" }, ", timestamp, "Fake keyboard event") )
	events = append(events, fmt.Sprintf("{ \"timestamp\": \"%v\", \"event\" : \"%v\" } ", timestamp, "Fake USB event") )
	events = append(events, "]")
	return events
}

func getPictureData(timestamp string, hostname string) []string {
	json := make([]string, 1)
	json = append(json, "{")
	json = append(json, fmt.Sprintf("\"timestamp\": \"%v\", \"computer\" : \"%v\", \"idleMilliSecs\": \"%v\", \"intervalMilliSecs\": \"%v\"",
	timestamp, hostname, 3451234, 899611) )
	json = append(json, "}")
	return json
}

func main() {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Microservice\Netpipe`, registry.QUERY_VALUE)
	check(err)

    var target_type string
    target_type, _, err = k.GetStringValue("TARGET_SERVER_TYPE")
	
    if target_type == "-d" {
        username, _, err = k.GetStringValue("D_USERNAME")
        check(err)
		password, _, err = k.GetStringValue("D_PASSWORD")
		check(err)
		uploadURL, _, err = k.GetStringValue("D_UPLOADURL")
		check(err)
        client = gowebdav.NewClient(uploadURL, username, password)
        takeScreenshotSendToWebdav()
	} else if target_type == "-s" {
		username, _, err = k.GetStringValue("S_USERNAME")
		check(err)
		password, _, err = k.GetStringValue("S_PASSWORD")
		check(err)
		uploadURL, _, err = k.GetStringValue("S_UPLOADURL")
		check(err)
        takeScreenshotAndSend()
	}
}
