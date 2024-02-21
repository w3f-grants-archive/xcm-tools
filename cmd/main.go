package main

import (
	"fmt"
	"github.com/gmajor-encrypt/xcm-tools/parse"
	"github.com/gmajor-encrypt/xcm-tools/tracker"
	"github.com/gmajor-encrypt/xcm-tools/tx"
	"github.com/itering/scale.go/utiles"
	"github.com/shopspring/decimal"
	"github.com/urfave/cli"
	"log"
	"os"
)

func main() {
	app := setup()
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func setup() *cli.App {
	app := cli.NewApp()
	app.Name = "Xcm tools"
	app.Usage = "Xcm tools"
	app.Action = func(*cli.Context) error { return nil }
	app.Flags = []cli.Flag{}
	app.Commands = subCommands()
	return app
}

func subCommands() []cli.Command {
	var sendFlag = []cli.Flag{
		cli.StringFlag{
			Name:     "dest",
			Usage:    "dest address",
			Required: true,
		},
		cli.StringFlag{
			Name:     "amount",
			Usage:    "send xcm transfer amount",
			Required: true,
		},
		cli.StringFlag{
			Name:     "keyring",
			Usage:    "Set sr25519 secret key",
			EnvVar:   "SK",
			Required: true,
		},
		cli.StringFlag{
			Name:     "endpoint",
			Usage:    "Set substrate endpoint, only support websocket protocol, like ws:// or wss://",
			EnvVar:   "ENDPOINT",
			Required: true,
		},
	}
	return []cli.Command{
		{
			Name:  "send",
			Usage: "send xcm message",
			Subcommands: []cli.Command{
				{
					Name:  "UMP",
					Usage: "send ump message",
					Flags: sendFlag,
					Action: func(c *cli.Context) error {
						fmt.Println("xxx")
						client := tx.NewClient(c.String("endpoint"))
						defer client.Close()
						client.SetKeyRing(c.String("keyring"))
						beneficiary := c.String("dest")
						transferAmount := decimal.RequireFromString(c.String("amount"))
						txHash, err := client.SendUmpTransfer(beneficiary, transferAmount)
						if err != nil {
							return err
						}
						log.Print("send UMP message success, tx hash: ", txHash)
						return nil
					},
				},
				{
					Name:  "HRMP",
					Usage: "send hrmp message",
					Flags: append([]cli.Flag{
						cli.IntFlag{
							Name:     "paraId",
							Usage:    "dest para id",
							Required: true,
						},
					}, sendFlag...),
					Action: func(c *cli.Context) error {
						client := tx.NewClient(c.String("endpoint"))
						defer client.Close()
						client.SetKeyRing(c.String("keyring"))
						beneficiary := c.String("dest")
						destParaId := c.Int("paraId")
						transferAmount := decimal.RequireFromString(c.String("amount"))
						txHash, err := client.SendHrmpTransfer(uint32(destParaId), beneficiary, transferAmount)
						if err != nil {
							return err
						}
						log.Print("send HRMP message success, tx hash: ", txHash)
						return nil
					},
				},
				{
					Name:  "DMP",
					Usage: "send dmp message",
					Flags: append([]cli.Flag{
						cli.IntFlag{
							Name:     "paraId",
							Usage:    "dest para id",
							Required: true,
						},
					}, sendFlag...),
					Action: func(c *cli.Context) error {
						client := tx.NewClient(c.String("endpoint"))
						defer client.Close()
						client.SetKeyRing(c.String("keyring"))
						beneficiary := c.String("dest")
						destParaId := c.Int("paraId")
						transferAmount := decimal.RequireFromString(c.String("amount"))
						txHash, err := client.SendDmpTransfer(uint32(destParaId), beneficiary, transferAmount)
						if err != nil {
							return err
						}
						log.Print("send HRMP message success, tx hash: ", txHash)
						return nil
					},
				},
			},
		},
		{
			Name:  "parse",
			Usage: "parse xcm message",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:     "message",
					Usage:    "xcm message raw data",
					Required: true,
				},
				cli.StringFlag{
					Name:     "endpoint",
					Usage:    "Set substrate endpoint, only support websocket protocol, like ws:// or wss://",
					EnvVar:   "ENDPOINT",
					Required: true,
				},
			},
			Action: func(c *cli.Context) error {
				client := tx.NewClient(c.String("endpoint"))
				defer client.Close()
				p := parse.New(client.Metadata())
				xcm, err := p.ParseXcmMessageInstruction(c.String("message"))
				if err != nil {
					return err
				}
				log.Println("parse xcm message success: ")
				utiles.Debug(xcm)
				return nil
			},
		},
		{
			Name:  "tracker",
			Usage: "tracker xcm message transaction",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:     "extrinsicIndex",
					Usage:    "xcm message extrinsicIndex",
					Required: true,
				},
				cli.StringFlag{
					Name:     "protocol",
					Usage:    "xcm protocol, such as UMP,HRMP,DMP",
					Required: true,
				},
				cli.StringFlag{
					Name:     "destEndpoint",
					Usage:    "dest endpoint, only support websocket protocol, like ws:// or wss://",
					Required: true,
				},
				cli.StringFlag{
					Name:     "relaychainEndpoint",
					Usage:    "relay chain endpoint, only support websocket protocol, like ws:// or wss://",
					Required: false,
				},
				cli.StringFlag{
					Name:     "endpoint",
					Usage:    "Set substrate endpoint, only support websocket protocol, like ws:// or wss://",
					EnvVar:   "ENDPOINT",
					Required: true,
				},
			},
			Action: func(c *cli.Context) error {
				_, err := tracker.TrackXcmMessage(c.String("extrinsicIndex"), tx.Protocol(c.String("protocol")), c.String("endpoint"), c.String("destEndpoint"), c.String("relaychainEndpoint"))
				if err != nil {
					return err
				}
				return nil
			},
		},
	}
}
