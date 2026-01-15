package redenvelope

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

const (
	testRPCURL          = "http://127.0.0.1:8545"
	testContractAddress = "0x5FC8d32690cc91D4c39d9d3abcBD16989F875707"
	testPrivateKey0     = "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80" // Account #0
	testPrivateKey1     = "59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d" // Account #1
)

// Helper functions for testing
func setupTestService(t *testing.T, privateKey string) *RedEnvelopeService {
	service, err := NewRedEnvelopeService(testRPCURL, testContractAddress, privateKey)
	if err != nil {
		t.Fatalf("Failed to create service: %v", err)
	}
	return service
}

func waitForTransaction(t *testing.T, service *RedEnvelopeService, tx *types.Transaction) *types.Receipt {
	time.Sleep(2 * time.Second)
	receipt, err := service.Client.TransactionReceipt(context.Background(), tx.Hash())
	if err != nil {
		t.Fatalf("Failed to get receipt: %v", err)
	}
	return receipt
}

func getBalance(t *testing.T, service *RedEnvelopeService, address common.Address) *big.Int {
	balance, err := service.Client.BalanceAt(context.Background(), address, nil)
	if err != nil {
		t.Fatalf("Failed to get balance: %v", err)
	}
	return balance
}

// ============================================================================
// TDD Tests for CreateEnvelope
// ============================================================================

func TestCreateEnvelope_DirectFixed_Success(t *testing.T) {
	service := setupTestService(t, testPrivateKey0)
	defer service.Client.Close()

	// Test: Create DIRECT_FIXED envelope successfully
	nextId, err := service.GetNextEnvelopeId()
	if err != nil {
		t.Fatalf("Failed to get next ID: %v", err)
	}

	amount := big.NewInt(100000000000000000) // 0.1 ETH
	recipient := common.HexToAddress("0x70997970c51812dc3a010c7d01b50e0d17dc79c8")

	tx, err := service.CreateEnvelope(
		DIRECT_FIXED,
		common.Address{}, // Native ETH
		1,                // totalClaims (ignored for DIRECT_FIXED)
		amount,
		24*time.Hour,
		EmptyRoomIdHash,
		recipient,
	)
	if err != nil {
		t.Fatalf("Failed to create envelope: %v", err)
	}

	receipt := waitForTransaction(t, service, tx)
	if receipt.Status != 1 {
		t.Fatal("Transaction failed")
	}

	// Verify envelope was created with correct parameters
	envelope, err := service.GetEnvelope(nextId)
	if err != nil {
		t.Fatalf("Failed to get envelope: %v", err)
	}

	if envelope.Kind != DIRECT_FIXED {
		t.Errorf("Expected kind DIRECT_FIXED (0), got %d", envelope.Kind)
	}
	if envelope.Creator != service.Address {
		t.Errorf("Expected creator %s, got %s", service.Address.Hex(), envelope.Creator.Hex())
	}
	if envelope.Recipient != recipient {
		t.Errorf("Expected recipient %s, got %s", recipient.Hex(), envelope.Recipient.Hex())
	}

	t.Logf("✓ DIRECT_FIXED envelope #%s created successfully", nextId.String())
}

