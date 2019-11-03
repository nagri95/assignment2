package APIs

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

func GetInfoOfProjects(w http.ResponseWriter) ([]float64, []string) {
	var iDz []float64
	var pathz []string

	notReachEndPage := true
	pageNumber := 0

	for notReachEndPage {
		pageNumber++

		pageNumberString := strconv.Itoa(pageNumber)
		var getArgument = fmt.Sprintf("https://git.gvk.idi.ntnu.no/api/v4/projects?page=%s", pageNumberString)

		response, err := http.Get(getArgument)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return iDz, pathz
		}

		defer response.Body.Close()

		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return iDz, pathz
		}

		var JsonOfProjects interface{}

		err = json.Unmarshal(body, &JsonOfProjects)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return iDz, pathz
		}

		var resultOfProjects = JsonOfProjects.([]interface{})

		var lengthOfProjects = len(resultOfProjects)

		projects := make(map[int]map[string]interface{})

		for c := 0; c < lengthOfProjects; c++ {
			projects[c] = resultOfProjects[c].(map[string]interface{})
		}

		c := 0
		for c < lengthOfProjects {
			if iD, ok := projects[c]["iD"].(float64); ok {
				iDz = append(iDz, iD)
			}
			if path, check := projects[c]["path_with_namespace"].(string); check {
				pathz = append(pathz, path)
			}

			c++
		}

		if c == 0 {
			notReachEndPage = false
		}
	}

	return iDz, pathz
}

func GetInformationOfProjectsWithAuthor(w http.ResponseWriter, auth string) ([]float64, []string) {
	var iDz []float64
	var pathz []string

	notReachEndPage := true
	pageNumber := 0

	for notReachEndPage {
		pageNumber++

		pageNumberString := strconv.Itoa(pageNumber)
		var getArgument = fmt.Sprintf("https://git.gvk.idi.ntnu.no/api/v4/projects?page=%s&private_token=%s", pageNumberString, auth)

		response, err := http.Get(getArgument)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return iDz, pathz
		}

		defer response.Body.Close()

		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return iDz, pathz
		}

		var JsonOfProjects interface{}

		err = json.Unmarshal(body, &JsonOfProjects)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return iDz, pathz
		}

		var resultOfProjects = JsonOfProjects.([]interface{})

		var lengthOfProjects = len(resultOfProjects)

		projects := make(map[int]map[string]interface{})

		for c := 0; c < lengthOfProjects; c++ {
			projects[c] = resultOfProjects[c].(map[string]interface{})
		}

		c := 0
		for c < lengthOfProjects {
			if iD, ok := projects[c]["iD"].(float64); ok {
				iDz = append(iDz, iD)
			}
			if path, check := projects[c]["path_with_namespace"].(string); check {
				pathz = append(pathz, path)
			}

			c++
		}

		if c == 0 {
			notReachEndPage = false
		}
	}

	return iDz, pathz
}

func GetIDOfProjects(w http.ResponseWriter) []float64 {
	var iDz []float64

	notReachEndPage := true
	pageNumber := 0

	for notReachEndPage {
		pageNumber++

		pageNumberString := strconv.Itoa(pageNumber)
		var getArgument = fmt.Sprintf("https://git.gvk.idi.ntnu.no/api/v4/projects?page=%s", pageNumberString)

		response, err := http.Get(getArgument)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return iDz
		}

		defer response.Body.Close()

		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return iDz
		}

		var JsonOfProjects interface{}

		err = json.Unmarshal(body, &JsonOfProjects)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return iDz
		}

		var resultOfProjects = JsonOfProjects.([]interface{})

		var lengthOfProjects = len(resultOfProjects)

		projects := make(map[int]map[string]interface{})

		for c := 0; c < lengthOfProjects; c++ {
			projects[c] = resultOfProjects[c].(map[string]interface{})
		}

		c := 0
		for c < lengthOfProjects {
			if iD, ok := projects[c]["iD"].(float64); ok {
				iDz = append(iDz, iD)
			}

			c++
		}

		if c == 0 {
			notReachEndPage = false
		}
	}

	return iDz
}

