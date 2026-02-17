package main

import (
	"fmt"
	// "github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestSavePages(t *testing.T) {
	urlsWithPageranks := map[string]float64{
		"http://joie.com/1": 0.533,
		"http://joie.com/2": 0.455,
	}
	pagesString := getJSONPages("site-id", urlsWithPageranks)
	fmt.Println(pagesString)
}

