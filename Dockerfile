FROM docker.io/golang:1.16 as build
ARG VERSION

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o ktr

FROM gcr.io/distroless/static:nonroot

COPY --from=build /app/ktr /usr/local/bin/

ENTRYPOINT [ "/usr/local/bin/ktr" ]

