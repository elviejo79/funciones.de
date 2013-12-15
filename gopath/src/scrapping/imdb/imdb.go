package imdb

import (
	"scrapping/cinebase"
	"fmt"
	"net/url"
	"regexp"
)

func Search(title string) (results ImdbResults, err error) {
	url := fmt.Sprintf("http://www.imdb.com/find?q=%s&s=tt",
		url.QueryEscape(title))
	xpath := "id('main')//table//tr[1]/td[2]"
	xNodes := cinebase.NodesExtractor(url, xpath)
	re, err := regexp.Compile(`^/title/(.+?)/.ref.*$`)

	for _, m := range xNodes {
		results.Title = cinebase.NodeContent("a", m)
		results.Href = cinebase.NodeContent("a/@href", m)
		res := re.FindAllStringSubmatch(results.Href, -1)
		results.TT = res[0][1]
	}

	return
}

func SearchImdbTT(title string) (string) {
	tmp,_ := Search(title)
	return tmp.TT
}

type ImdbEntry struct {
	Id                string
	Title             string
	Name              string
	Title_description string
	Episodo_title     string
	Description       string
}

type ImdbResults struct {
	Title string
	Href  string
	TT    string
}
