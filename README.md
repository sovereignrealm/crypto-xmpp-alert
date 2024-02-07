
# Crypto-xmpp-alert

This a golang project to send a xmpp message if % gained on a crypto asset reaches certain boundary.

You should replace input/data.json with your data.

The project uses interfaces as dependency injection and has its unit tests.

**IMPORTANT:** it only works by getting the current price from the following crytos: Bitcoin, Ethereum, Cardano, Polkadot.


run on development

```sh
go run main.go
```

run tests:

```sh
go test
```

compile and run build:

```sh
go build -o mybinary
./mybinary
```

The `scripts/write_true.sh` is designed to reset to true the content inside input files so it only sends 1 xmpp message per day in case your gains hit the boundary.
You can execute `scripts/write_true.sh` inside a cronjob to reset it to true once per day as an example.
