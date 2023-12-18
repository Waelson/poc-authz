package authz

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Resource struct {
	Namespace string `json:"namespace"`
	ID        string `json:"id"`
}

type Subject struct {
	Namespace string `json:"namespace"`
	ID        string `json:"id"`
}

type AuthzRequest struct {
	Resource Resource `json:"resource"`
	Relation string   `json:"relation"`
	Subject  Subject  `json:"subject"`
}

func WriteRelationship(authzReq AuthzRequest) error {
	requestBody, err := json.Marshal(authzReq)
	fmt.Println(string(requestBody))
	if err != nil {
		log.Fatalf("Error occurred during marshaling. Error: %s", err.Error())
		return err
	}

	resp, err := http.Post("http://localhost:7070/authz/relationship", "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		log.Fatalf("Error occurred during making a request. Error: %s", err.Error())
		return err
	}
	defer resp.Body.Close()

	// Lendo a resposta
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Erro ao ler o corpo da resposta. Erro: %s", err.Error())
	}

	// Tratando a resposta com base no status
	switch resp.StatusCode {
	case http.StatusAccepted:
		log.Println("Requisição bem-sucedida. Resposta:", string(body))
		return nil
	case http.StatusUnauthorized, http.StatusForbidden:
		log.Println("Acesso negado. Status Code:", resp.StatusCode)
		return fmt.Errorf("acesso negado. Status Code: %d", resp.StatusCode)
	default:
		return fmt.Errorf("falha na requisição. Status Code: %d. Resposta: %s", resp.StatusCode, string(body))
	}

}
