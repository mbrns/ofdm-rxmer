# ofdm-rxmer

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
