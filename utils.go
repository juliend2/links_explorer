package main

import (
	"strconv"
	"strings"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func hasKey(m map[string][]string, key string) bool {
	_, ok := m[key]
	return ok
}

func getStatusCodeFromError(err string) (int, error) {
	errParts := strings.SplitN(err, " ", 2)
	errCode, err2 := strconv.Atoi(errParts[0])
	return errCode, err2
}

func ValuesStringsMap(erroneousPages ErroneousPages) [][]string {
	var ret [][]string
	for _, erroneousPage := range erroneousPages.Pages {
		var urls []string
		for _, ce := range erroneousPage.ClientErrors {
			urls = append(urls, ce.Url)
		}
		ret = append(ret, urls)
	}
	return ret
}
