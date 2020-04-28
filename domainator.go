/**
 * Domain Name Super Fast Searcher
 */

package main

import (
	"bufio"
	// "context"
	"flag"
	"fmt"
	"log"
	// "net"
	"os"
	"sync"
	"regexp"
	"strings"
	"time"
	// "container/list"
	"github.com/miekg/dns"
)

type Name_Server struct {
	pri int
	tld string
	host string
	ipv4 string
}

// var rns_list = make([]Name_Server, 8192)
var rns_list = []Name_Server{}
var tld_list = map[string]string{}

func load_root_zone() {

	file := "root.zone"
	fh, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}
	defer fh.Close()

	var idx = 0
	var line string;
	scanner := bufio.NewScanner(fh)

	reNS := regexp.MustCompile(`^([\w\-\.]+)\s+\d+\s+IN\s+NS\s+(.+)$`) // NS Line
	reIP := regexp.MustCompile(`^([\w+\-\.]+)\s+\d+\s+IN\s+A\s+(.+)$`) // NS Line

	ip4_list := map[string]string{}

	// Scan for next token.
	for scanner.Scan() {

		line = scanner.Text()

		if (reNS.MatchString(line)) {

			// fmt.Println("NS:" + line)

			m1 := reNS.FindAllStringSubmatch(line, -1)

			nsX := Name_Server{
				pri: idx,
				tld: m1[0][1],
				host: m1[0][2],
			}

			rns_list = append(rns_list, nsX)
			tld_list[ nsX.tld ] = nsX.tld;

		}

		if (reIP.MatchString(line)) {

			mIP := reIP.FindAllStringSubmatch(line, -1)

			host := mIP[0][1]
			ipv4 := mIP[0][2]

			ip4_list[host] = ipv4;

		}

	}

	// Find something in this existing list to add the ipv4 information
	for i, rns := range rns_list {
		if (0 == len(rns.ipv4)) {
			// fmt.Println(host + " == " + rns.host + " == " + ipv4)
			rns_list[i].ipv4 = ip4_list[ rns.host ]
		}
	}


}


func main() {

	var name_base string
	var wg sync.WaitGroup

	flag.StringVar(&name_base, "name", "", "Domain Name to Search For")
	flag.Parse()

	if 0 == len(name_base) {
		panic("No Name was provided")
	}

	fmt.Printf("Searching for name: '%s'\n", name_base)

	load_root_zone()
	// fmt.Println(rns_list)

	for _, tld := range tld_list {
		for _, rns := range rns_list {
			if (tld == rns.tld) {
				wg.Add(1)
				go find_ns(&wg, rns.ipv4, rns.tld, name_base)
				time.Sleep(10 * time.Millisecond)
				break
			}
		}
	}

	wg.Wait()

}

func find_ns(wg *sync.WaitGroup, rns string, tld string, dom string) {

	defer wg.Done()

	var arg_note strings.Builder
	arg_note.Reset()
	arg_note.WriteString( "rns:" )
	arg_note.WriteString( rns )
	arg_note.WriteString( "; " )
	arg_note.WriteString( dom )
	arg_note.WriteString( "." )
	arg_note.WriteString( tld )
	arg_note.WriteString( "\n" )
	// fmt.Print( arg_note.String() )
	// return

	var dom_full strings.Builder
	var res_note strings.Builder
	
	c := dns.Client{}
	c.Net = "udp" // or udp4 or udp6 or tcp or tcp4 or tcp6

	dom_full.Reset()
	dom_full.WriteString( dom )
	dom_full.WriteString( "." )
	dom_full.WriteString( tld )

	// fmt.Print(name_full.String())
	res_note.Reset()
	// res_note.WriteString( fmt.Sprintf("%04d", idx) )
	// res_note.WriteString(" ")
	res_note.WriteString( dom_full.String() )

	m := dns.Msg{}
	m.SetQuestion( dom_full.String(), dns.TypeNS)
	m.RecursionDesired = true

	res, _, err := c.Exchange(&m, rns + ":53")

	if err != nil {
		res_note.WriteString(" ERROR ")
		res_note.WriteString(err.Error())
		fmt.Println( res_note.String() )
		return
	}

	if res.Rcode != dns.RcodeSuccess {
		res_note.WriteString(" ERROR ")
		res_note.WriteString(" RCODE ")
		fmt.Println( res_note.String() )
		return
	}

	if 0 == len(res.Answer) {
		res_note.WriteString(" NX_DOMAIN")
		fmt.Println( res_note.String() )
		return
	}

	for _, rec := range res.Answer {

		switch rec.(type) {
		case *dns.NS:
			r := rec.(*dns.NS)
			res_note.WriteString( " ns:" )
			res_note.WriteString( r.Ns )
		case *dns.CNAME:
			r := rec.(*dns.CNAME)
			res_note.WriteString( "  cname:" )
			res_note.WriteString( r.Target )
		default:
			res_note.WriteString(" UNKNOWN_TYPE ")
		}
	}

	fmt.Println( res_note.String() )

}
