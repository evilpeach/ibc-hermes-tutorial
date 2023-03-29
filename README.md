# Using Hermes

Hermes is a relayer for IBC (Inter-Blockchain Communication) in Rust and Golang. Here's how to use Hermes to transfer tokens between two chains.

## Prerequisites

Before using Hermes, you need to install Rust and Golang on your machine. Follow this guide from the Hermes documentation for detailed instructions on how to install the prerequisites: https://hermes.informal.systems/quick-start/pre-requisites.html

## Installation

To install Hermes, follow the instructions in this guide: https://hermes.informal.systems/quick-start/installation.html

## Configuration

After installing Hermes, create a `.hermes` directory in your root directory and copy `config.toml` to that folder. The path should look like this: `~/.hermes/config.toml`

You also need to add a relayer key to each chain. Follow the instructions below for adding the relayer key to Osmosis testnet and Beebchain.

### .ENV example

```
MNEMONIC=word1 word2 ... word24
RPC_PATH=https://rpc-test.osmosis.zone:443
GAS=1000000
FEE=25009uosmo
PREFIX=osmo
CHAIN_ID=osmo-test-4
```

### Adding a relayer key to Osmosis testnet

To add the relayer key to Osmosis testnet, run the following command:

```bash
hermes keys add --key-name relayer --chain osmo-test-4 --mnemonic-file mnemonic_file_hub.json
```

### Adding a relayer key to Beebchain

To add the relayer key to Beebchain, run the following command:

```bash
hermes keys add --key-name relayer --chain beebchain --mnemonic-file mnemonic_file_hub.json
```

## Health check

Run `hermes health-check` to check that both clients are ready to use.

## Starting Hermes

To start Hermes, run the following command:

```bash
hermes start
```

## Export variables

```bash
export OSMO=osmo-test-4
export BEEB=beebchain
```

## Creating new client, connection, and channel

To create a new client, connection, and channel on the transfer port in both chains, run the following command:

```bash
bash create_transfer_channel.sh $OSMO $BEEB
```

Get channel ID from result and transfer tokens through IBC, run the following command:

```bash
go run main.go transfer channel-3374
```

To create a new client, connection, and channel on the wasm port in both chains, run the following command:

```bash
bash create_wasm_channel.sh
```

Get channel ID from result and execute through IBC, run the following command:

```bash
go run main.go execute channel-3376
```

## Close channel

```bash
# Send MsgChannelCloseInit on Osmosis testnet
hermes tx chan-close-init --dst-chain osmo-test-4 --src-chain beebchain --dst-connection connection-3812 --dst-port wasm.osmo1tplpap2dze4yejnaecu6t9qt0azpvt4p6ttcxrmvlawlx9m7gdzslldfnp --src-port wasm.osmo14hj2tavq8fpesdwxxcu44rty3hh90vhujrvcmstl4zr3txmfvw9sq2r9g9 --dst-channel channel-3291 --src-channel channel-2

# Send MsgChannelCloseInit on Beebchain then Hermes will send MsgChannelCloseConfirm on Osmosis testnet
hermes tx chan-close-init --dst-chain beebchain --src-chain osmo-test-4 --dst-connection connection-3 --dst-port wasm.osmo14hj2tavq8fpesdwxxcu44rty3hh90vhujrvcmstl4zr3txmfvw9sq2r9g9 --src-port wasm.osmo1tplpap2dze4yejnaecu6t9qt0azpvt4p6ttcxrmvlawlx9m7gdzslldfnp --dst-channel channel-3 --src-channel channel-3292
```
