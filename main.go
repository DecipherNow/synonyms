package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/deciphernow/synonyms/wordnet"
)

// main sets up and launches server
func main() {
	http.HandleFunc("/", sentenceHandler)
	log.Println("Listening on *:8080")
	http.ListenAndServe(":8080", nil)
}

// sentenceHandler handles requests for a whole sentence.
// Prints all words in all synsets that contain each word in the sentence, in sentence order.
func sentenceHandler(w http.ResponseWriter, r *http.Request) {
	sentence := r.FormValue("q")
	tokens := wordnet.Tokenize(sentence)
	switch {
	case strings.HasSuffix(r.URL.Path, ".json"):
		w.Header().Set("Content-Type", "application/json")
		var synonyms syns
		for _, word := range tokens {
			synonyms = append(synonyms, syn{word, wordnet.SynonymsForWord(word)})
		}
		js, err := json.Marshal(synonyms)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.Write(js)
	case strings.HasSuffix(r.URL.Path, ".txt"):
		w.Header().Set("Content-Type", "text/plain")
		fallthrough
	default:
		for _, word := range tokens {
			fmt.Fprintf(w, "synonyms of '%s': %s\n", word, wordnet.SynonymsForWord(word))
		}
	}
}

type syns []syn

type syn struct {
	Word     string   `json:"word"`
	Synonyms []string `json:"synonyms"`
}
