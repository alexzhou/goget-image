package main

import (
	"errors"
	"fmt"
	"golang.org/x/net/html"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
)

func getImage(url string, dir string) error {
	fmt.Println("start download " + url)

	//get content
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return errors.New("Received non 200 response code")
	}
	//create dir
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		_ = os.MkdirAll(dir, os.ModePerm)
	}

	//create file path ans create file
	_, fName := path.Split(url)
	var fileName string
	if HasSuffix(dir, "/") {
		fileName = dir + fName
	} else {
		fileName = dir + "/" + fName
	}

	fmt.Println("create file: " + fileName)
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	//Write the bytes to the file
	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}

	return nil
}

func HasPrefix(s, prefix string) bool {
	return len(s) >= len(prefix) && s[0:len(prefix)] == prefix
}

func HasSuffix(s, suffix string) bool {
	return len(s) >= len(suffix) && s[len(s)-len(suffix):] == suffix
}

func Download(pageUrl string, dir string) {
	fmt.Println("start get images from " + pageUrl)
	resp, err := http.Get(pageUrl)
	// handle the error if there is one
	if err != nil {
		panic(err)
	}
	// do this now so it won't be forgotten
	defer resp.Body.Close()
	// reads html as a slice of bytes

	domes := html.NewTokenizer(resp.Body)

	for {
		tt := domes.Next()

		switch {
		case tt == html.ErrorToken:
			// End of the document, we're done
			return
		case tt == html.StartTagToken:
			t := domes.Token()

			isImage := t.Data == "img"
			if isImage {
				for _, value := range t.Attr {
					if value.Key == "src" {
						src := value.Val
						fmt.Println(src)
						var err2 error
						if strings.Contains(src, "http") {
							err2 = getImage(src, dir)
						} else {
							if HasPrefix(src, "/") {
								u, _ := url.Parse(pageUrl)
								err2 = getImage(u.Scheme+"://"+u.Host+src, dir)
							} else {
								err2 = getImage(pageUrl+"/"+src, dir)
							}
						}
						if err2 != nil {
							fmt.Println(src + " can't download")
						}
					}
				}
			}
		}
	}
}
