package main

import (
	"encoding/json"
	"database/sql"
	"fmt"
	"github.com/PuerkitoBio/gocrawl"
	"github.com/PuerkitoBio/goquery"
	"github.com/PuerkitoBio/purell"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	//"runtime"
	"strconv"
	"statusmachine.com/algos"
	"statusmachine.com/files"
	"statusmachine.com/crawling"
	"statusmachine.com/dao"
	"statusmachine.com/domain"
	"statusmachine.com/extensions"
	"statusmachine.com/s3"
	"time"
)

// Execute it like this:
//   envdir /etc/sm-env/ workers/links_explorer/build/main

// This software component depends on the following environment variables:
// - CRAWLER_USER_AGENT_STRING
// - DATABASE_URL
// - LINKSCRAWLER_TMP_DIR

var currentCrawl crawling.Crawl
var urlsMap map[string][]string = make(map[string][]string)
var errorPages ErroneousPages
var crawledPages []domain.Page

var linksCrawlTmpDir string = func(linksCrawlTmpDirString string) string {
	if linksCrawlTmpDirString == "" {
		return "/tmp"
	} else {
		return linksCrawlTmpDirString
	}
}(os.Getenv("LINKSCRAWLER_TMP_DIR"))

var db sql.DB

// See https://github.com/PuerkitoBio/purell/blob/master/purell.go for more
// info regarding this:
const FlagsUnsafeGreedyWithoutFlagRemoveDirectoryIndex = purell.FlagsUsuallySafeGreedy | purell.FlagRemoveFragment | purell.FlagRemoveDuplicateSlashes | purell.FlagSortQuery
const FlagsAllGreedyWithoutFlagsUnsafeGreedy = purell.FlagDecodeDWORDHost | purell.FlagDecodeOctalHost | purell.FlagDecodeHexHost | purell.FlagRemoveUnnecessaryHostDots | purell.FlagRemoveEmptyPortSeparator
const MyFlags = FlagsUnsafeGreedyWithoutFlagRemoveDirectoryIndex | FlagsAllGreedyWithoutFlagsUnsafeGreedy

// NOTE: According to https://github.com/PuerkitoBio/purell/blob/master/purell.go#L47
// we need to have either RemoveWWW or AddWWW. It can't handle both cases.

// Create the Extender implementation, based on the gocrawl-provided DefaultExtender,
// because we don't want/need to override all methods.
type LinksCrawlExtender struct {
	gocrawl.DefaultExtender // Will use the default implementation of all but Visit() and Filter()
}

func (this *LinksCrawlExtender) Error(err *gocrawl.CrawlError) {

	if err.Kind == gocrawl.CekParseRedirectURL {
		fmt.Println("++++++++++++++++++++++++++++++")
		fmt.Println("REDIRECT FROM:", err.Ctx.NormalizedSourceURL())
		fmt.Println("REDIRECT TO:", err.Ctx.NormalizedURL())
		fmt.Println("++++++++++++++++++++++++++++++")
	}

	if err.Kind == gocrawl.CekHttpStatusCode {
		normalizedSourceURL := err.Ctx.NormalizedSourceURL()
		erroneousURL := err.Ctx.NormalizedURL()
		var normalizedSourceURLString string
		fmt.Println("++++++++++++++++++++++++++++++")
		statusCode, _ := getStatusCodeFromError(err.Error())
		fmt.Println("Error:", statusCode, "SourceURL:", normalizedSourceURL, "ErroneousURL:", erroneousURL, "err.Kind:", err.Kind)
		fmt.Println("++++++++++++++++++++++++++++++")
		if statusCode >= 400 && statusCode <= 599 {
			if normalizedSourceURL == nil {
				normalizedSourceURLString = ""
			} else {
				normalizedSourceURLString = normalizedSourceURL.String()
			}
			if !extensions.IsEmptyString(erroneousURL.String()) {
				handleClientError(&errorPages, statusCode, normalizedSourceURLString, err.Ctx.NormalizedURL().String())
			}
			urlsMap[erroneousURL.String()] = make([]string, 0)
		}
	}
}

