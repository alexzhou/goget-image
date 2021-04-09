package main

import (
	"fmt"
	"golang.org/x/net/html"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"sync"
)

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

	//create dir
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, os.ModePerm)
		fmt.Println("create dir fail")
	}

	items := make([]string, 0)

	for {
		tt := domes.Next()
		if tt == html.ErrorToken {
			break
		} else if tt == html.StartTagToken {
			t := domes.Token()
			isImage := t.Data == "img"
			if isImage {
				for _, value := range t.Attr {
					if value.Key == "src" {
						src := value.Val
						var fullUrl string
						if strings.Contains(src, "http") {
							fullUrl = src
						} else {
							if HasPrefix(src, "/") {
								u, _ := url.Parse(pageUrl)
								fullUrl = u.Scheme + "://" + u.Host + src
							} else {
								fullUrl = pageUrl + "/" + src
							}
						}
						items = append(items, fullUrl)
					}
				}
			}
		}
	}

	var wg sync.WaitGroup
	nums := len(items)
	wg.Add(nums)
	for _, v := range items {
		go func(imageUrl string, imageDir string, wg *sync.WaitGroup) {
			defer wg.Done()
			fmt.Println("start download " + imageUrl)

			//get content
			response, err := http.Get(imageUrl)
			if err != nil {
				fmt.Println(err)
			}
			defer response.Body.Close()

			if response.StatusCode != 200 {
				fmt.Println("Received non 200 response code")
				return
			}

			//create file path ans create file
			_, fName := path.Split(imageUrl)
			var fileName string
			if HasSuffix(imageDir, "/") {
				fileName = imageDir + fName
			} else {
				fileName = imageDir + "/" + fName
			}

			file, err := os.Create(fileName)
			if err != nil {
				fmt.Println("create " + fileName + " fail")
				return
			}
			defer file.Close()

			//Write the bytes to the file
			_, err = io.Copy(file, response.Body)
			if err != nil {
				fmt.Println(err)
			}
			return
		}(v, dir, &wg)
	}
	wg.Wait()
	fmt.Println("all done")
}
