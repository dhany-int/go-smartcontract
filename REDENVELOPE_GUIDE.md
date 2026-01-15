# RedEnvelope Smart Contract - Go Integration Guide

Panduan lengkap untuk berinteraksi dengan RedEnvelope Smart Contract menggunakan Go.

## Setup

### 1. Deploy RedEnvelope Contract

Pertama, deploy contract di Hardhat:

```bash
# Terminal 1: Start Hardhat node
npx hardhat node

# Terminal 2: Deploy contract
npx hardhat run scripts/deploy.js --network localhost
```

Copy contract address yang muncul, misalnya:
```
RedEnvelope deployed to: 0x5FbDB2315678afecb367f032d93F642f64180aa3
```

### 2. Update Contract Address

Edit file [main.go](main.go) dan [examples/redenvelope_all.go](examples/redenvelope_all.go):

```go
contractAddress := "0x5FbDB2315678afecb367f032d93F642f64180aa3" // GANTI INI!
```

### 3. Install Dependencies

```bash
go mod download
```

### 4. Run Demo

```bash
# Demo utama (RPC + RedEnvelope)
go run main.go

# Demo semua jenis envelope
cd examples
go run redenvelope_all.go

# Demo RPC methods
go run rpc_methods.go
```

## API Reference

### Initialize Service

```go
import "rpcsol/redenvelope"

reService, err := redenvelope.NewRedEnvelopeService(
    "http://127.0.0.1:8545",              // RPC URL
    "0xYourContractAddress",               // Contract address
    "your_private_key_hex",                // Private key
)
defer reService.Client.Close()
```

### 1. Create DIRECT_FIXED Envelope

Angpao untuk 1 orang spesifik dengan jumlah tetap.

```go
recipient := common.HexToAddress("0x70997970C51812dc3A010C7d01b50e0d17dc79C8")
amount := big.NewInt(500000000000000000) // 0.5 ETH
roomIdHash := redenvelope.GenerateRoomIdHash("room-123")

tx, err := reService.CreateEnvelope(
    redenvelope.DIRECT_FIXED,
    common.Address{},     // Native token (ETH/BNB)
    1,                    // Total claims = 1 (hanya 1 orang)
    amount,               // Jumlah
    24*time.Hour,         // Expiry 24 jam
    roomIdHash,           // Room ID hash
    recipient,            // Penerima spesifik
)
```

**Use Case**: Gift untuk 1 orang, salary payment, bounty reward

### 2. Create GROUP_FIXED Envelope

Angpao grup dengan jumlah tetap per orang.

```go
amountPerClaim := big.NewInt(200000000000000000) // 0.2 ETH per orang
totalClaims := uint32(10)                         // 10 orang bisa klaim
roomIdHash := redenvelope.GenerateRoomIdHash("group-room-001")

tx, err := reService.CreateEnvelope(
    redenvelope.GROUP_FIXED,
    common.Address{},
    totalClaims,
    amountPerClaim,      // Jumlah PER KLAIM
    12*time.Hour,
    roomIdHash,
    common.Address{},    // Tidak perlu recipient spesifik
)
```

**Use Case**: Giveaway dengan hadiah sama rata, airdrop, reward distribution

### 3. Create GROUP_RANDOM Envelope

Angpao grup dengan distribusi random (luck-based).

```go
totalPot := big.NewInt(1000000000000000000) // 1.0 ETH total pot
totalClaims := uint32(8)                     // 8 orang bisa klaim
roomIdHash := redenvelope.GenerateRoomIdHash("lucky-room-001")

tx, err := reService.CreateEnvelope(
    redenvelope.GROUP_RANDOM,
    common.Address{},
    totalClaims,
    totalPot,            // TOTAL POT, bukan per klaim
    6*time.Hour,
    roomIdHash,
    common.Address{},
)
```

**Use Case**: Lucky draw, Chinese New Year angpao, gamified rewards

**Algorithm**: Setiap claim mendapat random amount antara 1 wei sampai `(remainingAmount * 2) / remainingClaims`. Klaim terakhir mendapat semua sisa.

### 4. Claim Envelope

```go
envelopeId := big.NewInt(1)

// Optional: Check if already claimed
hasClaimed, err := reService.HasClaimed(envelopeId, yourAddress)
if hasClaimed {
    fmt.Println("Already claimed!")
    return
}

// Claim
tx, err := reService.ClaimEnvelope(envelopeId)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Transaction: %s\n", tx.Hash().Hex())

// Wait for confirmation
time.Sleep(2 * time.Second)
receipt, _ := reService.Client.TransactionReceipt(context.Background(), tx.Hash())
if receipt.Status == 1 {
    fmt.Println("Claim successful!")
}
```

### 5. Get Envelope Information

```go
envelope, err := reService.GetEnvelope(envelopeId)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Creator: %s\n", envelope.Creator.Hex())
fmt.Printf("Kind: %d (0=DIRECT_FIXED, 1=GROUP_FIXED, 2=GROUP_RANDOM)\n", envelope.Kind)
fmt.Printf("Total Claims: %d\n", envelope.TotalClaims)
fmt.Printf("Remaining Claims: %d\n", envelope.RemainingClaims)
fmt.Printf("Remaining Amount: %s wei\n", envelope.RemainingAmount.String())
fmt.Printf("Expiry: %s\n", time.Unix(int64(envelope.Expiry), 0))
```

