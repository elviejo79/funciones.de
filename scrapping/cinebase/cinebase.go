package cinebase
import (
	"github.com/moovweb/gokogiri"
	"github.com/moovweb/gokogiri/xml"
	"io/ioutil"
	"net/http"
	"log"
)

var bodyCache = make(map[string][]byte)

func GetBody(url string) (body []byte, err error) {
	if page, ok := bodyCache[url]; ok {
		return page, nil
	}

	client := http.Client{}
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	bodyCache[url] = body
	return body, nil
}

func NodeContent(x_path string, m xml.Node) (result string) {
	ts, _ := m.Search(x_path)
	for _, e := range ts {
		result = e.Content()
	}
	return
}

func NodesExtractor(url string, xpath string) ([]xml.Node){
	html,err := GetBody(url)

	if err != nil {
		log.Fatal("NodesExtractor",err)
	}
	doc, err := gokogiri.ParseHtml(html)
	if err != nil {
		log.Fatal("NodesExtractor",err)
	}

	//defer doc.Free() //esto al por algún motivo cierra antes de y genera error de nulos

	movies, err := doc.Search(xpath)
	if err != nil {
		log.Fatal(err)
	}
	if len(movies) == 0 {
		log.Printf("\n xpath nodes not found \n\turl: %s \n\tXpath: %s \n\tExtractMovies: %#v \n\n", url, xpath, movies)
	}
	return movies
}

