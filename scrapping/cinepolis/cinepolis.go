package cinepolis

import (
	"../cinebase/"
	"fmt"
	"funciones.de/app/models"
	"strconv"
	"strings"
	"time"
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
		subtitulos := titulo[len(titulo)-3:]
		if subtitulos == "Sub" {
			subtitulos = "SUBTITULADA"
		} else {
			subtitulos = "ESPAÑOL"
		}

		sala := titulo[len(titulo)-7 : len(titulo)-4]
		if strings.Contains(titulo, " 4D") {
			titulo = titulo[:strings.Index(titulo, " 4D")]
			sala = "4D"
		} else if strings.Contains(titulo, " 3D ") {
			titulo = titulo[:strings.Index(titulo, " 3D ")]
			sala = "3D"
		} else if strings.Contains(titulo, " IMAX") {
			titulo = titulo[:strings.Index(titulo, " IMAX")]
			sala = "IMAX"
		} else if strings.Contains(titulo, " XE") {
			titulo = titulo[:strings.Index(titulo, " XE ")]
			sala = "XE"
		} else if strings.Contains(titulo, " Dig ") {
			titulo = titulo[:strings.Index(titulo, " Dig ")]
			sala = "Dig"
		} else {
			//titulo=titulo
			sala = ""
		}

		titulo = strings.ToUpper(titulo)
		titulo = strings.Replace(titulo, ":", "", -1)
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

		hours, err := m.Search("parent::*/parent::*//*[contains(@class,'horariosCartelera')]")
		if err != nil {
			fmt.Println(err)
		}

		var horas []cinebase.TimeLinks
		for _, e := range hours {
			t, _ := time.Parse("3:04pm", e.Content())
			horas = append(horas,
				cinebase.TimeLinks{
					t.Format("15:04"),
					"http://buySomeTicket.com??",//e.Attributes()["href"].Content(),
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
