package main

import (
	"fmt"
	"testing"
	"time"
	"github.com/redis/go-redis/v9"
)

func Test_timr(t *testing.T){
	fmt.Print(time.Until(time.Now().Add(time.Duration(5) * time.Second)).Milliseconds())
}

func Test_redis_integration_test(t *testing.T){
	client := redis.NewClient(&redis.Options{
        Addr:	  "localhost:6379",
        Password: "", 
        DB:		  0,  
        Protocol: 2,  
    })
	fmt.Println("Client :",client)

	str, err := create_taskid(client)
	if err != nil {
		t.Errorf("Error at create_taskid: %v\n", err)
	}
	fmt.Println("Task_id :", str)

	err = set_task_hset(client, str, "https://pinterest.com", 10)
	if err != nil {
		t.Errorf("Error at set task hset: %v", err)
	}

	err = set_sorted_set(client, str , 10)
	if err != nil {
		t.Errorf("Error at set sorted set: %v", err)
	}
}