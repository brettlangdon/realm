/*
Package realm implements a simple non-recursive DNS server.

INSTALLATION

To install realm:

    go get -u github.com/brettlangdon/realm/cmd/...

USAGE

Realm will parse your server configuration from a DNS zone file see https://en.wikipedia.org/wiki/Zone_file for more information.

To start a server:

    realm --zone ./domain.zone
    realm --zone ./first.domain.zone --zone ./second.domain.zone --bind "127.0.0.1:1053"

Full command usage:

    NAME:
       realm - A simple non-recursive DNS server

    USAGE:
       realm [global options] command [command options] [arguments...]

    VERSION:
       0.1.0

    COMMANDS:
       help, h	Shows a list of commands or help for one command

    GLOBAL OPTIONS:
       --zone, -z '--zone option --zone option'	location to DNS zone file [required] [$REALM_ZONE]
       --bind, -b ':53'				'[<host>]:<port>' to bind too [$REALM_BIND]
       --help, -h					show help
       --version, -v				print the version
*/
package realm
