# Test-Driven Development (TDD) Guide untuk RedEnvelope

## Overview

File ini menjelaskan implementasi TDD (Test-Driven Development) untuk fungsi-fungsi utama RedEnvelope service:
- `CreateEnvelope`
- `ClaimEnvelope`
- `RefundEnvelope`

## Struktur Test

### File Test
- **service_tdd_test.go** - Test suite lengkap dengan pendekatan TDD

### Helper Functions
```go
setupTestService(t, privateKey)     // Setup service untuk testing
waitForTransaction(t, service, tx)  // Tunggu konfirmasi transaksi
getBalance(t, service, address)     // Get balance address
```

## Test Cases

### 1. CreateEnvelope Tests

#### ✅ TestCreateEnvelope_DirectFixed_Success
**Tujuan**: Memverifikasi pembuatan envelope DIRECT_FIXED berhasil
- Setup: Create envelope untuk recipient tertentu
- Assert: 
  - Kind = DIRECT_FIXED (0)
  - Creator address benar
  - Recipient address benar
  - Transaction sukses

#### ✅ TestCreateEnvelope_GroupFixed_Success
**Tujuan**: Memverifikasi pembuatan envelope GROUP_FIXED berhasil
- Setup: Create envelope dengan 5 claims @ 0.05 ETH
- Assert:
  - Kind = GROUP_FIXED (1)
  - TotalClaims = 5
  - RemainingClaims = 5
  - AmountPerClaim = 0.05 ETH

#### ✅ TestCreateEnvelope_GroupRandom_Success
**Tujuan**: Memverifikasi pembuatan envelope GROUP_RANDOM berhasil
- Setup: Create envelope dengan 10 claims, total pot 0.5 ETH
- Assert:
  - Kind = GROUP_RANDOM (2)
  - TotalClaims = 10
  - Transaction sukses

#### ✅ TestCreateEnvelope_WithRoomIdHash
**Tujuan**: Memverifikasi envelope dengan room restriction
- Setup: Create envelope dengan roomIdHash
- Assert:
  - RoomIdHash tersimpan dengan benar
  - Envelope created successfully

---

### 2. ClaimEnvelope Tests

#### ✅ TestClaimEnvelope_FirstClaim_Success
**Tujuan**: Memverifikasi klaim pertama berhasil
- Setup: Create envelope dengan account #0
- Test: Claim dengan account #1
- Assert:
  - HasClaimed = false sebelum claim
  - Balance bertambah setelah claim (minimal 0.04 ETH)
  - HasClaimed = true setelah claim
  - RemainingClaims berkurang 1
  - Transaction sukses

#### ✅ TestClaimEnvelope_DoubleClaim_ShouldFail
**Tujuan**: Memverifikasi double claim ditolak
- Setup: Create dan claim envelope
- Test: Claim lagi dengan user yang sama
- Assert:
  - Error muncul
  - Double claim ditolak

#### ✅ TestClaimEnvelope_DirectFixedByRecipient_Success
**Tujuan**: Memverifikasi DIRECT_FIXED hanya bisa di-claim oleh recipient
- Setup: Create DIRECT_FIXED untuk account #1
- Test: Claim dengan account #1
- Assert:
  - Claim berhasil
  - Balance bertambah

#### ✅ TestClaimEnvelope_AllClaims_Success
**Tujuan**: Memverifikasi semua claims bisa diambil
- Setup: Create envelope dengan 2 claims
- Test: Claim dengan 2 accounts berbeda
- Assert:
  - RemainingClaims = 0
  - RemainingAmount = 0
  - Semua claims berhasil

---

### 3. RefundEnvelope Tests

#### ✅ TestRefundEnvelope_AfterExpiry_Success
**Tujuan**: Memverifikasi refund berhasil setelah expiry
- Setup: Create envelope dengan expiry 3 detik
- Test: Wait 5 detik, kemudian refund
- Assert:
  - Transaction sukses
  - Balance creator bertambah
  - Refund amount sesuai

