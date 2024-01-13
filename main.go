package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

func generator(done <-chan int, fn func() int) <-chan int {
	stream := make(chan int)

	go func() {
		defer close(stream)

		for {
			select {
			case <-done:
				return
			case stream <- fn():
			}
		}
	}()

	return stream
}

func take(done <-chan int, inStream <-chan int, N int) <-chan int {
	outStream := make(chan int)

	go func() {
		defer close(outStream)

		for i := 1; i <= N; i++ {
			select {
			case <-done:
				return
			case outStream <- <-inStream:
			}
		}
	}()

	return outStream
}

func primeFinder(done <-chan int, inStream <-chan int) <-chan int {
	primes := make(chan int)

	isPrime := func(n int) bool {
		for i := 2; i <= n-1; i++ {
			if n%i == 0 {
				return false
			}
		}
		return true
	}

	go func() {
		defer close(primes)

		for {
			select {
			case <-done:
				return
			case randInt := <-inStream:
				if isPrime(randInt) {
					primes <- randInt
				}
			}
		}
	}()

	return primes
}

func fanIn(done <-chan int, primeChannels ...<-chan int) <-chan int {
	stream := make(chan int)
	var wg sync.WaitGroup

	transfer := func(c <-chan int) {
		defer wg.Done()

		for {
			select {
			case <-done:
				return
			case stream <- <-c:
			}
		}
	}

	for _, c := range primeChannels {
		wg.Add(1)
		go transfer(c)
	}

	go func() {
		wg.Wait()
		close(stream)
	}()

	return stream
}

func main() {
	start := time.Now()
	done := make(chan int)

	fn := func() int {
		return rand.Intn(1000000000)
	}

	randomStream := generator(done, fn)

	// !NAIVE
	// primeStream := primeFinder(done, randomStream)

	// outStream := take(done, primeStream, 5)

	// for n := range outStream {
	// 	fmt.Printf("%d   ", n)
	// }

	// ? EFFICIENT
	noOfCpus := runtime.NumCPU()

	primeChannels := make([]<-chan int, noOfCpus)

	for i := 0; i < noOfCpus; i++ {
		primeChannels[i] = primeFinder(done, randomStream)
	}

	fannedInStream := fanIn(done, primeChannels...)

	outStream := take(done, fannedInStream, 5)

	for n := range outStream {
		fmt.Printf("%d   ", n)
	}

	fmt.Println(time.Since(start))
}
