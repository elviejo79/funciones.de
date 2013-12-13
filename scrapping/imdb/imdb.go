package imdb

import (
	"encoding/json"
	"fmt"
	"../cinebase"
	"net/url"
)

func SearchOne(title string) (result ImdbEntry) {
	records := Search(title)
	if records.Title_approx != nil {
		result =records.Title_approx[0]
	}else if records.Title_popular != nil {
		result =records.Title_popular[0]
	} 
	return
}

func Search(title string) (results results, err error) {
	url := fmt.Sprintf("http://www.imdb.com/xml/find?json=1&q=%s&s=all",
		url.QueryEscape(title))
	html, err := cinebase.GetBody(url)
	fmt.Printf("el url usado es: %s y regroso %d bytes \n",url,len(html))
	err = json.Unmarshal(html, &results)
	return
}

type ImdbEntry struct {
	Id                string
	Title             string
	Name              string
	Title_description string
	Episodo_title     string
	Description       string
}

type results struct {
	Title_popular []ImdbEntry
	Title_approx []ImdbEntry
}
