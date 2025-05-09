#!/bin/bash
set -e

# Colors for terminal output
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Default number of blocks to generate
NUM_BLOCKS=5

# Check if a custom number of blocks was specified
if [ ! -z "$1" ] && [[ "$1" =~ ^[0-9]+$ ]]; then
    NUM_BLOCKS=$1
fi

echo -e "${BLUE}=== Generating $NUM_BLOCKS Blocks for Regtest Wallet ===${NC}"

# Bitcoin CLI command prefix
BITCOIN_CLI="docker exec bitcoind bitcoin-cli -regtest -rpcuser=bitcoin -rpcpassword=bitcoin -rpcport=8332"

# Check if wallet exists
if ! $BITCOIN_CLI listwallets | grep -q "regtest-wallet"; then
    echo -e "${YELLOW}Creating regtest-wallet...${NC}"
    $BITCOIN_CLI createwallet "regtest-wallet"
    echo -e "${GREEN}Wallet created!${NC}"
else
    echo -e "${GREEN}Using existing regtest-wallet${NC}"
fi

# Get an address from regtest-wallet
echo -e "${BLUE}Getting address from regtest-wallet...${NC}"
REGTEST_ADDR=$($BITCOIN_CLI -rpcwallet=regtest-wallet getnewaddress)
echo -e "${GREEN}Address: $REGTEST_ADDR${NC}"

# Check initial balance
INITIAL_BALANCE=$($BITCOIN_CLI -rpcwallet=regtest-wallet getbalance)
echo -e "${BLUE}Initial balance: $INITIAL_BALANCE BTC${NC}"

# Generate blocks
echo -e "${BLUE}Generating $NUM_BLOCKS blocks to address $REGTEST_ADDR...${NC}"
BLOCK_HASHES=$($BITCOIN_CLI -rpcwallet=regtest-wallet generatetoaddress $NUM_BLOCKS "$REGTEST_ADDR")

# Print block hashes
echo -e "${GREEN}Generated $NUM_BLOCKS blocks:${NC}"
echo "$BLOCK_HASHES" | jq -r '.[]' | while read -r hash; do
    echo -e "  ${GREEN}→ $hash${NC}"
done

# Wait a moment for the wallet to register the new blocks
sleep 1

# Check new balance
NEW_BALANCE=$($BITCOIN_CLI -rpcwallet=regtest-wallet getbalance)
echo -e "${BLUE}New balance: $NEW_BALANCE BTC${NC}"

# Calculate reward
REWARD=$(echo "$NEW_BALANCE - $INITIAL_BALANCE" | bc)
echo -e "${GREEN}Mining reward: $REWARD BTC${NC}"

echo -e "${GREEN}Done!${NC}"