#### ✅ TestRefundEnvelope_BeforeExpiry_ShouldFail
**Tujuan**: Memverifikasi refund ditolak sebelum expiry
- Setup: Create envelope dengan expiry 24 jam
- Test: Langsung coba refund
- Assert:
  - Error muncul
  - Refund ditolak

#### ✅ TestRefundEnvelope_PartialClaims_Success
**Tujuan**: Memverifikasi refund dengan partial claims
- Setup: Create envelope dengan 3 claims, 1 di-claim
- Test: Wait expiry, kemudian refund
- Assert:
  - Refund berhasil
  - Amount = 2 remaining claims worth

#### ✅ TestRefundEnvelope_NonCreator_ShouldFail
**Tujuan**: Memverifikasi hanya creator yang bisa refund
- Setup: Create envelope dengan account #0
- Test: Coba refund dengan account #1
- Assert:
  - Error muncul
  - Unauthorized refund ditolak

#### ✅ TestRefundEnvelope_AllClaimsExhausted_ShouldFail
**Tujuan**: Memverifikasi refund ketika semua claims sudah diambil
- Setup: Create envelope 1 claim, claim it
- Test: Wait expiry, coba refund
- Assert:
  - Error atau 0 amount
  - No refund available

---

## Cara Menjalankan Tests

### Jalankan Semua TDD Tests
```bash
cd redenvelope
go test -v -run "TestCreateEnvelope|TestClaimEnvelope|TestRefundEnvelope"
```

### Jalankan Tests Per Fungsi

**CreateEnvelope tests:**
```bash
go test -v -run TestCreateEnvelope
```

**ClaimEnvelope tests:**
```bash
go test -v -run TestClaimEnvelope
```

**RefundEnvelope tests:**
```bash
go test -v -run TestRefundEnvelope
```

### Jalankan Test Spesifik
```bash
go test -v -run TestClaimEnvelope_DoubleClaim_ShouldFail
```

---

## Prasyarat

1. **Local Blockchain Running** (Hardhat/Anvil)
```bash
# Terminal 1 - Run local node
npx hardhat node
# atau
anvil
```

2. **Contract Deployed**
- Contract address: `0x5FbDB2315678afecb367f032d93F642f64180aa3`
- Account #0: `0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266`
- Account #1: `0x70997970C51812dc3A010C7d01b50e0d17dc79C8`

3. **Funded Accounts**
- Hardhat local node sudah menyediakan accounts dengan balance

---

## Metodologi TDD

### Red-Green-Refactor Cycle

1. **RED** - Write Failing Test
   - Tulis test case yang menguji behavior yang diinginkan
   - Test akan fail karena implementasi belum ada/belum benar

2. **GREEN** - Make Test Pass
   - Implementasi minimal untuk membuat test pass
   - Focus on making it work, not perfect

3. **REFACTOR** - Improve Code
   - Clean up implementation
   - Improve readability
   - Optimize performance
   - Tests tetap pass

### Contoh Aplikasi TDD

#### Cycle 1: CreateEnvelope Basic
1. **RED**: Write `TestCreateEnvelope_DirectFixed_Success` - FAIL
2. **GREEN**: Implement basic `CreateEnvelope` - PASS
3. **REFACTOR**: Clean up parameter handling

#### Cycle 2: CreateEnvelope Edge Cases
1. **RED**: Write `TestCreateEnvelope_GroupFixed_Success` - FAIL
2. **GREEN**: Handle GROUP_FIXED calculation - PASS
3. **REFACTOR**: Extract grossPot calculation logic

#### Cycle 3: ClaimEnvelope
1. **RED**: Write `TestClaimEnvelope_FirstClaim_Success` - FAIL
2. **GREEN**: Implement `ClaimEnvelope` - PASS
3. **REFACTOR**: Add better error handling

---

## Best Practices

