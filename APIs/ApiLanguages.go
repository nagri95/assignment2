package APIs

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

type LanguageInfoStruct struct {
	Name  string
	Count int
}

// Sorting by counting of languages number
type ByCountingLanguages []LanguageInfoStruct

func (a ByCountingLanguages) Len() int           { return len(a) }
func (a ByCountingLanguages) Less(i, j int) bool { return a[i].Count > a[j].Count } // which is More for our purpose
func (a ByCountingLanguages) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

type LanguagesResponseStruct struct {
	Languages []string `json:"languages"`
	Author      bool     `json:"author"`
}

func LanguagesHandler(w http.ResponseWriter, r *http.Request) {
	var LanguagesResponse = LanguagesResponseStruct{}

	http.Header.Add(w.Header(), "content-type", "application/json")
	parts := strings.Split(r.URL.Path, "/")

	if len(parts) != 4 || parts[3] != "languages" {
		http.Error(w, "Malformed URL", http.StatusBadRequest)
		return
	}

	requestLimit := 5

	//  Define limit number if it exists
	limit, check1 := r.URL.Query()["limit"]

	if check1 {
		limit, err := strconv.Atoi(limit[0])

		if err != nil {
			http.Error(w, "Internal error", http.StatusInternalServerError)
			return
		}

		requestLimit = limit
	}

	authorRequest := ""
	withAuthor := false

	//authentification if found
	author, check2 := r.URL.Query()["author"]

	if check2 {
		authorRequest = author[0]
		withAuthor = true
	}

	// do the webhook call
	var params []string
	if check1 {
		params = append(params, "limit")
	}
	if withAuthor {
		params = append(params, "authentification")
	}

	WebhookCheckfunc(w, "languages", params)


	// get the payload
	var projectsName []string
	err := json.NewDecoder(r.Body).Decode(&projectsName)

	// determine the ids we will work with
	var projectsID []float64
	if err != nil { // in case we have no payload
		if withAuthor {
			projectsID = GetIDOfProjectsWithAuthor(w, authorRequest)
		} else {
			projectsID = GetIDOfProjects(w)
		}
	} else { // in case we have payload
		var  mapOfNameID map[string]float64
		 mapOfNameID = make(map[string]float64)

		if withAuthor {
			 mapOfNameID = GetProjectsNameIdMapWithAuth(w, authorRequest)
		} else {
			 mapOfNameID = GetIDMapOfProjectsName(w)
		}

		for i := 0; i < len(projectsName); i++ {
			projectsID = append(projectsID,  mapOfNameID[projectsName[i]])
		}
	}

	var languageCounterMap map[string]int
	languageCounterMap = make(map[string]int)

	// update languageCounterMap by languages used
	projectsLengthID := len(projectsID)
	for i := 0; i < projectsLengthID; i++ {
		id := strconv.Itoa(int(projectsID[i]))

		var getArgument string
		if withAuthor {
			getArgument = fmt.Sprintf("https://git.gvk.idi.ntnu.no/api/v4/projects/%s/languages?private_token=%s", id, authorRequest)
		} else {
			getArgument = fmt.Sprintf("https://git.gvk.idi.ntnu.no/api/v4/projects/%s/languages", id)
		}

		resp, err := http.Get(getArgument)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var jsonOfLanguages interface{}

		err = json.Unmarshal(body, &jsonOfLanguages)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		mapOfLanguages, ok := jsonOfLanguages.(map[string]interface{})

		if ok {
			for lang := range mapOfLanguages {
				languageCounterMap[lang] += languageCounterMap[lang] + 1
			}
		} else {
			continue
		}
	}

	// determinante the new limit
	if requestLimit > len(languageCounterMap) {
		requestLimit = len(languageCounterMap)
	}

	var infoLanguages []LanguageInfoStruct

	for lang := range languageCounterMap {
		infoLanguages = append(infoLanguages, LanguageInfoStruct{lang, languageCounterMap[lang]})
	}

	sort.Sort(ByCountingLanguages(infoLanguages))

	for i := 0; i < requestLimit; i++ {
		LanguagesResponse.Languages = append(LanguagesResponse.Languages, infoLanguages[i].Name)
	}

	// encoding the answer
	LanguagesResponse.Author = withAuthor
	json.NewEncoder(w).Encode(LanguagesResponse)
}
