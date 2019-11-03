package APIs

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	firebase "firebase.google.com/go"
	"golang.org/x/net/context"

	"google.golang.org/api/option"
)

var StartingTime = time.Now()

type StatusRspnsStruct struct {
	Gitlab   int     `json:"gitlab"`
	Database int     `json:"database"`
	Uptime   float64 `json:"timeUP"`
	Version  string  `json:"version"`
}

func StatusHandler(w http.ResponseWriter, r *http.Request) {
	http.Header.Add(w.Header(), "content-type", "application/json")

	var params []string
	WebhookCheckfunc(w, "languages", params)

	// check the availablity of gitlab
	gitResponse, gitError := http.Get("https://git.gvk.idi.ntnu.no/api/v4/projects")

	if gitError != nil {
		http.Error(w, gitError.Error(), http.StatusBadRequest)
		return
	}
	// check availablity of database
	dataBaseAvailablity := 200

	ctx := context.Background()

	serviceAccount := option.WithCredentialsFile("/home/nabil/Downloads/assignment2v1-firebase.json")
	app, err := firebase.NewApp(ctx, nil, serviceAccount)
	if err != nil {
		log.Fatalln(err)
		dataBaseAvailablity = 404
	}

	clnt, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
		dataBaseAvailablity = 404
	}

	defer clnt.Close()

	timeUP := time.Since(StartingTime).Seconds()

	statusResponse := StatusRspnsStruct{gitResponse.StatusCode, dataBaseAvailablity, timeUP, "v1"}

	json.NewEncoder(w).Encode(statusResponse)
}
