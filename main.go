package main

import (
	"errors"
	"fmt"
	"github.com/akamensky/argparse"
	"github.com/go-playground/validator/v10"
	"log"
	"net"
	"net/netip"
	"os"
	"strings"
	"sync"
)

type ipHost struct {
	ipaddr  string
	dnsName []string
}

type options struct {
	resolverIP        string
	resolverPort      int
	useCustomResolver bool

	querys  []string
	timeout int //seconds
}

func main() {
	parser := argparse.NewParser("Jumpscan", "A quick, concurrent Reverse DNS scanner. Used PTR records to discover and map the domain. that can be dropped onto a victim machine and ran. NOTE: Custom Resolvers will not work on windows. Golang will not be able to change the default resolver.")

	hostarg := parser.String("t", "target", &argparse.Options{Required: true, Help: "IPv4 to target. Single, CIDR, comma seperated"})

	err := parser.Parse(os.Args)

	valHost, err := ParseHost(*hostarg)
	if err != nil {
		log.Fatal(err)
	}

	args := options{
		querys: valHost}
	startScan(&args)
}

func startScan(opt *options) {

	var ipArray []ipHost
	var wg sync.WaitGroup

	for _, host := range opt.querys {
		wg.Add(1)
		go func(query string) {
			ahost := ipHost{ipaddr: query}
			results, err := doLookup(query)
			if err == nil {

				ahost.dnsName = results
				ipArray = append(ipArray, ahost)
			}
			wg.Done()

		}(host)
	}
	wg.Wait()

	for r := range ipArray {
		if len(ipArray[r].dnsName) == 0 {
			continue
		}
		fmt.Printf("%s - ", ipArray[r].ipaddr)
		for dnsresult := range ipArray[r].dnsName {
			fmt.Printf("%s\n", ipArray[r].dnsName[dnsresult])
		}
	}

}

func doLookup(query string) ([]string, error) {
	ip, err := net.LookupAddr(query)
	return ip, err
}

func Cidr(cidr string) ([]string, error) {
	prefix, err := netip.ParsePrefix(cidr)
	if err != nil {
		return nil, err
	}

	var ips []netip.Addr
	for addr := prefix.Addr(); prefix.Contains(addr); addr = addr.Next() {
		ips = append(ips, addr)
	}

	var ipstring []string
	for ipClass := range ips {
		ipstring = append(ipstring, ips[ipClass].String())
	}

	return ipstring, nil
}

func ParseHost(host string) ([]string, error) {
	v := validator.New()
	var p []string

	//Single Host
	err := v.Var(host, "ip")
	if err == nil {
		p = append(p, host)
		return p, nil
	}

	//Cidr
	err = v.Var(host, "cidr")
	if err == nil {
		p, err = Cidr(host)
		if err != nil {
			return p, err
		}
		return p, nil
	}

	//Comma seperated
	if strings.Contains(host, ",") {
		potentialTarget := strings.Split(host, ",")
		for ahost := range potentialTarget {
			potentialTarget[ahost] = strings.TrimSpace(potentialTarget[ahost])
			err = v.Var(potentialTarget[ahost], "ip4_addr")
			if err != nil {
				continue
			}
			p = append(p, potentialTarget[ahost])

		}
		return p, err
	}
	return p, errors.New("could not parse IP address")

}
