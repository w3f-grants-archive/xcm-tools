# XCM Tools

---
XCM Tools is a tool and an SDK. This library is written in golang.
It provides the following functions: sending xcm messages, parsing xcm message instructions and getting the execution
result of the execution after sending xcm.

## Features

- [x] Send VMP(UMP & HRMP) message
- [x] Send HRMP message
- [x] Parse xcm message
- [x] Tracer xcm message result
- [x] Cli Support

## Get Start

### Requirement

1. Go 1.18+
2. docker(optional)

### Build

#### Build Docker Image
```bash 
docker build -f Dockerfile-build -t xcm-tools .
docker run -it xcm-tools -h
```

#### Build Binary
```bash
cd cmd && go build -o xcm-tools .
```

### Usage

#### Installation

```bash 
go install github.com/gmajor-encrypt/xcm-tools/cmd@latest 
```
You can find binary file(cmd) in $GOPATH/bin


### CLI Usage

XCM tools also support sending xcm messages, parsing messages, and tracking xcm transaction results through cli
commands.

```bash
go run .  -h
```

#### Commands

```
NAME:
   Xcm tools - Xcm tools

USAGE:
   cmd [global options] command [command options] [arguments...]

COMMANDS:
   send     send xcm message
   parse    parse xcm message
   tracker  tracker xcm message transaction
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h  show help
```

#### args

| Name               | Describe                                                                      | Suitable |
|--------------------|-------------------------------------------------------------------------------|----------|
| dest               | Dest address                                                                  | SendXCM  |
| amount             | Send xcm transfer amount                                                      | SendXCM  |
| keyring            | Set sr25519 secret key                                                        | SendXCM  |
| endpoint           | Set substrate endpoint, only support websocket protocol, like ws:// or wss:// | ALL      |
| paraId             | Send xcm transfer amount                                                      | SendXCM  |
| message            | Parsed xcm message raw data                                                   | Parse    |
| extrinsicIndex     | Xcm message extrinsicIndex                                                    | Tracker  |
| protocol           | Xcm protocol, such as UMP,HRMP,DMP                                            | Tracker  |
| destEndpoint       | Dest chain endpoint, only support websocket protocol, like ws:// or wss://    | Tracker  |
| relaychainEndpoint | Relay chain endpoint, only support websocket protocol, like ws:// or wss://   | Tracker  |

#### example

```bash
#ump
go run . send UMP --dest 0xd43593c715fdd31c61141abd04a99fd6822c8558854ccde39a5684e7a56da27d --amount 10 --keyring 0xe5be9a5092b81bca64be81d212e7f2f9eba183bb7a90954f7b76361f6edb5c0a --endpoint wss://rococo-asset-hub-rpc.polkadot.io
#dmp
go run . send DMP --dest 0xd43593c715fdd31c61141abd04a99fd6822c8558854ccde39a5684e7a56da27d --amount 10 --keyring 0xe5be9a5092b81bca64be81d212e7f2f9eba183bb7a90954f7b76361f6edb5c0a --endpoint wss://rococo-rpc.polkadot.io --paraId 1000
#hrmp
go run . send HRMP --dest 0xd43593c715fdd31c61141abd04a99fd6822c8558854ccde39a5684e7a56da27d --amount 10 --keyring 0xe5be9a5092b81bca64be81d212e7f2f9eba183bb7a90954f7b76361f6edb5c0a --endpoint wss://rococo-asset-hub-rpc.polkadot.io --paraId 2087
```


### Parse Xcm Message

We provide a function to parse xcm transaction instructions and deserialize the encoded raw message into readable JSON. Support XCM V0,V1,V2,V3.

```go
package example

import (
	"fmt"
	"github.com/gmajor-encrypt/xcm-tools/tx"
)

func ParseMessage() {
	client := tx.NewClient("wss://rococo-rpc.polkadot.io")
	defer client.Close()
	instruction, err := client.ParseXcmMessageInstruction("0x031000040000000007f5c1998d2a0a130000000007f5c1998d2a000d01020400010100ea294590dbcfac4dda7acd6256078be26183d079e2739dd1e8b1ba55d94c957a")
	fmt.Println(instruction, err)
}

```

### Tracker Xcm Message

We provide a function to track xcm transaction results. Support protocol UMP,HRMP,DMP.

```go
package example

import (
	"fmt"
	"github.com/gmajor-encrypt/xcm-tools/tracker"
	"github.com/gmajor-encrypt/xcm-tools/tx"
)

// TrackerMessage Tracker UMP Message with extrinsic_index 4310901-13
func TrackerMessage() {
	event, err := tracker.TrackXcmMessage("4310901-13", tx.UMP, "wss://moonbeam.api.onfinality.io/public-ws", "wss://polkadot.api.onfinality.io/public-ws", "")
	fmt.Println(event, err)
}
```


