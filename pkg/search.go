package pkg

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
)

func HandleSearch(w http.ResponseWriter, r *http.Request) {

	search := r.FormValue("search")

	Artsearch := []Artist{}

	tabId := []int{}
	temp := template.Must(template.ParseFiles("templates/resultsearch.html", "templates/navbar.html", "templates/form.html"))
	var err error
	if !CheckURL(tab_Location, search) {
		MyResult := ResultLocation(search)
		err = temp.Execute(w, MyResult)
	} else {

		for _, Art := range Artist_Data {

			if strings.ToUpper(search) == strings.ToUpper(Art.Aname) {
				tabId = append(tabId, Art.Aid)
			} else if search == strconv.Itoa(Art.Acread) {
				tabId = append(tabId, Art.Aid)
			} else if search == Art.Afalbum {
				tabId = append(tabId, Art.Aid)
			} else {
				for _, v := range Art.Amember {
					if strings.ToUpper(v) == strings.ToUpper(search) {
						tabId = append(tabId, Art.Aid)
					}
				}
			}
		}
		for _, id := range tabId {
			Artsearch = append(Artsearch, Artist_Data[id-1])
		}

		if len(Artsearch) == 0 {
			NoResultSearch := Filter{
				Boul:      true,
				NoResult:  "NO RESULTS FOR THE INFORMATION ENTERED",
				LocatFilt: tab_Location,
			}
			err = temp.Execute(w, NoResultSearch)
		} else {
			ResultSearch := Filter{
				Artboul:   true,
				Artists:   Artsearch,
				LocatFilt: tab_Location,
			}
			err = temp.Execute(w, ResultSearch)
		}
	}

	if err != nil {
		fmt.Println("Erreur lors de l'execution du template", err)
		Error500Handler(w, r)
		return
	}
}