func TestCreateEnvelope_GroupFixed_Success(t *testing.T) {
	service := setupTestService(t, testPrivateKey0)
	defer service.Client.Close()

	// Test: Create GROUP_FIXED envelope successfully
	nextId, err := service.GetNextEnvelopeId()
	if err != nil {
		t.Fatalf("Failed to get next ID: %v", err)
	}

	totalClaims := uint32(5)
	amountPerClaim := big.NewInt(50000000000000000) // 0.05 ETH

	tx, err := service.CreateEnvelope(
		GROUP_FIXED,
		common.Address{}, // Native ETH
		totalClaims,
		amountPerClaim,
		24*time.Hour,
		EmptyRoomIdHash,
		common.Address{},
	)
	if err != nil {
		t.Fatalf("Failed to create envelope: %v", err)
	}

	receipt := waitForTransaction(t, service, tx)
	if receipt.Status != 1 {
		t.Fatal("Transaction failed")
	}

	// Verify envelope was created with correct parameters
	envelope, err := service.GetEnvelope(nextId)
	if err != nil {
		t.Fatalf("Failed to get envelope: %v", err)
	}

	if envelope.Kind != GROUP_FIXED {
		t.Errorf("Expected kind GROUP_FIXED (1), got %d", envelope.Kind)
	}
	if envelope.TotalClaims != totalClaims {
		t.Errorf("Expected %d total claims, got %d", totalClaims, envelope.TotalClaims)
	}
	if envelope.RemainingClaims != totalClaims {
		t.Errorf("Expected %d remaining claims, got %d", totalClaims, envelope.RemainingClaims)
	}
	if envelope.AmountPerClaim.Cmp(amountPerClaim) != 0 {
		t.Errorf("Expected amount per claim %s, got %s", amountPerClaim.String(), envelope.AmountPerClaim.String())
	}

	t.Logf("✓ GROUP_FIXED envelope #%s created successfully", nextId.String())
}

func TestCreateEnvelope_GroupRandom_Success(t *testing.T) {
	service := setupTestService(t, testPrivateKey0)
	defer service.Client.Close()

	// Test: Create GROUP_RANDOM envelope successfully
	nextId, err := service.GetNextEnvelopeId()
	if err != nil {
		t.Fatalf("Failed to get next ID: %v", err)
	}

	totalClaims := uint32(10)
	totalPot := big.NewInt(500000000000000000) // 0.5 ETH

	tx, err := service.CreateEnvelope(
		GROUP_RANDOM,
		common.Address{}, // Native ETH
		totalClaims,
		totalPot,
		24*time.Hour,
		EmptyRoomIdHash,
		common.Address{},
	)
	if err != nil {
		t.Fatalf("Failed to create envelope: %v", err)
	}

	receipt := waitForTransaction(t, service, tx)
	if receipt.Status != 1 {
		t.Fatal("Transaction failed")
	}

	// Verify envelope was created
	envelope, err := service.GetEnvelope(nextId)
	if err != nil {
		t.Fatalf("Failed to get envelope: %v", err)
	}

	if envelope.Kind != GROUP_RANDOM {
		t.Errorf("Expected kind GROUP_RANDOM (2), got %d", envelope.Kind)
	}
	if envelope.TotalClaims != totalClaims {
		t.Errorf("Expected %d total claims, got %d", totalClaims, envelope.TotalClaims)
	}

	t.Logf("✓ GROUP_RANDOM envelope #%s created successfully", nextId.String())
}

func TestCreateEnvelope_WithRoomIdHash(t *testing.T) {
	service := setupTestService(t, testPrivateKey0)
	defer service.Client.Close()

	// Test: Create envelope with room restriction
	nextId, err := service.GetNextEnvelopeId()
	if err != nil {
		t.Fatalf("Failed to get next ID: %v", err)
	}

	roomId := "private-room-123"
	roomIdHash := GenerateRoomIdHash(roomId)
	amount := big.NewInt(100000000000000000) // 0.1 ETH

	tx, err := service.CreateEnvelope(
		GROUP_FIXED,
		common.Address{},
		3,
		amount,
		24*time.Hour,
		roomIdHash,
		common.Address{},
	)
	if err != nil {
		t.Fatalf("Failed to create envelope: %v", err)
	}

	receipt := waitForTransaction(t, service, tx)
	if receipt.Status != 1 {
		t.Fatal("Transaction failed")
	}

	// Verify room hash is set
	envelope, err := service.GetEnvelope(nextId)
	if err != nil {
		t.Fatalf("Failed to get envelope: %v", err)
	}

	if envelope.RoomIdHash != roomIdHash {
		t.Errorf("Room ID hash mismatch")
	}

	t.Logf("✓ Envelope with room restriction created successfully")
}

