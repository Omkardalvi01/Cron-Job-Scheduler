package main

import (
	"fmt"
	"time"
	"github.com/redis/go-redis/v9"
)

type Task struct {
	taskid string
	content string 
	exec_time     time.Time 
}
type Result struct {
	workerid int
	status int
}

type worker struct {
	workerid   int
	task       <-chan Task
	resultchan chan<- Result
}

func (w worker) Start() {
	for t := range w.task {
		fmt.Println("Worker has started working")
		success, err := get_data(t.content)
		if err != nil{
			w.resultchan <- Result{workerid : w.workerid , status: success}
			fmt.Print(err)
			return
		}
		w.resultchan <- Result{workerid: w.workerid, status: success}
	}
}

type workerpool struct {
	worker_num int
	taskqueue  chan Task
	result     chan Result
}

func Newpool(num int) *workerpool {
	return &workerpool{
		worker_num: num,
		taskqueue:  make(chan Task),
		result:     make(chan Result),
	}
}

func (wp workerpool) Start() {
	for i := 0; i < wp.worker_num; i++ {
		w := worker{workerid: i, task: wp.taskqueue, resultchan: wp.result}
		go w.Start()
	}
}

func (wp workerpool)  Submit(cancel chan struct{} , new_entry chan struct{}, r *redis.Client){

	t, err := get_top(r)
	if err == redis.Nil{
		<-new_entry
		go wp.Submit(cancel , new_entry, r)
		return 
	}
	if err != nil{
		fmt.Println(err)
	}
	

	fmt.Println("Duration until execution:", time.Until(t.exec_time))
	for {
		select{
		case <-cancel:
			fmt.Println("Context Cancelled")
			return
		default:
			if time.Now().After(t.exec_time) {
				wp.taskqueue <- t

				err := remove_from_db(r , t)
				if err != nil {
					fmt.Println("Error at remove from db: ", err)
				}

				go wp.Submit(cancel , new_entry, r)
				return
			}
			
		}
	}
	
}

func (wp workerpool) Cancel_Submit(cancel chan struct{} ,r *redis.Client ,delay int, new_entry chan struct{}){
	t, err := get_top(r)
	if err != nil {
		fmt.Print("Eror: ",err)
	}
	top_delay := t.exec_time
	new_exec_time := time.Now().Add(time.Duration(delay) * time.Second)
	if top_delay.After(new_exec_time)  {
		fmt.Println("Task switched")
		cancel <- struct{}{}
		go wp.Submit(cancel, new_entry, r)
	}

}
