package pkg

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type Artist struct {
	Aid        int      `json:"id"`
	AImg       string   `json:"image"`
	Aname      string   `json:"name"`
	Amember    []string `json:"members"`
	Acread     int      `json:"creationDate"`
	Afalbum    string   `json:"firstAlbum"`
	Aloc       string   `json:"locations"`
	Aconcerdat string   `json:"concertDates"`
	Arelat     string   `json:"relations"`
}

type OneArtist struct {
	Id         int
	Img        string
	Name       string
	Member     []string
	CreatDat   int
	FirstAlbum string
	DatLoc     map[string][]string
	Geoloc     []Cordonnate
	ConcertDat []string
	Artists    []Artist
	LocatFilt  []string
}

var Artist_Data = GetArtistData("https://groupietrackers.herokuapp.com/api/artists")

func GetArtistData(api string) []Artist {

	content, err := http.Get(api)
	if err != nil {
		fmt.Println("Erreur de recuperation des donnees", err)
	}
	defer content.Body.Close()

	scan, er := ioutil.ReadAll(content.Body)
	if er != nil {
		fmt.Println("Erreur lors de la lecture des donnees ", er)
	}
	var Artists []Artist
	erreur := json.Unmarshal([]byte(scan), &Artists)
	if erreur != nil {
		fmt.Println("Erreur lors du decodage", erreur)
	}

	return Artists
}

func HandleArtist(w http.ResponseWriter, r *http.Request) {

	FilterArtist := Filter{
		Artists:   Artist_Data,
		LocatFilt: tab_Location,
	}
	temp := template.Must(template.ParseFiles("templates/artist.html", "templates/form.html", "templates/navbar.html"))
	err := temp.Execute(w, FilterArtist)
	if err != nil {
		fmt.Println("Erreur lors de l'execution du template", err)
	}

}

func HandleArtistDeatail(w http.ResponseWriter, r *http.Request) {

	artistID := r.URL.Query().Get("Id")
	artistid, _ := strconv.Atoi(artistID)

	artist := Artist{}
	if artistid < 1 || artistid > 52 {
		Error404Handler(w, r)
	} else {

		artist = Artist_Data[artistid-1]

		data_dat := Dat_Artist(artist.Aconcerdat)
		newdat := []string{}

		for _, v := range data_dat.DAdat {
			newdat = append(newdat, strings.ReplaceAll(v, "*", ""))
		}

		data_rel := Rel_Artist(artist.Arelat) //Recuperation des donnees de relation de l'artist
		loc := []string{}
		newrel := map[string][]string{}
		for key, val := range data_rel.RAdatloc {
			loc = append(loc, key)
			newkey := strings.ReplaceAll(key, "-", "\n")
			newrel[newkey] = val
		}

		//Recuperation des longitude latitude et date de concert pour chaque ville de l'artist
		var cordonnates []Cordonnate
		my_map := map[string][]string{}

		for city, dat := range data_rel.RAdatloc {
			stock := GetCityCoordinates(city)
			my_map[city] = dat
			small_cordonnate := Cordonnate{
				City: my_map,
				Lat:  stock.MyResults[0].MyLocation[0].MyLatLng.Lat,
				Long: stock.MyResults[0].MyLocation[0].MyLatLng.Lng,
			}

			cordonnates = append(cordonnates, small_cordonnate)
			//Reinitiialisation
			my_map = map[string][]string{}
			small_cordonnate = Cordonnate{}
		}

		NewArtist := OneArtist{
			Img:        artist.AImg,
			Name:       artist.Aname,
			Member:     artist.Amember,
			CreatDat:   artist.Acread,
			FirstAlbum: artist.Afalbum,
			DatLoc:     newrel,
			Geoloc:     cordonnates,
			ConcertDat: newdat,
			LocatFilt:  tab_Location,
		}

		temp := template.Must(template.ParseFiles("templates/artist_detail.html", "templates/form.html", "templates/navbar.html"))
		temp.Execute(w, NewArtist)
	}
}
