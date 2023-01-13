# -------------
# build stage
# -------------
FROM golang:alpine AS build

# Attach sources
WORKDIR /src
ADD . /src

# System deps
RUN apk add build-base git npm

# Build
RUN go build -o oxigen
RUN (cd assets; npm i; npm run build)

# -------------
# runtime stage
# -------------
FROM alpine

# Copy app
WORKDIR /app
COPY --from=build /src/oxigen /app/

# Entrypoint
ENTRYPOINT ./oxigen
# Command
CMD -http 0.0.0.0:80
