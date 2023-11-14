// package customerimporter reads from the given customers.csv file and returns a
// sorted (data structure of your choice) of email domains along with the number
// of customers with e-mail addresses for each domain.  Any errors should be
// logged (or handled). Performance matters (this is only ~3k lines, but *could*
// be 1m lines or run on a small machine).
package customerimporter

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"sort"
	"strings"
	"time"
)

type Domain struct {
	Domain string
	Count  uint
}

type DomainCounter struct {
	Path    string
	Domains []Domain
}

type SortType int

const (
	SORT_ASCEND SortType = iota
	SORT_DESCEND
)

// Initializes DomainCounter structure with the given path to the file.
func InitDomainCounter(path string) *DomainCounter {
	return &DomainCounter{Path: path}
}

// Reads entire csv file but writes to slice only the given column and returns it.
func (d *DomainCounter) readCSVFileColumn(columnNumber int) ([]string, error) {
	file, err := os.Open(d.Path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	parser := csv.NewReader(file)

	// Skips headers.
	_, err = parser.Read()
	if err != nil {
		return nil, err
	}

	var records []string
	// Reads csv file record by record. I considered using goroutines but honestly this was faster.
	for {
		record, err := parser.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		records = append(records, record[columnNumber])
	}

	return records, nil
}

// Counts emails by domain.
func (d *DomainCounter) countDomains(emails []string) {
	emailMap := make(map[string]uint)

	for _, email := range emails {
		splittedEmail := strings.Split(email, "@")
		// Must have 2 length.
		if len(splittedEmail) != 2 {
			log.Printf("found invalid email: %v", email)
			continue
		}
		emailMap[splittedEmail[1]]++
	}

	for domain, count := range emailMap {
		d.Domains = append(d.Domains, Domain{Domain: domain, Count: count})
	}
}

// Sorts domains from structure ascending.
func (d *DomainCounter) sortAscend() {
	sort.Slice(d.Domains, func(i, j int) bool {
		return d.Domains[i].Count < d.Domains[j].Count
	})
}

// Sorts domains from structure descending.
func (d *DomainCounter) sortDescend() {
	sort.Slice(d.Domains, func(i, j int) bool {
		return d.Domains[i].Count > d.Domains[j].Count
	})
}

// Processes all the work from reading csv file to counting and sorting the domains.
// Returns the domain structure, which consists of the name and quantity for each domain.
// To call this method, specify the column number of the email (counted from 0) and the sort type  (SORT.ASCEND/SORT_DESCEND).
func (d *DomainCounter) CountEmailsByDomains(columnNumber int, sortType SortType) ([]Domain, error) {
	start := time.Now()

	emails, err := d.readCSVFileColumn(columnNumber)
	if err != nil {
		return nil, err
	}

	d.countDomains(emails)

	if sortType == SORT_ASCEND {
		d.sortAscend()
	}
	if sortType == SORT_DESCEND {
		d.sortDescend()
	}

	elapsed := time.Since(start)
	log.Printf("process took %s", elapsed)

	return d.Domains, nil
}
