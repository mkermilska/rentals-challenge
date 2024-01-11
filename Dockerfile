FROM golang:1.21

RUN mkdir /app
ADD . /app
WORKDIR /app/cmd/rentals-challenge

RUN CGO_ENABLED=0 go build -ldflags "-s -w" -o /rentals-challenge

EXPOSE 59191

CMD [ "/rentals-challenge" ]