#!/bin/bash
set -e

# Configuration
# Amount in BTC
AMOUNT_BTC=1.0
# Amount in satoshis (100000000 satoshis = 1 BTC)
AMOUNT_SATS=100000000
# Fee in BTC
FEE=0.0001

# Bitcoin CLI command prefix
BITCOIN_CLI="docker exec bitcoind bitcoin-cli -regtest -rpcuser=bitcoin -rpcpassword=bitcoin -rpcport=8332 -rpcwallet=regtest-wallet"
# Electrs API URL
ELECTRS_URL="http://localhost:50001"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color
BOLD='\033[1m'

# Header function
header() {
  echo -e "${BOLD}${PURPLE}== $1 ==${NC}"
}

# Success function
success() {
  echo -e "${GREEN}✓ $1${NC}"
}

# Info function
info() {
  echo -e "${BLUE}ℹ $1${NC}"
}

# Warning function
warning() {
  echo -e "${YELLOW}⚠ $1${NC}"
}

# Error function
error() {
  echo -e "${RED}✗ $1${NC}"
}

# Value function for displaying important values
value() {
  echo -e "${CYAN}$1:${NC} $2"
}

# JSON Function to post data to Electrs
electrs_post() {
  curl -s -X POST -H "Content-Type: application/json" -d "$2" "$ELECTRS_URL$1"
}

# Function to try broadcasting using Electrs
broadcast_transaction() {
  local tx_hex="$1"

  # Method 1: Try Electrs with "tx" property
  info "Attempting to broadcast via Electrs (JSON format)..."
  local payload="{\"tx\": \"$tx_hex\"}"
  local result=$(electrs_post "/tx" "$payload")

  if ! echo "$result" | grep -q "error"; then
    success "Transaction broadcast successfully via Electrs API (JSON format)"
    value "Transaction ID" "$result"
    echo "$result"
    return
  fi

  warning "Electrs JSON broadcasting format failed"
  echo "Error: $result"

  # Method 2: Try Electrs with raw hex
  info "Attempting to broadcast via Electrs (raw hex format)..."
  local result=$(curl -s -X POST -H "Content-Type: text/plain" -d "$tx_hex" "$ELECTRS_URL/tx")

  if ! echo "$result" | grep -q "error"; then
    success "Transaction broadcast successfully via Electrs API (raw hex format)"
    value "Transaction ID" "$result"
    echo "$result"
    return
  fi

  error "Failed to broadcast transaction via Electrs. Error: $result"
  exit 1
}

header "Bitcoin Transaction Using Electrs API"

# Step 1: Make sure we have a wallet
info "Checking for wallet..."
if ! $BITCOIN_CLI listwallets | grep -q "regtest-wallet"; then
  warning "Creating wallet 'regtest-wallet'..."
  $BITCOIN_CLI createwallet "regtest-wallet"
  success "Wallet created"
else
  success "Wallet found"
fi

# Step 2: Create two addresses - one for sending and one for receiving
info "Creating addresses..."
SENDER_ADDRESS=$($BITCOIN_CLI getnewaddress "sender" "bech32")
RECEIVER_ADDRESS=$($BITCOIN_CLI getnewaddress "receiver" "bech32")

value "Sender address" "$SENDER_ADDRESS"
value "Receiver address" "$RECEIVER_ADDRESS"

# Step 3: Fund the sender address with some bitcoin
header "Funding Sender Address"
info "Funding sender address with 10 BTC..."
TXID=$($BITCOIN_CLI sendtoaddress "$SENDER_ADDRESS" 10.0)
success "Funding transaction sent"
value "Funding transaction ID" "$TXID"

# Step: Generate a block to confirm the funding transaction
info "Generating block to confirm funding transaction..."
BLOCK_HASH=$($BITCOIN_CLI generatetoaddress 1 "$($BITCOIN_CLI getnewaddress)")
success "Block generated: ${BLOCK_HASH[0]}"

