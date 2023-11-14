package customerimporter

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitDomainCounter(t *testing.T) {
	path := "./test.csv"
	domainStruct := InitDomainCounter(path)
	assert.Equal(t, &DomainCounter{Path: path}, domainStruct, "struct not created correctly")
}

func TestReadCsvFileColumn_NoFile(t *testing.T) {
	path := "test.csv"
	columnNum := 0
	d := DomainCounter{Path: path}
	_, err := d.readCSVFileColumn(columnNum)
	assert.Equal(t, fmt.Sprintf("open %v: no such file or directory", path), err.Error())

}

func TestReadCsvFileColumn(t *testing.T) {
	tests := []struct {
		name      string
		path      string
		columnNum int
		data      []byte
		want      []string
		wantErr   string
	}{
		{"Empty file", "test.txt", 0, []byte(""), nil, "EOF"},
		{"OK", "test.txt", 2, []byte("first_name,last_name,email,gender,ip_address\nMildred,Hernandez,mhernandez0@github.io,Female,38.194.51.128"), []string{"mhernandez0@github.io"}, "EOF"},
		{"OK", "test.txt", 1, []byte("first_name,last_name,email,gender,ip_address\nMildred,Hernandez,mhernandez0@github.io,Female,38.194.51.128"), []string{"Hernandez"}, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Temporary file (just for tests).
			file, err := os.CreateTemp("", tt.path)
			assert.Nil(t, err)
			err = os.WriteFile(file.Name(), []byte(tt.data), 0666)
			assert.Nil(t, err)
			defer os.Remove(file.Name())

			d := DomainCounter{Path: file.Name()}
			got, err := d.readCSVFileColumn(tt.columnNum)
			if err != nil {
				assert.Equal(t, tt.wantErr, err.Error())
			} else {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestCountDomains(t *testing.T) {
	tests := []struct {
		name   string
		emails []string
		want   DomainCounter
	}{
		{"One wrong email",
			[]string{"mhernandez0@github.io", "email"},
			DomainCounter{
				Domains: []Domain{
					{Domain: "github.io", Count: 1}}}},
		{"All good emails",
			[]string{"mhernandez0@github.io", "tst@github.io"},
			DomainCounter{Domains: []Domain{
				{Domain: "github.io", Count: 2}}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := DomainCounter{}
			d.countDomains(tt.emails)
			assert.Equal(t, tt.want, d)
		})
	}
}

func TestSortAscend(t *testing.T) {
	domains := []Domain{
		{Domain: "gmail.com", Count: 2},
		{Domain: "outlook.com", Count: 1},
	}
	d := DomainCounter{Domains: domains}

	d.sortAscend()

	want := []Domain{
		{Domain: "outlook.com", Count: 1},
		{Domain: "gmail.com", Count: 2},
	}
	assert.Equal(t, want, d.Domains, "ascend sorting erorr")
}

func TestSortDescend(t *testing.T) {
	domains := []Domain{
		{Domain: "outlook.com", Count: 1},
		{Domain: "gmail.com", Count: 2},
	}
	d := DomainCounter{Domains: domains}

	d.sortDescend()

	want := []Domain{
		{Domain: "gmail.com", Count: 2},
		{Domain: "outlook.com", Count: 1},
	}
	assert.Equal(t, want, d.Domains, "descend sorting error")
}

func TestCountEmailsByDomains(t *testing.T) {
	tests := []struct {
		name      string
		path      string
		columnNum int
		sortType  SortType
		data      []byte
		want      []Domain
		wantErr   string
	}{
		{"Empty file", "test.txt", 0, SORT_ASCEND, []byte(""), nil, "EOF"},
		{"OK", "test.txt", 2, SORT_DESCEND, []byte("first_name,last_name,email,gender,ip_address\nMildred,Hernandez,mhernandez0@github.io,Female,38.194.51.128\nMildred,Hernandez,sdfs@github.io,Female,38.194.51.128\nMildred,Hernandez,sdfs@gmail.com,Female,38.194.51.128"), []Domain{{Domain: "github.io", Count: 2}, {Domain: "gmail.com", Count: 1}}, ""},
		{"OK", "test.txt", 2, SORT_ASCEND, []byte("first_name,last_name,email,gender,ip_address\nMildred,Hernandez,mhernandez0@github.io,Female,38.194.51.128\nMildred,Hernandez,sdfs@github.io,Female,38.194.51.128\nMildred,Hernandez,sdfs@gmail.com,Female,38.194.51.128"), []Domain{{Domain: "gmail.com", Count: 1}, {Domain: "github.io", Count: 2}}, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Temporary file (just for tests).
			file, err := os.CreateTemp("", tt.path)
			assert.Nil(t, err)
			err = os.WriteFile(file.Name(), []byte(tt.data), 0666)
			assert.Nil(t, err)
			defer os.Remove(file.Name())

			d := DomainCounter{Path: file.Name()}
			got, err := d.CountEmailsByDomains(tt.columnNum, tt.sortType)
			if err != nil {
				assert.Equal(t, tt.wantErr, err.Error())
			} else {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
