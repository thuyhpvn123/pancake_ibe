module main

go 1.18

replace gitlab.com/meta-node/client => ./

replace gitlab.com/meta-node/core => ./core 

// replace testing1/testing/network/messages => ./network/messages
require (
	github.com/ethereum/go-ethereum v1.10.25
	github.com/gorilla/websocket v1.5.0
	github.com/holiman/uint256 v1.2.1
	github.com/sirupsen/logrus v1.9.0
	gitlab.com/meta-node/client v0.0.0-00010101000000-000000000000
	gitlab.com/meta-node/core v0.0.0-00010101000000-000000000000
	google.golang.org/protobuf v1.28.1
)

require (
	github.com/gofiber/fiber/v2 v2.39.0
	go.mongodb.org/mongo-driver v1.10.4
)

require (
	github.com/andybalholm/brotli v1.0.4 // indirect
	github.com/btcsuite/btcd/btcec/v2 v2.2.0 // indirect
	github.com/decred/dcrd/dcrec/secp256k1/v4 v4.0.1 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/gorilla/mux v1.8.0 // indirect
	github.com/klauspost/compress v1.15.12 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.16 // indirect
	github.com/mattn/go-runewidth v0.0.14 // indirect
	github.com/montanaflynn/stats v0.6.6 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/rivo/uniseg v0.2.0 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasthttp v1.40.0 // indirect
	github.com/valyala/tcplisten v1.0.0 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.1.1 // indirect
	github.com/xdg-go/stringprep v1.0.3 // indirect
	github.com/youmark/pkcs8 v0.0.0-20201027041543-1326539a0a0a // indirect
	golang.org/x/crypto v0.1.0 // indirect
	golang.org/x/sync v0.1.0 // indirect
	golang.org/x/sys v0.1.0 // indirect
	golang.org/x/text v0.4.0 // indirect
)
