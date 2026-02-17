package main

import (
	//"fmt"
	"reflect"
	"strconv"
	"strings"
)

type ClientError struct {
	Code int    `json:"code"`
	Url  string `json:"url"`
}

func (this *ClientError) String() string {
	return "ClientError{ Code: " + strconv.Itoa(this.Code) + ", Url: '" + this.Url + "' }"
}

func NewClientError(statusCode int, url string) *ClientError {
	ret := new(ClientError)
	ret.Code = statusCode
	ret.Url = url
	return ret
}

// ErroneousPage

type ErroneousPage struct {
	Url          string        `json:"url"`
	ClientErrors []ClientError `json:"client_errors"`
}

func (this *ErroneousPage) String() string {
	var clientErrs []string
	for _, ce := range this.ClientErrors {
		clientErrs = append(clientErrs, ce.String())
	}
	return "ErroneousPage{ Url: '" + this.Url + "', ClientErrors: " + strings.Join(clientErrs, ", ") + " }"
}

func NewErroneousPage(url string, clientErrs []ClientError) *ErroneousPage {
	ret := new(ErroneousPage)
	ret.Url = url
	ret.ClientErrors = clientErrs
	return ret
}

func (this *ErroneousPage) HasClientError(url string) bool {
	for _, clientErr := range this.ClientErrors {
		if clientErr.Url == url {
			return true
		}
	}
	return false
}

func (this *ErroneousPage) AppendClientError(clientError ClientError) {
	this.ClientErrors = append(this.ClientErrors, clientError)
}

// ErroneousPages

type ErroneousPages struct {
	Pages map[string]ErroneousPage
}

/*
	ErroneousPages:
		Pages: map[string]ErroneousPage
			ErroneousPage:
				Url          string
				ClientErrors []ClientError
					Code int
					Url  string
*/
func (this *ErroneousPages) PreMarshal() map[string][]ClientError {
	var ret map[string][]ClientError = make(map[string][]ClientError)
	for url, erroneousPage := range this.Pages {
		cliErrors := make([]ClientError, len(erroneousPage.ClientErrors))
		for i, cliError := range erroneousPage.ClientErrors {
			if cliError.Code == 404 {
				cliErrors[i] = cliError
			}
		}
		ret[url] = cliErrors
	}
	return ret
}

func typeof(v interface{}) string {
	return reflect.TypeOf(v).String()
}

// Returns the status code of the error, or -1 if none found
func (this *ErroneousPages) HasPageByUrl(givenUrl string) int {
	// fmt.Printf("HasPageByUrl url: '%v' \n", givenUrl)
	for _, erroneousPage := range this.Pages {
		for _, clientErr := range erroneousPage.ClientErrors {
			// fmt.Println("Type of elements: ", typeof(clientErr.Url), typeof(givenUrl))
			if clientErr.Url == givenUrl {
				return clientErr.Code
			}
		}
	}
	return -1
}

func (this *ErroneousPages) Reset() {
	this.Pages = make(map[string]ErroneousPage)
}

func (this *ErroneousPages) String() string {
	var errPagesString []string
	for _, errPage := range this.Pages {
		errPagesString = append(errPagesString, errPage.String())
	}
	return "ErroneousPages{ Pages: " + strings.Join(errPagesString, ", ") + " }"
}

func (this *ErroneousPages) HasKey(key string) bool {
	_, ok := this.Pages[key]
	return ok
}

func (this *ErroneousPages) Length() int {
	return len(this.Pages)
}
