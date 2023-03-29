package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/joho/godotenv"
	"github.com/osmosis-labs/osmosis/v15/app"

	ibctypes "github.com/cosmos/ibc-go/v4/modules/apps/transfer/types"
)

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found")
	}
}

func main() {
	// Set prefix
	sdk.GetConfig().SetBech32PrefixForAccount(os.Getenv("PREFIX"), "")
	sendTx()
}

func sendTx() {
	MNEMONIC := os.Getenv("MNEMONIC")
	RPC_PATH := os.Getenv("RPC_PATH")
	CHAIN_ID := os.Getenv("CHAIN_ID")
	GAS := os.Getenv("GAS")
	FEE := os.Getenv("FEE")

	gasInt, err := strconv.ParseInt(GAS, 10, 64)
	if err != nil {
		panic(err)
	}

	// Setup keyring
	kb := keyring.NewInMemory()
	path := sdk.GetConfig().GetFullBIP44Path()
	key, err := kb.NewAccount("alice", MNEMONIC, "", path, hd.Secp256k1)
	if err != nil {
		panic(err)
	}

	println("Account:", key.GetAddress().String())

	// Create client
	clientNode, err := client.NewClientFromNode(RPC_PATH)
	if err != nil {
		panic(err)
	}
	clientCtx := client.Context{
		Client:            clientNode,
		ChainID:           CHAIN_ID,
		NodeURI:           RPC_PATH,
		InterfaceRegistry: app.MakeEncodingConfig().InterfaceRegistry,
		TxConfig:          app.MakeEncodingConfig().TxConfig,
		Keyring:           kb,
	}

	// Retrieve account info
	accountRetriever := authtypes.AccountRetriever{}
	acc, err := accountRetriever.GetAccount(clientCtx, key.GetAddress())
	if err != nil {
		panic(err)
	}

	var sdkMsg sdk.Msg

	// Crafting message
	switch os.Args[1] {
		case "transfer":
			channelId := os.Args[2]

			if channelId == "" {
				panic("Channel ID is required when use transfer command")
			}

			msg := ibctypes.MsgTransfer{
				SourcePort:    "transfer",
				SourceChannel: channelId,
				Token:         sdk.NewCoin("uosmo", sdk.NewInt(10000)),
				Sender:        key.GetAddress().String(),
				Receiver:      "osmo1cvve4mvnxhhqd6ywgfs9047gyg4xkvfgdk05la",
				TimeoutTimestamp: uint64(time.Now().Add(time.Hour * 24).UnixNano()),
			} 
			sdkMsg = sdk.Msg(&msg)
		case "execute":
			channelId := os.Args[2]

			if channelId == "" {
				panic("Channel ID is required when use ibc execute command")
			}

			executeMsg := fmt.Sprintf(`{"send_packet":{"channel_id":"%s","function":"execute"}}`, channelId)

			msg := wasmtypes.MsgExecuteContract{
				Sender: key.GetAddress().String(),
				Contract: "osmo1tplpap2dze4yejnaecu6t9qt0azpvt4p6ttcxrmvlawlx9m7gdzslldfnp",
				Msg: []byte(executeMsg),
			}
			sdkMsg = sdk.Msg(&msg)
		default:
			panic("First argument must be 'transfer' or 'execute'")
	}
	
	// Create transaction factory
	txf := tx.Factory{}.
		WithKeybase(kb).
		WithTxConfig(app.MakeEncodingConfig().TxConfig).
		WithAccountRetriever(clientCtx.AccountRetriever).
		WithAccountNumber(acc.GetAccountNumber()).
		WithSequence(acc.GetSequence()).
		WithGas(uint64(gasInt)).WithGasAdjustment(2).
		WithChainID(CHAIN_ID).
		WithMemo("").
		WithFees(FEE).
		WithSignMode(signing.SignMode_SIGN_MODE_DIRECT)

	txb, err := tx.BuildUnsignedTx(txf, sdkMsg)
	if err != nil {
		panic(err)
	}

	err = tx.Sign(txf, key.GetName(), txb, true)
	if err != nil {
		panic(err)
	}

	txBytes, err := clientCtx.TxConfig.TxEncoder()(txb.GetTx())
	if err != nil {
		panic(err)
	}

	res, err := clientCtx.BroadcastTxSync(txBytes)
	if err != nil {
		panic(err)
	}

	fmt.Println("response", res)
}