### 6. Check Claim Status

```go
user := common.HexToAddress("0x...")
hasClaimed, err := reService.HasClaimed(envelopeId, user)

if hasClaimed {
    fmt.Println("User already claimed")
} else {
    fmt.Println("User can claim")
}
```

### 7. Refund Expired Envelope

```go
// Only creator can refund after expiry
tx, err := reService.RefundEnvelope(envelopeId)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Refund transaction: %s\n", tx.Hash().Hex())
```

### 8. Get Next Envelope ID

```go
nextId, err := reService.GetNextEnvelopeId()
fmt.Printf("Next envelope ID will be: %s\n", nextId.String())
```

## Helper Functions

### Generate Room ID Hash

```go
roomIdHash := redenvelope.GenerateRoomIdHash("my-room-123")
// Returns [32]byte hash using Keccak256
```

### Wei to Ether Conversion

```go
func weiToEther(wei *big.Int) string {
    ether := new(big.Float).SetInt(wei)
    ether = ether.Quo(ether, big.NewFloat(1e18))
    return ether.Text('f', 4)
}

amount := big.NewInt(1000000000000000000) // 1 ETH in wei
fmt.Println(weiToEther(amount))           // "1.0000"
```

## Error Handling

```go
tx, err := reService.CreateEnvelope(...)
if err != nil {
    if strings.Contains(err.Error(), "insufficient funds") {
        log.Println("Not enough balance")
    } else if strings.Contains(err.Error(), "already claimed") {
        log.Println("Already claimed this envelope")
    } else {
        log.Printf("Error: %v", err)
    }
    return
}

// Wait and check receipt
time.Sleep(2 * time.Second)
receipt, err := reService.Client.TransactionReceipt(context.Background(), tx.Hash())
if err != nil {
    log.Printf("Failed to get receipt: %v", err)
    return
}

if receipt.Status == 0 {
    log.Println("Transaction failed!")
} else {
    log.Println("Transaction successful!")
}
```

## Complete Examples

Lihat file-file berikut untuk contoh lengkap:

1. **[main.go](main.go)** - Demo RPC operations + RedEnvelope basic
2. **[examples/redenvelope_all.go](examples/redenvelope_all.go)** - Demo semua jenis envelope
3. **[examples/rpc_methods.go](examples/rpc_methods.go)** - Demo semua RPC methods

## Testing

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Build binary
go build -o redenvelope-demo main.go
./redenvelope-demo
```

## Troubleshooting

### "no contract code at given address"
**Problem**: Contract belum di-deploy atau address salah

**Solution**: 
1. Deploy contract terlebih dahulu
2. Copy address yang benar
3. Update di main.go

### "insufficient funds"
**Problem**: Balance tidak cukup untuk create envelope + gas

**Solution**: Pastikan balance > (envelope amount + gas fee)

### "already claimed"
**Problem**: User sudah claim envelope ini

**Solution**: Check dengan `HasClaimed()` sebelum claim

### "envelope expired"
**Problem**: Envelope sudah kadaluarsa

**Solution**: Creator bisa refund dengan `RefundEnvelope()`

## Fee Calculation

Contract menggunakan fee system:

```
feeAmount = (totalAmount * feeBps) / 10000
netAmount = totalAmount - feeAmount
```

Contoh dengan fee 2.5% (250 bps):
- Deposit: 1.0 ETH
- Fee: 0.025 ETH (ke treasury)
- Net: 0.975 ETH (untuk klaim)

## Security Notes

1. **Private Keys**: Jangan hardcode private key di production. Gunakan environment variable atau secret manager.

2. **Random Distribution**: Implementasi saat ini menggunakan pseudo-random. Untuk production, gunakan Chainlink VRF.

3. **Gas Limits**: Sudah di-set reasonable default, tapi bisa disesuaikan jika perlu.

4. **Expiry Time**: Pastikan expiry time reasonable (tidak terlalu pendek atau panjang).

5. **Room ID**: Room ID hash hanya untuk tracking, tidak enforce access control di contract.

## Production Checklist

- [ ] Deploy contract ke testnet dulu
- [ ] Test semua fungsi (create, claim, refund)
- [ ] Verify contract di block explorer
- [ ] Set treasury address yang benar
- [ ] Configure fee (feeBps)
- [ ] Monitor gas usage
- [ ] Implement proper error handling
- [ ] Use environment variables untuk sensitive data
- [ ] Add logging untuk audit trail
- [ ] Test dengan multiple accounts
- [ ] Consider Chainlink VRF untuk random

## Resources

- [Ethereum go-ethereum Documentation](https://geth.ethereum.org/docs)
- [Hardhat Documentation](https://hardhat.org/)
- [ABI Specification](https://docs.soliditylang.org/en/latest/abi-spec.html)
- [EIP-155 (Transaction Signing)](https://eips.ethereum.org/EIPS/eip-155)

## License

MIT
