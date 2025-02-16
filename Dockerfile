FROM golang:1.24.0 as base
LABEL authors="nawaf"
FROM base as dev
WORKDIR /src/
COPY . .
RUN curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s -- -b $(go env GOPATH)/bin
# Copy the Go application source code into the container

# Build the Go application
RUN go build -o bin/cmd/orders cmd/orders/main.go

ENV GOOGLE_APPLICATION_CREDENTIALS=/src/google_service_account.json
# Expose on Port 9003
EXPOSE 9000
CMD ["air","./bin/cmd/orders"]