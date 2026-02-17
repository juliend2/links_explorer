package main

import (
	"encoding/json"
	_ "github.com/lib/pq"
	"path/filepath"
	"statusmachine.com/extensions"
)

var siteURLsCache map[string]string = make(map[string]string)

const dirSeparator string = string(filepath.Separator)

/*
started_at = now
check_id = worker.append_check(server_error_urls, SERVER_ERRORS_CHECK, started_at)
if server_error_urls.size > 0
	worker.append_issue_to_raise(SERVER_ERRORS_ISSUE, check_id: check_id, concerned_pages_count: server_error_urls.size)
else
	worker.append_issue_to_resolve(SERVER_ERRORS_ISSUE, concerned_pages_count: 0)
end
*/

func generateUrlsJsonSlice(urlsMap map[string][]string) []string {
	var urlsSlice []string
	for key, value := range urlsMap {
		theMap := map[string][]string{key: value}
		jsonBytes, _ := json.Marshal(theMap)
		urlsSlice = append(urlsSlice, string(jsonBytes))
	}
	return urlsSlice
}

func handleClientError(clientErrorPages *ErroneousPages, statusCode int, sourceURL string, erroneousURL string) {
	// source url exists in collection:
	if clientErrorPages.HasKey(sourceURL) {
		erroneousPage := clientErrorPages.Pages[sourceURL]
		// so if this error page is not already present in the collection, append it:
		if !erroneousPage.HasClientError(erroneousURL) {
			erroneousPage.AppendClientError(*NewClientError(statusCode, erroneousURL))
		}
		clientErrorPages.Pages[sourceURL] = erroneousPage
	} else {
		// source url doesn't exists in collection:
		// create the first one-element list:
		var list []ClientError = []ClientError{*NewClientError(statusCode, erroneousURL)}
		// if the collection is empty:
		if clientErrorPages.Length() == 0 {
			// create a new map of ErroneousPages that we'll assign to it:
			clientErrorPages.Pages = make(map[string]ErroneousPage)
		}
		// assign the ErroneousPage to the already-existing-map:
		clientErrorPages.Pages[sourceURL] = *NewErroneousPage(sourceURL, list)
	}
}

func handleServerError(serverErrors *[]string, statusCode int, erroneousURL string) {
	if extensions.IndexOfString(*serverErrors, erroneousURL) == -1 {
		*serverErrors = append(*serverErrors, erroneousURL)
	}
}
