# Section IP Controller

HTTP Listener that receives a payload of bad IP addresses and updates application blocklists via the Aperture API.

## Installation

### Docker

Run directly from a docker image:
```sh
docker run --rm ghcr.io/dpc-sdp/bay-section-ip-controller:main bay-section-ip-controller <flags>
```

Or add to your docker image:

```Dockerfile
COPY --from=ghcr.io/dpc-sdp/bay-section-ip-controller:main /usr/local/bin/bay-section-ip-controller /usr/local/bin/bay-section-ip-controller
```

## Usage

```
$ bay-section-ip-controller -h

Usage of bay-section-ip-controller:
  -a string
        Comma separate list of applications to update
  -b string
        Comma separated list of IPs to always include in the blocklist
  -debug
        Sets log level to debug
  -e string
        Comma separated list of environments to update (default "Develop")
  -i string
        Account ID for Section API (default os.Getenv("SECTION_IO_ACCOUNT_ID"))
  -p string
        TCP listen port (default "80")
  -t string
        Token for Section API (default os.Getenv("SECTION_IO_TOKEN"))
  -u string
        User for Section API (default os.Getenv("SECTION_IO_USERNAME"))
```

## Local development

### Build

```sh
git clone git@github.com:dpc-sdp/bay-section-ip-controller.git && cd bay-section-ip-controller
go generate ./...
go build -ldflags="-s -w" -o build/bay-section-ip-controller .
go run . -h
```