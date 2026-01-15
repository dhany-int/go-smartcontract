# Go RPC Service untuk RedEnvelope Smart Contract

Service Go untuk berinteraksi dengan **RedEnvelope Smart Contract** (Angpao Digital) di Hardhat local node atau testnet.

## Tentang RedEnvelope Contract

RedEnvelope adalah smart contract untuk sistem Angpao Digital di blockchain dengan 3 mode berbeda:

- **DIRECT_FIXED** - Angpao langsung untuk 1 penerima spesifik
- **GROUP_FIXED** - Angpao grup dengan jumlah tetap per klaim  
- **GROUP_RANDOM** - Angpao grup dengan distribusi random (luck-based)

### Fitur Contract
- ✅ Support native token (BNB/ETH) dan ERC20
- ✅ Fee system yang dapat dikonfigurasi
- ✅ Expiry time untuk setiap envelope
- ✅ Refund otomatis setelah expiry
- ✅ Room ID hash untuk verifikasi grup

## Struktur Project

```
go-rpc-sol/
├── main.go                    # Demo RPC operations + RedEnvelope
├── go.mod                     # Go module dependencies
├── redenvelope/
│   ├── abi.go                # RedEnvelope contract ABI
│   └── service.go            # Service untuk interact dengan contract
├── service/
│   └── ethereum.go           # General Ethereum service
└── examples/
    ├── README.md             # Dokumentasi examples
    ├── rpc_methods.go        # Demo semua RPC methods
    └── redenvelope_all.go    # Demo semua jenis envelope
```

## FKoneksi ke Hardhat local node (http://127.0.0.1:8545)
- ✅ Load private key dan generate address
- ✅ Get balance dari address
- ✅ Get chain ID
- ✅ Get block number terkini
- ✅ Get block information
- ✅ Create, sign, dan send transaction (transfer ETH)
- ✅ Wait for transaction receipt
- ✅ Verify transaction status

### RedEnvelope Contract Operations
- ✅ Create DIRECT_FIXED envelope
- ✅ Create GROUP_FIXED envelope
- ✅ Create GROUP_RANDOM envelope
- ✅ Claim envelope
- ✅ Get envelope information
- ✅ CRedEnvelope Contract** sudah di-deploy. Ganti contract address di kode:
   ```go
   contractAddress := "0xYourContractAddress"
   ```

3. **heck claim status
- ✅ Refund expired envelope
- ✅ Generate room ID hash
- ✅ Get block information
- ✅**Update contract address** di [main.go](main.go):
   ```go
   contractAddress := "0x5FbDB2315678afecb367f032d93F642f64180aa3" // Ganti dengan address Anda
   ```

3. Jalankan aplikasi:
   ```bash
   go run main.go
   ```

## Quick Start - RedEnvelope

### 1. Initialize Service

```go
import "rpcsol/redenvelope"

reService, err := redenvelope.NewRedEnvelopeService(
    "http://127.0.0.1:8545",
    "0xYourContractAddress",
    "your_private_key_hex",
)
```

### 2. Create DIRECT_FIXED Envelope

Angpao untuk 1 orang spesifik:

```go
roomIdHash := redenvelope.GenerateRoomIdHash("room-123")
recipient := common.HexToAddress("0x70997970C51812dc3A010C7d01b50e0d17dc79C8")
amount := big.NewInt(500000000000000000) // 0.5 ETH

tx, err := reService.CreateEnvelope(
    redenvelope.DIRECT_FIXED,
    common.Address{},     // Native token
    1,                    // Total claims = 1
    amount,
    24*time.Hour,         // Expiry
    roomIdHash,
    recipient,            // Specific recipient
)
```

### 3. Create GROUP_FIXED Envelope

Angpao grup dengan jumlah tetap per orang:

```go
amountPerClaim := big.NewInt(200000000000000000) // 0.2 ETH per orang
totalClaims := uint32(10)

tx, err := reService.CreateEnvelope(
    redenvelope.GROUP_FIXED,
    common.Address{},
    totalClaims,
    amountPerClaim,
    12*time.Hour,
    roomIdHash,
    common.Address{},     // No specific recipient
)
```

### 4. Create GROUP_RANDOM Envelope

Angpao grup dengan distribusi random:

```go
totalPot := big.NewInt(1000000000000000000) // 1.0 ETH total
totalClaims := uint32(8)

tx, err := reService.CreateEnvelope(
    redenvelope.GROUP_RANDOM,
    common.Address{},
    totalClaims,
    totalPot,             // Total pot, bukan per claim
    6*time.Hour,
    roomIdHash,
    common.Address{},
)
```

### 5. Claim Envelope

```go
envelopeId := big.NewInt(1)

// Check if already claimed
hasClaimed, err := reService.HasClaimed(envelopeId, yourAddress)

if !hasClaimed {
    tx, err := reService.ClaimEnvelope(envelopeId)
}
```

