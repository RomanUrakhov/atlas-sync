# syntax=docker/dockerfile:1

FROM golang:1.26.3-alpine AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# CGO off so the binary runs on the distroless/static base.
ENV CGO_ENABLED=0
RUN go build -trimpath -ldflags="-s -w" -o /atlas-sync ./cmd/atlas-sync

# distroless/static-debian12 ships CA certs for outbound HTTPS to Atlassian and GCP.
FROM gcr.io/distroless/static-debian12

COPY --from=build /atlas-sync /atlas-sync

USER nonroot:nonroot

ENTRYPOINT ["/atlas-sync"]
CMD ["sync"]
