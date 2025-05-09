#!/bin/bash
set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[0;33m'
NC='\033[0m' # No Color

echo -e "${BLUE}=== Funding Test Wallet ===${NC}"

# File paths for wallets
MINING_WALLET_FILE="./mining-wallet.txt"
REGTEST_WALLET_FILE="./regtest-wallet.txt"

# Bitcoin CLI command prefix
BITCOIN_CLI="docker exec bitcoind bitcoin-cli -regtest -rpcuser=bitcoin -rpcpassword=bitcoin -rpcport=8332"

# Function to load or create a wallet
load_or_create_wallet() {
    local wallet_name=$1
    
    # First try to load the wallet
    if $BITCOIN_CLI loadwallet "$wallet_name" 2>/dev/null; then
        echo -e "${GREEN}Loaded existing $wallet_name${NC}"
    else
        # If loading fails, try to create it
        if $BITCOIN_CLI createwallet "$wallet_name" 2>/dev/null; then
            echo -e "${GREEN}Created new $wallet_name${NC}"
        else
            # If both fail, check if it's already loaded
            if $BITCOIN_CLI listwallets | grep -q "$wallet_name"; then
                echo -e "${GREEN}Wallet $wallet_name is already loaded${NC}"
            else
                echo -e "${RED}Failed to load or create $wallet_name${NC}"
                return 1
            fi
        fi
    fi
    return 0
}

# Load or create wallets
echo -e "${BLUE}Setting up regtest-wallet...${NC}"
load_or_create_wallet "regtest-wallet"

echo -e "${BLUE}Setting up mining-wallet...${NC}"
load_or_create_wallet "mining-wallet"

# Create or load mining wallet address
if [ -f "$MINING_WALLET_FILE" ]; then
    MINING_ADDR=$(cat "$MINING_WALLET_FILE")
    echo -e "${GREEN}Loaded mining wallet address: $MINING_ADDR${NC}"
else
    # Get a new address from the mining wallet
    MINING_ADDR=$($BITCOIN_CLI -rpcwallet=mining-wallet getnewaddress)
    echo "$MINING_ADDR" > "$MINING_WALLET_FILE"
    echo -e "${GREEN}Created mining wallet address: $MINING_ADDR${NC}"
fi

# Create or load regtest wallet address
if [ -f "$REGTEST_WALLET_FILE" ]; then
    REGTEST_ADDR=$(cat "$REGTEST_WALLET_FILE")
    echo -e "${GREEN}Loaded regtest wallet address: $REGTEST_ADDR${NC}"
else
    REGTEST_ADDR=$($BITCOIN_CLI -rpcwallet=regtest-wallet getnewaddress)
    echo "$REGTEST_ADDR" > "$REGTEST_WALLET_FILE"
    echo -e "${GREEN}Created regtest wallet address: $REGTEST_ADDR${NC}"
fi

# Generate 101 blocks to mining wallet for coinbase maturity
echo -e "${BLUE}Generating 101 blocks to mining wallet ($MINING_ADDR)...${NC}"
$BITCOIN_CLI -rpcwallet=mining-wallet generatetoaddress 101 "$MINING_ADDR"

# Check mining wallet balance before transfer
echo -e "${BLUE}Checking mining wallet balance before transfer...${NC}"
MINING_BALANCE=$($BITCOIN_CLI -rpcwallet=mining-wallet getbalance)
echo -e "${GREEN}Mining wallet balance: $MINING_BALANCE BTC${NC}"

# Check if we have sufficient funds
if (( $(echo "$MINING_BALANCE < 20" | bc -l) )); then
    echo -e "${YELLOW}Insufficient funds in mining wallet. Generating more blocks...${NC}"
    $BITCOIN_CLI -rpcwallet=mining-wallet generatetoaddress 50 "$MINING_ADDR"
    
    # Check balance again
    MINING_BALANCE=$($BITCOIN_CLI -rpcwallet=mining-wallet getbalance)
    echo -e "${GREEN}New mining wallet balance: $MINING_BALANCE BTC${NC}"
    
    if (( $(echo "$MINING_BALANCE < 20" | bc -l) )); then
        echo -e "${RED}Still insufficient funds after generating more blocks.${NC}"
        echo -e "${RED}Try resetting the regtest environment or check for other issues.${NC}"
        exit 1
    fi
fi

# Transfer funds from mining wallet to regtest wallet
TRANSFER_AMOUNT=20
echo -e "${BLUE}Transferring $TRANSFER_AMOUNT BTC from mining wallet to regtest wallet ($REGTEST_ADDR)...${NC}"
TXID=$($BITCOIN_CLI -rpcwallet=mining-wallet sendtoaddress "$REGTEST_ADDR" $TRANSFER_AMOUNT)
echo -e "${GREEN}Transaction ID: $TXID${NC}"

# Generate 1 block to confirm the transaction
echo -e "${BLUE}Generating block to confirm the transaction...${NC}"
BLOCK=$($BITCOIN_CLI -rpcwallet=mining-wallet generatetoaddress 1 "$MINING_ADDR")
echo -e "${GREEN}Block generated: $BLOCK${NC}"

# Check balance of regtest-wallet
echo -e "${BLUE}Checking balance of regtest-wallet...${NC}"
BALANCE=$($BITCOIN_CLI -rpcwallet=regtest-wallet getbalance)
echo -e "${GREEN}Balance: $BALANCE BTC${NC}"

# Final balance checks
mining_balance=$($BITCOIN_CLI -rpcwallet=mining-wallet getbalance)
regtest_balance=$($BITCOIN_CLI -rpcwallet=regtest-wallet getbalance)

echo -e "${GREEN}Mining wallet balance: $mining_balance${NC}"
echo -e "${GREEN}Regtest wallet balance: $regtest_balance${NC}"

echo -e "${GREEN}Done!${NC}"
