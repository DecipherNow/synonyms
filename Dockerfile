FROM iron/go

ENV tmpdir /tmp/synonyms-service-wordnet-db
ADD http://wordnetcode.princeton.edu/wn3.1.dict.tar.gz $tmpdir/dict.tar.gz
RUN chown -R daemon: $tmpdir

USER daemon

WORKDIR /app/
COPY synonyms .

ENTRYPOINT ["./synonyms"]
CMD ["8080"]

EXPOSE 8080
