package main

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"statusmachine.com/extensions"
	"testing"
	// "time"
)

func TestErroneousPages(t *testing.T) {
	clientErrors := make([]ClientError, 2)
	clientErrors[0] = ClientError{
		Url:  "http://cool1.com",
		Code: 404,
	}
	clientErrors[1] = ClientError{
		Url:  "http://cool2.com",
		Code: 404,
	}
	brokenLinks := ErroneousPages{
		Pages: map[string]ErroneousPage{
			"http://joie1.com": ErroneousPage{
				Url:          "http://joie1.com",
				ClientErrors: clientErrors,
			},
			"http://joie2.com": ErroneousPage{
				Url: "http://joie1.com",
				ClientErrors: []ClientError{
					ClientError{
						Url:  "http://poulet.com",
						Code: 404,
					},
				},
			},
		},
	}
	var numBrokenLinks int = len(extensions.UniqStrings(extensions.FlattenStringSlices(ValuesStringsMap(brokenLinks))))
	assert.Equal(t, numBrokenLinks, 3, "Should be 3 broken links")
	brokenLinksJson, _ := json.Marshal(brokenLinks.PreMarshal())
	fmt.Println(string(brokenLinksJson))
	clientError := ClientError{
		Code: 404,
		Url:  "http://joie.com",
	}
	clientErrorJson, _ := json.Marshal(clientError)
	fmt.Println(string(clientErrorJson))
}