### 1. Test Naming
- Format: `Test<Function>_<Scenario>_<ExpectedResult>`
- Contoh: `TestClaimEnvelope_DoubleClaim_ShouldFail`

### 2. Test Structure (AAA Pattern)
```go
func TestSomething(t *testing.T) {
    // Arrange - Setup
    service := setupTestService(t, privateKey)
    
    // Act - Execute
    result, err := service.DoSomething()
    
    // Assert - Verify
    if err != nil {
        t.Fatalf("Expected no error, got: %v", err)
    }
}
```

### 3. Test Independence
- Setiap test harus independen
- Tidak bergantung pada urutan eksekusi
- Setup dan cleanup yang jelas

### 4. Descriptive Assertions
```go
// ❌ Bad
if result != expected {
    t.Error("wrong")
}

// ✅ Good
if result != expected {
    t.Errorf("Expected %v, got %v", expected, result)
}
```

### 5. Use Helper Functions
- Reduce code duplication
- Improve readability
- Easier maintenance

---

## Coverage Goals

Target coverage untuk TDD:
- ✅ **Happy Path**: Normal successful operations
- ✅ **Edge Cases**: Boundary conditions
- ✅ **Error Cases**: Expected failures
- ✅ **Security**: Access control, double-claim, etc.

### Coverage Breakdown

**CreateEnvelope**: 4 tests
- DIRECT_FIXED success ✅
- GROUP_FIXED success ✅
- GROUP_RANDOM success ✅
- With room restriction ✅

**ClaimEnvelope**: 4 tests
- First claim success ✅
- Double claim prevention ✅
- DIRECT_FIXED recipient ✅
- All claims exhausted ✅

**RefundEnvelope**: 5 tests
- After expiry success ✅
- Before expiry fail ✅
- Partial claims refund ✅
- Non-creator fail ✅
- All claims exhausted ✅

**Total**: 13 comprehensive test cases

---

## Troubleshooting

### Tests Failing - "Connection Refused"
**Problem**: Local blockchain tidak running
**Solution**: 
```bash
npx hardhat node
```

### Tests Failing - "Insufficient Funds"
**Problem**: Account tidak punya cukup ETH
**Solution**: Gunakan Hardhat local node yang sudah pre-funded

### Tests Failing - "Contract Not Found"
**Problem**: Contract address salah atau belum deployed
**Solution**: 
1. Deploy contract
2. Update `testContractAddress` di test file

### Tests Timing Out
**Problem**: Transaction tidak dikonfirmasi
**Solution**: 
- Increase wait time di `waitForTransaction`
- Check blockchain logs

---

## Next Steps

### Improvements
1. ✅ Add more edge case tests
2. ✅ Add performance benchmarks
3. ✅ Add integration tests
4. ✅ Add contract event verification
5. ✅ Add gas optimization tests

### Advanced TDD
1. **Property-Based Testing**: Using `rapid` or similar
2. **Fuzz Testing**: Random input generation
3. **Mutation Testing**: Verify test quality
4. **Coverage Analysis**: Ensure high coverage

---

## Referensi

- [Test-Driven Development by Example](https://www.amazon.com/Test-Driven-Development-Kent-Beck/dp/0321146530)
- [Go Testing Best Practices](https://go.dev/doc/tutorial/add-a-test)
- [Ethereum Testing Guide](https://ethereum.org/en/developers/docs/testing/)

---

## Kesimpulan

TDD memberikan:
- ✅ **Confidence**: Tests verify functionality
- ✅ **Documentation**: Tests explain behavior
- ✅ **Safety**: Catch regressions early
- ✅ **Design**: Better code structure
- ✅ **Speed**: Faster debugging

File `service_tdd_test.go` berisi 13 comprehensive test cases yang mengcover:
- Happy paths
- Edge cases
- Error handling
- Security constraints

Tests ini bisa dijalankan kapan saja untuk verify bahwa fungsi-fungsi core bekerja dengan benar! 🚀
