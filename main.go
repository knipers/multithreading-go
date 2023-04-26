package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type Cep struct {
	Cep    string
	City   string
	State  string
	Origin string
}

type ViaCep struct {
	Cep   string `json:"cep"`
	City  string `json:"localidade"`
	State string `json:"uf"`
}

type ApiCep struct {
	Cep   string `json:"code"`
	City  string `json:"city"`
	State string `json:"state"`
}

func main() {
	zipCode := os.Args[1]

	viaCepCh := make(chan Cep)
	apiCepCh := make(chan Cep)

	go getViaCep(zipCode, viaCepCh)
	go getApiCep(zipCode, apiCepCh)

	select {
	case cep := <-viaCepCh:
		fmt.Printf("The data for cep was found by %s: %v", cep.Origin, cep)

	case cep := <-apiCepCh:
		fmt.Printf("The data for cep was found by %s: %v", cep.Origin, cep)

	case <-time.After(time.Second):
		println("timeout")
	}

}

func getViaCep(zipCode string, ch chan Cep) {
	req, err := http.Get("http://viacep.com.br/ws/" + zipCode + "/json/")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao fazer requisição: %v\n", err)
	}
	defer req.Body.Close()
	res, err := io.ReadAll(req.Body)
	var data ViaCep
	err = json.Unmarshal(res, &data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao fazer parse da resposta: %v\n", err)
	}
	cep := Cep{data.Cep, data.City, data.State, "ViaCEP"}

	ch <- cep
}

func getApiCep(zipCode string, ch chan Cep) {
	req, err := http.Get("https://cdn.apicep.com/file/apicep/" + zipCode + ".json")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao fazer requisição: %v\n", err)
	}
	defer req.Body.Close()
	res, err := io.ReadAll(req.Body)
	var data ApiCep
	err = json.Unmarshal(res, &data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao fazer parse da resposta: %v\n", err)
	}
	cep := Cep{data.Cep, data.City, data.State, "apiCEP"}

	ch <- cep
}