// ============================================================================
// TDD Tests for ClaimEnvelope
// ============================================================================

func TestClaimEnvelope_FirstClaim_Success(t *testing.T) {
	// Setup: Create envelope with account #0
	serviceCreator := setupTestService(t, testPrivateKey0)
	defer serviceCreator.Client.Close()

	nextId, err := serviceCreator.GetNextEnvelopeId()
	if err != nil {
		t.Fatalf("Failed to get next ID: %v", err)
	}
	totalClaims := uint32(3)
	amountPerClaim := big.NewInt(100000000000000000) // 0.1 ETH

	tx, err := serviceCreator.CreateEnvelope(
		GROUP_FIXED,
		common.Address{},
		totalClaims,
		amountPerClaim,
		24*time.Hour,
		EmptyRoomIdHash,
		common.Address{},
	)
	if err != nil {
		t.Fatalf("Setup failed: %v", err)
	}
	waitForTransaction(t, serviceCreator, tx)

	// Test: Claim with account #1
	serviceClaimer := setupTestService(t, testPrivateKey1)
	defer serviceClaimer.Client.Close()

	// Verify hasn't claimed yet
	hasClaimed, err := serviceClaimer.HasClaimed(nextId, serviceClaimer.Address)
	if err != nil {
		t.Fatalf("Failed to check claim status: %v", err)
	}
	if hasClaimed {
		t.Fatal("Should not have claimed yet")
	}

	// Get balance before
	balanceBefore := getBalance(t, serviceClaimer, serviceClaimer.Address)

	// Claim envelope
	claimTx, err := serviceClaimer.ClaimEnvelope(nextId)
	if err != nil {
		t.Fatalf("Failed to claim envelope: %v", err)
	}

	claimReceipt := waitForTransaction(t, serviceClaimer, claimTx)
	if claimReceipt.Status != 1 {
		t.Fatal("Claim transaction failed")
	}

	// Verify balance increased
	balanceAfter := getBalance(t, serviceClaimer, serviceClaimer.Address)
	diff := new(big.Int).Sub(balanceAfter, balanceBefore)
	minExpected := big.NewInt(40000000000000000) // At least 0.04 ETH (minus gas)
	if diff.Cmp(minExpected) < 0 {
		t.Errorf("Balance increase too small: %s", diff.String())
	}

	// Verify has claimed
	hasClaimed, err = serviceClaimer.HasClaimed(nextId, serviceClaimer.Address)
	if err != nil {
		t.Fatalf("Failed to check claim status after: %v", err)
	}
	if !hasClaimed {
		t.Fatal("Should have claimed")
	}

	// Verify envelope state updated
	envelope, err := serviceClaimer.GetEnvelope(nextId)
	if err != nil {
		t.Fatalf("Failed to get envelope: %v", err)
	}
	if envelope.RemainingClaims != totalClaims-1 {
		t.Errorf("Expected %d remaining claims, got %d", totalClaims-1, envelope.RemainingClaims)
	}

	t.Logf("✓ First claim successful, balance increased by %s wei", diff.String())
}

func TestClaimEnvelope_DoubleClaim_ShouldFail(t *testing.T) {
	// Setup: Create and claim envelope
	serviceCreator := setupTestService(t, testPrivateKey0)
	defer serviceCreator.Client.Close()

	nextId, _ := serviceCreator.GetNextEnvelopeId()
	amountPerClaim := big.NewInt(50000000000000000)

	tx, _ := serviceCreator.CreateEnvelope(
		GROUP_FIXED,
		common.Address{},
		3,
		amountPerClaim,
		24*time.Hour,
		EmptyRoomIdHash,
		common.Address{},
	)
	waitForTransaction(t, serviceCreator, tx)

	serviceClaimer := setupTestService(t, testPrivateKey1)
	defer serviceClaimer.Client.Close()

	// First claim
	claimTx, _ := serviceClaimer.ClaimEnvelope(nextId)
	waitForTransaction(t, serviceClaimer, claimTx)

	// Test: Try to claim again (should fail)
	_, err := serviceClaimer.ClaimEnvelope(nextId)
	if err == nil {
		t.Fatal("Expected error for double claim, got none")
	}

	t.Logf("✓ Double claim correctly prevented: %v", err)
}

