# -------------
# build stage
# -------------
FROM golang:alpine AS build

# System deps
RUN apk add build-base git npm

# Attach sources
WORKDIR /src
ADD . /src

# Build
RUN (cd assets; npm i; npm run build)
RUN go build -o oxigen

# -------------
# runtime stage
# -------------
FROM alpine

# Copy app
WORKDIR /app
COPY --from=build /src/oxigen /app/oxigen
COPY --from=build /src/assets/dist /app/assets/dist

# Command
CMD ["./oxigen", "-http", "0.0.0.0:80"]
