package main

import (
	"./cinepolis"
	"./cinemex"
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
	cinepolis.Company.Save()
	cinemex.Company.Save()

	CreateCities(&cinemex.Company)
	CreateCities(&cinepolis.Company)

}

func CreateCities(co *models.Company){
	var cities []models.City

	if co.Name == "cinepolis" {
		cities = cinepolis.ExtractCities()
	} else {
		cities = cinemex.ExtractCities()
	}
	
	for _, c := range cities {
		c.Save()
		co.Node().Relate("operates_in", c.Node().Id(), neoism.Props{})
		CreateTheaters(co,c)
	}

}

///////////////////////////////////////// anidados
func CreateTheaters(co *models.Company, newCity models.City) {

	var theaters []models.Theater
	if co.Name == "cinepolis" {
		theaters = cinepolis.ExtractTheaters(newCity)
	} else {
		theaters = cinemex.ExtractTheaters(newCity)
	}


	fmt.Printf("For City %s found %d theaters\n",newCity.Name,len(theaters))
	for _, t := range theaters {
		t.Save()
		newCity.Node().Relate("has_a", t.Node().Id(), neoism.Props{})

		CreateShowtimes(co, newCity, t)
	}
}

func CreateShowtimes(co *models.Company, newCity models.City, newTheater models.Theater) {
	var showtimes []models.Showtime
	if co.Name == "cinepolis" {
		showtimes,_ = cinepolis.ExtractShowtimes(newCity, newTheater)
	} else {
		showtimes,_ = cinemex.ExtractShowtimes(newCity, newTheater)
	}

	for _, s := range showtimes {
		s.Save()
		newTheater.Node().Relate("presents_at", s.Node().Id(), neoism.Props{})
		movie := models.NewMovie(s.IdMovie)
		movie.Save()
		s.Node().Relate("exhibits", movie.Node().Id(), neoism.Props{})

		//agarcia: showtimes as relationship not nodes
		//newTheater.Node().Relate("showtime",movie.Node().Id(),models.StructToMap(s))
	}
}