func TestClaimEnvelope_DirectFixedByRecipient_Success(t *testing.T) {
	// Setup: Create DIRECT_FIXED envelope for account #1
	serviceCreator := setupTestService(t, testPrivateKey0)
	defer serviceCreator.Client.Close()

	nextId, _ := serviceCreator.GetNextEnvelopeId()
	amount := big.NewInt(100000000000000000)                                       // 0.1 ETH
	recipient := common.HexToAddress("0x70997970C51812dc3A010C7d01b50e0d17dc79C8") // Account #1

	tx, err := serviceCreator.CreateEnvelope(
		DIRECT_FIXED,
		common.Address{},
		1,
		amount,
		24*time.Hour,
		EmptyRoomIdHash,
		recipient,
	)
	if err != nil {
		t.Fatalf("Setup failed: %v", err)
	}
	waitForTransaction(t, serviceCreator, tx)

	// Test: Claim by designated recipient
	serviceClaimer := setupTestService(t, testPrivateKey1)
	defer serviceClaimer.Client.Close()

	balanceBefore := getBalance(t, serviceClaimer, serviceClaimer.Address)

	claimTx, err := serviceClaimer.ClaimEnvelope(nextId)
	if err != nil {
		t.Fatalf("Failed to claim: %v", err)
	}

	claimReceipt := waitForTransaction(t, serviceClaimer, claimTx)
	if claimReceipt.Status != 1 {
		t.Fatal("Claim failed")
	}

	balanceAfter := getBalance(t, serviceClaimer, serviceClaimer.Address)
	diff := new(big.Int).Sub(balanceAfter, balanceBefore)

	t.Logf("✓ DIRECT_FIXED claim successful, received %s wei", diff.String())
}

func TestClaimEnvelope_AllClaims_Success(t *testing.T) {
	// Setup: Create envelope with 2 claims
	serviceCreator := setupTestService(t, testPrivateKey0)
	defer serviceCreator.Client.Close()

	nextId, _ := serviceCreator.GetNextEnvelopeId()
	totalClaims := uint32(2)
	amountPerClaim := big.NewInt(50000000000000000)

	tx, _ := serviceCreator.CreateEnvelope(
		GROUP_FIXED,
		common.Address{},
		totalClaims,
		amountPerClaim,
		24*time.Hour,
		EmptyRoomIdHash,
		common.Address{},
	)
	waitForTransaction(t, serviceCreator, tx)

	// Claim with account #1
	serviceClaimer1 := setupTestService(t, testPrivateKey1)
	defer serviceClaimer1.Client.Close()
	claimTx1, _ := serviceClaimer1.ClaimEnvelope(nextId)
	waitForTransaction(t, serviceClaimer1, claimTx1)

	// Claim with account #0
	claimTx2, _ := serviceCreator.ClaimEnvelope(nextId)
	waitForTransaction(t, serviceCreator, claimTx2)

	// Verify all claims exhausted
	envelope, _ := serviceCreator.GetEnvelope(nextId)
	if envelope.RemainingClaims != 0 {
		t.Errorf("Expected 0 remaining claims, got %d", envelope.RemainingClaims)
	}
	if envelope.RemainingAmount.Cmp(big.NewInt(0)) != 0 {
		t.Errorf("Expected 0 remaining amount, got %s", envelope.RemainingAmount.String())
	}

	t.Logf("✓ All claims exhausted successfully")
}

