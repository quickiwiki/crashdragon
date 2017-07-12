FROM golang:latest

ENV GOPATH /go
ENV PATH $PATH:$GOPATH/bin

ADD . $GOPATH/src/git.1750studios.com/GSoC/CrashDragon
WORKDIR $GOPATH/src/git.1750studios.com/GSoC/CrashDragon

RUN apt-get update && apt-get -y install libcurl4-gnutls-dev rsync postgresql sassc

RUN go get -u github.com/kardianos/govendor
RUN govendor sync

RUN make

RUN /etc/init.d/postgresql start && su postgres -c 'createuser -w crashdragon' && su postgres -c 'createdb -w -O crashdragon crashdragon'
RUN echo "local all all trust" > /etc/postgresql/9.4/main/pg_hba.conf
RUN echo "host all all all trust" >> /etc/postgresql/9.4/main/pg_hba.conf

RUN cd assets/stylesheets && sassc -t compressed app.scss > app.css

EXPOSE 8080
CMD /etc/init.d/postgresql start && sleep 15 && ./bin/CrashDragon