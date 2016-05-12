package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/deciphernow/synonyms/wordnet"
)

var synsets *wordnet.Synsets

// main sets up and launches server with word and sentence endpoints
func main() {
	synsets = wordnet.Initialize()
	http.HandleFunc("/", sentenceHandler)
	log.Println("Listening on *:8080")
	http.ListenAndServe(":8080", nil)
}

// sentenceHandler handles requests for a whole sentence.
// Prints all words in all synsets that contain each word in the sentence, in sentence order.
func sentenceHandler(w http.ResponseWriter, r *http.Request) {
	sentence := r.FormValue("q")
	for _, word := range wordnet.Tokenize(sentence) {
		fmt.Fprintf(w, "synonyms of '%s': %s\n", word, wordnet.SynonymsForWord(word))
	}
}