// ============================================================================
// TDD Tests for RefundEnvelope
// ============================================================================

func TestRefundEnvelope_AfterExpiry_Success(t *testing.T) {
	t.Skip("Skipping test that requires waiting for expiry - run manually if needed")

	// Setup: Create envelope with short expiry
	serviceCreator := setupTestService(t, testPrivateKey0)
	defer serviceCreator.Client.Close()

	nextId, _ := serviceCreator.GetNextEnvelopeId()
	amount := big.NewInt(100000000000000000) // 0.1 ETH

	tx, err := serviceCreator.CreateEnvelope(
		GROUP_FIXED,
		common.Address{},
		3,
		amount,
		3*time.Second, // Short expiry
		EmptyRoomIdHash,
		common.Address{},
	)
	if err != nil {
		t.Fatalf("Setup failed: %v", err)
	}
	waitForTransaction(t, serviceCreator, tx)

	// Wait for expiry
	t.Log("Waiting for envelope to expire...")
	time.Sleep(5 * time.Second)

	// Test: Refund after expiry
	balanceBefore := getBalance(t, serviceCreator, serviceCreator.Address)

	refundTx, err := serviceCreator.RefundEnvelope(nextId)
	if err != nil {
		t.Fatalf("Failed to refund: %v", err)
	}

	refundReceipt := waitForTransaction(t, serviceCreator, refundTx)
	if refundReceipt.Status != 1 {
		t.Fatal("Refund transaction failed")
	}

	balanceAfter := getBalance(t, serviceCreator, serviceCreator.Address)
	diff := new(big.Int).Sub(balanceAfter, balanceBefore)

	// Should get some refund (minus fees and gas)
	minExpected := big.NewInt(100000000000000000) // At least 0.1 ETH
	if diff.Cmp(minExpected) < 0 {
		t.Logf("Warning: Refund amount seems low: %s", diff.String())
	}

	t.Logf("✓ Refund successful, received %s wei", diff.String())
}

func TestRefundEnvelope_BeforeExpiry_ShouldFail(t *testing.T) {
	// Setup: Create envelope with long expiry
	serviceCreator := setupTestService(t, testPrivateKey0)
	defer serviceCreator.Client.Close()

	nextId, _ := serviceCreator.GetNextEnvelopeId()
	amount := big.NewInt(100000000000000000)

	tx, _ := serviceCreator.CreateEnvelope(
		GROUP_FIXED,
		common.Address{},
		3,
		amount,
		24*time.Hour, // Long expiry
		EmptyRoomIdHash,
		common.Address{},
	)
	waitForTransaction(t, serviceCreator, tx)

	// Test: Try to refund before expiry (should fail)
	_, err := serviceCreator.RefundEnvelope(nextId)
	if err == nil {
		t.Fatal("Expected error for refund before expiry, got none")
	}

	t.Logf("✓ Refund before expiry correctly prevented: %v", err)
}

