package cinepolis

import(
	"io/ioutil"
	"fmt"
	"github.com/moovweb/gokogiri"
	//	"github.com/gregjones/httpcache"
	"net/http"
	"strconv"
)

func ExtractCities(url string) map[int]string {

	html,_ := getBody(url)
	doc, err := gokogiri.ParseHtml(html)
	if err !=nil {
		fmt.Println(err)
	}
	defer doc.Free()

	options,_ := doc.Search("id('ctl00_ddlCiudad')/option");
	cities := make(map[int]string)
	for _, o := range options{
		val,_ := strconv.Atoi(o.Attributes()["value"].Content())
		if val > 0 {
			cities[val] = o.Content()
		}
	}
	
	return cities

}

func getBody(url string) (body []byte, err error) {
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
	return body,nil
}
