package wordnet

import (
	"archive/tar"
	"bufio"
	"compress/gzip"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

var synsets *Synsets
var nonAlphanumeric *regexp.Regexp
var stopwords *regexp.Regexp

// Synsets contains maps from string to byte offset, and byte offset to tokens
type Synsets struct {
	Index map[string][]uint64
	Data  map[uint64][]string
}

// TODO Download WordNet database if it doesn't exist where expected

// Initialize pulls WordNet's synset database into memory
func Initialize() *Synsets {
	nonAlphanumeric = regexp.MustCompile("\\W+")
	stopwords = regexp.MustCompile("\\s+(I|a|an|as|at|by|he|she|his|hers|it|its|me|or|thou|us|who)\\s")

	currentDirectory, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	downloadDatabase(currentDirectory)

	synsets = &Synsets{Index: make(map[string][]uint64), Data: make(map[uint64][]string)}

	for _, filename := range []string{"index.adj", "index.adv", "index.noun", "index.verb", "data.adj", "data.adv", "data.noun", "data.verb"} {
		isIndexFile := strings.HasPrefix(filename, "index")
		isDataFile := !isIndexFile
		f, err := os.Open(os.ExpandEnv(path.Join(currentDirectory, "dict", filename)))
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, "  ") {
				continue
			}

			fields := strings.Fields(line)
			switch {
			case isIndexFile: // load index files containing word->synset byte offset
				word := fields[0]
				synsetCnt, err := strconv.ParseUint(fields[2], 10, 32)
				ptrCnt, err := strconv.ParseUint(fields[3], 10, 32)
				if err != nil {
					log.Fatal(err)
				}
				synsetIndices := make([]uint64, synsetCnt)
				for i, synsetOffset := range fields[6+ptrCnt:] {
					wideIndex, err := strconv.ParseUint(synsetOffset, 10, 32)
					if err != nil {
						log.Fatal(err)
					}
					synsetIndices[i] = wideIndex
				}
				synsets.Index[word] = synsetIndices

			case isDataFile: // load data files containing the synset byte offset -> synonyms
				synsetOffset, err := strconv.ParseUint(fields[0], 10, 32)
				if err != nil {
					log.Fatal(err)
				}
				wordCnt, err := strconv.ParseUint(fields[3], 16, 32)
				if err != nil {
					log.Fatal(err)
				}
				// TODO finish parsing out the data file <---
				words := make([]string, wordCnt)
				for i, word := range fields[4 : 3+(2*wordCnt)] {
					if i%2 == 0 {
						words[i/2] = word
					}
				}
				synsets.Data[synsetOffset] = words
			}
		}
	}
	return synsets
}

// Tokenize is a relatively naive sentence tokenizer
// TODO further normalize
func Tokenize(sentence string) []string {
	alphanumericOnly := nonAlphanumeric.ReplaceAllString(sentence, " ")
	withoutStopwords := stopwords.ReplaceAllString(" "+alphanumericOnly+" ", " ")
	fields := strings.Fields(withoutStopwords)
	for i, field := range fields {
		fields[i] = strings.ToLower(field)
	}
	return fields
}

// SynonymsForWord returns the synonyms of a given word.
func SynonymsForWord(word string) []string {
	sliceOfOffsets := synsets.Index[word]
	var allOffsets []string
	for _, offset := range sliceOfOffsets {
		allOffsets = append(allOffsets, synsets.Data[offset]...)
	}
	return deDup(allOffsets)
}

func deDup(s []string) []string {
	result := []string{}
	seen := map[string]bool{}
	for _, val := range s {
		if _, ok := seen[val]; !ok {
			result = append(result, val)
			seen[val] = true
		}
	}
	return result
}

// downloadDatabase downloads the WordNet 3.1 database if necessary
func downloadDatabase(intoDirectory string) {
	// Don't already have it? Download it.
	if _, err := os.Stat(path.Join(intoDirectory, "dict")); os.IsNotExist(err) {
		log.Println("Downloading WordNet 3.1 database files...")
		out, err := os.Create(path.Join(intoDirectory, "dict.tar.gz"))
		if err != nil {
			log.Fatal(err)
		}
		resp, err := http.Get("http://wordnetcode.princeton.edu/wn3.1.dict.tar.gz")
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()
		_, err = io.Copy(out, resp.Body)
		if err != nil {
			log.Fatal(err)
		}

		// extract it
		if err := ungzip(path.Join(intoDirectory, "dict.tar.gz"), intoDirectory); err != nil {
			log.Fatal(err)
		}
		if err := untar(path.Join(intoDirectory, "dict.tar"), intoDirectory); err != nil {
			log.Fatal(err)
		}
	}
}

// ungzip from http://blog.ralch.com/tutorial/golang-working-with-tar-and-gzip/

func ungzip(source, target string) error {
	reader, err := os.Open(source)
	if err != nil {
		return err
	}
	defer reader.Close()

	archive, err := gzip.NewReader(reader)
	if err != nil {
		return err
	}
	defer archive.Close()

	target = filepath.Join(target, archive.Name)
	writer, err := os.Create(target)
	if err != nil {
		return err
	}
	defer writer.Close()

	_, err = io.Copy(writer, archive)
	return err
}

// untar from http://blog.ralch.com/tutorial/golang-working-with-tar-and-gzip/
func untar(tarball, target string) error {
	reader, err := os.Open(tarball)
	if err != nil {
		return err
	}
	defer reader.Close()
	tarReader := tar.NewReader(reader)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		path := filepath.Join(target, header.Name)
		info := header.FileInfo()
		if info.IsDir() {
			if err = os.MkdirAll(path, info.Mode()); err != nil {
				return err
			}
			continue
		}

		file, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, info.Mode())
		if err != nil {
			return err
		}
		defer file.Close()
		_, err = io.Copy(file, tarReader)
		if err != nil {
			return err
		}
	}
	return nil
}
