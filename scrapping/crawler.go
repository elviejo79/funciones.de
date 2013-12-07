package main

import (
	"./cinepolis"
	"fmt"
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

	root, err := neoism.Connect("http://localhost:7474/db/data")
	if err != nil {
		log.Fatal("failed to open default database", err)
	}

	cleanup(root)

	company, _ := root.CreateNode(map[string]interface{}{"url": "http://cinepolis.com", "name": "cinepolis"})
	company.AddLabel("Company")

	cities := cinepolis.ExtractCities("http://cinepolis.com")
	var newCity *neoism.Node
	for i, c := range cities {
		newCity, err = root.CreateNode(map[string]interface{}{"idCity": i, "name": c})
		if err != nil {
			log.Fatal("failed to create node:", err)
		}
		newCity.AddLabel("City")
		_, err = company.Relate("operates_in", newCity.Id(), neoism.Props{})
		if err != nil {
			log.Fatal("failed to create relationship:", err)
		}

		CreateTheaters(root, newCity)

	}


}

type Movie struct {
	NeoID int    `json:"neoId"`
	Title string `json:"n.title"`
}

///////////////////////////////////////// anidados
func CreateTheaters(root *neoism.Database, newCity *neoism.Node) {
	properties, _ := newCity.Properties()
	fmt.Printf("%#v \n", properties)
	urlTheater := fmt.Sprintf("http://cinepolis.com/_CARTELERA/cartelera.aspx?ic=%d", int(properties["idCity"].(float64)))
	theaters := cinepolis.ExtractTheaters(urlTheater)

	for it, t := range theaters {
		newTheater, err := root.CreateNode(map[string]interface{}{"idTheater": it, "name": t})
		if err != nil {
			log.Fatal("failed to create Theater:", err)
		}
		newTheater.AddLabel("Theater")
		newCity.Relate("has_a", newTheater.Id(), neoism.Props{})
		
		CreateShowtimes(root, newCity, newTheater)
	}


}

func CreateShowtimes (root *neoism.Database, newCity *neoism.Node, newTheater *neoism.Node) {
	properties, _ := newCity.Properties()
	fmt.Printf("%#v \n", properties)
	urlTheater := fmt.Sprintf("http://cinepolis.com/_CARTELERA/cartelera.aspx?ic=%d", int(properties["idCity"].(float64)))

	propertiesT, _ := newTheater.Properties()
	showtimes, err := cinepolis.ExtractMovies(urlTheater, int(propertiesT["idTheater"].(float64)))
	fmt.Printf("Showtimes: %#v \n", showtimes)
	if err != nil {
		log.Fatal("Problemas al sacar las funciones de cinepolis:", err)
	}
	for _, s := range showtimes {
		newNode, err := root.CreateNode(map[string]interface{}{
			"IdTheater": s.IdTheater,
			"IdMovie":   s.IdMovie,
			"Language":  s.Language,
			"RoomType":  s.RoomType,
			"Date":      s.Date,
			"Time":      s.Time,
		})

		if err != nil {
			log.Fatal("failed to create Showtime:", err)
		}

		newNode.AddLabel("Showtime")
		newTheater.Relate("performs_at", newNode.Id(), neoism.Props{})
		movieID, err := GetMovieID(root, s.IdMovie)
		if err != nil {
			movieID, _ = CreateMovie(root, s.IdMovie, "Movie")
		}
		newNode.Relate("presents", movieID, neoism.Props{})

	}

}

func CreateMovie(db *neoism.Database, title string, label string) (int, error) {
	newNode, err := db.CreateNode(map[string]interface{}{
		"title": title,
	})

	if err != nil {
		log.Fatal("failed to create %s:", label, err)
		return 0, err
	}
	newNode.AddLabel(label)
	return newNode.Id(), nil
}

//gets the id of a movie
func GetMovieID(db *neoism.Database, title string) (movieID int, err error) {
	//
	// Query with integer parameter and string results
	//
	var result []Movie
	cq := neoism.CypherQuery{
		Statement: "MATCH (n:Movie) where n.title = '" + title + "'  RETURN ID(n) as neoId, n.title",
		Parameters: map[string]interface{}{
			"title": title,
		},
		Result: &result,
	}


	err = db.Cypher(&cq)
	if err != nil {
		fmt.Printf("%#v", err)
		log.Fatal(err)
	}

	if len(result) > 0 {
		movieID = result[0].NeoID
	} else {
		err = fmt.Errorf("movie with title %s not found", title)
	}

	return
}
