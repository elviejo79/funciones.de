package controllers

import (
	"github.com/robfig/revel"
	"funciones.de/app/models"
	"runtime/pprof"
	"os"
)

type App struct {
	*revel.Controller
}

func (c App) TheatersByCity(cityName string) revel.Result {
	//profiling
	f, _ := os.Create("my_profile.file")
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()
	city:= models.CityByName(cityName)


	theaters := models.TheatersByCity(city)
	movies := models.MoviesByTheaters(theaters)
	showtimes := models.ShowtimesForTheaters(theaters)
	return c.Render(city,theaters,movies,showtimes)
}


func (c App) Buy(showtimeKey string) revel.Result {
	s:=models.ShowtimeByKey(showtimeKey)
	return c.Render(s)
}

func (c App) Index(strCity string) revel.Result {
	cities := models.AllCities()
	return c.Render(cities)
}
