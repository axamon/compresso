package main

//Accesslog tipo per transaction
type Accesslog struct {
	Hash      string
	Type      string
	Time      string
	TTS       int
	SEIp      string
	Clientip  string
	Request   string
	Bytes     int
	Method    string
	URL       string
	Urlschema string
	Urlhost   string
	Urlpath   string
	Urlquery  string
	Mime      string
	Ua        string
}

//Fruizioni keeps data relavant to fruitions
type Fruizioni struct {
	Hashfruizione map[string]bool
	Clientip      map[string]string
	Idvideoteca   map[string]string
	Idaps         map[string]string
	Edgeip        map[string]string
	Giorno        map[string]string
	Orario        map[string]string
	Details       map[string][]float64 `json:"-"`
}
