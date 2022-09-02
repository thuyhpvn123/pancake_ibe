module testin

go 1.18

replace gitlab.com/meta-node/client => ./

replace gitlab.com/meta-node/core => ./core

require (
	github.com/ethereum/go-ethereum v1.10.22
	github.com/gorilla/websocket v1.5.0
	gitlab.com/meta-node/client v0.0.0-00010101000000-000000000000
	gitlab.com/meta-node/core v0.0.0-00010101000000-000000000000
	google.golang.org/protobuf v1.28.1
)

require (
	github.com/holiman/uint256 v1.2.0 // indirect
	github.com/sirupsen/logrus v1.8.1 // indirect
	golang.org/x/crypto v0.0.0-20210921155107-089bfa567519 // indirect
	golang.org/x/sys v0.0.0-20220520151302-bc2c85ada10a // indirect
)
