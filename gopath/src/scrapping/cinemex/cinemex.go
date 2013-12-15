package cinemex

import (
	"../cinebase/"
	"bytes"
	"encoding/json"
	"fmt"
	"funciones.de/app/models"
	"log"
	"strconv"
	"strings"
	"time"
)

var Company models.Company = models.NewCompany("http://cinemex.com", "cinemex")

func ExtractCities() (cities []models.City) {
	as, _ := extractJsonArea()
	for _, a := range as {
		cities = append(cities, models.NewCity(a.Id, a.Name))
	}
	return cities
}

func ExtractTheaters(c models.City) (theaters []models.Theater) {
	as, _ := extractJsonArea()
	for _, a := range as {
		if a.Id == c.IdCity {
			for _, c := range a.Cinemas {
				theaters = append(theaters, models.NewTheater(c.Id, c.Name))
			}
		}
	}
	return
}

func ExtractShowtimes(c models.City, t models.Theater) (res []models.Showtime, err error) {
	log.Printf("http://cinemex.com/partials/sidebarCinema/%d/date-%s", t.IdTheater, time.Now().Add(2*time.Hour).Format("20060102"))
	url := fmt.Sprintf("http://cinemex.com/partials/sidebarCinema/%d/date-%s", t.IdTheater, time.Now().Add(2*time.Hour).Format("20060102"))
	movies := cinebase.NodesExtractor(url, "id('sidebar-mycinema')/div/div")
	for _, m := range movies {
		cineId := strconv.Itoa(t.IdTheater)
		titulo := cinebase.NodeContent("a", m)
		subtitulos := cinebase.NodeContent("div/div/span[1]", m)
		if subtitulos == "Inglés" {
			subtitulos = "SUBTITULADA"
		}
		subtitulos = strings.ToUpper(subtitulos)

		sala := cinebase.NodeContent("div/div/span[2]", m)
		if strings.Contains(sala, "Digital") {
			sala = "Dig"
		}

		titulo = strings.Replace(strings.ToUpper(titulo), ":", "", -1)
		t := time.Now().Format("20060102")

		row := models.NewShowtime(
			cineId, //cineID
			titulo,
			subtitulos,
			sala,
			t,
			"",
			"00:00",
		)

		hours, err := m.Search("div/div/a")
		if err != nil {
			log.Fatal(err)
		}

		var horas []cinebase.TimeLinks
		for _, e := range hours {
			t, _ := time.Parse("3:04 PM", e.Content())
			horas = append(horas,
				cinebase.TimeLinks{
					t.Format("15:04"),
					e.Attributes()["href"].Content(),
				})

			for _, h := range horas {
				row.Time = h.T
				row.BuyLink = h.BuyLink
				res = append(res, row)
			}
		}
	}
	return
}
func ExtractMovies()(results []models.Movie){
	var tmp []movie
	html := extractJson([]byte("var movies = [{"),[]byte("}];"))
	json.Unmarshal(html,&tmp)
	for _,m := range tmp {
		results = append(models.NewMovie(
			m.Name,
			m.Cover,
			m.Info.Country,
			m.Info.Director,
			m.Info.Genre[0],
			m.Info.Original_title,
			m.Info.Sinopsis,
			m.Info.Trailer,
			m.Info.Duration,
			m.Info.Year,
		))
	}
	return
}

func extractJsonArea() (as []area, err error) {
	html := extractJson([]byte("var areas       = [{"), []byte("}}]}];"))
	err = json.Unmarshal(html, &as)
	return
}

func extractJson(begin_area, end_area []byte)(html []byte){
	html, err := cinebase.GetBody(Company.Url)
	if err != nil {
		log.Fatal(" %#v", err)
	}
	html = html[bytes.Index(html, begin_area)+len(begin_area)-2:]
	html = html[:bytes.Index(html,end_area)+len(end_area)-1]
	return
}

type cinema struct {
	Area_id  int
	Id       int
	Info     interface{}
	Lat      float32
	Lng      float32
	Name     string
	Platinum bool
	State_id int
}

type area struct {
	Cinemas  []cinema
	Id       int
	Name     string
	State_id int
}
type info struct {
	Cast string
	Country string
	Director string
	Duration string
	Genre []string
	Original_title string
	Rating string
	Sinopsis string
	Trailer string
	Year string
}
type movie struct{
	Attributes []interface{}
	Cover string
	Data string
	Id int
	Info info
	Name string
	Score float32
	Type []string
	Url string
	Versions []interface{}
	Youtube_id string
}
