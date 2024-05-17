package main

import (
	"bufio"
	"fmt"
	"net"
	"sync"
)

func dial(wait *sync.WaitGroup) {
	conn, err := net.Dial("tcp", "127.0.0.1:8000")
	if err != nil {
		fmt.Println("err conncting: ", err)
		return
	}

	defer func() {
		conn.Close()
		wait.Done()
	}()

	scanner := bufio.NewScanner(conn)
	scanner.Scan()
	err = scanner.Err()
	if err != nil {
		fmt.Println("err scanning: ", err)
	}
	fmt.Println(scanner.Text())
}

func main() {
	var w sync.WaitGroup
	done := make(chan struct{})

	listner, err := net.Listen("tcp", "127.0.0.1:8000")
	if err != nil {
		fmt.Println("Err listening to connection: ", err)
		return
	}

	var numberOfConnections int
	go func() {
		defer func() {
			done <- struct{}{} // write to the channel
		}()

		for {
			conn, err := listner.Accept()
			if err != nil {
				return
			}

			w.Add(1)
			go func(c net.Conn, connNum int) {
				defer func() {
					c.Close()
					w.Done() // decrementing
				}()

				_, err := c.Write([]byte(fmt.Sprintf("writing to connection number %d\n", connNum)))
				if err != nil {
					fmt.Println("Err writing: ", err)
					return
				}
			}(conn, numberOfConnections)

			numberOfConnections++
		}
	}()

	for i := 0; i < 100; i++ {
		w.Add(1)
		go dial(&w)
	}

	w.Wait()
	listner.Close() // closing the listener
	<-done          // wait for accept to return
}
