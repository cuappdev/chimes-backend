FROM golang:1.25 AS build-stage
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /chimes-backend

FROM gcr.io/distroless/base-debian11
WORKDIR /
COPY --from=build-stage /chimes-backend /chimes-backend
EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/chimes-backend"]