// This sample program demonstrates how to create goroutines and
// how the goroutine scheduler behaves with three logical processors.
 
package main
 
import (
    "fmt"
    "runtime"
    "sync"
)
 
func main() {
    // Allocate three logical processors for the scheduler to use.
    runtime.GOMAXPROCS(3)
 
    // processTest is used to wait for the program to finish.
    var processTest sync.WaitGroup
    // Add a count of three, one for each goroutine.
    processTest.Add(3)
     
    // Declaration of three anonymous function and create a goroutine.
    go func() {
        defer processTest.Done()
        for i := 0; i < 30; i++ {
			fmt.Println("proc a: 51<num<100 round: ", i)
			fmt.Printf("proc a: ")
            for j := 51; j <= 100; j++ {
                fmt.Printf(" %d", j)
                if j == 100{
                    fmt.Println()
                }
            }
        }
    }()
    go func() {
        defer processTest.Done()
        for j := 0; j < 10; j++ {
			fmt.Println("proc b: letters round: ", j)
			fmt.Printf("proc b: ")
            for char := 'A'; char < 'A'+26; char++ {
                fmt.Printf("%c ", char)
                if char == 'Z' {
                    fmt.Println()
                }
 
            }
        }
    }()
    go func() {
        defer processTest.Done()
        for i := 0; i < 30; i++ {
			fmt.Println("proc c: 0<num<50 round: ", i)

			fmt.Printf("proc c: ")
            for j := 0; j <= 50; j++ {
                fmt.Printf(" %d", j)
                if j == 50 {
                    fmt.Println()
                }
            }
        }
    }()
 
    // Wait for the goroutines to finish.
    processTest.Wait()  
}
