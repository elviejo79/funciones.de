package main

import (
	"./cinepolis"
	"fmt"
	"funciones.de/app/models"
	"github.com/jmcvetta/neoism"
	"log"
)

func cleanup(db *neoism.Database) {
	qs := []*neoism.CypherQuery{
		&neoism.CypherQuery{
			Statement: `START r=rel(*) DELETE r`,
		},
		&neoism.CypherQuery{
			Statement: `START n=node(*) DELETE n`,
		},
	}
	err := db.CypherBatch(qs)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {

	cleanup(models.GlobalDb())

	company := models.NewCompany("http://cinepolis.com", "cinepolis")
	company.Save()

	cities := cinepolis.ExtractCities(company.Url)

	for _, c := range cities {
		c.Save()
		company.Node().Relate("operates_in", c.Node().Id(), neoism.Props{})
		CreateTheaters(c)
	}

}

///////////////////////////////////////// anidados
func CreateTheaters(newCity models.City) {
	urlTheater := fmt.Sprintf("http://cinepolis.com/_CARTELERA/cartelera.aspx?ic=%d", newCity.IdCity)
	theaters := cinepolis.ExtractTheaters(urlTheater)
	fmt.Printf("theatre: %#v\n",theaters)
	for _, t := range theaters {
		t.Save()
		newCity.Node().Relate("has_a", t.Node().Id(), neoism.Props{})

		CreateShowtimes(newCity, t)
	}
}

func CreateShowtimes(newCity models.City, newTheater models.Theater) {
	urlTheater := fmt.Sprintf("http://cinepolis.com/_CARTELERA/cartelera.aspx?ic=%d", newCity.IdCity)
	showtimes, _ := cinepolis.ExtractMovies(urlTheater, newTheater.IdTheater)


	for _, s := range showtimes {
		s.Save()
		newTheater.Node().Relate("presents_at", s.Node().Id(), neoism.Props{})
		movie := models.NewMovie(s.IdMovie)
		movie.Save()
		s.Node().Relate("exhibits", movie.Node().Id(), neoism.Props{})
	}
}
