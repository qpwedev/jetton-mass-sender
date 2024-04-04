<p align="center" style="color: #444">
  <h1 align="center">🍭 Jetton Mass Sender</h1>
</p>
<p align="center" style="font-size: 1.2rem;">Send TON Jettons Easily!</p>

1. Prepare json file with the following structure:

```json
[
  {
    "amount": "65089.333268244",
    "address": "UQAZF-cErbXnXbSTDJCFM3k5GI4dqFh5NSzgIT7tIMK5rSOX"
  },
  {
    "amount": "1731376.26493529",
    "address": "UQBCSonWNi0mVwHdxord0JnUfjzDCAJlCIsas2fBmR1p00tk"
  },
  {
    "amount": "2169644.44227224",
    "address": "UQDXOSJeAPITOr7NrdmqXRALMzskzKo087qCb-0V3ZrKFrjs"
  }
]
```

2. Clone Repository

```bash
git clone git@github.com:qpwedev/jetton-mass-sender.git; cd jetton-mass-sender;
```

3. [Install golang](https://go.dev/doc/install)
4. Setup Jetton Sender

```bash
go run src/main.go src/cli.go src/highload.go setup
```

5. Run Jetton Sender

```bash
go run src/main.go src/cli.go src/highload.go
```