# Wait for Electrs to index the new block - use UTXO endpoint directly instead of address balance
info "Waiting for Electrs to index the funding transaction..."
wait_for_utxo=0
for i in {1..30}; do
  UTXOS=$(curl -s $ELECTRS_URL/address/$SENDER_ADDRESS/utxo)
  if [ -n "$UTXOS" ] && [ "$UTXOS" != "[]" ]; then
    wait_for_utxo=1
    break
  fi
  spinner="⠋⠙⠹⠸⠼⠴⠦⠧⠇⠏"
  echo -ne "\r${YELLOW}Waiting for Electrs to index UTXOs ($i/30) ${spinner:i%10:1}${NC}"
  sleep 2
done

if [ $wait_for_utxo -eq 0 ]; then
  error "Failed to find UTXOs in Electrs after 30 attempts. Try increasing the timeout."
  exit 1
fi
success "Successfully found UTXOs using Electrs"

# Step: Check balances after funding
header "Balances After Funding (Electrs)"

# Get UTXOs for sender address using Electrs API
info "Getting UTXOs for sender address using Electrs API..."
UTXOS=$(curl -s $ELECTRS_URL/address/$SENDER_ADDRESS/utxo)

if [ -z "$UTXOS" ] || [ "$UTXOS" = "[]" ]; then
  error "No UTXOs found with Electrs even after waiting. Something is wrong."
  exit 1
fi

# Parse UTXO from Electrs
TXID=$(echo $UTXOS | jq -r '.[0].txid')
VOUT=$(echo $UTXOS | jq -r '.[0].vout')
AMOUNT_IN_SATOSHIS=$(echo $UTXOS | jq -r '.[0].value')
AMOUNT_IN_BTC=$(bc <<< "scale=8; $AMOUNT_IN_SATOSHIS / 100000000")

value "Using UTXO" "$TXID:$VOUT"
value "UTXO amount" "$AMOUNT_IN_BTC BTC ($AMOUNT_IN_SATOSHIS satoshis)"

# Calculate change amount
CHANGE_AMOUNT_BTC=$(bc <<< "scale=8; $AMOUNT_IN_BTC - $AMOUNT_BTC - $FEE")
CHANGE_AMOUNT_SATOSHIS=$(bc <<< "scale=0; $CHANGE_AMOUNT_BTC * 100000000 / 1")

value "Change amount" "$CHANGE_AMOUNT_BTC BTC"
value "Amount to send" "$AMOUNT_SATS satoshis ($AMOUNT_BTC BTC)"
value "Change amount" "$CHANGE_AMOUNT_SATOSHIS satoshis"

# Step 9: Create a raw transaction
header "Creating Transaction"
info "Creating raw transaction..."
# Print the exact command for debugging
CREATE_TX_CMD="$BITCOIN_CLI createrawtransaction '[{\"txid\":\"$TXID\",\"vout\":$VOUT}]' '{\"$RECEIVER_ADDRESS\":$AMOUNT_BTC, \"$SENDER_ADDRESS\":$CHANGE_AMOUNT_BTC}'"
info "Running command: $CREATE_TX_CMD"

RAW_TX=$($BITCOIN_CLI createrawtransaction \
  "[{\"txid\":\"$TXID\",\"vout\":$VOUT}]" \
  "{\"$RECEIVER_ADDRESS\":$AMOUNT_BTC, \"$SENDER_ADDRESS\":$CHANGE_AMOUNT_BTC}")

success "Raw transaction created"
echo -e "${CYAN}Transaction (preview):${NC} ${RAW_TX:0:60}..."

# Show decoded raw transaction to verify it has inputs
info "Verifying transaction structure:"
$BITCOIN_CLI decoderawtransaction "$RAW_TX" | jq '{txid: .txid, vin: .vin | length, vout: .vout | length}'

