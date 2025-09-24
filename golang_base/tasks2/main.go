package main

import (
	"fmt"
	"math"
	"sync"
	"sync/atomic"
	"time"
)

// ✅指针
// 题目 ：编写一个Go程序，定义一个函数，该函数接收一个整数指针作为参数，在函数内部将该指针指向的值增加10，然后在主函数中调用该函数并输出修改后的值。
// 考察点 ：指针的使用、值传递与引用传递的区别。
// 题目 ：实现一个函数，接收一个整数切片的指针，将切片中的每个元素乘以2。
// 考察点 ：指针运算、切片操作。
// ✅Goroutine
// 题目 ：编写一个程序，使用 go 关键字启动两个协程，一个协程打印从1到10的奇数，另一个协程打印从2到10的偶数。
// 考察点 ： go 关键字的使用、协程的并发执行。
// 题目 ：设计一个任务调度器，接收一组任务（可以用函数表示），并使用协程并发执行这些任务，同时统计每个任务的执行时间。
// 考察点 ：协程原理、并发任务调度。
// ✅面向对象
// 题目 ：定义一个 Shape 接口，包含 Area() 和 Perimeter() 两个方法。然后创建 Rectangle 和 Circle 结构体，实现 Shape 接口。在主函数中，创建这两个结构体的实例，并调用它们的 Area() 和 Perimeter() 方法。
// 考察点 ：接口的定义与实现、面向对象编程风格。
// 题目 ：使用组合的方式创建一个 Person 结构体，包含 Name 和 Age 字段，再创建一个 Employee 结构体，组合 Person 结构体并添加 EmployeeID 字段。为 Employee 结构体实现一个 PrintInfo() 方法，输出员工的信息。
// 考察点 ：组合的使用、方法接收者。
// ✅Channel
// 题目 ：编写一个程序，使用通道实现两个协程之间的通信。一个协程生成从1到10的整数，并将这些整数发送到通道中，另一个协程从通道中接收这些整数并打印出来。
// 考察点 ：通道的基本使用、协程间通信。
// 题目 ：实现一个带有缓冲的通道，生产者协程向通道中发送100个整数，消费者协程从通道中接收这些整数并打印。
// 考察点 ：通道的缓冲机制。
// ✅锁机制
// 题目 ：编写一个程序，使用 sync.Mutex 来保护一个共享的计数器。启动10个协程，每个协程对计数器进行1000次递增操作，最后输出计数器的值。
// 考察点 ： sync.Mutex 的使用、并发数据安全。
// 题目 ：使用原子操作（ sync/atomic 包）实现一个无锁的计数器。启动10个协程，每个协程对计数器进行1000次递增操作，最后输出计数器的值。
// 考察点 ：原子操作、并发数据安全。
func HandlerArrayPointer(num *int) {
	fmt.Println("HandlerArrayPointer before num:", *num)
	*num += 10
	fmt.Println("HandlerArrayPointer after num:", *num)
}

func AnswerPointer() {
	fmt.Println(" < ---------- AnswerPointer1() ------------ >")
	// 题目 ：编写一个Go程序，定义一个函数，该函数接收一个整数指针作为参数，在函数内部将该指针指向的值增加10，然后在主函数中调用该函数并输出修改后的值。
	// 考察点 ：指针的使用、值传递与引用传递的区别。
	num := 5
	fmt.Println("AnswerPointer before num:", num)
	HandlerArrayPointer(&num)
	fmt.Println("AnswerPointer after num:", num)
}

func HandlerArrayPointer2(nums *[]int) {
	fmt.Println("HandlerArrayPointer2 before num:", *nums)
	for i, v := range *nums {
		(*nums)[i] = v * 2
	}
	fmt.Println("HandlerArrayPointer2 after num:", *nums)

}
func AnswerPointer2() {
	fmt.Println(" < ---------- AnswerPointer2() ------------ >")
	// 题目 ：实现一个函数，接收一个整数切片的指针，将切片中的每个元素乘以2。
	// 考察点 ：指针运算、切片操作。
	num := []int{1, 2, 3}
	fmt.Println("AnswerPointer2 before num:", num)
	HandlerArrayPointer2(&num)
	fmt.Println("AnswerPointer2 after num:", num)
}

// ✅Goroutine
func ScheduleTasks(tasks []func()) {
	var wg sync.WaitGroup
	wg.Add(len(tasks))
	// 为每个任务启动一个协程
	for i, task := range tasks {
		go func(index int, task func()) {
			defer wg.Done()

			// 记录开始时间
			start := time.Now()

			// 执行任务
			task()

			// 记录结束时间并计算耗时
			end := time.Now()
			duration := end.Sub(start)

			fmt.Printf("任务 %d 执行完毕，耗时：%v\n", index+1, duration)
		}(i, task) // 传递索引和任务函数，避免闭包问题
	}

	// 等待所有任务完成
	wg.Wait()
	fmt.Println("所有任务执行完成")
}

