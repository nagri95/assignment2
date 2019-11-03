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


type ReposInfoStruct struct {
	Repository string `json:"repository"`
	Commits    int    `json:"commits"`
}

// Sorting by counting of commits number
type ByCountingCommits []ReposInfoStruct
func (x ByCountingCommits) Len() int           { return len(x) }
func (x ByCountingCommits) Less(i, j int) bool { return x[i].Commits > x[j].Commits }
func (x ByCountingCommits) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }

type CommitsAnswerStruct struct {
	Repos []ReposInfoStruct `json:"repos"`
	Author  bool        `json:"auth"`
}

func CommitsHandler(w http.ResponseWriter, r *http.Request) {
	var commitsResponse = CommitsAnswerStruct{}

	http.Header.Add(w.Header(), "content-type", "application/json")
	parts := strings.Split(r.URL.Path, "/")

	if len(parts) != 4 || parts[3] != "commits" {
		http.Error(w, "Malformed URL", http.StatusBadRequest)
		return
	}

	requestLimit := 5

	// Define limit number if it exists
	limit, check1 := r.URL.Query()["limit"]

	if check1 {
		limit, err := strconv.Atoi(limit[0])

		if err != nil {
			http.Error(w, "Internal error", http.StatusInternalServerError)
			return
		}

		requestLimit = limit
	}

	authRqst := ""
	withAuthor := false

	//authentification if found
	auth, check2 := r.URL.Query()["auth"]

	if check2 {
		authRqst = auth[0]
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

	WebhookCheckfunc(w, "commits", params)

	// get projects
	var projectsID []float64
	var paths []string
	if withAuthor {
		projectsID, paths = GetInformationOfProjectsWithAuthor(w, authRqst)
	} else {
		projectsID, paths = GetInfoOfProjects(w)
	}

	

	// get the number of commits for each project
	var commitsCounts []int
	lenProjectsID := len(projectsID)
	for i := 0; i < lenProjectsID; i++ {
		commitsCount := 0

		lastPageNoReached := true
		j := 0

		for lastPageNoReached {
			j++

			jString := strconv.Itoa(j)
			iD := strconv.Itoa(int(projectsID[i]))

			var getArgument string
			if withAuthor {
				getArgument = fmt.Sprintf("https://git.gvk.idi.ntnu.no/api/v4/projects/%s/repository/commits?page=%s&private_token=%s", iD, jString, authRqst)
			} else {
				getArgument = fmt.Sprintf("https://git.gvk.idi.ntnu.no/api/v4/projects/%s/repository/commits?page=%s", iD, jString)
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

			var jsonCommits interface{}

			err = json.Unmarshal(body, &jsonCommits)

			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			mapOfCommits, ok := jsonCommits.([]interface{})

			if ok {
				var lenCommitsMap = len(mapOfCommits)

				commitsCount += lenCommitsMap

				if lenCommitsMap == 0 {
					lastPageNoReached = false
				}
			} else {
				lastPageNoReached = false
			}
		}

		commitsCounts = append(commitsCounts, commitsCount)
	}

	// determinante the new limit
	if requestLimit > lenProjectsID {
		requestLimit = lenProjectsID
	}

	var repos []ReposInfoStruct

	for i := 0; i < lenProjectsID; i++ {
		repos = append(repos, ReposInfoStruct{paths[i], commitsCounts[i]})
	}

	sort.Sort(ByCountingCommits(repos))

	commitsResponse.Repos = repos[0:requestLimit]

	// Response encoding
	commitsResponse.Author = withAuthor
	json.NewEncoder(w).Encode(commitsResponse)
}
