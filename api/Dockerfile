FROM --platform=linux/amd64 golang:1.21.5-bullseye AS build
WORKDIR /opt/build
COPY . .
RUN curl -LO https://github.com/tinygo-org/tinygo/releases/download/v0.32.0/tinygo_0.32.0_amd64.deb && dpkg -i tinygo_0.32.0_amd64.deb
RUN tinygo build -target=wasi -o main.wasm main.go

FROM scratch
COPY --from=build /opt/build/spin.toml .
COPY --from=build /opt/build/main.wasm .
