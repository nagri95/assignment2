package APIs

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"google.golang.org/api/option"

	"google.golang.org/api/iterator"

	firebase "firebase.google.com/go"
)

type Registration struct {
	Event string `json:"event"`
	Url   string `json:"url"`
}

type RegistrationResponse struct {
	ID string `json:"id"`
}

type RegisteredWebhook struct {
	ID    string `json:"id"`
	Event string `json:"event"`
	Time  string `json:"time"`
}

var keyAddress = "/home/nabil/Downloads/assignment2v1-firebase.json"

func WebhookHandler(w http.ResponseWriter, r *http.Request) {

	// connection to the DB
	contxt := context.Background()

	serviceAccount := option.WithCredentialsFile(keyAddress)
	application1, err := firebase.NewApp(contxt, nil, serviceAccount)
	if err != nil {
		log.Fatalln(err)
	}

	clnt1, err := application1.Firestore(contxt)
	if err != nil {
		log.Fatalln(err)
	}

	defer clnt1.Close()

	switch r.Method {
	case http.MethodPost:
		var registratn Registration
		err := json.NewDecoder(r.Body).Decode(&registratn)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			fmt.Println("Decoding: " + err.Error())
			return
		}

		fmt.Println("Adding to db ...")

		ID, _, err := clnt1.Collection("webhooks").Add(contxt, map[string]interface{}{
			"event": registratn.Event,
			"url":   registratn.Url,
		})
		if err != nil {
			log.Fatalf("Failed adding alovelace: %v", err)
		}

		http.Header.Add(w.Header(), "content-type", "application/json")

		var registratnResponse = RegistrationResponse{}
		registratnResponse.ID = ID.ID

		json.NewEncoder(w).Encode(registratnResponse)

		return

	case http.MethodGet:
		http.Header.Add(w.Header(), "content-type", "application/json")

		var registerdWebhoks = []RegisteredWebhook{}

		itertr := clnt1.Collection("webhooks").Documents(contxt)
		for {
			document, err := itertr.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				log.Fatalf("Failed to iterate: %v", err)
			}

			var registeredWebhook = RegisteredWebhook{}
			registeredWebhook.ID = document.Ref.ID
			registeredWebhook.Event = document.Data()["event"].(string)
			registeredWebhook.Time = document.CreateTime.String()

			registerdWebhoks = append(registerdWebhoks, registeredWebhook)
		}

		json.NewEncoder(w).Encode(registerdWebhoks)

		return

	default:
		http.Error(w, "not implemented yet", http.StatusNotImplemented)
		return
	}
}

func WebhookWithIdHandler(w http.ResponseWriter, r *http.Request) {

	// connecting to the DataBase
	contxt := context.Background()

	serviceAccount := option.WithCredentialsFile(keyAddress)
	application1, err := firebase.NewApp(contxt, nil, serviceAccount)
	if err != nil {
		log.Fatalln(err)
	}

	clnt1, err := application1.Firestore(contxt)
	if err != nil {
		log.Fatalln(err)
	}

	defer clnt1.Close()

	http.Header.Add(w.Header(), "content-type", "application/json")

	parts := strings.Split(r.URL.Path, "/")

	switch r.Method {
	case http.MethodGet:

		var registerdWebhoks = map[string]RegisteredWebhook{}

		itertr := clnt1.Collection("webhooks").Documents(contxt)
		for {
			document, err := itertr.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				log.Fatalf("Failed to iterate: %v", err)
			}

			var registeredWebhook = RegisteredWebhook{}
			registeredWebhook.ID = document.Ref.ID
			registeredWebhook.Event = document.Data()["event"].(string)
			registeredWebhook.Time = document.CreateTime.String()

			registerdWebhoks[document.Ref.ID] = registeredWebhook
		}

		json.NewEncoder(w).Encode(registerdWebhoks[parts[4]])

		return

	case http.MethodDelete:

		_, err := clnt1.Collection("webhooks").Doc(parts[4]).Delete(contxt)
		if err != nil {
			log.Printf("An error has occurred: %s", err)
		}

		return
	}
}

type InformationSend struct {
	Event  string   `json:"event"`
	Params []string `json:"params"`
	Time   string   `json:"time"`
}

type WebhookRegisteredWithUrl struct {
	ID    string `json:"id"`
	Event string `json:"event"`
	Url   string `json:"url"`
}

func WebhookCheckfunc(w http.ResponseWriter, eventType string, parameters []string) {
	// creating pyld that will be sent
	infoSend := InformationSend{eventType, parameters, time.Now().String()}

	// connecting to the DataBase
	contxt := context.Background()

	serviceAccount := option.WithCredentialsFile(keyAddress)
	application1, err := firebase.NewApp(contxt, nil, serviceAccount)
	if err != nil {
		log.Fatalln(err)
	}

	clnt1, err := application1.Firestore(contxt)
	if err != nil {
		log.Fatalln(err)
	}

	defer clnt1.Close()

	// get the webhooks in the database
	var registerdWebhoks = map[string]WebhookRegisteredWithUrl{}

	itertr := clnt1.Collection("webhooks").Documents(contxt)
	for {
		document, err := itertr.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Iteration failed: %v", err)
		}

		var registeredWebhook = WebhookRegisteredWithUrl{}
		registeredWebhook.ID = document.Ref.ID
		registeredWebhook.Event = document.Data()["event"].(string)
		registeredWebhook.Url = document.Data()["url"].(string)

		registerdWebhoks[document.Ref.ID] = registeredWebhook
	}

	// get only eventType related
	var ids []string

	for webhock := range registerdWebhoks {
		if registerdWebhoks[webhock].Event == eventType {
			ids = append(ids, registerdWebhoks[webhock].ID)
		}
	}

	// sending pyld Get request

	lengthofIDs := len(ids)

	for counter := 0; counter < lengthofIDs; counter++ {
		url := registerdWebhoks[ids[counter]].Url
		requestBody, err := json.Marshal(infoSend)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			continue
		}

		pyld := bytes.NewBuffer(requestBody)

		request, err1 := http.NewRequest(http.MethodGet, url, pyld)

		if err1 != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			continue
		}

		_, _ = http.DefaultClient.Do(request)
	}
}
