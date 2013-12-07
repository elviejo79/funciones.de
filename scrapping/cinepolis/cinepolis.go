package cinepolis

import(
	"io/ioutil"
	"fmt"
	"github.com/moovweb/gokogiri"
	"github.com/moovweb/gokogiri/xml"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Showtime struct{
	IdTheater string
	IdMovie string
	Language string
	RoomType string
	Date string
	Time string
}

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

func ExtractTheaters(url string) map[int]string {
	result := make(map[int]string)
	html,err := getBody(url)
	if err != nil {
		fmt.Printf("%#v",err)
	}
	doc, _ := gokogiri.ParseHtml(html)
	defer doc.Free()
	//theaters,_ := doc.Search("//a[ends-with(@id,'306')]");
	movies,_ := doc.Search("//a[contains(@id, 'idPelCine')]")
	for _, m := range movies{ 
		cineId := nodeContent("@id",m)[14:]
		cineName := nodeContent("//select[@name='cartelera"+cineId+"']/parent::*/parent::*/parent::*//span[@class='TitulosBlanco']",m)
		intCineId, _ := strconv.Atoi(cineId)
		result[intCineId] = cineName
	}
	return result
}

func ExtractMovies(url string, idCine int) (res []Showtime, err error) {
	html,err := getBody(url)
	if err != nil {
		fmt.Printf("%#v",err)
	}
	doc, _ := gokogiri.ParseHtml(html)
	defer doc.Free()
	//theaters,_ := doc.Search("//a[ends-with(@id,'306')]");
	xpath := fmt.Sprintf("//a[contains(@id, 'idPelCine') and (substring(@id, string-length(@id) -2)=%d)]",idCine)
	movies,_ := doc.Search(xpath)
	fmt.Printf("url: %s \nXpath: %s \nExtractMovies: %#v \n\n", url,xpath,movies)
	for _, m := range movies{ 
		cineId := nodeContent("@id",m)[14:]
		titulo := nodeContent("parent::*//a[@class='peliculaCartelera']",m)
		subtitulos := titulo[len(titulo)-3:]
		if subtitulos == "Sub" {
			subtitulos = "SUBTITULADA"
                } else {
			subtitulos = "ESPAÑOL"
		}

		sala := titulo[len(titulo)-7:len(titulo)-4]
		if strings.Contains(titulo," 4D") {
			titulo=titulo[:strings.Index(titulo," 4D")]
			sala = "4D"
		} else if strings.Contains(titulo," 3D ") {
			titulo=titulo[:strings.Index(titulo," 3D ")]
			sala = "3D"
		} else if strings.Contains(titulo," IMAX") {
			titulo=titulo[:strings.Index(titulo," IMAX")]
			sala = "IMAX"
		} else if strings.Contains(titulo," XE") {
			titulo=titulo[:strings.Index(titulo," XE ")]
			sala = "XE"
		} else if strings.Contains(titulo," Dig ") {
			titulo=titulo[:strings.Index(titulo," Dig ")]
			sala = "Dig"
		} else {
			//titulo=titulo
			sala = ""
		}
		sala = "R"+sala
		titulo = strings.ToUpper(titulo)
		titulo = strings.Replace(titulo,":","", -1)
		t := time.Now().Format("20060102")
		
		row := Showtime{
			cineId, //cineID
			titulo,
			subtitulos,
			sala,
			t,
			"00:00",
		}
		

		hours,err := m.Search("parent::*/parent::*//*[contains(@class,'horariosCartelera')]")
		if err != nil {
			fmt.Println(err)
		}

		horas := []string{}
		for _,e := range hours {
			t,_:= time.Parse("3:04pm",e.Content())
			horas = append(horas,t.Format("15:04"))
		}

		//if row != nil {
			for _,h := range horas {
				row.Time = h
				res = append(res,row)
			}
		//}

	}

	return
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

func nodeContent(x_path string,m xml.Node) (result string){
	ts,_ := m.Search(x_path)
	for _,e := range ts{
		result = e.Content()
	}
	return 
}