### 6. Get Envelope Info

```go
envelope, err := reService.GetEnvelope(envelopeId)

fmt.Printf("Creator: %s\n", envelope.Creator.Hex())
fmt.Printf("Kind: %d\n", envelope.Kind)
fmt.Printf("Total Claims: %d\n", envelope.TotalClaims)
fmt.Printf("Remaining Claims: %d\n", envelope.RemainingClaims)
fmt.Printf("Remaining Amount: %s\n", envelope.RemainingAmount.String())
```

### 7. Refund (after expiry)

```go
// Only creator can refund after expiry
tx, err := reService.RefundEnvelope(envelopeId)

## Prerequisites

1. **Hardhat node** harus berjalan:
   ```bash
   npx hardhat node
   ```

2. **Go** versi 1.21 atau lebih tinggi

## Installation

1. Install dependencies:
   ```bash
   go mod download
   ```

2. Jalankan aplikasi:
   ```bash
   go run main.go
   ```

## Output Example

```
=== Ethereum RPC Service Demo ===
Connected to: http://127.0.0.1:8545

Your address: 0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266

=== Get Balance ===
Balance: 9999.9937 ETH

=== Chain Information ===
Chain ID: 31337
Current block number: 3

=== Latest Block Info ===
Block Hash: 0x52b7b85700ef5e4fd364185a7abeb9da470e5cb6aaf86b49e9766f4d6a273f6c
Block Time: 1768458943
Transactions: 1

=== Sending ETH Transaction ===
From: 0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266
To: 0x70997970C51812dc3A010C7d01b50e0d17dc79C8
Amount: 1.0000 ETH
Transaction Hash: 0x3c3543d5ee6964a3aedcafe360bd1040f4d17c58cb792c3ca29750e8621effb8

=== Transaction Receipt ===
Waiting for transaction to be mined...
Status: 1 (1 = success, 0 = failed)
Block Number: 4
Gas Used: 21000

=== Updated Balances ===
Your balance: 9998.9936 ETH
Recipient balance: 10001.0000 ETH

=== Demo Completed Successfully! ===
```

## Cara Penggunaan

### 1. Connect ke Ethereum Node

```go
client, err := ethclient.Dial("http://127.0.0.1:8545")
if err != nil {
    log.Fatal(err)
}
defer client.Close()
```

### 2. Load Private Key

```go
privateKey, err := crypto.HexToECDSA("your_private_key_here")
if err != nil {
    log.Fatal(err)
}

publicKey := privateKey.Public()
publicKeyECDSA := publicKey.(*ecdsa.PublicKey)
address := crypto.PubkeyToAddress(*publicKeyECDSA)
```

### 3. Get Balance

```go
balance, err := client.BalanceAt(context.Background(), address, nil)
```

### 4. Get Block Information

```go
// Get latest block number
blockNumber, err := client.BlockNumber(context.Background())

// Get block by number
block, err := client.BlockByNumber(context.Background(), big.NewInt(blockNumber))
```

### 5. Send Transaction

```go
// Get nonce
nonce, err := client.PendingNonceAt(context.Background(), fromAddress)

// Get gas price
gasPrice, err := client.SuggestGasPrice(context.Background())

// Create transaction
tx := types.NewTransaction(nonce, toAddress, amount, gasLimit, gasPrice, nil)

// Sign transaction
signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)

