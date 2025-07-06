package main

import (
	"context"
	"fmt"
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
	execTime := time.Now().Add(time.Duration(delay) * time.Second).Unix()
	
	ctx  := context.Background()
	redis_key := "task:" + taskid

	err := r.HSet(ctx, redis_key ,
		"taskid" , taskid,
		"url" , url	,
		"delay", int64(execTime)).Err()
	if err != nil{
		return err
	}

	return nil
}

func set_sorted_set(r *redis.Client, taskid string, delay int) error {
	execTime := time.Now().Add(time.Duration(delay) * time.Second).Unix()
	ctx := context.Background()

	err := r.ZAdd(ctx, "tasks_schedular", redis.Z{
		Score: float64(execTime),
		Member: taskid,
	}).Err()
	if err != nil{
		return err
	}

	return nil
} 

func get_top(r *redis.Client) (Task , error) {
	ctx := context.Background()
	task_id, err := r.ZRange(ctx, "tasks_schedular", 0 , 0).Result()
	if err != nil {
		return Task{} , err
	} 
	if len(task_id) == 0{
		return Task{} , redis.Nil
	}
	fmt.Println("Task id: ", task_id)

	redis_key := "task:" + task_id[0] 
	t , err := r.HGetAll(ctx, redis_key).Result()
	if(err != nil){
		return Task{} , err
	}

	t_id , t_url, t_delay_str := t["taskid"] , t["url"], t["delay"]
	
	t_delay_int , err := strconv.ParseInt(t_delay_str, 10, 64)
	t_delay := time.Unix(t_delay_int, 0)
	if err != nil {
		return Task{} , err
	}

	task := Task{taskid: t_id, content: t_url, exec_time: t_delay}
	return task , nil
}

func remove_from_db(r *redis.Client, t Task) error{
	ctx := context.Background()
	
	redis_key := "task:" + t.taskid
	_ , err := r.Del(ctx, redis_key).Result()
	if err != nil {
		return err
	}

	_ , err = r.ZRem(ctx, "tasks_schedular" , t.taskid).Result()
	if err != nil {
		return err
	}

	return nil
}