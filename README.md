# Providence

![logo](panel/source/favicon.png)

A general-purpose, free and open-source command and control pannel on top gRPC written in Go and Mint.

Currently under experimental stage and likely many breaking change, prototype stage.

## Progress

- [x] Bootstraping data types (such Task,Stock, etc).
- [x] Basic gRPC communication.
- [ ] Panel View.
- [ ] Revisit JWT with JWE.
- [ ] Breaking Full mode into proxy and panel mode.
- [ ] Implementing HTTP/S RPC as alternative comm channel.
- [ ] TBD

## Mode

Providence is sparated in three mode:

- Full Mode

Works as command and control server with panel.

- Proxy Mode

C2 server only interaction without panel.

- Panel Mode

web panel that connected to proxy mode providence instance 
and manage it.

## How To run

### First timer

- Install MariaDB and create a database.
- Change the `DBDSN` variable inside Makefile and `db.config.yml.bak` file according your database and server connection.
- rename file `db.config.yml.bak` to `db.config.yml`
- Run migration, `make up-migrate`, to down migration do `make down-migrate`.
- Make keypair for `PUBKEYPATH` and `PRIVKEYPATH`, run `make keypair`.

To start the panel, `make run-panel`.

To start the server, `make run-server`.

## Dependecies

### Runtime

- MariaDB > 10.x.x
- GNU make
- openssl

### Development

- Go compiler > 1.18.0
- GNU make
- richgo 
- mockgen
- sql-migrate
- protoc

For front end (inside `panel` folder):

- [mint-lang](https://mint-lang.com/install)

## LICENSE

[UNLICENSE](/UNLICENSE)

[Contributor License Agreement](/CONTRIBUTOR_LICENSE_AGREEMENT.md)