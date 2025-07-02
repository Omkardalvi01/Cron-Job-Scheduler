package main

import (
	"time"
	"fmt"
)

type Task struct {
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

func (wp workerpool)  Submit(c string, exec time.Time){
	t := Task{content: c, exec_time: exec }
	fmt.Println("Duration until execution:", time.Until(t.exec_time))

	time.AfterFunc(time.Until(t.exec_time), func(){ wp.taskqueue <- t })
}

func (wp workerpool) getResult() Result{
	return <-wp.result
}