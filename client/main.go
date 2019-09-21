package main

import (
	"flag"
	"fmt"
	"math"
	"net/http"
	"os"
	"regexp"
	"time"

	"github.com/eiannone/keyboard"
)

var requestsBySecond = 1

func main() {
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("No url provided")
		os.Exit(1)
	}

	if !validateURL(args[0]) {
		fmt.Println("Invalid url provided")
		os.Exit(1)
	}

	count := make(chan int)

	go process(args[0], count, &requestsBySecond)

	go measure(args[0], count)

	count <- 0

	err := keyboard.Open()
	if err != nil {
		panic(err)
	}
	defer keyboard.Close()

	for {
		char, key, err := keyboard.GetKey()
		if err != nil {
			panic(err)
		}

		if key == keyboard.KeyEsc {
			os.Exit(0)
		}

		switch char {
		case '+', '=':
			requestsBySecond++
			fmt.Printf("Speed up requests\n")
		case '-':
			requestsBySecond = int(math.Max(1, float64(requestsBySecond-1)))
			fmt.Printf("Slow down requests\n")
		}
	}
}

func process(url string, count chan int, requestsBySecond *int) {
	for {
		go order(url, count)

		time.Sleep(time.Duration(1000000 / *requestsBySecond) * time.Microsecond)
	}
}

func measure(url string, count chan int) {
	fmt.Printf("\033[H\033[2J 0 requests/second to %s\n\n", url)

	for {
		time.Sleep(1 * time.Second)

		fmt.Printf("\033[H\033[2J %d requests/second to %s\n\n", <-count, url)
		count <- 0
	}
}

func validateURL(url string) bool {
	r := regexp.MustCompile("^https?://")

	return r.MatchString(url)
}

func order(url string, count chan int) {
	v := <-count

	count <- v + 1

	resp, err := http.Get(url)

	if err != nil {
		fmt.Println(err)
		return
	}

	resp.Body.Close()
}