# Step 10: Sign the transaction
info "Signing transaction..."
SIGNED_RESULT=$($BITCOIN_CLI signrawtransactionwithwallet "$RAW_TX")
SIGNED_TX=$(echo "$SIGNED_RESULT" | jq -r '.hex')
COMPLETE=$(echo "$SIGNED_RESULT" | jq -r '.complete')

if [ "$COMPLETE" != "true" ]; then
  error "Transaction signing incomplete!"
  echo "$SIGNED_RESULT" | jq
  exit 1
fi

success "Transaction signed successfully"
value "Signed transaction hex" "${SIGNED_TX:0:50}... (truncated)"

# Verify signed transaction
info "Verifying signed transaction structure:"
$BITCOIN_CLI decoderawtransaction "$SIGNED_TX" | jq '{txid: .txid, vin: .vin | length, vout: .vout | length}'

# Step 11: Broadcast the transaction using Electrs
header "Broadcasting Transaction"
TXID_RESULT=$(broadcast_transaction "$SIGNED_TX")
TXID_BITCOIND=$(echo "$TXID_RESULT" | tail -n 1)

# Verify we got a valid transaction ID
if [ -z "$TXID_BITCOIND" ] || [[ "$TXID_BITCOIND" == *"error"* ]]; then
  error "Failed to broadcast transaction: $TXID_BITCOIND"
  exit 1
fi

value "Transaction ID for confirmation" "$TXID_BITCOIND"

# Step 12: Generate a block to confirm the transaction
info "Generating block to confirm the transaction..."
CONFIRM_BLOCK=$($BITCOIN_CLI generatetoaddress 1 "$($BITCOIN_CLI getnewaddress)")
success "Block generated: ${CONFIRM_BLOCK[0]}"

# Wait for Electrs to index the new block with robust retry mechanism
info "Waiting for transaction to be confirmed in Electrs..."
wait_for_confirmation=0
for i in {1..40}; do
  TX_DETAILS=$(curl -s $ELECTRS_URL/tx/$TXID_BITCOIND)
  if [ -n "$TX_DETAILS" ] && [ "$TX_DETAILS" != "null" ]; then
    # Check if confirmed
    CONFIRMED=$(echo $TX_DETAILS | jq -r '.status.confirmed')
    if [ "$CONFIRMED" = "true" ]; then
      wait_for_confirmation=1
      success "Transaction confirmed after $i attempts"
      break
    fi
  fi
  spinner="⠋⠙⠹⠸⠼⠴⠦⠧⠇⠏"
  echo -ne "\r${YELLOW}Waiting for transaction confirmation ($i/40) ${spinner:i%10:1}${NC}"
  sleep 3
done

if [ $wait_for_confirmation -eq 0 ]; then
  error "Failed to confirm transaction in Electrs after 40 attempts."
  warning "Will continue anyway and check balances..."
fi

# Step 13: Check transaction confirmation status via Electrs
header "Transaction Status"
info "Checking transaction confirmation status via Electrs..."
CONFIRMED_TX=$(curl -s $ELECTRS_URL/tx/$TXID_BITCOIND)
if [ -z "$CONFIRMED_TX" ] || [ "$CONFIRMED_TX" = "null" ]; then
  warning "Transaction not found in Electrs yet. It may still be indexing."
