// Go 1.2
// go run threads_in_go.go

package main

import 
(
    ."fmt"
    "runtime"
    "time"
)

//This is a global variable
var i int = 0
var availible = make(chan int, 1)


func thread_1(){
    var local_i int
    for j := 0; j<1000000; j++{
		local_i = <-availible
        local_i++
        availible <-local_i
	}
}

func thread_2(){
    var local_i int
    for j := 0; j<1000000; j++{
        local_i = <-availible
        local_i--
        availible <-local_i
    }
}


func main(){

    availible <- i
    runtime.GOMAXPROCS(runtime.NumCPU())    // I guess this is a hint to what GOMAXPROCS does...
                                            // Try doing the exercise both with and without it!
   
    go thread_1()                      // This spawns someGoroutine() as a goroutine
    go thread_2()



    // We have no way to wait for the completion of a goroutine (without additional syncronization of some sort)
    // We'll come back to using channels in Exercise 2. For now: Sleep.
    time.Sleep(100*time.Millisecond)
    Println("Inside main! This is i: ",i);
}