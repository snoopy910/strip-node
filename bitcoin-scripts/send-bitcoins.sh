#!/bin/bash

# Script to send bitcoins from the mining wallet to a specified address
# Usage: ./send-bitcoins.sh <address> <amount>

set -e

# Define colors
GREEN='\033[0;32m'
RED='\033[0;31m'
BLUE='\033[0;34m'
YELLOW='\033[0;33m'
NC='\033[0m' # No Color

# Check if both arguments are provided
if [ $# -ne 2 ]; then
    echo -e "${RED}Error: Missing required arguments${NC}"
    echo -e "Usage: $0 <address> <amount>"
    exit 1
fi

ADDRESS=$1
AMOUNT=$2

# Validate the amount is a valid number
if ! [[ $AMOUNT =~ ^[0-9]+(\.[0-9]+)?$ ]]; then
    echo -e "${RED}Error: Amount must be a valid number${NC}"
    exit 1
fi

# Check wallet balance before sending
echo -e "${BLUE}Checking wallet balance...${NC}"
BALANCE=$(docker exec bitcoind bitcoin-cli -regtest -rpcuser=bitcoin -rpcpassword=bitcoin -rpcport=8332 -rpcwallet=regtest-wallet getbalance)

echo -e "${BLUE}Current wallet balance: ${YELLOW}$BALANCE BTC${NC}"

# Check if there are sufficient funds
if (( $(echo "$BALANCE < $AMOUNT" | bc -l) )); then
    echo -e "${RED}Error: Insufficient funds. You are trying to send ${YELLOW}$AMOUNT BTC${RED} but only have ${YELLOW}$BALANCE BTC${RED} available.${NC}"
    exit 1
fi

echo -e "${BLUE}Sending ${YELLOW}$AMOUNT BTC${BLUE} to address ${YELLOW}$ADDRESS${BLUE}...${NC}"

# Execute the bitcoin-cli command to send the bitcoins
TXID=$(docker exec bitcoind bitcoin-cli -regtest -rpcuser=bitcoin -rpcpassword=bitcoin -rpcport=8332 -rpcwallet=regtest-wallet sendtoaddress "$ADDRESS" "$AMOUNT")

if [ $? -eq 0 ]; then
    echo -e "${GREEN}Transaction successful!${NC}"
    echo -e "${BLUE}Transaction ID: ${YELLOW}$TXID${NC}"

    # Check the new balance
    NEW_BALANCE=$(docker exec bitcoind bitcoin-cli -regtest -rpcuser=bitcoin -rpcpassword=bitcoin -rpcport=8332 -rpcwallet=regtest-wallet getbalance)
    echo -e "${BLUE}New wallet balance: ${YELLOW}$NEW_BALANCE BTC${NC}"

    # Calculate and display the difference
    DIFF=$(echo "$BALANCE - $NEW_BALANCE" | bc -l)
    echo -e "${BLUE}Total spent (including fees): ${YELLOW}$DIFF BTC${NC}"
else
    echo -e "${RED}Transaction failed!${NC}"
    exit 1
fi