// Send transaction
err = client.SendTransaction(context.Background(), signedTx)
```

### 6. Wait for Transaction Receipt

```go
receipt, err := client.TransactionReceipt(context.Background(), txHash)
fmt.Printf("Status: %d\n", receipt.Status) // 1 = success, 0 = failed
```

## Accounts Hardhat yang Digunakan

**Account #0** (Sender):
- Address: `0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266`
- Private Key: `0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80`
- Balance: 10000 ETH

**Account #1** (Recipient):
- Address: `0x70997970C51812dc3A010C7d01b50e0d17dc79C8`
- Private Key: `0x59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d`
- Balance: 10000 ETH

## Dependencies

```go
github.com/ethereum/go-ethereum v1.16.8
```

Package ini menyediakan:
- `ethclient` - Ethereum JSON-RPC client
- `crypto` - Cryptographic operations
- `common` - Common types (Address, Hash, etc.)
- `types` - Transaction and block types
- `accounts/abi` - ABI encoding/decoding

## Tips & Best Practices

### 1. Ganti Private Key
Untuk menggunakan account lain, ganti private key di [main.go](main.go):
```go
privateKeyHex := "your_private_key_here"
```

### 2. Error Handling
Selalu handle error dengan baik:
```go
if err != nil {
    log.Printf("Error: %v", err)
    // Handle error appropriately
}
```

### 3. Gas Management
- Transfer ETH: 21000 gas
- Smart contract deployment: 3000000 gas (estimasi)
- Smart contract interaction: varies (estimasi dengan `EstimateGas`)

### 4. Transaction Confirmation
Tunggu transaction receipt sebelum melanjutkan:
```go
receipt, err := client.TransactionReceipt(ctx, txHash)
if receipt.Status == 1 {
    // Success
} else {
    // Failed
}
```

### 5. Wei to Ether Conversion
```go
func weiToEther(wei *big.Int) string {
    ether := new(big.Float).SetInt(wei)
    ether = ether.Quo(ether, big.NewFloat(1e18))
    return ether.Text('f', 4)
}
```

## Advanced: Smart Contract Interaction

Untuk berinteraksi dengan smart contract, lihat implementasi di:
- [service/ethereum.go](service/ethereum.go) - Service layer untuk contract interaction
- [contract/contract.go](contract/contract.go) - ABI dan bytecode definition
- [contracts/SimpleStorage.sol](contracts/SimpleStorage.sol) - Solidity contract

## Troubleshooting

### Connection refused
**Problem**: `Failed to connect to ethereum node: dial tcp 127.0.0.1:8545: connect: connection refused`

**Solution**: Pastikan Hardhat node sudah running:
```bash
npx hardhat node
```

### Insufficient funds
**Problem**: Transaction gagal karena insufficient funds

**Solution**: 
- Check balance dengan `GetBalance`
- Pastikan gas price + amount tidak melebihi balance
- Gunakan account Hardhat yang memiliki 10000 ETH

### Transaction timeout
**Problem**: Transaction tidak ter-mine setelah beberapa saat

**Solution**:
- Hardhat akan auto-mine transactions
- Jika menggunakan manual mining, jalankan `await network.provider.send("evm_mine")`

### Wrong chain ID
**Problem**: `invalid sender` error

**Solution**: Pastikan menggunakan chain ID yang benar (Hardhat default: 31337)
```go
chainID, err := client.ChainID(context.Background())
```

## Test-Driven Development (TDD)

Project ini menggunakan **TDD approach** untuk memastikan kualitas code dan functionality yang reliable.

### 📊 Test Coverage

**13 comprehensive test cases** untuk 3 fungsi utama:

#### CreateEnvelope: 4 Tests ✅
- ✅ DIRECT_FIXED envelope creation
- ✅ GROUP_FIXED envelope creation
- ✅ GROUP_RANDOM envelope creation
- ✅ Envelope with room restriction

#### ClaimEnvelope: 4 Tests ✅
- ✅ First claim success
- ✅ Double claim prevention
- ✅ DIRECT_FIXED recipient claim
- ✅ All claims exhausted

#### RefundEnvelope: 5 Tests ⚠️
- ✅ Refund before expiry prevention
- ⏭️ Refund after expiry (manual test)
- ⏭️ Partial claims refund (manual test)
- ⏭️ Non-creator prevention (manual test)
- ⏭️ All claims exhausted (manual test)

**Success Rate: 100% (9/9 executed tests passed)** 🎉

### 🚀 Running Tests

Jalankan semua TDD tests:
```bash
cd redenvelope
go test -v -run "TestCreateEnvelope|TestClaimEnvelope|TestRefundEnvelope"
```

Jalankan tests per fungsi:
```bash
# CreateEnvelope tests
go test -v -run TestCreateEnvelope

# ClaimEnvelope tests
go test -v -run TestClaimEnvelope

# RefundEnvelope tests
go test -v -run TestRefundEnvelope
```

Jalankan test spesifik:
```bash
go test -v -run TestClaimEnvelope_DoubleClaim_ShouldFail
```

### 📚 TDD Documentation

Lihat dokumentasi lengkap:
- **[TDD_GUIDE.md](TDD_GUIDE.md)** - Metodologi, best practices, troubleshooting
- **[TDD_SUMMARY.md](TDD_SUMMARY.md)** - Hasil implementasi dan summary

### 💡 Benefits TDD

- ✅ **Confidence**: Semua fungsi core terverifikasi
- ✅ **Documentation**: Tests menjelaskan expected behavior
- ✅ **Safety**: Catch regressions early
- ✅ **Design**: Better code structure
- ✅ **Speed**: Faster debugging

## Next Steps

Untuk development lebih lanjut, Anda bisa:

1. **Add Smart Contract Integration**
   - Deploy custom contracts
   - Call contract functions
   - Listen to contract events

2. **Implement REST API**
   - Wrap functionality dalam HTTP server
   - Expose endpoints untuk RPC operations

3. **Add Database**
   - Store transaction history
   - Track account balances

4. **Add Monitoring**
   - Log all transactions
   - Monitor gas usage
   - Alert on failed transactions

## Resources

- [go-ethereum Documentation](https://geth.ethereum.org/docs)
- [Ethereum JSON-RPC](https://ethereum.org/en/developers/docs/apis/json-rpc/)
- [Hardhat Documentation](https://hardhat.org/docs)
- [Solidity Documentation](https://docs.soliditylang.org/)

