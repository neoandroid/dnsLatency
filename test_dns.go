package main

import (
	"flag"
	"log"
	"net"
	"math/rand"
	"sync"
	"time"
	"github.com/miekg/dns"
)

var wg sync.WaitGroup
var servers []string

func random(min, max int) int {
    rand.Seed(time.Now().Unix())
    return rand.Intn(max - min) + min
}


func checkDns(connections int) {
	defer wg.Done()
	for i := 0; i<connections ; i++ {
		var status string
		start := time.Now()
		ip, err := net.ResolveIPAddr("ip4", "messaging.schibsted.io")
		if err != nil {
			log.Println("Resolution error:", err)
			status = "Unreachable"
		} else {
			status = "Ok"
		}
		elapsed := time.Since(start)
		log.Println(ip, status + " Time:", elapsed)
	}
}

func checkDns2(connections int) {
	defer wg.Done()
	for i := 0; i<connections ; i++ {
		target := "messaging.schibsted.io"
		c := dns.Client{}
		m := dns.Msg{}
		m.SetQuestion(target+".", dns.TypeA)
		r, t, err := c.Exchange(&m, servers[random(0,3)]+":53")
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Took %v", t)
		if len(r.Answer) == 0 {
			log.Fatal("No results")
		}
		for _, ans := range r.Answer {
			Arecord := ans.(*dns.A)
			log.Printf("%s", Arecord.A)
		}
	}
}

func main() {
	concurrency := flag.Int("concurrency", 2, "goroutines")
	connections := flag.Int("connections", 1, "connections per goroutine")
	flag.Parse()
	log.Println("Concurrency:", *concurrency)
	log.Println("Connections:", *connections)

	wg.Add(*concurrency)
	servers = []string{ "205.251.197.250", "205.251.198.57", "205.251.195.51", "205.251.192.148"}

	for routine := 0 ; routine < *concurrency ; routine++ {
		go checkDns2(*connections)
	}

	wg.Wait()
}
