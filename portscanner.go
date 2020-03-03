package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/sync/semaphore"
)

type PortScanner struct {
	domain string
	ip     string
	lock   *semaphore.Weighted
}

var top_20_tcp_ports = map[int]string{
	21:   "ftp",
	22:   "ssh",
	23:   "telnet",
	25:   "smtp",
	53:   "domain",
	80:   "http",
	110:  "pop3",
	111:  "rpcbind",
	135:  "msrpc",
	139:  "netbios-ssn",
	143:  "imap",
	443:  "https",
	445:  "microsoft-ds",
	993:  "imaps",
	995:  "pop3s",
	1723: "pptp",
	3306: "mysql",
	3389: "ms-wbt-server",
	5900: "vnc",
	8080: "http-proxy",
}

func Ulimit() int64 {
	out, err := exec.Command("ulimit", "-n").Output()
	if err != nil {
		panic(err)
	}

	s := strings.TrimSpace(string(out))

	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		panic(err)
	}

	return i
}

func ScanPort(domain string, ip string, port int, timeout time.Duration) {
	target := fmt.Sprintf("%s:%d", ip, port)
	conn, err := net.DialTimeout("tcp", target, timeout)

	if err != nil {
		if strings.Contains(err.Error(), "too many open files") {
			time.Sleep(timeout)
			ScanPort(domain, ip, port, timeout)
		} else {
			//fmt.Println(port, "closed")
		}
		return
	}

	conn.Close()
	fmt.Println("\t", ip, "->", port, "open")

}

func (ps *PortScanner) Start(f, l int, timeout time.Duration) {
	wg := sync.WaitGroup{}
	defer wg.Wait()
	if f == 0 && l == 0 { // change this to outside func, its unnecessary check here and extra step
		for port := range top_20_tcp_ports {
			ps.lock.Acquire(context.TODO(), 1)
			wg.Add(1)
			go func(port int) {
				defer ps.lock.Release(1)
				defer wg.Done()
				ScanPort(ps.domain, ps.ip, port, timeout)
			}(port)
		}

	}

	for port := f; port <= l; port++ {
		ps.lock.Acquire(context.TODO(), 1)
		wg.Add(1)
		go func(port int) {
			defer ps.lock.Release(1)
			defer wg.Done()
			ScanPort(ps.domain, ps.ip, port, timeout)
		}(port)
	}

}

func main() {
	interactiveTest()
}

func resolveDomain(target string) net.IP {
	var target_ip net.IP
	addr, err := net.LookupIP(target)
	if err != nil {
		fmt.Println("Unknown host")
	} else {
		target_ip = addr[0] // returns ipv4
	}
	return target_ip
}

func getDomainsIPs(filename string) map[string]string {
	domain_ips := make(map[string]string)

	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		slice := strings.Split(line, " ")
		domain := slice[0]
		ip := slice[1]
		domain_ips[domain] = ip
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return domain_ips
}

func interactiveTest() {
	filename := "domains.txt"
	domains := getDomainsIPs(filename)
	var i int
	for dom, ip := range domains {
		fmt.Println(dom + ": " + strconv.Itoa(i))
		i = i + 1
		if ip != "" {
			ps := &PortScanner{
				domain: dom,
				ip:     ip,
				lock:   semaphore.NewWeighted(Ulimit()),
			}
			ps.Start(0, 0, 75*time.Millisecond) // approx 75 is enough to test all 20 ports

		}
	}

}