// Get the path into which we'll save the data of the site that we're crawling.
func siteCrawlDir(baseTmpDir string, siteID int) string {
	return baseTmpDir + "/link_crawls/site_" + strconv.Itoa(siteID)
}

// Generates the path of the file in which we'll save the HTML content of the
// page we crawl.
func crawlFilePath(baseTmpDir string, siteID int, url string) string {
	return siteCrawlDir(baseTmpDir, siteID) + "/" + algos.MurmurHash(url) + ".html"
}

func pageHeadersFilePath(baseTmpDir string, siteID int, url string) string {
	return siteCrawlDir(baseTmpDir, siteID) + "/" + algos.MurmurHash(url) + ".json"
}

func (this *LinksCrawlExtender) Visit(context *gocrawl.URLContext, response *http.Response, doc *goquery.Document) (interface{}, bool) {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	check(err)
	defer db.Close()
	var normalizedURL string = context.NormalizedURL().String()
	headersJSON, _ := json.Marshal(response.Header)
	headersFilename := pageHeadersFilePath(linksCrawlTmpDir, currentCrawl.SiteId, normalizedURL)
	files.WriteInSubdir(headersFilename, string(headersJSON))
	defer response.Body.Close()
	docString, _ := ioutil.ReadAll(response.Body)
	filename := crawlFilePath(linksCrawlTmpDir, currentCrawl.SiteId, normalizedURL)
	files.WriteInSubdir(filename, string(docString))

	var anchors []string = crawling.GetAnchorsFor(string(docString), context)
	var internalAnchors []string
	for _, anchor := range anchors {
		url, _ := context.URL().Parse(anchor)
		if crawling.GetIsInternal(url.Host, currentCrawl.URL) {
			internalAnchors = append(internalAnchors, anchor)
		}
	}
	urlsMap[normalizedURL] = internalAnchors
	// TODO: fill out the gaps in this map, but also adding the sub links. This
	// makes sure the lenght of the slice we'll save, accounts for all the links
	// loop into the anchors.
	// 	if there doesn't exist a key with the current anchor, add it, with an empty value
	return nil, true
}

func (this *LinksCrawlExtender) Filter(ctx *gocrawl.URLContext, isVisited bool) bool {
	url := *ctx.NormalizedURL()
	var isInternal bool = crawling.GetIsInternal(url.Host, currentCrawl.URL)
	regex := `(?i)\.(flv|swf|png|jpg|jpeg|gif|asx|zip|rar|tar|7z|gz|jar|js|css|dtd|xsd|ico|raw|mp3|mp4|wav|wmv|ape|aac|ac3|wma|aiff|mpg|mpeg|avi|mov|ogg|mkv|mka|asx|asf|mp2|m1v|m3u|f4v|pdf|doc|xls|ppt|pps|bin|exe|rss|xml)$`
	isMediaFile, _ := regexp.MatchString(regex, url.Path)
	// fmt.Println("isVisited", isVisited, "and it should be false")
	// fmt.Println("isInternal", isInternal, "and it should be true")
	// fmt.Println("isMediaFile", isMediaFile, "and it should be false")
	var condition bool = !isVisited && isInternal && !isMediaFile
	// fmt.Println("So:", condition, "for", url.Host, url.Path, "based on:", currentCrawl.URL)
	// fmt.Println("--------------------------------------------------------------------------")
	return condition
}

func (this *LinksCrawlExtender) RequestRobots(ctx *gocrawl.URLContext, robotAgent string) (data []byte, request bool) {
	if currentCrawl.ObeyRobots {
		fmt.Println("Obey robots")
		return nil, true // will obey the robots.txt
	} else {
		fmt.Println("Disobey robots")
		return []byte(""), false // as if there was no robots.txt at all to obey
	}
}

