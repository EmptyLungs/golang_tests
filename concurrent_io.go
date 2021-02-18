package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/google/uuid"
)

var imageUrl = "https://loremflickr.com/g/1920/1080/cat"

func MakeRequest(wg *sync.WaitGroup) {
	resp, _ := http.Get(imageUrl)
	body, _ := ioutil.ReadAll(resp.Body)
	fileName := fmt.Sprintf("./images/%s.jpg", uuid.New().String())
	err := ioutil.WriteFile(fileName, body, 0644)
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Printf("done getting file %s\n", fileName)
	wg.Done()
}

func DownloadFiles(iter int) {
	wg := sync.WaitGroup{}
	for i := 1; i <= iter; i++ {
		wg.Add(1)
		go MakeRequest(&wg)
	}
	wg.Wait()

}

func checkOrCreateImgDir() {
	if _, err := os.Stat("./images"); os.IsNotExist(err) {
		err := os.Mkdir("./images", os.FileMode(0777))
		if err != nil {
			log.Fatal(err.Error())
		}
	}
}

func Concurrent() {
	start := time.Now()
	checkOrCreateImgDir()
	DownloadFiles(10)
	fmt.Printf("%.2fs elapsed \n", time.Since(start).Seconds())
}
