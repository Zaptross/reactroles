FROM golang:alpine as build

RUN apk --no-cache add ca-certificates

WORKDIR /build
COPY . .
RUN go mod download
ENV CGO_ENABLED=0
RUN cd /build/cmd/reactroles && go build -a -o /build/reactroles

# Create the output container from the built image.
FROM scratch
COPY --from=build /build/reactroles /reactroles
COPY --from=build /etc/ssl/certs/ /etc/ssl/certs

ENTRYPOINT [ "/reactroles" ]