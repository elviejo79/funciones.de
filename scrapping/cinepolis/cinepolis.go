package cinepolis

import (
	"../cinebase/"
	"fmt"
	"funciones.de/app/models"
	"strconv"
	"strings"
	"time"
	"log"
)

var Company models.Company = models.NewCompany("http://cinepolis.com", "cinepolis")

func ExtractCities() (cities []models.City) {
	options := cinebase.NodesExtractor(Company.Url, "id('ctl00_ddlCiudad')/option")

	for _, o := range options {
		if nil != o.Attributes()["value"] {

			val, _ := strconv.Atoi(o.Attributes()["value"].Content())
			cityName := o.Content()
			df_am := "D.F. y A.M. ("
			if strings.Contains(cityName, df_am) {
				cityName = cityName[len(df_am) : len(cityName)-1]
			}
			if val > 0 {
				cities = append(cities, models.NewCity(val, cityName))
			}
		}
	}

	return cities
}

func ExtractTheaters(c models.City) (theaters []models.Theater) {
	tmp := make(map[int]models.Theater)
	url := fmt.Sprintf("http://cinepolis.com/_CARTELERA/cartelera.aspx?ic=%d", c.IdCity)
	movies := cinebase.NodesExtractor(url, "//a[contains(@id, 'idPelCine')]")
	for _, m := range movies {
		cineId := cinebase.NodeContent("@id", m)[14:]
		cineName := cinebase.NodeContent("//select[@name='cartelera"+cineId+"']/parent::*/parent::*/parent::*//span[@class='TitulosBlanco']", m)
		intCineId, _ := strconv.Atoi(cineId)
		if cineName != "" {
			tmp[intCineId] = models.NewTheater(intCineId, cineName)
		}
	}

	for _, t := range tmp {
		theaters = append(theaters, t)
	}

	return
}

func ExtractShowtimes(c models.City, t models.Theater) (res []models.Showtime, err error) {
	url := fmt.Sprintf("http://cinepolis.com/_CARTELERA/cartelera.aspx?ic=%d", c.IdCity)
	len_idCine := len(strconv.Itoa(t.IdTheater)) - 1
	xpath := fmt.Sprintf("//a[contains(@id, 'idPelCine') and (substring(@id, string-length(@id) -%d)=%d)]", len_idCine, t.IdTheater)

	movies := cinebase.NodesExtractor(url, xpath)

	for _, m := range movies {
		cineId := cinebase.NodeContent("@id", m)[14:]
		titulo := cinebase.NodeContent("parent::*//a[@class='peliculaCartelera']", m)
		pT := parseTitle(titulo)
		
		t := time.Now().Format("20060102")

		row := models.NewShowtime(
			cineId, //cineID
			pT.Title,
			pT.Subtitle,
			pT.Other,
			t,
			"",
			"00:00",
		)

		hours, err := m.Search("parent::*/parent::*//*[contains(@class,'horariosCartelera')]")
		if err != nil {
			fmt.Println(err)
		}

		var horas []cinebase.TimeLinks
		for _, e := range hours {
			t, _ := time.Parse("3:04pm", e.Content())
			log.Printf("e.Attributes(): %#v \n",e.Attributes() )
			horas = append(horas,
				cinebase.TimeLinks{
					t.Format("15:04"),
					e.Attributes()["href"].Content(),
				})
		}

		//if row != nil {
		for _, h := range horas {
			row.Time = h.T
			row.BuyLink = h.BuyLink
			res = append(res, row)
		}
		//}

	}
	return
}

func ExtractMovies()(results []models.Movie){

	xNodes := cinebase.NodesExtractor("http://cinepolis.com/index.asp", "id('ctl00_ddlPelicula')/option")

	for _, x := range xNodes {
		rawTitle := x.Content()
		pT := parseTitle(rawTitle)
		results = append(results,models.NewMovieWithTitle(pT.Title))
	}
	return
}

type parsedTitle struct {
	Title string
	Subtitle string
	Other string
}
func parseTitle(complexTitle string)(parsed parsedTitle){
	complexTitle = strings.Replace(complexTitle,":"," ", -1)

	languages := []string{"SUB","ESP"}
	roomTypes := []string{"3D","4DX","IMAX","XE","DIG",}
	parts := strings.Fields(strings.ToUpper(complexTitle))

	langPart := ""
	if stringInSlice(parts[len(parts)-1],languages){
		langPart = parts[len(parts)-1]
		parts = parts[:len(parts)-1]
	}
	if langPart == "SUB" {
		parsed.Subtitle = "SUBTITULADA"
	} else {
		parsed.Subtitle  = "ESPAÑOL"
	}

	var arrTitle []string
	i := -1
	found:=true
	for found && i < len(parts)-1 {		
		i=i+1
		if stringInSlice(parts[i],roomTypes) {
			parsed.Other = strings.Join(parts[i:]," ")
			found = false
		}else{
			arrTitle=append(arrTitle,parts[i])
		}
	}
	parsed.Title = strings.Join(arrTitle," ")

	return
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
