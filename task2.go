package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

func main() {
	// 指针 题目1
	num := 20
	increaseTen(&num)
	fmt.Println("Modified value:", num)

	//指针 题目2
	nums := []int{1, 2, 3, 4, 5}
	doubleSlice(&nums)
	fmt.Println(nums)

	//Goroutine 题目1
	go func() {
		for i := 1; i <= 10; i += 2 {
			fmt.Println("Odd:", i)
		}
	}()

	go func() {
		for i := 2; i <= 10; i += 2 {
			fmt.Println("Even:", i)
		}
	}()
	time.Sleep(100 * time.Millisecond)

	// Goroutine 题目2
	tasks := []Task{
		{
			Name: "task1",
			Content: func() {
				time.Sleep(3000 * time.Millisecond)
			},
		},
		{
			Name: "task2",
			Content: func() {
				time.Sleep(2000 * time.Millisecond)
			},
		},
		{
			Name: "task3",
			Content: func() {
				time.Sleep(3000 * time.Millisecond)
			},
		},
		{
			Name: "task4",
			Content: func() {
				time.Sleep(2000 * time.Millisecond)
			},
		},
		{
			Name: "task5",
			Content: func() {
				time.Sleep(1000 * time.Millisecond)
			},
		},
	}
	runTask(tasks)
	time.Sleep(5000 * time.Millisecond)

	//面向对象 题目1
	r := Rectangle{}
	r.Area()
	r.Perimeter()

	c := Circle{}
	c.Area()
	c.Perimeter()
	// 面向对象 题目2
	e := Employee{EmployeeID: 1, Person: Person{Name: "Tom", Age: 18}}
	e.printInfo()

	//Channel 题目1
	ch := make(chan int)
	go sendNumbers(ch)
	go receiveNumbers(ch)
	time.Sleep(100 * time.Millisecond)
	//Channel 题目2
	ch1 := make(chan int, 20)
	go produce(ch1)
	time.Sleep(1 * time.Second)
	go func() {
		for {
			select {
			case v, ok := <-ch1:
				if !ok {
					fmt.Println("Channel 通道已关闭")
					return
				}
				fmt.Println("go1接收到数据:", v)
			default:
				fmt.Println("g1没有数据，等待中...")
				return
			}
		}
	}()
	go func() {
		for {
			select {
			case v, ok := <-ch1:
				if !ok {
					fmt.Println("Channel 通道已关闭")
					return
				}
				fmt.Println("go2接收到数据：", v)
			default:
				fmt.Println("go2没有数据，等待中...")
				time.Sleep(1 * time.Second)
				return
			}
		}
	}()
	time.Sleep(1 * time.Second)

	//锁机制题目1
	unsafeCounter := SafeCounter{}
	for i := 0; i < 1000; i++ {
		go func() {
			for j := 0; j < 10; j++ {
				unsafeCounter.addCount()
			}
		}()
	}
	time.Sleep(time.Second)
	fmt.Println("最终的计数为：", unsafeCounter.count)

	//锁机制题目2
	counter := Counter{}
	for i := 0; i < 1000; i++ {
		go func() {
			for j := 0; j < 10; j++ {
				counter.atomicAddCount()
			}
		}()
	}
	time.Sleep(time.Second)
	fmt.Println("原子操作最终的计数为：", counter.count)
}

func increaseTen(num *int) {
	*num += 10
}

func doubleSlice(slice *[]int) {
	for i := 0; i < len(*slice); i++ {
		(*slice)[i] *= 2
	}
}

type Task struct {
	Name    string
	Content func()
}

func runTask(tasks []Task) {
	for index, task := range tasks {
		go func(idx int, t Task) {
			start := time.Now()
			t.Content()
			during := time.Since(start)
			fmt.Println(t.Name, "执行完毕, during: ", during)
		}(index, task)
	}
}

type Shape interface {
	Area()
	Perimeter()
}

type Rectangle struct {
}

type Circle struct {
}

func (r *Rectangle) Area() {
	fmt.Println("Rectangle Area method")
}

func (r *Rectangle) Perimeter() {
	fmt.Println("Rectangle Perimeter method")
}

func (r *Circle) Area() {
	fmt.Println("Circle Area method")
}

func (r *Circle) Perimeter() {
	fmt.Println("Circle Perimeter method")
}

type Person struct {
	Name string
	Age  int
}

type Employee struct {
	EmployeeID int
	Person
}

func (e *Employee) printInfo() {
	employeeId := e.EmployeeID
	employeeName := e.Person.Name
	employeeAge := e.Person.Age
	fmt.Printf("员工ID为：%d, 员工Name为：%s, 员工Age为：%d", employeeId, employeeName, employeeAge)
}

func sendNumbers(ch chan<- int) {
	for i := 1; i <= 10; i++ {
		ch <- i
	}
	close(ch)
}

func receiveNumbers(ch <-chan int) {
	for v := range ch {
		fmt.Println(v)
	}
}

func produce(ch chan<- int) {
	for i := 1; i <= 100; i++ {
		ch <- i
	}
	close(ch)
}

type SafeCounter struct {
	mu    sync.Mutex
	count int
}

func (c *SafeCounter) addCount() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.count++
}

type Counter struct {
	count int64
}

func (c *Counter) atomicAddCount() {
	atomic.AddInt64(&c.count, 1)
}
