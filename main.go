package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
	"github.com/redis/go-redis/v9"
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

	client := redis.NewClient(&redis.Options{
        Addr:	  "localhost:6379",
        Password: "", 
        DB:		  0,  
        Protocol: 2,  
    })

	wp := Newpool(3)
	wp.Start()

	http.HandleFunc("/serve", func(w http.ResponseWriter, r *http.Request) {

		var m map[string]interface{}
		body , _ := io.ReadAll(r.Body)
		json.Unmarshal(body, &m)

		task_id, err:= create_taskid(client)
		if err != nil{
			fmt.Print("Error at create_task_id", err)
		}

		urlStr, _ := m["url"].(string)
		delayFloat , _ := m["delay"].(float64)
		delay := int(delayFloat)
		
		err = set_task_hset(client , task_id, urlStr, delay)
		if err != nil{
			fmt.Print("Error at set_task_hset", err)
		}

		err = set_sorted_set(client, task_id, delay)
		if err != nil{
			fmt.Print("Errpr at sorted task set", err)
		}

		wp.Submit(urlStr, time.Now().Add(time.Duration(delay) * time.Second))
	})

	go func(){
		for r := range wp.result{
		fmt.Printf("workerid : %d success: %v\n", r.workerid, r.status)
		}
	}()

	http.ListenAndServe(":5000", nil)

	
}	