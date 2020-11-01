package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"sort"
	"strings"
	"time"
)

var (
	searchURI = "https://www.finn.no/car/used/search.html?location=20061&location=22030&location=22038&make=0.749&model=1.749.2000264&price_from=100000&q=120+fully&sort=0&year_from=2019"
)

func main() {
	r, err := http.Get(searchURI)
	check(err)
	
	today := time.Now().Format("2006-01-02")
	fmt.Println(today)	

	b, err := ioutil.ReadAll(r.Body)
	check(err)

	foundBytes := regexp.MustCompile(`finnkode=[0-9]{9}`).FindAll(b, -1)

	adIDs := []string{}
	for _, v := range foundBytes {
		adIDs = append(adIDs, string(v))
	}

	adList := removeDuplicates(adIDs)
	trimText(adList)
	sort.Strings(adList)
	
	listings := map[string](string){}
	for _, kode := range adList {
		listings[kode] = getPrice(kode)
		// fmt.Println(kode + " " + getPrice(kode))
	}
	
}


func trimText(list []string) {
	for i, v := range list {
		list[i] = strings.Split(v, "=")[1]
	}
}

func removeDuplicates(xs []string) []string {
	unique := map[string](string){}
	for _, v := range xs {
		if _, ok := unique[v]; !ok {
			unique[v] = "present"
		}
	}
	
	list := []string{}
	for k := range unique {
		list = append(list, k)
	}
	return list
}

func getPrice(finnKode string) string {
	adURL := fmt.Sprintf("https://www.finn.no/car/used/ad.html?finnkode=%s", finnKode)
	
	r, err := http.Get(adURL)
	if err != nil {
		log.Fatalf("error getting ad: %e", err)
	}
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatalf("error reading body: %e", err)
	}
	
	if strings.Contains(string(b), "SOLGT") {
		return "SOLGT"
	}
	
	re := regexp.MustCompile(`[0-9]{3}.[0-9]{3}\skr`)
	
	price := re.FindAllString(string(b), 1)[0]
	
	return trimWhiteSpaceAndKR(price)
}

func trimWhiteSpaceAndKR(s string) string {
	s = strings.TrimSpace(s)
	s = strings.Replace(s, string('\u00A0'), "", -1) // forbanna irriterende NO-BREAK SPACE
	s = strings.Replace(s, "kr", "", -1)
	return s
}

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}