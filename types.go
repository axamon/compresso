package main

type Fruizioni struct {
	Hashfruizione map[string]bool
	Clientip      map[string]string
	Idvideoteca   map[string]string
	Details       map[string][]float64 `json:"-"`
}

type Fruizioniexport struct {
	Hashfruizione string `json:"-"` //Non permette a json di esportare il campo
	Clientip      string
	Idvideoteca   string
	Edgeip        string
	Giorno        string
	Orario        string
	Errori        int
}