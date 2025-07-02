package main

import (
	"fmt"
	"testing"
	"time"
)

func Test_timr(t *testing.T){
	fmt.Print(time.Until(time.Now().Add(time.Duration(5) * time.Second)).Milliseconds())
}