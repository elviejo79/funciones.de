package cinemex

import (
	"fmt"
	"funciones.de/app/models"
	"encoding/json"
	"strings"
	"time"
	"log"
	"strconv"
	"bytes"
	"../cinebase/"
)

var Company models.Company = models.NewCompany("http://cinemex.com", "cinemex")

func ExtractCities() (cities []models.City) {
	as, _ := extractJsonArea()
	for _,a := range as {
		cities= append(cities, models.NewCity(a.Id,a.Name))
	}
	return cities
}

func ExtractTheaters(c models.City) ( theaters []models.Theater) {
	as,_ := extractJsonArea()
	for _,a := range as {
		if a.Id == c.IdCity {
			for _,c := range a.Cinemas {
				theaters = append(theaters,models.NewTheater(c.Id,c.Name))
			}
		}
	}
	return
}

func ExtractShowtimes(c models.City, t models.Theater) (res []models.Showtime, err error) {
	url := fmt.Sprintf("http://cinemex.com/partials/sidebarCinema/%d/date-%s",t.IdTheater,time.Now().Format("20060102"))
	movies := cinebase.NodesExtractor(url,"id('sidebar-mycinema')/div/div")
	for _, m := range movies{
		cineId := strconv.Itoa(t.IdTheater)
		titulo := cinebase.NodeContent("a", m)
		subtitulos := cinebase.NodeContent("div/div/span[1]",m) 
		if subtitulos == "Inglés" {
			subtitulos = "SUBTITULADA"
		} 
		subtitulos = strings.ToUpper(subtitulos)

		sala := cinebase.NodeContent("div/div/span[2]",m) 
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
			"00:00",
		)


		hours, err := m.Search("div/div/a")
		if err != nil {
			log.Fatal(err)
		}

		horas := []string{}
		for _, e := range hours {
			t, _ := time.Parse("3:04 PM", e.Content())
			horas = append(horas, t.Format("15:04"))
		}

		for _, h := range horas {
			row.Time = h
			res = append(res, row)
		}
	}

	return
}

func extractJsonArea() (as []area, err error) {
	html, err := cinebase.GetBody(Company.Url)
	if err != nil {
		log.Fatal(" %#v", err)
	}
	begin_area := []byte("var areas       = [{") 
	html = html[bytes.Index(html,begin_area)+len(begin_area)-2:]
	html = html[:bytes.Index(html,[]byte("}}]}];"))+5]

	err = json.Unmarshal(html, &as)
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

