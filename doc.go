/*
Package realm implements a simple non-recursive DNS server.

INSTALLATION

To install realm:

    go get -u github.com/brettlangdon/realm/cmd/...

USAGE

Realm will parse your server configuration from a DNS zone file see https://en.wikipedia.org/wiki/Zone_file for more information.

To start a server:

    realm ./domain.zone
    realm --bind "127.0.0.1:1053" ./first.domain.zone ./second.domain.zone

Full command usage:

    usage: realm [--bind BIND] [ZONE [ZONE ...]]

    positional arguments:
      zone                   DNS zone files to serve from this server

    options:
      --bind BIND            [<host>]:<port> to bind too [default: :53]
      --help, -h             display this help and exit
*/
package realm
