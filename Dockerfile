FROM golang:1.21 as build

WORKDIR /app

COPY go.* ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o /server


FROM gcr.io/distroless/static-debian11

COPY --from=build /server /server

USER nonroot:nonroot

EXPOSE 8000

ENTRYPOINT [ "/server" ]