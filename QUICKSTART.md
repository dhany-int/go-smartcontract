# Quick Start Guide

Panduan cepat untuk memulai menggunakan Go RPC Service dengan Hardhat.

## 🚀 Quick Setup (5 menit)

### 1. Start Hardhat Node

Di terminal pertama, jalankan Hardhat node:

```bash
npx hardhat node
```

Anda akan melihat output seperti ini:
```
Started HTTP and WebSocket JSON-RPC server at http://127.0.0.1:8545/

Accounts
========
Account #0: 0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266 (10000 ETH)
Private Key: 0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80
...
```

### 2. Install Go Dependencies

Di terminal kedua, masuk ke project directory dan install dependencies:

```bash
cd go-rpc-sol
go mod download
```

### 3. Run Demo

Jalankan aplikasi:

```bash
go run main.go
```

Output yang diharapkan:
```
=== Ethereum RPC Service Demo ===
Connected to: http://127.0.0.1:8545

Your address: 0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266

=== Get Balance ===
Balance: 10000.0000 ETH

=== Sending ETH Transaction ===
From: 0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266
To: 0x70997970C51812dc3A010C7d01b50e0d17dc79C8
Amount: 1.0000 ETH
Transaction Hash: 0x...

=== Demo Completed Successfully! ===
```

## 📚 Explore Examples

### Test All RPC Methods
```bash
cd examples
go run rpc_methods.go
```

### Smart Contract Interaction (Advanced)
```bash
cd examples
go run contract_interaction.go
```

## 🎯 Common Use Cases

### 1. Check Account Balance

```go
package main

import (
    "context"
    "fmt"
    "log"
    "math/big"
    
    "github.com/ethereum/go-ethereum/common"
    "github.com/ethereum/go-ethereum/ethclient"
)

func main() {
    client, _ := ethclient.Dial("http://127.0.0.1:8545")
    address := common.HexToAddress("0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266")
    
    balance, err := client.BalanceAt(context.Background(), address, nil)
    if err != nil {
        log.Fatal(err)
    }
    
    ether := new(big.Float).Quo(new(big.Float).SetInt(balance), big.NewFloat(1e18))
    fmt.Printf("Balance: %s ETH\n", ether.Text('f', 4))
}
```

### 2. Send ETH Transaction

Lihat [main.go](main.go) untuk implementasi lengkap send transaction.

### 3. Get Block Information

```go
blockNumber, _ := client.BlockNumber(context.Background())
block, _ := client.BlockByNumber(context.Background(), big.NewInt(int64(blockNumber)))

fmt.Printf("Block: %d\n", block.Number().Uint64())
fmt.Printf("Hash: %s\n", block.Hash().Hex())
fmt.Printf("Transactions: %d\n", len(block.Transactions()))
```

## 🔧 Configuration

### Change RPC URL

Edit di [main.go](main.go):
```go
rpcURL := "http://127.0.0.1:8545"  // Hardhat default
// rpcURL := "http://127.0.0.1:8545"  // Ganache
// rpcURL := "https://mainnet.infura.io/v3/YOUR-API-KEY"  // Mainnet
```

### Use Different Account

Edit private key di [main.go](main.go):
```go
// Account #0 (default)
privateKeyHex := "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"

// Account #1
// privateKeyHex := "59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d"
```

## 🐛 Troubleshooting

### Error: "connection refused"
**Problem**: Hardhat node belum berjalan

**Solution**:
```bash
npx hardhat node
```

### Error: "insufficient funds"
**Problem**: Account tidak memiliki ETH yang cukup

**Solution**: Gunakan account Hardhat yang memiliki 10000 ETH atau kurangi jumlah transfer

### Error: "nonce too low"
**Problem**: Nonce conflict

**Solution**: Restart Hardhat node untuk reset state
```bash
# Ctrl+C untuk stop
npx hardhat node
```

## 📖 Next Steps

1. ✅ Baca [README.md](README.md) untuk dokumentasi lengkap
2. ✅ Explore [examples/](examples/) untuk contoh lebih lanjut
3. ✅ Check [service/ethereum.go](service/ethereum.go) untuk advanced features
4. ✅ Lihat [contracts/SimpleStorage.sol](contracts/SimpleStorage.sol) untuk smart contract example

## 🎓 Learning Resources

- [Go Ethereum Documentation](https://geth.ethereum.org/docs)
- [Ethereum JSON-RPC Spec](https://ethereum.org/en/developers/docs/apis/json-rpc/)
- [Hardhat Documentation](https://hardhat.org/docs)
- [Solidity by Example](https://solidity-by-example.org/)

## 💡 Tips

1. **Monitor Hardhat console** untuk melihat transactions real-time
2. **Use Hardhat accounts** yang sudah pre-funded dengan 10000 ETH
3. **Check transaction receipt** untuk verify transaction status
4. **Handle errors properly** untuk production code
5. **Test di local network** sebelum deploy ke testnet/mainnet

## ✨ What's Next?

Setelah familiar dengan basics, coba:

- Deploy smart contract sendiri
- Implement event listening
- Build REST API wrapper
- Add database untuk transaction history
- Integrate dengan front-end application

---

**Happy Coding! 🚀**

Jika ada pertanyaan, check [README.md](README.md) atau examples di [examples/](examples/).
