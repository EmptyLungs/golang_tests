package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

var imageUrl = "https://loremflickr.com/g/1920/1080/cat"

func MakeRequest(ch chan<- string) {
	resp, _ := http.Get(imageUrl)
	body, _ := ioutil.ReadAll(resp.Body)
	fileName := fmt.Sprintf("./images/%s.jpg", uuid.New().String())
	err := ioutil.WriteFile(fileName, body, 0644)
	if err != nil {
		log.Fatal(err.Error())
	}
	ch <- fileName
}

func Concurrent() {
	start := time.Now()
	ch := make(chan string)
	for i := 0; i <= 10; i++ {
		go MakeRequest(ch)
	}
	for i := 0; i <= 10; i++ {
		fmt.Println(<-ch)
	}

	fmt.Printf("%.2fs elapsed \n", time.Since(start).Seconds())
}
