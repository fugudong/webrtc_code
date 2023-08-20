package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func main2() {
	/*
	例子中我们通过select判断stop是否接受到值，如果接受到值就表示可以推出停止了，
	如果没有接受到，就会执行default里面的监控逻辑，继续监控，直到收到stop的通知

	以上控制goroutine的方式在大多数情况下可以满足我们的使用，但是也存在很多局限性，
	比如有很多goroutiine，并且这些goroutine还衍生了其他goroutine，此时chan就比较困难解决这样的问题了
	*/
	stop := make(chan bool)

	go func() {
		for {
			select {
			case <-stop:
				fmt.Println("监控退出，停止了...")
				return
			default:
				fmt.Println("goroutine监控中...")
				time.Sleep(2 * time.Second)
			}
		}
	}()

	time.Sleep(10 * time.Second)
	fmt.Println("可以了，通知监控停止")
	stop<- true
	//为了检测监控过是否停止，如果没有监控输出，就表示停止了
	time.Sleep(5 * time.Second)
}

func main3() {
	//以上例子一定要等到两个goroutine同时做完才会全部完成，
	// 这种控制并发方式尤其适用于多个goroutine协同做一件事情的时候。
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		fmt.Println("1号开始")
		time.Sleep(2*time.Second)
		fmt.Println("1号完成")
		wg.Done()
	}()
	go func() {
		fmt.Println("2号开始")
		time.Sleep(2*time.Second)
		fmt.Println("2号完成")
		wg.Done()
	}()
	wg.Wait()
	fmt.Println("好了，大家都干完了，放工")
}

func main4() {
	ctx, cancel := context.WithTimeout(context.Background(),5*time.Second)
	go watch(ctx,"【监控1】")
	go watch(ctx,"【监控2】")
	go watch(ctx,"【监控3】")
	time.Sleep(10 * time.Second)
	fmt.Println("可以了，通知监控停止")
	cancel()
	//为了检测监控过是否停止，如果没有监控输出，就表示停止了
	time.Sleep(5 * time.Second)
}
func watch(ctx context.Context, name string) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println(name,"监控退出，停止了...,原因："+ ctx.Err().Error())
			return
		default:
			deadTime,ok := ctx.Deadline()
			ok = ok
			fmt.Println(name,"goroutine监控中..." + deadTime.Format("2006-01-02 15:04:05"))
			time.Sleep(2 * time.Second)
		}
	}
}
func main() {
	main4()
}