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
)
type node struct{
	end bool
	folw map[byte] *node
}
var (
	DomainRegex = `[-a-z0-9]+\.(cat|goog|sohu|ac|academy|ac\.cn|accountant|accountants|actor|ad|adult|ae|aero|af|ag|agency|ah\.cn|ai|airforce|al|am|amsterdam|an|ao|apartments|app|aq|ar|archi|army|art|as|asia|associates|at|attorney|au|auction|auto|autos|aw|az|ba|baby|band|bar|barcelona|bargains|bayern|bb|bd|be|beauty|beer|berlin|best|bet|bf|bg|bh|bi|bid|bike|bingo|bio|biz|biz\.pl|bj|bj\.cn|black|blog|blue|bm|bn|bo|boats|boston|boutique|br|bs|bt|build|builders|business|buzz|bv|bw|by|bz|ca|cab|cafe|camera|camp|capital|car|cards|care|careers|cars|casa|cash|casino|catering|cc|cd|center|ceo|cf|cg|ch|charity|chat|cheap|church|ci|city|ck|cl|claims|cleaning|clinic|clothing|cloud|club|cm|cn|co|coach|codes|coffee|co\.in|co\.jp|co\.kr|college|com|com\.ag|com\.au|com\.br|com\.bz|com\.cn|com\.co|com\.es|com\.ky|community|com\.mx|company|com\.pe|com\.ph|com\.pl|computer|com\.tw|condos|construction|consulting|contact|contractors|co\.nz|cooking|cool|coop|co\.uk|country|coupons|courses|co\.za|cq\.cn|cr|credit|creditcard|cricket|cruises|cu|cv|cx|cy|cymru|cz|dance|date|dating|de|deals|degree|delivery|democrat|dental|dentist|design|dev|diamonds|digital|direct|directory|discount|dj|dk|dm|do|doctor|dog|domains|download|dz|earth|ec|edu|education|ee|eg|eh|email|energy|engineer|engineering|enterprises|equipment|er|es|estate|et|eu|events|exchange|expert|exposed|express|fail|faith|family|fan|fans|farm|fashion|fi|film|finance|financial|firm\.in|fish|fishing|fit|fitness|fj|fj\.cn|fk|flights|florist|fm|fo|football|forsale|foundation|fr|fun|fund|furniture|futbol|fyi|ga|gallery|games|garden|gay|gd|gd\.cn|ge|gen\.in|gf|gg|gh|gi|gifts|gives|gl|glass|global|gm|gmbh|gn|gold|golf|gov|gov\.cn|gp|gq|gr|graphics|gratis|green|gripe|group|gs|gs\.cn|gt|gu|guide|guru|gw|gx\.cn|gy|gz\.cn|ha\.cn|hair|haus|hb\.cn|health|healthcare|he\.cn|hi\.cn|hk|hk\.cn|hl\.cn|hm|hn|hn\.cn|hockey|holdings|holiday|homes|horse|hospital|host|house|hr|ht|hu|icu|id|idv|idv\.tw|ie|il|im|immo|immobilien|in|inc|ind\.in|industries|info|info\.pl|ink|institute|insure|int|international|investments|io|iq|ir|irish|is|ist|istanbul|it|je|jetzt|jewelry|jl\.cn|jm|jo|jobs|jp|js\.cn|jx\.cn|kaufen|ke|kg|kh|ki|kim|kitchen|kiwi|km|kn|kp|kr|kw|ky|kz|la|land|law|lawyer|lb|lc|lease|legal|lgbt|li|life|lighting|limited|limo|link|live|lk|llc|ln\.cn|loan|loans|london|love|lr|ls|lt|ltd|ltda|lu|luxury|lv|ly|ma|maison|makeup|management|market|marketing|mba|mc|md|me|media|melbourne|memorial|men|menu|me\.uk|mg|mh|miami|mil|mk|ml|mm|mn|mo|mobi|mo\.cn|moda|moe|money|monster|mortgage|motorcycles|movie|mp|mq|mr|ms|mt|mu|museum|mv|mw|mx|my|mz|na|nagoya|name|navy|nc|ne|ne\.kr|net|net\.ag|net\.au|net\.br|net\.bz|net\.cn|net\.co|net\.in|net\.ky|net\.nz|net\.pe|net\.ph|net\.pl|network|news|nf|ng|ni|ninja|nl|nm\.cn|no|nom\.co|nom\.es|nom\.pe|np|nr|nrw|nu|nx\.cn|nyc|nz|okinawa|om|one|onl|online|org|org\.ag|org\.au|org\.cn|org\.es|org\.in|org\.ky|org\.nz|org\.pe|org\.ph|org\.pl|org\.uk|pa|page|paris|partners|parts|party|pe|pet|pf|pg|ph|photography|photos|pictures|pink|pizza|pk|pl|place|plumbing|plus|pm|pn|poker|porn|pr|press|pro|productions|promo|properties|protection|ps|pt|pub|pw|py|qa|qh\.cn|quebec|quest|racing|re|realestate|recipes|red|rehab|reise|reisen|re\.kr|ren|rent|rentals|repair|report|republican|rest|restaurant|review|reviews|rich|rip|ro|rocks|rodeo|ru|run|rw|ryukyu|sa|sale|salon|sarl|sb|sc|sc\.cn|school|schule|science|sd|sd\.cn|se|security|services|sex|sg|sh|sh\.cn|shiksha|shoes|shop|shopping|show|si|singles|site|sj|sk|ski|skin|sl|sm|sn|sn\.cn|so|soccer|social|software|solar|solutions|space|/span|sr|st|storage|store|stream|studio|study|style|supplies|supply|support|surf|surgery|sv|sx\.cn|sy|sydney|systems|sz|tax|taxi|tc|td|team|tech|technology|tel|tennis|tf|tg|th|theater|theatre|tienda|tips|tires|tj|tj\.cn|tk|tl|tm|tn|to|today|tokyo|tools|top|tours|town|toys|tp|tr|trade|training|travel|tt|tube|tv|tw|tw\.cn|tz|ua|ug|uk|um|university|uno|us|uy|uz|va|vacations|vc|ve|vegas|ventures|vet|vg|vi|viajes|video|villas|vin|vip|vision|vn|vodka|vote|voto|voyage|vu|wales|wang|watch|webcam|website|wedding|wf|wiki|win|wine|work|works|world|ws|wtf|xin|xj\.cn|xxx|xyz|xz\.cn|yachts|ye|yn\.cn|yoga|yokohama|yr|yt|yu|za|zj\.cn|zm|zone|zw|中国|中文网|企业|佛山|信息|公司|商城|商店|商标|在线|娱乐|广东|我爱你|手机|招聘|游戏|移动|网址|网络|集团|餐厅)$`
	Tree map[byte] *node = make(map[byte] *node)
	domainRegex = regexp.MustCompile(DomainRegex)
	MaxDepth int
	Parallelism int
	Home string
	HomeDomain string
	HomePrefix string
	mutex sync.Mutex
	domaMu sync.Mutex
	visted map[string] bool = map[string]bool{}
	domainFile string
)
func searchDomain(domain string) string  {
	//doubleSlah := strings.Index(domain, "//")
	//if doubleSlah != -1 {
	//	domain = domain[doubleSlah + 2:]
	//}
	shortDomain := domainRegex.FindAllStringSubmatch(domain, -1)
	len := len(shortDomain)
	if len == 0 {
		return ""
	}
	if isDomain(shortDomain[0][0]) {
		return shortDomain[0][0]
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
				c.Visit(absUrl)
			}
		}

	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	// Start scraping on https://hackerspaces.org
	c.Visit(Home)
	c.Wait()
}