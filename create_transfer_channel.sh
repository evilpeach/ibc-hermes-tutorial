# !/bin/bash

echo "Create channel between A chain: $1 and B chain: $2"

hermes create channel --a-chain $1 --b-chain $2 --a-port transfer --b-port transfer --new-client-connection --yes
