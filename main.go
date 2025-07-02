package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

func get_data(u string) (int , error) {
	resp , err := http.Get(u)
	if err != nil{
		return -1 , err
	}
	fmt.Println(resp.StatusCode)
	return resp.StatusCode , nil
}

func main(){
	wp := Newpool(3)
	wp.Start()

	http.HandleFunc("/serve", func(w http.ResponseWriter, r *http.Request) {
		var m map[string]interface{}
		body , _ := io.ReadAll(r.Body)
		json.Unmarshal(body, &m)
		urlStr, _ := m["url"].(string)
		delayFloat , _ := m["delay"].(float64)
		delay := int(delayFloat)
		wp.Submit(urlStr, time.Now().Add(time.Duration(delay) * time.Second))
	})

	go func(){
		for r := range wp.result{
		fmt.Printf("workerid : %d success: %v\n", r.workerid, r.status)
		}
	}()
	http.ListenAndServe(":5000", nil)

	
}	