// 示例任务函数
func task1() {
	time.Sleep(2 * time.Second)
	fmt.Println("任务1完成")
}

func task2() {
	time.Sleep(1 * time.Second)
	fmt.Println("任务2完成")
}

func task3() {
	time.Sleep(3 * time.Second)
	fmt.Println("任务3完成")
}
func AswerGoroutine() {
	// 题目 ：设计一个任务调度器，接收一组任务（可以用函数表示），并使用协程并发执行这些任务，同时统计每个任务的执行时间。
	// 考察点 ：协程原理、并发任务调度。
	fmt.Println(" < ---------- AnswerGoroutine ---------- >")
	// 调用任务调度器
	tasks := []func(){
		task1,
		task2,
		task3,
	}
	ScheduleTasks(tasks)
}

func AswerGoroutine2() {

	// 题目 ：编写一个程序，使用 go 关键字启动两个协程，一个协程打印从1到10的奇数，另一个协程打印从2到10的偶数。
	// 考察点 ： go 关键字的使用、协程的并发执行。
	fmt.Println(" < ---------- AnswerGoroutine2 ---------- >")
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 1; i <= 10; i++ {
			if i%2 != 0 {
				fmt.Println("奇数 goroutine:", i)
				time.Sleep(time.Millisecond * 100)
			}

		}
	}()
	go func() {
		defer wg.Done()
		for i := 1; i <= 10; i++ {
			if i%2 == 0 {
				fmt.Println("偶数 goroutine:", i)
				time.Sleep(time.Millisecond * 100)
			}

		}
	}()
	wg.Wait()
	fmt.Println("所有协程执行完成")
}

// ✅面向对象

type Shape interface {
	Area() float64
	Perimeter() float64
}

type Rectangle struct {
	width  float64
	height float64
}

func (r *Rectangle) Area() float64 {
	// 矩形的面积
	return r.width * r.height
}

func (r *Rectangle) Perimeter() float64 {
	// 矩形的周长
	return 2 * (r.width + r.height)
}

type Circle struct {
	// 圆形的半径
	radius float64
}

func (c *Circle) Area() float64 {
	// 圆的面积
	fmt.Printf("r = %f 计算圆的面积 = π * r * r = %f\n", c.radius, math.Pi*c.radius*c.radius)
	return math.Pi * c.radius * c.radius
}
func (c *Circle) Perimeter() float64 {
	// 圆的周长
	fmt.Printf("r = %f 圆的周长公式 = 2 * π * r = %f\n", c.radius, 2*math.Pi*c.radius)
	return 2 * math.Pi * c.radius
}

func AswerObjectOriented() {
	// 题目 ：定义一个 Shape 接口，包含 Area() 和 Perimeter() 两个方法。然后创建 Rectangle 和 Circle 结构体，实现 Shape 接口。
	// 在主函数中，创建这两个结构体的实例，并调用它们的 Area() 和 Perimeter() 方法。
	// 考察点 ：接口的定义与实现、面向对象编程风格。
	fmt.Println(" < ---------- AnswerObjectOriented() ---------- >")
	RectangleShape := &Rectangle{height: 5, width: 3}
	fmt.Println("Rectangle Area:", RectangleShape.Area())
	fmt.Println("Rectangle Perimeter:", RectangleShape.Perimeter())
	fmt.Println(" --------------------------------------------- ")
	CircleShape := &Circle{radius: 4}
	fmt.Println("Circle Area:", CircleShape.Area())
	fmt.Println("Circle Perimeter:", CircleShape.Perimeter())

}

type Person struct {
	Name string
	Age  int
}

type Employee struct {
	Person
	EmployeeID int
}

func (p *Employee) PrintInfo() {
	fmt.Printf("Name: %s, Age: %d, EmployeeID: %d\n", p.Name, p.Age, p.EmployeeID)
}

func AswerObjectInfo() {
	// 题目 ：使用组合的方式创建一个 Person 结构体，包含 Name 和 Age 字段，再创建一个 Employee 结构体，组合 Person 结构体并添加 EmployeeID 字段。为 Employee 结构体实现一个 PrintInfo() 方法，输出员工的信息。
	// 考察点 ：组合的使用、方法接收者。
	fmt.Println(" < ---------- AnswerObjectInfo ---------- >")
	emp := &Employee{
		EmployeeID: 1001,
		Person: Person{
			Name: "张三",
			Age:  18,
		},
	}
	emp.PrintInfo()
}

