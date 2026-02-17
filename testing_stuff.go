package main

import (
	"fmt"
	"log"
	"statusmachine.com/crawling"
	"statusmachine.com/files"
)

func failIf(err error, msg string) {
	if err != nil {
		log.Fatalf("error "+msg+": %v", err)
	}
}

func failIfFalse(truth bool, msg string) {
	if truth != true {
		log.Fatalf("error " + msg)
	}
}

func failIfNotEqual(val1 string, val2 string) {
	if val1 != val2 {
		log.Fatalf("NOT EQUAL: %v != %v", val1, val2)
	}
}

func main() {
	fmt.Println("Testing...")
	failIfFalse(crawling.GetIsInternal("https://www.juliendesrosiers.com/2020/06/05/html-css-courses.php", "https://www.juliendesrosiers.com/"), "URL is not Internal")
	failIfFalse("https://www.domain.com" == crawling.GetSiteAddress("https://www.domain.com/"), "Fail")
	failIfFalse("https://domain.com" == crawling.GetSiteAddress("https://domain.com/"), "Fail")
	failIfFalse("https://domain.co.uk" == crawling.GetSiteAddress("https://domain.co.uk/"), "Fail")
	failIfFalse("https://domain.co.uk" == crawling.GetSiteAddress("https://domain.co.uk/lol/yolo/trololol/index.php?d=1&var=uno"), "Fail")
	failIfNotEqual("https://québec.ca", crawling.GetSiteAddress("https://québec.ca/lol/yolo/trololol/index.php?d=1&var=uno"))
	failIfNotEqual("https://québec.ca", crawling.GetSiteAddress("https://québec.ca"))
	failIfNotEqual("https://www1.bell.ca", crawling.GetSiteAddress("https://www1.bell.ca/buy/some/phones/"))
	failIfNotEqual("https://julien@luge.com", crawling.GetSiteAddress("https://julien@luge.com/buy/"))

	failIfNotEqual("/", crawling.GetUrlSuffix("https://google.com/"))
	failIfNotEqual("/search/index.php?query=bob", crawling.GetUrlSuffix("https://google.com/search/index.php?query=bob"))
	failIfNotEqual("/folder1;1/folder2:2/page.htm", crawling.GetUrlSuffix("https://google.com/folder1;1/folder2:2/page.htm"))

	// Files:

	files.Write("./testtext.txt", `
	J'aime bouger et la joie me le permet.
	Je l'mérite!
	`)

	content, err := files.Read("./testtext.txt")
	failIf(err, "Oops!")
	fmt.Println(content)


	files.WriteInSubdir("./tests/je/vais/marcher/bientot.txt", `
	<html>
		<body><h1>Marcher bientot</h1></body>
	</html>
	`)

	fmt.Println("Done.")
}
