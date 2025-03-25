#!/bin/bash

# Script to check the balance of a Bitcoin address using Electrs REST API
# Usage: ./check-address-balance.sh <address>

set -e

# Define colors
GREEN='\033[0;32m'
RED='\033[0;31m'
BLUE='\033[0;34m'
YELLOW='\033[0;33m'
NC='\033[0m' # No Color

# Check if address argument is provided
if [ $# -ne 1 ]; then
    echo -e "${RED}Error: Missing required argument${NC}"
    echo -e "Usage: $0 <address>"
    exit 1
fi

ADDRESS=$1

# Improved validation for different Bitcoin address formats
# Support legacy, P2SH, and Bech32 address formats
if [[ "$ADDRESS" =~ ^([123]|bc1|tb1|bcrt1) ]]; then
    # Minimum length check - most addresses are at least 26 chars
    if [ ${#ADDRESS} -lt 26 ]; then
        echo -e "${RED}Warning: Address seems unusually short${NC}"
        echo -e "${RED}Continuing anyway...${NC}"
    fi
else
    echo -e "${RED}Warning: Address doesn't start with known Bitcoin prefixes${NC}"
    echo -e "${RED}Known prefixes: 1, 2, 3 (legacy/P2SH), bc1/tb1/bcrt1 (Bech32)${NC}"
    echo -e "${RED}Continuing anyway...${NC}"
fi

echo -e "${BLUE}Checking balance for address: ${YELLOW}$ADDRESS${NC}"

# Electrs REST API URL (standard Electrum protocol port)
ELECTRUM_API="http://localhost:50001"

# Direct way to get address balance in Electrs
echo -e "${BLUE}Retrieving balance...${NC}"
BALANCE_RESPONSE=$(curl -s -X GET "${ELECTRUM_API}/address/${ADDRESS}")

# Check if the API call was successful
if [ $? -ne 0 ] || [[ "$BALANCE_RESPONSE" == *"error"* ]]; then
    echo -e "${RED}Error: Failed to retrieve balance information.${NC}"
    echo -e "${RED}Response: $BALANCE_RESPONSE${NC}"
    
    # Try JSON-RPC method as a fallback
    echo -e "${BLUE}Attempting fallback method...${NC}"
    FALLBACK_RESPONSE=$(curl -s -X POST \
        -H "Content-Type: application/json" \
        -d "{\"jsonrpc\":\"2.0\",\"id\":\"check-balance\",\"method\":\"blockchain.address.get_balance\",\"params\":[\"$ADDRESS\"]}" \
        $ELECTRUM_API)
    
    if [ $? -ne 0 ] || [[ "$FALLBACK_RESPONSE" == *"error"* ]]; then
        echo -e "${RED}Fallback method also failed.${NC}"
        echo -e "${RED}Response: $FALLBACK_RESPONSE${NC}"
        exit 1
    else
        # Try to extract confirmed and unconfirmed balance from the response
        CONFIRMED=$(echo "$FALLBACK_RESPONSE" | grep -o '"confirmed":[0-9]*' | sed 's/"confirmed"://')
        UNCONFIRMED=$(echo "$FALLBACK_RESPONSE" | grep -o '"unconfirmed":[0-9]*' | sed 's/"unconfirmed"://')
        
        if [ -n "$CONFIRMED" ]; then
            CONFIRMED_BTC=$(echo "scale=8; $CONFIRMED / 100000000" | bc -l)
            CONFIRMED_BTC=$(echo $CONFIRMED_BTC | sed 's/^\./0./')
            echo -e "${GREEN}Confirmed balance: ${YELLOW}$CONFIRMED_BTC BTC${NC}"
        fi
        
        if [ -n "$UNCONFIRMED" ]; then
            UNCONFIRMED_BTC=$(echo "scale=8; $UNCONFIRMED / 100000000" | bc -l)
            UNCONFIRMED_BTC=$(echo $UNCONFIRMED_BTC | sed 's/^\./0./')
            echo -e "${BLUE}Unconfirmed balance: ${YELLOW}$UNCONFIRMED_BTC BTC${NC}"
        fi
        
        exit 0
    fi
fi

# Extract values from chain_stats and mempool_stats
CHAIN_FUNDED=$(echo "$BALANCE_RESPONSE" | grep -o '"chain_stats":{[^}]*"funded_txo_sum":[0-9]*' | grep -o '"funded_txo_sum":[0-9]*' | sed 's/"funded_txo_sum"://')
CHAIN_SPENT=$(echo "$BALANCE_RESPONSE" | grep -o '"chain_stats":{[^}]*"spent_txo_sum":[0-9]*' | grep -o '"spent_txo_sum":[0-9]*' | sed 's/"spent_txo_sum"://')
CHAIN_FUNDED_COUNT=$(echo "$BALANCE_RESPONSE" | grep -o '"chain_stats":{[^}]*"funded_txo_count":[0-9]*' | grep -o '"funded_txo_count":[0-9]*' | sed 's/"funded_txo_count"://')
CHAIN_SPENT_COUNT=$(echo "$BALANCE_RESPONSE" | grep -o '"chain_stats":{[^}]*"spent_txo_count":[0-9]*' | grep -o '"spent_txo_count":[0-9]*' | sed 's/"spent_txo_count"://')
CHAIN_TX_COUNT=$(echo "$BALANCE_RESPONSE" | grep -o '"chain_stats":{[^}]*"tx_count":[0-9]*' | grep -o '"tx_count":[0-9]*' | sed 's/"tx_count"://')

MEMPOOL_FUNDED=$(echo "$BALANCE_RESPONSE" | grep -o '"mempool_stats":{[^}]*"funded_txo_sum":[0-9]*' | grep -o '"funded_txo_sum":[0-9]*' | sed 's/"funded_txo_sum"://')
MEMPOOL_SPENT=$(echo "$BALANCE_RESPONSE" | grep -o '"mempool_stats":{[^}]*"spent_txo_sum":[0-9]*' | grep -o '"spent_txo_sum":[0-9]*' | sed 's/"spent_txo_sum"://')
MEMPOOL_TX_COUNT=$(echo "$BALANCE_RESPONSE" | grep -o '"mempool_stats":{[^}]*"tx_count":[0-9]*' | grep -o '"tx_count":[0-9]*' | sed 's/"tx_count"://')

# Calculate confirmed and unconfirmed balances
if [ -n "$CHAIN_FUNDED" ] && [ -n "$CHAIN_SPENT" ]; then
    CONFIRMED_SATS=$((CHAIN_FUNDED - CHAIN_SPENT))
    CONFIRMED_BTC=$(echo "scale=8; $CONFIRMED_SATS / 100000000" | bc -l)
    CONFIRMED_BTC=$(echo $CONFIRMED_BTC | sed 's/^\./0./')
    UTXO_COUNT=$((CHAIN_FUNDED_COUNT - CHAIN_SPENT_COUNT))
    
    echo -e "${GREEN}Confirmed balance: ${YELLOW}$CONFIRMED_BTC BTC${NC}"
    echo -e "${BLUE}Confirmed UTXOs: ${YELLOW}$UTXO_COUNT${NC}"
    echo -e "${BLUE}Total transactions: ${YELLOW}$CHAIN_TX_COUNT${NC}"
    
    # Check if there are unconfirmed transactions
    if [ -n "$MEMPOOL_FUNDED" ] && [ -n "$MEMPOOL_SPENT" ] && [ "$MEMPOOL_TX_COUNT" -gt 0 ]; then
        UNCONFIRMED_SATS=$((MEMPOOL_FUNDED - MEMPOOL_SPENT))
        if [ $UNCONFIRMED_SATS -ne 0 ]; then
            UNCONFIRMED_BTC=$(echo "scale=8; $UNCONFIRMED_SATS / 100000000" | bc -l)
            UNCONFIRMED_BTC=$(echo $UNCONFIRMED_BTC | sed 's/^\./0./')
            TOTAL_BTC=$(echo "scale=8; ($CONFIRMED_SATS + $UNCONFIRMED_SATS) / 100000000" | bc -l)
            TOTAL_BTC=$(echo $TOTAL_BTC | sed 's/^\./0./')
            
            echo -e "${BLUE}Unconfirmed balance: ${YELLOW}$UNCONFIRMED_BTC BTC${NC}"
            echo -e "${BLUE}Pending transactions: ${YELLOW}$MEMPOOL_TX_COUNT${NC}"
            echo -e "${GREEN}Total balance (including unconfirmed): ${YELLOW}$TOTAL_BTC BTC${NC}"
        fi
    fi
    
    # If there are UTXOs, get their details
    if [ $UTXO_COUNT -gt 0 ]; then
        echo -e "${BLUE}Getting UTXO details...${NC}"
        UTXO_RESPONSE=$(curl -s -X GET "${ELECTRUM_API}/address/${ADDRESS}/utxo")
        
        if [ $? -eq 0 ] && [ -n "$UTXO_RESPONSE" ] && [[ "$UTXO_RESPONSE" != *"error"* ]]; then
            echo -e "${BLUE}UTXO Details:${NC}"
            echo "$UTXO_RESPONSE" | sed 's/\[//' | sed 's/\]//' | sed 's/},{/}\n{/g' | while read -r UTXO; do
                TXID=$(echo "$UTXO" | grep -o '"txid":"[^"]*' | sed 's/"txid":"//')
                VOUT=$(echo "$UTXO" | grep -o '"vout":[0-9]*' | sed 's/"vout"://')
                VALUE=$(echo "$UTXO" | grep -o '"value":[0-9]*' | sed 's/"value"://')
                STATUS=$(echo "$UTXO" | grep -o '"status":{[^}]*}' | sed 's/"status":{//' | sed 's/}//')
                BLOCK_HEIGHT=$(echo "$STATUS" | grep -o '"block_height":[0-9]*' | sed 's/"block_height"://')
                
                if [ -n "$VALUE" ]; then
                    VALUE_BTC=$(echo "scale=8; $VALUE / 100000000" | bc -l)
                    VALUE_BTC=$(echo $VALUE_BTC | sed 's/^\./0./')
                    if [ -n "$BLOCK_HEIGHT" ]; then
                        echo -e "${YELLOW}  - ${VALUE_BTC} BTC${NC} (txid: ${TXID}:${VOUT}, height: ${BLOCK_HEIGHT})"
                    else
                        echo -e "${YELLOW}  - ${VALUE_BTC} BTC${NC} (txid: ${TXID}:${VOUT}, unconfirmed)"
                    fi
                fi
            done
        else
            echo -e "${RED}Could not retrieve detailed UTXO information${NC}"
        fi
    fi
else
    echo -e "${BLUE}No balance information found for this address${NC}"
    echo -e "${BLUE}Balance: ${YELLOW}0 BTC${NC}"
fi