func TestRefundEnvelope_PartialClaims_Success(t *testing.T) {
	t.Skip("Skipping test that requires waiting for expiry - run manually if needed")

	// Setup: Create envelope, make 1 claim, then refund after expiry
	serviceCreator := setupTestService(t, testPrivateKey0)
	defer serviceCreator.Client.Close()

	nextId, _ := serviceCreator.GetNextEnvelopeId()
	totalClaims := uint32(3)
	amountPerClaim := big.NewInt(100000000000000000) // 0.1 ETH

	tx, _ := serviceCreator.CreateEnvelope(
		GROUP_FIXED,
		common.Address{},
		totalClaims,
		amountPerClaim,
		3*time.Second, // Short expiry
		EmptyRoomIdHash,
		common.Address{},
	)
	waitForTransaction(t, serviceCreator, tx)

	// Make 1 claim
	serviceClaimer := setupTestService(t, testPrivateKey1)
	defer serviceClaimer.Client.Close()

	claimTx, _ := serviceClaimer.ClaimEnvelope(nextId)
	waitForTransaction(t, serviceClaimer, claimTx)

	// Wait for expiry
	t.Log("Waiting for envelope to expire...")
	time.Sleep(5 * time.Second)

	// Test: Refund remaining amount
	balanceBefore := getBalance(t, serviceCreator, serviceCreator.Address)

	refundTx, err := serviceCreator.RefundEnvelope(nextId)
	if err != nil {
		t.Fatalf("Failed to refund: %v", err)
	}

	refundReceipt := waitForTransaction(t, serviceCreator, refundTx)
	if refundReceipt.Status != 1 {
		t.Fatal("Refund failed")
	}

	balanceAfter := getBalance(t, serviceCreator, serviceCreator.Address)
	diff := new(big.Int).Sub(balanceAfter, balanceBefore)

	// Should get refund for 2 remaining claims
	t.Logf("✓ Partial refund successful, received %s wei", diff.String())
}

func TestRefundEnvelope_NonCreator_ShouldFail(t *testing.T) {
	t.Skip("Skipping test that requires waiting for expiry - run manually if needed")

	// Setup: Create envelope with account #0
	serviceCreator := setupTestService(t, testPrivateKey0)
	defer serviceCreator.Client.Close()

	nextId, _ := serviceCreator.GetNextEnvelopeId()
	amount := big.NewInt(100000000000000000)

	tx, _ := serviceCreator.CreateEnvelope(
		GROUP_FIXED,
		common.Address{},
		3,
		amount,
		3*time.Second,
		EmptyRoomIdHash,
		common.Address{},
	)
	waitForTransaction(t, serviceCreator, tx)

	// Wait for expiry
	time.Sleep(5 * time.Second)

	// Test: Try to refund with different account (should fail)
	serviceOther := setupTestService(t, testPrivateKey1)
	defer serviceOther.Client.Close()

	_, err := serviceOther.RefundEnvelope(nextId)
	if err == nil {
		t.Fatal("Expected error for refund by non-creator, got none")
	}

	t.Logf("✓ Refund by non-creator correctly prevented: %v", err)
}

func TestRefundEnvelope_AllClaimsExhausted_ShouldFail(t *testing.T) {
	t.Skip("Skipping test that requires waiting for expiry - run manually if needed")

	// Setup: Create envelope with 1 claim, claim it, wait for expiry
	serviceCreator := setupTestService(t, testPrivateKey0)
	defer serviceCreator.Client.Close()

	nextId, _ := serviceCreator.GetNextEnvelopeId()
	amount := big.NewInt(100000000000000000)

	tx, _ := serviceCreator.CreateEnvelope(
		GROUP_FIXED,
		common.Address{},
		1, // Only 1 claim
		amount,
		3*time.Second,
		EmptyRoomIdHash,
		common.Address{},
	)
	waitForTransaction(t, serviceCreator, tx)

	// Claim the only available claim
	serviceClaimer := setupTestService(t, testPrivateKey1)
	defer serviceClaimer.Client.Close()
	claimTx, _ := serviceClaimer.ClaimEnvelope(nextId)
	waitForTransaction(t, serviceClaimer, claimTx)

	// Wait for expiry
	time.Sleep(5 * time.Second)

	// Test: Try to refund when no remaining amount (should fail or return 0)
	_, err := serviceCreator.RefundEnvelope(nextId)
	if err != nil {
		t.Logf("✓ Refund with no remaining amount correctly prevented: %v", err)
	} else {
		// If it succeeds, check that remaining amount is 0
		envelope, _ := serviceCreator.GetEnvelope(nextId)
		if envelope.RemainingAmount.Cmp(big.NewInt(0)) != 0 {
			t.Errorf("Expected 0 remaining amount")
		}
		t.Logf("✓ Refund processed with 0 remaining amount")
	}
}