#### Xcm Client

```go
package example

import (
	"github.com/gmajor-encrypt/xcm-tools/tx"
	"github.com/itering/substrate-api-rpc/keyring"
)

func New() *tx.Client {
	endpoint := "endpoint" // need set endpoint
	client := tx.NewClient(endpoint)
	client.SetKeyRing("...") //  need set secret key
	return client
}
```

#### Send Xcm Transfer(simplified)

We provide the following methods to simplify the parameters of sending xcm transfer, so that there is no need to
construct complex **multiLocation** and **multiAssets**.

```go
package example

import (
	"fmt"
	"github.com/gmajor-encrypt/xcm-tools/tx"
	"github.com/itering/substrate-api-rpc/keyring"
	"github.com/shopspring/decimal"
)

// SendUmpMessage
// Send ump message from asset-hub to rococo relay chain
func SendUmpMessage() {
	client := New()
	beneficiary := "beneficiary Account id"
	transferAmount := decimal.New(10000, 0) // transfer amount 
	txHash, err := client.SendUmpTransfer(beneficiary, transferAmount)
	fmt.Println(txHash, err)
}

// SendHrmpMessage
// Send hrmp message from rococo asset-hub to picasso-rococo testnet
func SendHrmpMessage() {
	client := New()
	destParaId := 2087 // destination parachain id
	beneficiary := "beneficiary Account id"
	transferAmount := decimal.New(10000, 0) // transfer amount 
	txHash, err := client.SendHrmpTransfer(destParaId, beneficiary, transferAmount)
	fmt.Println(txHash, err)
}

// SendDmpMessage 
// Send dmp message from rococo to asset-hub 
func SendDmpMessage() {
	client := New()
	destParaId := 1000 // destination parachain id, rococo asset hub
	beneficiary := "beneficiary Account id"
	transferAmount := decimal.New(10000, 0) // transfer amount 
	txHash, err := client.SendDmpTransfer(destParaId, beneficiary, transferAmount)
	fmt.Println(txHash, err)
}

```

### Send Xcm Message

#### Support Methods

* `LimitedReserveTransferAssets`
* `LimitedTeleportAssets`
* `TeleportAssets`
* `ReserveTransferAssets`
* `Send`

#### Example

limited_reserve_transfer_assets, transfer assets from asset-hub to relay chain(rococo)

```go
package main

import (
	. "github.com/gmajor-encrypt/xcm-tools/tx"
	"github.com/shopspring/decimal"
	"log"
)

func main() {
	endpoint := "endpoint" // you need set an endpoint

	client := NewClient(endpoint)
	client.SetKeyRing(".....") // you need set a sr25519 secret key

	// dest is a relay chain
	dest := VersionedMultiLocation{V2: &V1MultiLocation{Interior: V0MultiLocation{Here: "NULL"}, Parents: 1}}

	// beneficiary is a relay chain account id
	beneficiary := VersionedMultiLocation{V2: &V1MultiLocation{Interior: V0MultiLocation{
		X1: &XCMJunction{AccountId32: &XCMJunctionAccountId32{Network: Enum{"Any": "NULL"}, Id: "0xd43593c715fdd31c61141abd04a99fd6822c8558854ccde39a5684e7a56da27d"}}},
		Parents: 0,
	}}

	// transfer amount
	amount := decimal.New(1, 0)

	// multi assets
	assets := MultiAssets{V2: []V2MultiAssets{{
		Id:  AssetsId{Concrete: &V1MultiLocation{Interior: V0MultiLocation{Here: "NULL"}, Parents: 1}},
		Fun: AssetsFun{Fungible: &amount},
	}}}

	// weight set
	weight := Weight{Limited: &WeightLimited{ProofSize: 0, RefTime: 4000000000}}

	// send an ump message use limited_reserve_transfer_assets
	callName, args := client.Ump.LimitedTeleportAssets(&dest, &beneficiary, &assets, 0, &weight)

	// sign the extrinsic
	signed, err := client.Conn.SignTransaction(client.Ump.GetModuleName(), callName, args...)

	if err != nil {
		log.Panic(err)
	}
	_, err = client.Conn.SendAuthorSubmitAndWatchExtrinsic(signed)
	if err != nil {
		log.Panic(err)
	}
}

```

More examples can be found in the [example](./example) or [xcm_test](./tx/xcm_test.go) directory.

### Test

```bash
go test -v ./...
```

docker

```bash
docker build -t xcm-tools-test .
docker run -it --rm xcm-tools-test
```

## License

The package is available as open source under the terms of
the [Apache LICENSE-2.0](https://www.apache.org/licenses/LICENSE-2.0)