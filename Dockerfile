FROM golang:1.17.7-alpine3.15

ENV SHORT_LINK_BASE=http://localhost:15001/ 

ENV POSTGRES_CONN_STRING=postgres://user:password@host:port/databaseName 

RUN GOROOT=/go

WORKDIR /go/src/ozonLinks

ADD . .

RUN go mod tidy

RUN CGO_ENABLED=0 GOOS=linux go build -a -o ozonLinks .

EXPOSE 15001

RUN chmod 777 /go/src/ozonLinks

ENTRYPOINT [ "/go/src/ozonLinks/ozonLinks" ]
