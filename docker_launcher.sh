#!/bin/sh
COMMAND="galtcoin --data-dir $DATA_DIR --wallet-dir $WALLET_DIR $@"

adduser -D -u 10000 galtcoin

if [[ \! -d $DATA_DIR ]]; then
    mkdir -p $DATA_DIR
fi
if [[ \! -d $WALLET_DIR ]]; then
    mkdir -p $WALLET_DIR
fi

chown -R galtcoin:galtcoin $( realpath $DATA_DIR )
chown -R galtcoin:galtcoin $( realpath $WALLET_DIR )

su galtcoin -c "$COMMAND"
