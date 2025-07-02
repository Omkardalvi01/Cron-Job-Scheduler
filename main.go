package main

import (
	"fmt"
	"net/http"
	"time"
)

func get_data(u string) (int , error) {
	resp , err := http.Get(u)
	if err != nil{
		return 0 , err
	}
	fmt.Println(resp.StatusCode)
	return resp.StatusCode , nil
}

func main(){
	wp := Newpool(3)
	wp.Start()
	urls := []string{"https://www.google.com/","https://www.pinterest.com/"}
	for i, u := range urls{
		go wp.Submit(u , time.Now().Add(time.Duration(i + 2) * time.Second))
	}
	for i := 0; i < len(urls); i++{
		r := wp.getResult()
		fmt.Printf("workerid : %d success: %v\n", r.workerid, r.status)
	}
}	