package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

type Test struct {
	Message string `json:"message,omitempty"`
	Status  int    `json:"status,omitempty"`
}

type Log struct {
	Date     string `json:"date,omitempty"`
	Username string `json:"username,omitempty"`
	// Id       string `json:"id,omitempty"`
}

var ctx = context.Background()

var rdb = redis.NewClient(&redis.Options{
	Addr: "localhost:6379",
	DB:   3,
})

func createLog(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	data := Log{}
	err = json.Unmarshal(reqBody, &data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	data.Date = time.Now().String()
	logId := uuid.New().String()
	fmt.Println(data)
	rawData, _ := json.Marshal(data)
	err = rdb.Set(ctx, logId, rawData, 0).Err()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(fmt.Sprintf("{\"log_id\": \"%s\"}", logId)))
}

func ping(w http.ResponseWriter, r *http.Request) {
	test_data := Test{
		"message",
		200,
	}
	json_data, err := json.Marshal(test_data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println(json_data)
	w.Header().Set("Content-Type", "application/json")
	w.Write(json_data)
}

func main() {
	http.HandleFunc("/ping", ping)
	http.HandleFunc("/log", createLog)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
