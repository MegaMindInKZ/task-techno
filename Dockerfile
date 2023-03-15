FROM golang

WORKDIR /app


COPY . .

RUN go mod download
RUN go build -o task-techno .

EXPOSE 8000
CMD ["./task-techno"]

