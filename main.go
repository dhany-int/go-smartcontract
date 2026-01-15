package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"rpcsol/redenvelope"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	// Konfigurasi
	rpcURL := "http://127.0.0.1:8545"
	// Private key dari Account #0 Hardhat
	privateKeyHex := "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"

	// !! GANTI DENGAN CONTRACT ADDRESS ANDA !!
	contractAddress := "0x5FbDB2315678afecb367f032d93F642f64180aa3"

	// Connect ke Ethereum node
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		log.Fatalf("Failed to connect to ethereum node: %v", err)
	}
	defer client.Close()

	fmt.Println("=== Ethereum RPC Service Demo ===")
	fmt.Printf("Connected to: %s\n", rpcURL)
	fmt.Println()

	// Load private key
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		log.Fatalf("Failed to load private key: %v", err)
	}

	// Get address dari private key
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatalf("Failed to cast public key to ECDSA")
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	fmt.Printf("Your address: %s\n", fromAddress.Hex())
	fmt.Println()

	// 1. Get Balance
	fmt.Println("=== Get Balance ===")
	balance, err := client.BalanceAt(context.Background(), fromAddress, nil)
	if err != nil {
		log.Fatalf("Failed to get balance: %v", err)
	}
	fmt.Printf("Balance: %s ETH\n", weiToEther(balance))
	fmt.Println()

	// 2. Get Chain ID
	fmt.Println("=== Chain Information ===")
	chainID, err := client.ChainID(context.Background())
	if err != nil {
		log.Fatalf("Failed to get chain ID: %v", err)
	}
	fmt.Printf("Chain ID: %s\n", chainID.String())

	// 3. Get Block Number
	blockNumber, err := client.BlockNumber(context.Background())
	if err != nil {
		log.Fatalf("Failed to get block number: %v", err)
	}
	fmt.Printf("Current block number: %d\n", blockNumber)
	fmt.Println()

	// 4. Get Block by Number
	if blockNumber > 0 {
		fmt.Println("=== Latest Block Info ===")
		block, err := client.BlockByNumber(context.Background(), big.NewInt(int64(blockNumber)))
		if err != nil {
			log.Fatalf("Failed to get block: %v", err)
		}
		fmt.Printf("Block Hash: %s\n", block.Hash().Hex())
		fmt.Printf("Block Time: %d\n", block.Time())
		fmt.Printf("Transactions: %d\n", len(block.Transactions()))
		fmt.Println()
	}

	// 5. Send Transaction (transfer ETH)
	fmt.Println("=== Sending ETH Transaction ===")
	// Kirim ke Account #1
	toAddress := common.HexToAddress("0x70997970C51812dc3A010C7d01b50e0d17dc79C8")
	amount := big.NewInt(1000000000000000000) // 1 ETH dalam wei

	// Get nonce
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatalf("Failed to get nonce: %v", err)
	}

	// Get gas price
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatalf("Failed to get gas price: %v", err)
	}

	// Create transaction
	tx := types.NewTransaction(
		nonce,
		toAddress,
		amount,
		21000, // gas limit untuk transfer ETH
		gasPrice,
		nil,
	)

	// Sign transaction
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Fatalf("Failed to sign transaction: %v", err)
	}

	// Send transaction
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatalf("Failed to send transaction: %v", err)
	}

	fmt.Printf("From: %s\n", fromAddress.Hex())
	fmt.Printf("To: %s\n", toAddress.Hex())
	fmt.Printf("Amount: %s ETH\n", weiToEther(amount))
	fmt.Printf("Transaction Hash: %s\n", signedTx.Hash().Hex())
	fmt.Println()

	// 6. Get Transaction Receipt
	fmt.Println("=== Transaction Receipt ===")
	fmt.Println("Waiting for transaction to be mined...")
	receipt, err := waitForTransactionReceipt(client, signedTx.Hash())
	if err != nil {
		log.Printf("Warning: Failed to get receipt: %v", err)
	} else {
		fmt.Printf("Status: %d (1 = success, 0 = failed)\n", receipt.Status)
		fmt.Printf("Block Number: %d\n", receipt.BlockNumber.Uint64())
		fmt.Printf("Gas Used: %d\n", receipt.GasUsed)
		fmt.Println()
	}

	// 7. Get Updated Balances
	fmt.Println("=== Updated Balances ===")
	newBalance, err := client.BalanceAt(context.Background(), fromAddress, nil)
	if err != nil {
		log.Printf("Failed to get new balance: %v", err)
	} else {
		fmt.Printf("Your balance: %s ETH\n", weiToEther(newBalance))
	}

	toBalance, err := client.BalanceAt(context.Background(), toAddress, nil)
	if err != nil {
		log.Printf("Failed to get recipient balance: %v", err)
	} else {
		fmt.Printf("Recipient balance: %s ETH\n", weiToEther(toBalance))
	}
	fmt.Println()

	fmt.Println("=== Demo RPC Completed Successfully! ===")
	fmt.Println()

	// ========================================
	// DEMO REDENVELOPE CONTRACT
	// ========================================
	fmt.Println("\n========================================")
	fmt.Println("=== RedEnvelope Contract Demo ===")
	fmt.Println("========================================")

	demoRedEnvelope(rpcURL, contractAddress, privateKeyHex)
}

