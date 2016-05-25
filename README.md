Synonyms Service
================

A small service for returning the synonyms for each meaningful word in a sentence using the WordNet 3.1 synset database.

Building and Installing
-----------------------

### Manually

The synonyms service is written in [golang](https://golang.org/). If Go is correctly installed, we should be able to acquire and install the source in the usual way with

```{bash}
go get github.com/deciphernow/synonyms
go install github.com/deciphernow/synonyms
```

### With Docker

For quickest results, we can use the prebuilt image on Docker Hub:

```{bash}
docker pull deciphernow/synonyms
```

Alternatively, to build synonyms in a container:

```{bash}
# from the repo root:
docker run --rm \
           -v "$PWD":/go/src/github.com/deciphernow/synonyms \
           -w /go/src/github.com/deciphernow/synonyms \
           iron/go:dev \
           go build -o synonyms
```

Which produces a "synonyms" Linux binary in the repo root and terminates.

We then build the `deciphernow/synonyms` production image with

```{bash}
# still from the repo root:
docker build -t deciphernow/synonyms .
```

Starting
--------

### Manual Start

If $GOPATH/bin is in $PATH, we can launch the synonyms service with

```{bash}
synonyms [port]
```

Otherwise, run with

```{bash}
$GOPATH/bin/synonyms [port]
```

### Docker Start

Once the `deciphernow/synonyms` image exists locally (see above), we can run it with, e.g.:

```{bash}
docker run --publish 8080:8080 --rm deciphernow/synonyms
```

This will run the `deciphernow/synonyms` image in a container, publishing internal port 8080 on external port 8080, and cleaning up the container filesystem upon exit.

Usage
-----

The service first checks for the presence of the WordNet database files in $TMPDIR/synonyms-service-wordnet-db/dict, and if not yet present, it downloads them (16Mb download, 53Mb uncompressed) to this location. Thereafter, the service listens on port 8080 (or first argument `port` if provided).

Synonyms supports text and JSON output. A simple GET to `localhost:8080?q=Hello, World!` will return

```
synonyms of 'hello': [hello hullo hi howdy how-do-you-do]
synonyms of 'world': [universe existence creation world cosmos macrocosm domain reality Earth earth globe populace public worldly_concern earthly_concern human_race humanity humankind human_beings humans mankind man]
```

and a GET to `localhost:8080/synonyms.json?q=Hello, World!` will return

```{json}

[{"word":"hello","synonyms":["hello","hullo","hi","howdy","how-do-you-do"]},{"word":"world","synonyms":["universe","existence","creation","world","cosmos","macrocosm","domain","reality","Earth","earth","globe","populace","public","worldly_concern","earthly_concern","human_race","humanity","humankind","human_beings","humans","mankind","man"]}]
```

The synonyms service also supports querying instead by header. A GET to `localhost:8080/synonyms.json` with the header `Q: Hello, World!` returns the same JSON output as above. Similarly, GET to `localhost:8080` with that header returns the same text output as above.

Roadmap
-------

A list of features and changes we'd like to make:

* ! Improve error handling (pass to requestor with accurate error code instead of dying)
* Load a third-party stop words list
* More sophisticated tokenization (e.g., better supporting contractions like "It's", which currently tokenizes as "it" and "s")
* A separate endpoint that applies heuristics to guess the *intended* WordNet sense of each word in a sentence, and return only the synset for that particular sense

Acknowledgements
----------------

* Princeton University "About WordNet." WordNet. Princeton University. 2010. <http://wordnet.princeton.edu>

* http://blog.ralch.com/tutorial/golang-working-with-tar-and-gzip/

License
-------

© Copyright 2016 Decipher Technology Studios

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
