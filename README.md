Realm
=====
[![GoDoc](https://godoc.org/github.com/brettlangdon/realm?status.svg)](https://godoc.org/github.com/brettlangdon/realm)

A simple non-recursive DNS server written in [go](https://golang.org).

## Installation
### Go get
```
go get -u github.com/brettlangdon/realm/cmd/...
```

### Build
```
git clone https://github.com/brettlangdon/realm
cd ./realm
make
```

## Usage
### Zone file
To run a server you must have a DNS [zone file](https://en.wikipedia.org/wiki/Zone_file).

A simple example looks like the following:

```
$ORIGIN example.com.     ; designates the start of this zone file in the namespace
$TTL 1h                  ; default expiration time of all resource records without their own TTL value
example.com.  IN  SOA   ns.example.com. username.example.com. ( 2007120710 1d 2h 4w 1h )
example.com.  IN  NS    ns                    ; ns.example.com is a nameserver for example.com
example.com.  IN  NS    ns.somewhere.example. ; ns.somewhere.example is a backup nameserver for example.com
example.com.  IN  MX    10 mail.example.com.  ; mail.example.com is the mailserver for example.com
@             IN  MX    20 mail2.example.com. ; equivalent to above line, "@" represents zone origin
@             IN  MX    50 mail3              ; equivalent to above line, but using a relative host name
example.com.  IN  A     192.0.2.1             ; IPv4 address for example.com
              IN  AAAA  2001:db8:10::1        ; IPv6 address for example.com
ns            IN  A     192.0.2.2             ; IPv4 address for ns.example.com
              IN  AAAA  2001:db8:10::2        ; IPv6 address for ns.example.com
www           IN  CNAME example.com.          ; www.example.com is an alias for example.com
wwwtest       IN  CNAME www                   ; wwwtest.example.com is another alias for www.example.com
mail          IN  A     192.0.2.3             ; IPv4 address for mail.example.com
mail2         IN  A     192.0.2.4             ; IPv4 address for mail2.example.com
mail3         IN  A     192.0.2.5             ; IPv4 address for mail3.example.com
```

Example taken from [here](https://en.wikipedia.org/wiki/Zone_file#File_format).

### Starting the server
```
realm --zone ./domain.zone
```

By default `realm` binds to port `53`, which usually requires root, so you may need to run `sudo realm --zone ./domain.zone`.

### Options
* `--zone, -z` - the file file to load (e.g. `./domain.zone`), this argument is required
    * You may instead specify the environment variable `REALM_ZONE="./domain.zone"`
* `--bind, -b` - the `[<host>]:<port>` to bind the server to (e.g. `0.0.0.0:53`), default is `:53`
    * You may instead specify the environment variable `REALM_BIND=":53"`
* `--help, -h` - show help message
* `--version, -v` - show version information

To see the latest command usage run `realm --help`.