func main() {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	check(err)
	defer db.Close()
	fmt.Println("Starting Links Explorer. Will crawl sites.")
	// Infinite Loop:
	for {
		currentCrawl = crawling.GetNextCrawl(db)
		if currentCrawl.SiteId == 0 {
			fmt.Println("Nothing to crawl. (looping)")
		} else {
			var url string = currentCrawl.URL
			fmt.Println("Going to crawl:", url)

			// Add a trailing slash to the starting URL, because on creation, we
			// chomp it. So it makes sense to add it when starting the crawl, to
			// make sure it doesn't get flagged as 'isVisited' right from the
			// beginning, for example on lafamilledulait.com/fr, we have a 301
			// redirect to lafamilledulait.com/fr/ , but it was flagged as
			// isVisited because of the way the crawler normalizes the URLs
			// internally:
			// Set custom options
			opts := gocrawl.NewOptions(new(LinksCrawlExtender))
			opts.CrawlDelay = 1 * time.Second
			opts.LogFlags = gocrawl.LogAll
			opts.UserAgent = os.Getenv("CRAWLER_USER_AGENT_STRING")
			opts.SameHostOnly = false
			opts.URLNormalizationFlags = MyFlags
			crawler := gocrawl.NewCrawlerWithOptions(opts)

			startErr := dao.UpdateCrawlStatus(db, currentCrawl.LinkcrawlID, "started")
			check(startErr)
			crawling.PerformCrawl(crawler, currentCrawl, url+"/")

			// NOTE: HEADS UP: PageRank is taking a lot of RAM (relative to N of URLs):
			var urlsWithPageranks map[string]float64 = algos.GeneratePageranks(urlsMap)
			pagesJSON := domain.GetJSONPages(currentCrawl.SiteId, currentCrawl.LinkcrawlID, urlsWithPageranks)
			s3.UploadFile(pagesJSON, "sm-dev", crawling.LinksCrawlFileName(currentCrawl.SiteId, currentCrawl.LinkcrawlID))
			// For each page, insert or update it with the proper values
			for url, pageRank := range urlsWithPageranks {
				fmt.Println("URL:", url)
				urlSuffix := crawling.GetUrlSuffix(url)
				upsertErr := dao.UpsertEnvSitePage(db, currentCrawl.SiteId, urlSuffix, pageRank)
				if upsertErr != nil {
					log.Fatalf("upsert error %v: ", upsertErr)
				}
				htmlFilePath := crawlFilePath(linksCrawlTmpDir, currentCrawl.SiteId, url)
				headersFilePath := pageHeadersFilePath(linksCrawlTmpDir, currentCrawl.SiteId, url)
				var html string
				var headers string
				if files.FileExists(htmlFilePath) {
					html, _ = files.Read(htmlFilePath)
					headers, _ = files.Read(headersFilePath)
				} else {
					html = ""
					headers = "{}"
				}

				// NOTE: HEADS UP: This part does a lot of round trips to the DB (1
				// select + 1 insert for EACH page)
				env_site_page_id, selectErr := dao.SelectEnvSitePage(db, currentCrawl.SiteId, urlSuffix)
				check(selectErr)
				var statusCode int = 200

				if code := errorPages.HasPageByUrl(url); code != -1 {
					statusCode = code
				}
				fmt.Println("STATUS CODE:", statusCode)

				insertErr := dao.InsertPage(db, domain.Page{
					EnvSitePageID: env_site_page_id,
					LinkCrawlID: currentCrawl.LinkcrawlID,
					Url: url,
					StatusCode: statusCode,
					Headers: headers,
					Body: html,
				})

				if insertErr != nil {
					fmt.Printf("insert error: %v", insertErr)
				}
			}

			finishErr := dao.UpdateCrawlStatus(db, currentCrawl.LinkcrawlID, "finished")
			check(finishErr)
		}

		urlsMap = make(map[string][]string) // VERY IMPORTANT to reset this data structure after EACH site
		errorPages.Reset()
		time.Sleep(10 * time.Second)
	}
}
