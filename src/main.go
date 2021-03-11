package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

const baseURL = "https://groupietrackers.herokuapp.com/api"

type Page struct {
	Title string
	Body  []byte
}
type MyArtist struct {
	ID           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
	Locations    []string `json:"locations"`
	ConcertDates []string `json:"concertDates"`
	Relations    string   `json:"relations"`
}
type MyDate struct {
	ID    int      `json:"id"`
	Dates []string `json:"dates"`
}
type MyDatess struct {
	Index []MyDate `json:"index"`
}
type MyRelations struct {
	Index []MyRelation `json:"index"`
}
type MyRelation struct {
	ID             int                 `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}

type DatesLocations struct {
	Id             int
	DatesLocations map[string][]string
}
type MyLocation struct {
	ID        int      `json:"id"`
	Locations []string `json:"locations"`
	Dates     []string `json:"dates"`
}
type MyLocations struct {
	Index []MyLocation `json:"index"`
}
type MyDates struct {
	ID    int      `json:"id"`
	Dates []string `json:"dates"`
}
type MyArtistFull struct {
	ID             int                 `json:"id"`
	Image          string              `json:"image"`
	Name           string              `json:"name"`
	Members        []string            `json:"members"`
	CreationDate   int                 `json:"creationDate"`
	FirstAlbum     string              `json:"firstAlbum"`
	Locations      []string            `json:"locations"`
	ConcertDates   []string            `json:"concertDates"`
	DatesLocations map[string][]string `json:"datesLocations"`
	Location       string
}

var ArtistFull []MyArtistFull
var Artists []MyArtist
var Dates MyDatess
var Locations MyLocations
var Relations MyRelations

func (p *Page) save() error {
	filename := p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}
func GetFullDataById(id int) (MyArtistFull, error) {
	for _, artist := range ArtistFull {
		if artist.ID == id {
			return artist, nil
		}
	}
	return MyArtistFull{}, errors.New("Not found")
}
func GetDateByID(id int) (MyDate, error) {
	for _, date := range Dates.Index {
		if date.ID == id {
			return date, nil
		}
	}
	return MyDate{}, errors.New("Not found")
}
func GetArtistByID(id int) (MyArtist, error) {
	for _, artist := range Artists {

		if artist.ID == id {
			return artist, nil
		}
	}
	return MyArtist{}, errors.New("Not found")
}
func GetArtistsData() error {
	resp, err := http.Get(baseURL + "/artists")
	if err != nil {
		return errors.New("Error by get")
	}
	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.New("Error by ReadAll")
	}
	json.Unmarshal(bytes, &Artists)
	return nil
}

func GetDatesData() error {

	resp, err := http.Get(baseURL + "/dates")
	if err != nil {
		return errors.New("Error by get")
	}
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.New("Error by ReadAll")
	}
	json.Unmarshal(bytes, &Dates)
	return nil
}

func GetLocationsData() error {
	resp, err := http.Get(baseURL + "/locations")
	if err != nil {
		return errors.New("Error by get")
	}
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.New("Error by ReadAll")
	}
	json.Unmarshal(bytes, &Locations)
	return nil
}

func GetRelationsData() error {
	resp, err := http.Get(baseURL + "/relation")
	if err != nil {
		return errors.New("Error by get")
	}
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.New("Error by ReadAll")
	}
	json.Unmarshal(bytes, &Relations)
	return nil
}
func GetData() error {

	if len(ArtistFull) != 0 {
		return nil
	}
	err1 := GetArtistsData()
	err2 := GetLocationsData()
	err3 := GetDatesData()
	err4 := GetRelationsData()
	if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
		return errors.New("Error by get data artists, locations, dates")
	}
	for i := range Artists {
		var tmpl MyArtistFull
		tmpl.ID = i + 1
		tmpl.Image = Artists[i].Image
		tmpl.Name = Artists[i].Name
		tmpl.Members = Artists[i].Members
		tmpl.CreationDate = Artists[i].CreationDate
		tmpl.FirstAlbum = Artists[i].FirstAlbum
		tmpl.Locations = Locations.Index[i].Locations
		tmpl.ConcertDates = Dates.Index[i].Dates
		tmpl.DatesLocations = Relations.Index[i].DatesLocations
		ArtistFull = append(ArtistFull, tmpl)
	}

	return nil
}

func GetRelationsDataById(id int) DatesLocations {
	response, err := http.Get("https://groupietrackers.herokuapp.com/api/relation/" + strconv.Itoa(id))
	if err != nil {
		fmt.Printf("The http request failed with error %s\n", err)

	}
	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	var data DatesLocations

	json.Unmarshal(responseData, &ArtistFull)
	return data

}
func GetLocationsDataById(id int) MyLocation {
	response, err := http.Get("https://groupietrackers.herokuapp.com/api/locations/" + strconv.Itoa(id))
	if err != nil {
		fmt.Printf("The http request failed with error %s\n", err)

	}
	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	var data MyLocation

	json.Unmarshal(responseData, &data)
	return data

}
func GetLocationByID(id int) (MyLocation, error) {
	for _, location := range Locations.Index {
		if location.ID == id {
			return location, nil
		}
	}
	return MyLocation{}, errors.New("Not found")
}
func GetDatesById(id int) MyDates {
	response, err := http.Get("https://groupietrackers.herokuapp.com/api/dates/" + strconv.Itoa(id))
	if err != nil {
		fmt.Printf("The http request failed with error %s\n", err)

	}
	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	var data MyDates

	json.Unmarshal(responseData, &data)
	return data

}
func GetArtistDataById(id int) MyArtist {
	response, err := http.Get("https://groupietrackers.herokuapp.com/api/artists/" + strconv.Itoa(id))
	if err != nil {
		fmt.Printf("The http request failed with error %s\n", err)

	}
	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	var data MyArtist

	json.Unmarshal(responseData, &data)

	return data

}
func makeSearch(search string) []MyArtistFull {
	if search == "" {
		return ArtistFull
	}
	art, err := GetAllArtistStringData() // renvoie toutes les données de tous les artistes(name, memezbers, dates, creation date...)
	if err != nil {
		errors.New("Error by converter")
	}
	var search_artist []MyArtistFull

	for i, artist := range art { //art est un tablmeau qui rpz tous les données des artistes et i est l'index du tableau
		//On boucle sur charque artiste du tableau(un artiste = 315 lettres)
		lower_band := strings.ToLower(artist) //idem que art sans capitals
		// lower_band est chaque itération d'artiste ainsi que tous ses données
		for i_name, l_name := range []byte(lower_band) { //  i_name -> indice du tab et l_name -> données artistes en bytes
			lower_search := strings.ToLower(search) // on decapitalize la recherche | et on peut ondexer lower_search
			// fmt.Println("i_name", i_name)
			// fmt.Println("l_name", l_name)
			if lower_search[0] == l_name { // si le premier caract de search == à la 1ere lettre du nom l'artsite
				// fmt.Println("l_name", l_name)
				lenght_name := 0
				indx := i_name
				for _, l := range []byte(lower_search) {
					if l == lower_band[indx] {
						// fmt.Println("L", l, "lowerbandindx", lower_band[indx])
						// fmt.Println(len(lower_band))
						if indx+1 == len(lower_band) {
							break
						}
						indx++
						lenght_name++
					} else {
						break
					}
				}

				if len(search) == lenght_name {
					band, _ := GetFullDataById(i + 1)
					search_artist = append(search_artist, band)
					break
				}
			}
		}

	}

	return search_artist
}

func HomePage(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/"):] //Chemin de la page
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	t, _ := template.ParseFiles("index.html")
	t.Execute(w, p)
}

func MenuPage(w http.ResponseWriter, r *http.Request) {
	var Datadata []MyArtistFull

	err := GetData()
	search := r.FormValue("search")

	Datadata = makeSearch("")
	if !(search == "" && len(ArtistFull) != 0) {
		Datadata = makeSearch(search)
	}

	res, err := http.Get("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Bad Request", 400)
	}
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Internal Server Error", 500)
	}
	var Data []MyArtistFull

	json.Unmarshal(data, &Data)

	ts, err := template.ParseFiles("menu.html")
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Internal Server Error : check if template file exists", 500)
	}

	if err := ts.Execute(w, Datadata); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

}

func ConcertPage(w http.ResponseWriter, r *http.Request) {

	idStr := r.FormValue("id")
	id, _ := strconv.Atoi(idStr)
	relations := GetRelationsDataById(id)
	artist := GetArtistDataById(id)
	locations := GetLocationsDataById(id)
	concertDates := GetDatesById(id)
	loc := locations.Locations[0]
	var fullArtist MyArtistFull

	fullArtist.ID = artist.ID
	fullArtist.Image = artist.Image
	fullArtist.Name = artist.Name
	fullArtist.Members = artist.Members
	fullArtist.CreationDate = artist.CreationDate
	fullArtist.FirstAlbum = artist.FirstAlbum
	fullArtist.Locations = locations.Locations
	fullArtist.ConcertDates = concertDates.Dates
	fullArtist.DatesLocations = relations.DatesLocations
	fullArtist.Location = loc

	for key, value := range fullArtist.Locations {

		fmt.Println(key)
		fmt.Println(value)

	}

	tmpl, err := template.ParseFiles("pageConcert.html")
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	if err := tmpl.Execute(w, fullArtist); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

}

func GetAllArtistStringData() ([]string, error) { // recuperer tous les données
	var data []string

	for i := 1; i <= len(Artists); i++ {
		artist, err1 := GetArtistByID(i)
		locations, err2 := GetLocationByID(i)
		date, err3 := GetDateByID(i)
		if err1 != nil || err2 != nil || err3 != nil {
			return data, errors.New("Error by converter")
		}

		str := artist.Name + " "
		for _, member := range artist.Members {
			str += member + " "
		}
		str += strconv.Itoa(artist.CreationDate) + " "
		str += artist.FirstAlbum + " "
		for _, location := range locations.Locations {
			str += location + " "
		}
		for _, d := range date.Dates {
			str += d + " "
		}
		data = append(data, str)

	}

	return data, nil
}

func main() {
	response, err := http.Get("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		fmt.Printf("The http request failed with error %s\n", err)

	}
	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	var data []MyArtist
	json.Unmarshal(responseData, &data)

	http.HandleFunc("/", HomePage)

	http.HandleFunc("/Menu", MenuPage)
	http.HandleFunc("/ConcertInfos", ConcertPage)

	fmt.Printf("Servor is running on 8050 port")

	log.Fatal(http.ListenAndServe(":8050", nil))
}