else
  CONFIRMED=$(echo $CONFIRMED_TX | jq -r '.status.confirmed')
  if [ "$CONFIRMED" = "true" ]; then
    success "Transaction confirmed"
  else
    warning "Transaction not yet confirmed"
  fi

  # Generate one more block to ensure we have enough confirmations for Electrs
  info "Generating one more block to make sure Electrs has enough confirmations..."
  EXTRA_BLOCK=$($BITCOIN_CLI generatetoaddress 1 "$($BITCOIN_CLI getnewaddress)")
  success "Extra block generated: ${EXTRA_BLOCK[0]}"
  sleep 5  # Give Electrs some time to update

  # Step 14: Display final balances via Electrs - with advanced retry mechanism
  header "Final Balances (Electrs)"

  # Get the balance directly from Bitcoin Core to compare
  info "Getting balances directly from Bitcoin Core for comparison..."
  CORE_RECEIVER_UTXOS=$($BITCOIN_CLI listunspent 1 9999 "[\"$RECEIVER_ADDRESS\"]")
  CORE_RECEIVER_BALANCE=$(echo $CORE_RECEIVER_UTXOS | jq -r 'if length > 0 then .[0].amount * 100000000 else 0 end')
  value "Receiver balance according to Bitcoin Core" "$CORE_RECEIVER_BALANCE satoshis"

  # Wait for Electrs to update with more retries and a forceful approach
  info "Aggressively querying Electrs until receiver balance updates..."
  max_retries=40
  success=0

  for i in {1..40}; do
    # Query both the address endpoint and the UTXO endpoint
    RECEIVER_ADDRESS_DATA=$(curl -s $ELECTRS_URL/address/$RECEIVER_ADDRESS)
    RECEIVER_UTXOS=$(curl -s $ELECTRS_URL/address/$RECEIVER_ADDRESS/utxo)

    # Check UTXOs first (more reliable)
    if [ -n "$RECEIVER_UTXOS" ] && [ "$RECEIVER_UTXOS" != "[]" ]; then
      UTXO_AMOUNT=$(echo $RECEIVER_UTXOS | jq -r '.[0].value')
      if [ -n "$UTXO_AMOUNT" ] && [ "$UTXO_AMOUNT" -eq "$AMOUNT_SATS" ]; then
        success "Receiver UTXO found in Electrs with correct amount: $UTXO_AMOUNT satoshis"
        success=1
        break
      fi
    fi

    # Then check address stats
    if [ -n "$RECEIVER_ADDRESS_DATA" ]; then
      FUNDED=$(echo $RECEIVER_ADDRESS_DATA | jq -r '.chain_stats.funded_txo_sum')
      if [ -n "$FUNDED" ] && [ "$FUNDED" -gt 0 ]; then
        success "Receiver balance updated in Electrs: $FUNDED satoshis"
        success=1
        break
      fi
    fi

    spinner="⠋⠙⠹⠸⠼⠴⠦⠧⠇⠏"
    echo -ne "\r${YELLOW}Waiting for Electrs to update receiver balance ($i/$max_retries) ${spinner:i%10:1}${NC}"
    sleep 3
  done

  echo "" # New line after spinner

  # Display detailed information with more flexible error handling
  SENDER_FINAL=$(curl -s $ELECTRS_URL/address/$SENDER_ADDRESS)
  RECEIVER_FINAL=$(curl -s $ELECTRS_URL/address/$RECEIVER_ADDRESS)

  # Process sender data
  if [ -n "$SENDER_FINAL" ] && [ "$SENDER_FINAL" != "null" ]; then
    SENDER_FUNDED=$(echo $SENDER_FINAL | jq -r '.chain_stats.funded_txo_sum // 0')
    SENDER_SPENT=$(echo $SENDER_FINAL | jq -r '.chain_stats.spent_txo_sum // 0')
    SENDER_FINAL_BALANCE=$((SENDER_FUNDED - SENDER_SPENT))

    value "Sender funding transactions" "$SENDER_FUNDED satoshis"
    value "Sender spending transactions" "$SENDER_SPENT satoshis"
    value "Sender final balance" "$SENDER_FINAL_BALANCE satoshis"
  else
    warning "Could not retrieve sender balance from Electrs"
  fi

  # Process receiver data with UTXOs as backup
  if [ -n "$RECEIVER_FINAL" ] && [ "$RECEIVER_FINAL" != "null" ]; then
    RECEIVER_FUNDED=$(echo $RECEIVER_FINAL | jq -r '.chain_stats.funded_txo_sum // 0')
    RECEIVER_SPENT=$(echo $RECEIVER_FINAL | jq -r '.chain_stats.spent_txo_sum // 0')
    RECEIVER_FINAL_BALANCE=$((RECEIVER_FUNDED - RECEIVER_SPENT))

    value "Receiver funding transactions" "$RECEIVER_FUNDED satoshis"
    value "Receiver spending transactions" "$RECEIVER_SPENT satoshis"
    value "Receiver final balance" "$RECEIVER_FINAL_BALANCE satoshis"
  else
    warning "Could not retrieve receiver balance from Electrs address endpoint"

    # Check UTXOs as a backup method
    RECEIVER_UTXOS=$(curl -s $ELECTRS_URL/address/$RECEIVER_ADDRESS/utxo)
    if [ -n "$RECEIVER_UTXOS" ] && [ "$RECEIVER_UTXOS" != "[]" ]; then
      UTXO_AMOUNT=$(echo $RECEIVER_UTXOS | jq -r '[.[].value] | add // 0')
      value "Receiver balance from UTXOs" "$UTXO_AMOUNT satoshis"
      RECEIVER_FINAL_BALANCE=$UTXO_AMOUNT
    else
      warning "Could not retrieve receiver UTXOs from Electrs"
      if [ "$CORE_RECEIVER_BALANCE" -gt 0 ]; then
        info "Using Bitcoin Core balance as backup"
        RECEIVER_FINAL_BALANCE=$CORE_RECEIVER_BALANCE
      else
        RECEIVER_FINAL_BALANCE=0
      fi
    fi
  fi

  # Show balance changes - using 0 as initial since we checked at the beginning
  SENDER_CHANGE=$((SENDER_FINAL_BALANCE))
  RECEIVER_CHANGE=$((RECEIVER_FINAL_BALANCE))

  if [ $SENDER_CHANGE -lt 0 ]; then
    echo -e "${RED}Sender balance change: ${SENDER_CHANGE} satoshis${NC}"
  else
    echo -e "${GREEN}Sender balance change: +${SENDER_CHANGE} satoshis${NC}"
  fi

  if [ $RECEIVER_CHANGE -gt 0 ]; then
    echo -e "${GREEN}Receiver balance change: +${RECEIVER_CHANGE} satoshis${NC}"
  else
    echo -e "${RED}Receiver balance change: ${RECEIVER_CHANGE} satoshis${NC}"
  fi

  # Verify expected amounts with more diagnostics if there's a mismatch
  if [ "$RECEIVER_FINAL_BALANCE" -eq "$AMOUNT_SATS" ]; then
    success "Receiver received exact expected amount"
  else
    warning "Receiver did not receive expected amount in Electrs API. Expected: $AMOUNT_SATS satoshis, Actual: $RECEIVER_FINAL_BALANCE satoshis"

    if [ "$CORE_RECEIVER_BALANCE" -eq "$AMOUNT_SATS" ]; then
      info "However, Bitcoin Core shows the correct amount ($CORE_RECEIVER_BALANCE satoshis)"
      info "This is likely an Electrs indexing delay, not an actual transaction issue."
    fi
  fi

  # Step 15: Display transaction details
  header "Transaction Details (Electrs)"
  TX_DETAILS=$(curl -s $ELECTRS_URL/tx/$TXID_BITCOIND)

  if [ -n "$TX_DETAILS" ] && [ "$TX_DETAILS" != "null" ]; then
    echo $TX_DETAILS | jq -r '{
      txid: .txid,
      version: .version,
      size: .size,
      fee: .fee,
      inputs: .vin | length,
      outputs: .vout | length
    }'

    # Show outputs specifically to verify the amounts
    info "Transaction outputs:"
    echo $TX_DETAILS | jq -r '.vout[] | {address: .scriptpubkey_address, value: .value}'

    # Additionally, show detailed output structure
    info "Transaction confirmation status:"
    echo $TX_DETAILS | jq -r '.status'
  else
    warning "Could not retrieve detailed transaction information from Electrs"
  fi

  header "Transaction Complete"
fi