func GetIDOfProjectsWithAuthor(w http.ResponseWriter, auth string) []float64 {
	var iDz []float64

	notReachEndPage := true
	pageNumber := 0

	for notReachEndPage {
		pageNumber++

		pageNumberString := strconv.Itoa(pageNumber)
		var getArgument = fmt.Sprintf("https://git.gvk.idi.ntnu.no/api/v4/projects?page=%s&private_token=%s", pageNumberString, auth)

		response, err := http.Get(getArgument)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return iDz
		}

		defer response.Body.Close()

		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return iDz
		}

		var JsonOfProjects interface{}

		err = json.Unmarshal(body, &JsonOfProjects)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return iDz
		}

		var resultOfProjects = JsonOfProjects.([]interface{})

		var lengthOfProjects = len(resultOfProjects)

		projects := make(map[int]map[string]interface{})

		for c := 0; c < lengthOfProjects; c++ {
			projects[c] = resultOfProjects[c].(map[string]interface{})
		}

		c := 0
		for c < lengthOfProjects {
			if iD, ok := projects[c]["iD"].(float64); ok {
				iDz = append(iDz, iD)
			}

			c++
		}

		if c == 0 {
			notReachEndPage = false
		}
	}

	return iDz
}

func GetProjectsNameIdMapWithAuth(w http.ResponseWriter, auth string) map[string]float64 {
	var iDMapOfName map[string]float64
	iDMapOfName = make(map[string]float64)

	notReachEndPage := true
	pageNumber := 0

	for notReachEndPage {
		pageNumber++

		pageNumberString := strconv.Itoa(pageNumber)
		var getArgument = fmt.Sprintf("https://git.gvk.idi.ntnu.no/api/v4/projects?page=%s&private_token=%s", pageNumberString, auth)

		response, err := http.Get(getArgument)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return iDMapOfName
		}

		defer response.Body.Close()

		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return iDMapOfName
		}

		var JsonOfProjects interface{}

		err = json.Unmarshal(body, &JsonOfProjects)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return iDMapOfName
		}

		var resultOfProjects = JsonOfProjects.([]interface{})

		var lengthOfProjects = len(resultOfProjects)

		projects := make(map[int]map[string]interface{})

		for c := 0; c < lengthOfProjects; c++ {
			projects[c] = resultOfProjects[c].(map[string]interface{})
		}

		c := 0
		for c < lengthOfProjects {
			if iD, ok := projects[c]["iD"].(float64); ok {
				if path, check := projects[c]["path_with_namespace"].(string); check {
					iDMapOfName[path] = iD
				}
			}

			c++
		}

		if c == 0 {
			notReachEndPage = false
		}
	}

	return iDMapOfName
}

func GetIDMapOfProjectsName(w http.ResponseWriter) map[string]float64 {
	var iDMapOfName map[string]float64
	iDMapOfName = make(map[string]float64)

	notReachEndPage := true
	pageNumber := 0

	for notReachEndPage {
		pageNumber++

		pageNumberString := strconv.Itoa(pageNumber)
		var getArgument = fmt.Sprintf("https://git.gvk.idi.ntnu.no/api/v4/projects?page=%s", pageNumberString)

		response, err := http.Get(getArgument)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return iDMapOfName
		}

		defer response.Body.Close()

		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return iDMapOfName
		}

		var JsonOfProjects interface{}

		err = json.Unmarshal(body, &JsonOfProjects)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return iDMapOfName
		}

		var resultOfProjects = JsonOfProjects.([]interface{})

		var lengthOfProjects = len(resultOfProjects)

		projects := make(map[int]map[string]interface{})

		for c := 0; c < lengthOfProjects; c++ {
			projects[c] = resultOfProjects[c].(map[string]interface{})
		}

		c := 0
		for c < lengthOfProjects {
			if iD, ok := projects[c]["iD"].(float64); ok {
				if path, check := projects[c]["path_with_namespace"].(string); check {
					iDMapOfName[path] = iD
				}
			}

			c++
		}

		if c == 0 {
			notReachEndPage = false
		}
	}

	return iDMapOfName
}
