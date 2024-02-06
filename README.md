# ofdm-rxmer

[![Docker Image CI](https://github.com/mbrns/ofdm-rxmer/actions/workflows/docker-image.yml/badge.svg)](https://github.com/mbrns/ofdm-rxmer/actions/workflows/docker-image.yml)
[![Go](https://github.com/mbrns/ofdm-rxmer/actions/workflows/go.yml/badge.svg)](https://github.com/mbrns/ofdm-rxmer/actions/workflows/go.yml)

---


## Binary
#### Building
```
make
```

## Docker
#### Building 
```
make docker
```
#### Running
```
docker run --rm ofdm-rxmer:latest rxmer -cm 10.10.10.10 -comm private -tftp 10.11.11.11 -out json-pretty
```