// ✅Channel

func sendOnly(ch chan<- int, num int) {
	for i := 1; i <= num; i++ {
		ch <- i
		fmt.Printf("发送: %d\n", i)
		// time.Sleep(time.Millisecond * 1)
	}
	close(ch)
}
func receiveOnly(ch <-chan int) {
	for v := range ch {
		fmt.Printf("接收到: %d\n", v)
	}
}
func AswerOutputChannelInfo() {
	// 题目 ：编写一个程序，使用通道实现两个协程之间的通信。一个协程生成从1到10的整数，并将这些整数发送到通道中，另一个协程从通道中接收这些整数并打印出来。
	// 考察点 ：通道的基本使用、协程间通信。
	fmt.Println(" < ---------- AnswerOutputChannelInfo() ---------- >")
	ch := make(chan int, 10)
	go sendOnly(ch, 10)
	go receiveOnly(ch)
	time.Sleep(time.Second * 1)

}
func AswerOutputChannelInfo2() {
	// 题目 ：实现一个带有缓冲的通道，生产者协程向通道中发送100个整数，消费者协程从通道中接收这些整数并打印。
	// 考察点 ：通道的缓冲机制。
	fmt.Println(" < ---------- AnswerOutputChannelInfo2() ---------- >")

	// 创建一个缓冲通道，缓冲区大小为10（可以根据需要调整）
	ch := make(chan int, 10)

	// 使用 WaitGroup 等待所有协程完成
	var wg sync.WaitGroup
	wg.Add(2) // 一个生产者协程，一个消费者协程

	// 生产者协程：发送100个整数
	go func() {
		defer wg.Done()
		sendOnly(ch, 100)
		fmt.Println("生产者完成发送")
	}()

	// 消费者协程：接收并打印整数
	go func() {
		defer wg.Done()
		receiveOnly(ch)
		fmt.Println("消费者完成接收")
	}()

	// 等待所有协程完成
	wg.Wait()
	fmt.Println("生产者-消费者模式执行完成")

}

// ✅锁机制
type SafeCounter struct {
	mu    sync.Mutex
	count int
}

func (c *SafeCounter) Increment() {
	c.mu.Lock()
	c.count++
	c.mu.Unlock()
}

func (c *SafeCounter) GetCount() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.count
}

func AswerSafeIncrement() {
	// 题目 ：编写一个程序，使用 sync.Mutex 来保护一个共享的计数器。启动10个协程，每个协程对计数器进行1000次递增操作，最后输出计数器的值。
	// 考察点 ： sync.Mutex 的使用、并发数据安全。
	fmt.Println(" < ---------- AnswerSafeIncrement ---------- >")
	safeCounter := SafeCounter{}
	var wg sync.WaitGroup
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < 1000; j++ {
				safeCounter.Increment()
			}
		}()
	}
	wg.Wait()
	fmt.Printf("safeCounter Final count: %d\n", safeCounter.GetCount())

}

func AswerUnlockCounter() {
	// 题目 ：使用原子操作（ sync/atomic 包）实现一个无锁的计数器。启动10个协程，每个协程对计数器进行1000次递增操作，最后输出计数器的值。
	// 考察点 ：原子操作、并发数据安全。
	fmt.Println(" < ---------- AswerUnlockCounter ---------- >")
	var counter int64 // 使用int64类型，因为atomic包主要操作int64
	var wg sync.WaitGroup

	// 启动10个协程
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func(goroutineID int) {
			defer wg.Done()
			// 每个协程进行1000次递增操作
			for j := 0; j < 1000; j++ {
				atomic.AddInt64(&counter, 1) // 原子递增操作
			}
			fmt.Printf("Goroutine %d 完成计数\n", goroutineID)
		}(i) // 传递i值避免闭包问题；就是这里 i 传值给goroutineID
	}

	// 等待所有协程完成
	wg.Wait()

	// 输出最终计数器值
	fmt.Printf("原子操作最终计数: %d\n", counter)

	// 验证结果是否正确（应该是10 * 1000 = 10000）
	if counter == 10000 {
		fmt.Println("✅ 计数正确：10个协程 × 1000次递增 = 10000")
	} else {
		fmt.Printf("❌ 计数错误：期望10000，实际%d\n", counter)
	}

}
func main() {
	AnswerPointer()
	AnswerPointer2()
	AswerGoroutine()
	AswerGoroutine2()
	AswerObjectOriented()
	AswerObjectInfo()
	AswerOutputChannelInfo()
	AswerOutputChannelInfo2()
	AswerSafeIncrement()
	AswerUnlockCounter()
}
