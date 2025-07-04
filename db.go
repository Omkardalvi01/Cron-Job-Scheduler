package main

import (
	"context"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

func create_taskid(r *redis.Client) (string,  error) {
	var temp string
	ctx := context.Background()

	for{
		temp = uuid.NewString()
		_ , err := r.HGet(ctx, "Task_Hset", temp).Result()
		if err == redis.Nil {
			break
		}
		if err != nil {
			return "", err
		}
	}
	return temp , nil
}

func set_task_hset(r *redis.Client, taskid, url string, delay int ) error {
	ctx  := context.Background()
	values := []string{"taskid" , taskid , "url" , url, "delay" , strconv.Itoa(delay)}
	err := r.HSet(ctx, "Task_HSet" , values).Err()
	if err != nil{
		return err
	}
	return nil
}

func set_sorted_set(r *redis.Client, taskid string, delay int) error {
	execTime := time.Now().Add(time.Duration(delay) * time.Second).Unix()
	ctx := context.Background()
	err := r.ZAdd(ctx, "tasks_schedulae", redis.Z{
		Score: float64(execTime),
		Member: taskid,
	}).Err()

	if err != nil{
		return err
	}
	return nil
} 