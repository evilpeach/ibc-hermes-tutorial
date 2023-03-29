# !/bin/bash

hostChainId = "beebchain"
referenceChainId = "osmo-test-4"

# first way
hermes create client --host-chain $hostChainId --reference-chain $referenceChainId
hermes create client --host-chain $referenceChainId --reference-chain $hostChainId

hermes create connection --a-chain osmo-test-4 --a-client 07-tendermint-4379 --b-client 07-tendermint-392

hermes create channel --a-chain osmo-test-4 --a-connection connection-3795 --a-port transfer --b-port transfer

# wasm (old connection)
hermes create channel --a-chain beebchain  --a-port wasm.osmo14hj2tavq8fpesdwxxcu44rty3hh90vhujrvcmstl4zr3txmfvw9sq2r9g9 --b-port wasm.osmo1tplpap2dze4yejnaecu6t9qt0azpvt4p6ttcxrmvlawlx9m7gdzslldfnp --a-connection connection-3 --channel-version ics20-1 --yes

# second way(new connection)
# transfer
hermes create channel --a-chain osmo-test-4 --b-chain beebchain --a-port transfer --b-port transfer --new-client-connection

# wasm 
hermes create channel --a-chain beebchain --b-chain osmo-test-4  --a-port wasm.osmo14hj2tavq8fpesdwxxcu44rty3hh90vhujrvcmstl4zr3txmfvw9sq2r9g9 --b-port wasm.osmo1tplpap2dze4yejnaecu6t9qt0azpvt4p6ttcxrmvlawlx9m7gdzslldfnp --new-client-connection --channel-version ics20-1 --yes

# channel close init (send to A)
hermes tx chan-close-init --dst-chain osmo-test-4 --src-chain beebchain --dst-connection connection-3812 --dst-port wasm.osmo1tplpap2dze4yejnaecu6t9qt0azpvt4p6ttcxrmvlawlx9m7gdzslldfnp --src-port wasm.osmo14hj2tavq8fpesdwxxcu44rty3hh90vhujrvcmstl4zr3txmfvw9sq2r9g9 --dst-channel channel-3291 --src-channel channel-2
hermes tx chan-close-init --dst-chain beebchain --src-chain osmo-test-4 --dst-connection connection-3 --dst-port wasm.osmo14hj2tavq8fpesdwxxcu44rty3hh90vhujrvcmstl4zr3txmfvw9sq2r9g9 --src-port wasm.osmo1tplpap2dze4yejnaecu6t9qt0azpvt4p6ttcxrmvlawlx9m7gdzslldfnp --dst-channel channel-3 --src-channel channel-3292

# channel close confirm (send to B)
# hermes auto send this to B
hermes tx chan-close-confirm --dst-chain beebchain --src-chain osmo-test-4 --dst-connection connection-3 --dst-port wasm.osmo14hj2tavq8fpesdwxxcu44rty3hh90vhujrvcmstl4zr3txmfvw9sq2r9g9 --src-port wasm.osmo1tplpap2dze4yejnaecu6t9qt0azpvt4p6ttcxrmvlawlx9m7gdzslldfnp --dst-channel channel-2 --src-channel channel-3291