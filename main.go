package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"reflect"
	"regexp"
	"strings"
	"sync"

	"github.com/gocolly/colly/v2"
	"github.com/1121170088/find-domain/search"
)
type node struct{
	end bool
	folw map[byte] *node
}
var (
	Tree map[byte] *node = make(map[byte] *node)
	MaxDepth int
	Parallelism int
	Home string
	HomeDomain string
	HomePrefix string
	mutex sync.Mutex
	domaMu sync.Mutex
	visted map[string] bool = map[string]bool{}
	domainFile string
	domainSuffixFile string
)
func searchDomain(domain string) string  {
	//doubleSlah := strings.Index(domain, "//")
	//if doubleSlah != -1 {
	//	domain = domain[doubleSlah + 2:]
	//}
	shortDomain := search.Search(domain)

	if isDomain(shortDomain) {
		return shortDomain
	}
	return ""
}
func isDomain(domain string) bool  {
	match1, _ := regexp.Match(`^[A-Za-z0-9-.]{1,63}$`, []byte(domain))
	match2, _ := regexp.Match(`[A-Za-z0-9-.]{1,63}\.[A-Za-z0-9-.]{1,63}`, []byte(domain))
	return  match1 && match2 && []byte(domain)[0] != '-'

}
func reverse(s interface{}) {
	n := reflect.ValueOf(s).Len()
	swap := reflect.Swapper(s)
	for i, j := 0, n-1; i < j; i, j = i+1, j-1 {
		swap(i, j)
	}
}
func hasDomain(domain string) (has bool) {
	bytes := []byte(domain)
	reverse(bytes)
	var preNode *node = nil
	var ok bool = false
	has = true
	for _, b := range bytes {
		if preNode == nil {
			preNode, ok = Tree[b]
			if !ok {
				preNode = &node{
					end:  false,
					folw: make(map[byte] *node),
				}
				Tree[b] = preNode
				has = false
			}
		} else {
			preNode2, ok := preNode.folw[b]
			if !ok {
				preNode2 = &node{
					end:  false,
					folw: make(map[byte] *node),
				}
				preNode.folw[b]= preNode2
				has = false
			}
			preNode = preNode2
		}

	}
	if preNode != nil {
		if !preNode.end {
			preNode.end = true
			has = false
		}
	}
	return
}
func init() {
	flag.StringVar(&Home, "h", "", "visit it")
	flag.IntVar(&MaxDepth, "d", 1, "max depth")
	flag.IntVar(&Parallelism, "p", 2, "Parallelism")
	flag.StringVar(&domainFile, "dl", "", "domain file")
	flag.StringVar(&domainSuffixFile, "dsf", "", "domain suffix file")
	flag.Parse()
}

func prefix(url string) string {
	len := len(url)
	if len == 0 {
		return ""
	}
	slash1 := strings.Index(url, "/")
	if slash1 == -1 {
		return ""
	}
	from := slash1 + 1
	slash2 := strings.Index(url[from:], "/")
	if slash2 == -1 {
		return ""
	}
	from = from + slash2 + 1
	slash3 := strings.Index(url[from:], "/")
	quest := strings.Index(url[from:], "?")
	if slash3 == -1 && quest == -1 {
		return url
	}
	if slash3 == -1 {
		to := from + quest
		return url[:to]
	}
	if quest == -1 {
		to := from + slash3
		return url[:to]
	}
	min := slash3
	if min > quest {
		min = quest
	}
	to := from + min
	return url[:to]

}
func OpenFile(filename string) (*os.File, error) {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return os.Create(filename)
	}
	return os.OpenFile(filename, os.O_APPEND|os.O_RDWR, os.ModePerm)
}

func main() {
	search.Init(domainSuffixFile)
	if domainFile != "" {
		domainFilePtr, err := OpenFile(domainFile)
		if err != nil {
			log.Fatal(err)
		}
		rd := bufio.NewReader(domainFilePtr)

		for {
			line, err := rd.ReadString('\n')
			if err != nil || err == io.EOF {
				break
			}
			line = strings.Trim(line, "\n")
			line = strings.Trim(line, "\r")
			line = strings.Trim(line, " ")
			hasDomain(line)
		}
		domainFilePtr.Close()
	}

	HomePrefix = prefix(Home)
	HomeDomain = searchDomain(HomePrefix)
	if HomeDomain == "" {
		log.Panic("home domain is nil")
	}
	// Instantiate default collector
	c := colly.NewCollector(
		// Visit only domains: hackerspaces.org, wiki.hackerspaces.org
		colly.MaxDepth(MaxDepth),
		colly.Async(),
		//colly.AllowedDomains("hackerspaces.org", "wiki.hackerspaces.org"),
		colly.DisallowedDomains("www.beian.gov.cn", "beian.miit.gov.cn"),
	)
	c.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: Parallelism})
	// On every a element which has href attribute call callback
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		absUrl := e.Request.AbsoluteURL(link)
		if e.Request.Depth > MaxDepth {
			return
		}
		// Print link
		//fmt.Printf("Link found: %q -> %s\n", e.Text, absUrl)
		if absUrl == "" {
			return
		}
		if strings.Index(absUrl, "javascript") == 0 {
			return
		}
		prefix := prefix(absUrl)
		if prefix == HomePrefix {
			var ok bool
			mutex.Lock()
			_, ok = visted[absUrl]
			mutex.Unlock()
			if !ok {
				mutex.Lock()
				visted[absUrl] = true
				mutex.Unlock()
				c.Visit(absUrl)
			}
		} else {
			domain := searchDomain(prefix)
			if domain == "" {
				return
			}
			var has bool
			domaMu.Lock()
			has = hasDomain(domain)
			domaMu.Unlock()
			if !has {
				fmt.Println("Visiting", absUrl, " domain: ", domain)
				c.Visit(absUrl)
			}
		}

	})

	// Before making a request print "Visiting ..."
	//c.OnRequest(func(r *colly.Request) {
	//	fmt.Println("Visiting", r.URL.String())
	//})

	// Start scraping on https://hackerspaces.org
	c.Visit(Home)
	c.Wait()
	log.Printf("======================over======================")
}