func demoRedEnvelope(rpcURL, contractAddress, privateKeyHex string) {
	// Initialize RedEnvelope service
	reService, err := redenvelope.NewRedEnvelopeService(rpcURL, contractAddress, privateKeyHex)
	if err != nil {
		log.Printf("Failed to initialize RedEnvelope service: %v", err)
		log.Println("Pastikan contract address sudah benar!")
		return
	}
	defer reService.Client.Close()

	fmt.Printf("Connected to RedEnvelope contract: %s\n", contractAddress)
	fmt.Printf("Your address: %s\n", reService.Address.Hex())
	fmt.Println()

	// Get balance
	balance, err := reService.Client.BalanceAt(context.Background(), reService.Address, nil)
	if err != nil {
		log.Printf("Failed to get balance: %v", err)
		return
	}
	fmt.Printf("Balance: %s ETH\n", weiToEther(balance))
	fmt.Println()

	// 1. Get Next Envelope ID
	fmt.Println("=== Get Next Envelope ID ===")
	nextId, err := reService.GetNextEnvelopeId()
	if err != nil {
		log.Printf("Failed to get next envelope ID: %v", err)
	} else {
		fmt.Printf("Next Envelope ID: %s\n", nextId.String())
	}
	fmt.Println()

	// 2. Create GROUP_FIXED Envelope
	fmt.Println("=== Creating GROUP_FIXED Envelope ===")
	amountPerClaim := big.NewInt(100000000000000000) // 0.1 ETH per claim
	totalClaims := uint32(5)
	grossPot := new(big.Int).Mul(amountPerClaim, big.NewInt(int64(totalClaims)))            // 0.5 ETH total
	fee := new(big.Int).Div(new(big.Int).Mul(grossPot, big.NewInt(250)), big.NewInt(10000)) // 2.5% fee
	netPot := new(big.Int).Sub(grossPot, fee)

	fmt.Printf("Type: GROUP_FIXED\n")
	fmt.Printf("Total Claims: %d\n", totalClaims)
	fmt.Printf("Amount per Claim: %s ETH\n", weiToEther(amountPerClaim))
	fmt.Printf("Gross Pot (sent): %s ETH\n", weiToEther(grossPot))
	fmt.Printf("Fee (2.5%%): %s ETH\n", weiToEther(fee))
	fmt.Printf("Net Pot (escrowed): %s ETH\n", weiToEther(netPot))
	fmt.Printf("Expiry: 1 hour from now\n")
	fmt.Printf("RoomIdHash: Empty (no restriction)\n")

	// Untuk GROUP_FIXED: amount parameter adalah amountPerClaim (contract akan × totalClaims)
	// Gunakan EmptyRoomIdHash supaya siapa saja bisa claim
	tx, err := reService.CreateEnvelope(
		redenvelope.GROUP_FIXED,
		common.Address{}, // Native token (ETH/BNB)
		totalClaims,
		amountPerClaim, // Amount PER CLAIM (bukan total!)
		1*time.Hour,
		redenvelope.EmptyRoomIdHash, // Tidak ada room restriction
		common.Address{},            // No specific recipient
	)

	if err != nil {
		log.Printf("Failed to create envelope: %v", err)
	} else {
		fmt.Printf("✓ Transaction sent: %s\n", tx.Hash().Hex())
		fmt.Println("Waiting for confirmation...")
		time.Sleep(2 * time.Second)

		receipt, err := reService.Client.TransactionReceipt(context.Background(), tx.Hash())
		if err == nil && receipt.Status == 1 {
			fmt.Printf("✓ Envelope created successfully!\n")
			fmt.Printf("  Envelope ID: %s\n", nextId.String())
		}
	}
	fmt.Println()

	// 3. Get Envelope Info
	if nextId != nil {
		fmt.Println("=== Get Envelope Information ===")
		envelope, err := reService.GetEnvelope(nextId)
		if err != nil {
			log.Printf("Failed to get envelope: %v", err)
		} else {
			fmt.Printf("Envelope ID: %s\n", nextId.String())
			fmt.Printf("Creator: %s\n", envelope.Creator.Hex())
			fmt.Printf("Kind: %s\n", getEnvelopeKindName(envelope.Kind))
			fmt.Printf("Total Claims: %d\n", envelope.TotalClaims)
			fmt.Printf("Remaining Claims: %d\n", envelope.RemainingClaims)
			fmt.Printf("Amount Per Claim: %s ETH\n", weiToEther(envelope.AmountPerClaim))
			fmt.Printf("Remaining Amount: %s ETH\n", weiToEther(envelope.RemainingAmount))
			fmt.Printf("Expiry: %s\n", time.Unix(int64(envelope.Expiry), 0).Format("2006-01-02 15:04:05"))
		}
		fmt.Println()

		// 4. Check if already claimed
		fmt.Println("=== Check Claim Status ===")
		hasClaimed, err := reService.HasClaimed(nextId, reService.Address)
		if err != nil {
			log.Printf("Failed to check claim status: %v", err)
		} else {
			fmt.Printf("Has claimed: %v\n", hasClaimed)
		}
		fmt.Println()

		// 5. Claim Envelope
		if !hasClaimed {
			fmt.Println("=== Claiming Envelope ===")
			claimTx, err := reService.ClaimEnvelope(nextId)
			if err != nil {
				log.Printf("Failed to claim envelope: %v", err)
			} else {
				fmt.Printf("✓ Claim transaction sent: %s\n", claimTx.Hash().Hex())
				fmt.Println("Waiting for confirmation...")
				time.Sleep(2 * time.Second)

				receipt, err := reService.Client.TransactionReceipt(context.Background(), claimTx.Hash())
				if err == nil && receipt.Status == 1 {
					fmt.Printf("✓ Claim successful!\n")

					// Get updated envelope info
					updatedEnvelope, _ := reService.GetEnvelope(nextId)
					if updatedEnvelope != nil {
						fmt.Printf("  Remaining Claims: %d\n", updatedEnvelope.RemainingClaims)
						fmt.Printf("  Remaining Amount: %s ETH\n", weiToEther(updatedEnvelope.RemainingAmount))
					}
				}
			}
			fmt.Println()
		}
	}

	fmt.Println("=== RedEnvelope Demo Completed! ===")
	fmt.Println()
	fmt.Println("Fitur yang sudah didemonstrasikan:")
	fmt.Println("✓ Connect ke RedEnvelope contract")
	fmt.Println("✓ Get next envelope ID")
	fmt.Println("✓ Create GROUP_FIXED envelope")
	fmt.Println("✓ Get envelope information")
	fmt.Println("✓ Check claim status")
	fmt.Println("✓ Claim envelope")
}

func getEnvelopeKindName(kind uint8) string {
	switch kind {
	case redenvelope.DIRECT_FIXED:
		return "DIRECT_FIXED"
	case redenvelope.GROUP_FIXED:
		return "GROUP_FIXED"
	case redenvelope.GROUP_RANDOM:
		return "GROUP_RANDOM"
	default:
		return "UNKNOWN"
	}
}

// weiToEther converts wei to ether
func weiToEther(wei *big.Int) string {
	ether := new(big.Float).SetInt(wei)
	ether = ether.Quo(ether, big.NewFloat(1e18))
	return ether.Text('f', 4)
}

// waitForTransactionReceipt waits for transaction to be mined
func waitForTransactionReceipt(client *ethclient.Client, txHash common.Hash) (*types.Receipt, error) {
	for i := 0; i < 30; i++ { // max 30 attempts
		receipt, err := client.TransactionReceipt(context.Background(), txHash)
		if err == nil {
			return receipt, nil
		}
		// Wait 1 second before retry
		// Note: In production, you should use time.Sleep()
	}
	return nil, fmt.Errorf("transaction not mined after 30 attempts")
}
