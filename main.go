package main

import (
	"fmt"
	"log"
	"teamworkgotests/customerimporter"
)

// Basic variables.
var (
	path         = "./data/customers.csv"
	columnNumber = 2
	sortType     = customerimporter.SORT_DESCEND
)

func main() {
	domainCounter := customerimporter.InitDomainCounter(path)
	result, err := domainCounter.CountEmailsByDomains(columnNumber, sortType)
	if err != nil {
		log.Panic(err)
	}

	fmt.Println(result)
}
