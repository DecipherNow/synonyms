FROM iron/go:dev

ENV dbdir /tmp/synonyms-service-wordnet-db
ENV workdir /go/src/github.com/deciphernow/synonyms
ENV PATH $PATH:/go/bin

ADD http://wordnetcode.princeton.edu/wn3.1.dict.tar.gz $dbdir/dict.tar.gz

COPY . $workdir

WORKDIR $workdir

CMD go install && synonyms 8080

EXPOSE 8080
