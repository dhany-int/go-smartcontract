package redenvelope

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"reflect"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

type RedEnvelopeService struct {
	Client          *ethclient.Client
	ContractAddress common.Address
	PrivateKey      *ecdsa.PrivateKey
	Address         common.Address
	ChainID         *big.Int
	ABI             abi.ABI
}

// Envelope struct sesuai dengan contract
type Envelope struct {
	Creator         common.Address
	Token           common.Address
	Kind            uint8
	AmountPerClaim  *big.Int
	RemainingAmount *big.Int
	TotalClaims     uint32
	RemainingClaims uint32
	ClaimIndex      uint32
	Expiry          uint64
	RoomIdHash      [32]byte
	Recipient       common.Address
}

// NewRedEnvelopeService membuat instance baru RedEnvelope service
func NewRedEnvelopeService(rpcURL string, contractAddress string, privateKeyHex string) (*RedEnvelopeService, error) {
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to ethereum node: %v", err)
	}

	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return nil, fmt.Errorf("failed to load private key: %v", err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("failed to cast public key to ECDSA")
	}
	address := crypto.PubkeyToAddress(*publicKeyECDSA)

	chainID, err := client.ChainID(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to get chain ID: %v", err)
	}

	parsedABI, err := abi.JSON(strings.NewReader(RedEnvelopeABI))
	if err != nil {
		return nil, fmt.Errorf("failed to parse ABI: %v", err)
	}

	return &RedEnvelopeService{
		Client:          client,
		ContractAddress: common.HexToAddress(contractAddress),
		PrivateKey:      privateKey,
		Address:         address,
		ChainID:         chainID,
		ABI:             parsedABI,
	}, nil
}

// CreateEnvelope membuat envelope baru
// Parameter amount:
//   - DIRECT_FIXED: amount untuk penerima
//   - GROUP_FIXED: amount PER CLAIM (contract akan × totalClaims)
//   - GROUP_RANDOM: total pot
//
// msg.value (untuk native token):
//   - DIRECT_FIXED: amount
//   - GROUP_FIXED: amount × totalClaims
//   - GROUP_RANDOM: amount
//
// Fee diambil DARI grossPot, bukan ditambahkan!
func (s *RedEnvelopeService) CreateEnvelope(
	kind uint8,
	token common.Address,
	totalClaims uint32,
	amount *big.Int,
	expiryDuration time.Duration,
	roomIdHash [32]byte,
	recipient common.Address,
) (*types.Transaction, error) {
	expiry := uint64(time.Now().Add(expiryDuration).Unix())

	nonce, err := s.Client.PendingNonceAt(context.Background(), s.Address)
	if err != nil {
		return nil, fmt.Errorf("failed to get nonce: %v", err)
	}

	gasPrice, err := s.Client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to get gas price: %v", err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(s.PrivateKey, s.ChainID)
	if err != nil {
		return nil, fmt.Errorf("failed to create transactor: %v", err)
	}
	auth.Nonce = big.NewInt(int64(nonce))
	auth.GasLimit = uint64(500000)
	auth.GasPrice = gasPrice

	// Calculate grossPot (amount yang harus dikirim sebagai msg.value)
	var grossPot *big.Int
	if kind == GROUP_FIXED {
		// GROUP_FIXED: grossPot = amount × totalClaims
		grossPot = new(big.Int).Mul(amount, big.NewInt(int64(totalClaims)))
	} else {
		// DIRECT_FIXED & GROUP_RANDOM: grossPot = amount
		grossPot = amount
	}

	if token == (common.Address{}) {
		auth.Value = grossPot
	} else {
		auth.Value = big.NewInt(0)
	}

	boundContract := bind.NewBoundContract(s.ContractAddress, s.ABI, s.Client, s.Client, s.Client)

	tx, err := boundContract.Transact(auth, "createEnvelope", kind, token, totalClaims, amount, expiry, roomIdHash, recipient)
	if err != nil {
		return nil, fmt.Errorf("failed to create envelope: %v", err)
	}

	return tx, nil
}

// ClaimEnvelope klaim envelope
func (s *RedEnvelopeService) ClaimEnvelope(envelopeId *big.Int) (*types.Transaction, error) {
	nonce, err := s.Client.PendingNonceAt(context.Background(), s.Address)
	if err != nil {
		return nil, fmt.Errorf("failed to get nonce: %v", err)
	}

	gasPrice, err := s.Client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to get gas price: %v", err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(s.PrivateKey, s.ChainID)
	if err != nil {
		return nil, fmt.Errorf("failed to create transactor: %v", err)
	}
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)
	auth.GasLimit = uint64(300000)
	auth.GasPrice = gasPrice

	boundContract := bind.NewBoundContract(s.ContractAddress, s.ABI, s.Client, s.Client, s.Client)

	tx, err := boundContract.Transact(auth, "claimEnvelope", envelopeId)
	if err != nil {
		return nil, fmt.Errorf("failed to claim envelope: %v", err)
	}

	return tx, nil
}

// GetEnvelope mendapatkan informasi envelope
func (s *RedEnvelopeService) GetEnvelope(envelopeId *big.Int) (*Envelope, error) {
	boundContract := bind.NewBoundContract(s.ContractAddress, s.ABI, s.Client, s.Client, s.Client)

	var result []interface{}
	err := boundContract.Call(&bind.CallOpts{}, &result, "getEnvelope", envelopeId)
	if err != nil {
		return nil, fmt.Errorf("failed to get envelope: %v", err)
	}

	// Gunakan reflection untuk handle struct dengan atau tanpa JSON tags
	envelopeData := result[0]

	envelope := &Envelope{}

	// Extract data menggunakan reflection
	val := reflect.ValueOf(envelopeData)

	envelope.Creator = val.Field(0).Interface().(common.Address)
	envelope.Token = val.Field(1).Interface().(common.Address)
	envelope.Kind = val.Field(2).Interface().(uint8)
	envelope.AmountPerClaim = val.Field(3).Interface().(*big.Int)
	envelope.RemainingAmount = val.Field(4).Interface().(*big.Int)
	envelope.TotalClaims = val.Field(5).Interface().(uint32)
	envelope.RemainingClaims = val.Field(6).Interface().(uint32)
	envelope.ClaimIndex = val.Field(7).Interface().(uint32)
	envelope.Expiry = val.Field(8).Interface().(uint64)
	envelope.RoomIdHash = val.Field(9).Interface().([32]byte)
	envelope.Recipient = val.Field(10).Interface().(common.Address)

	return envelope, nil
}

// HasClaimed cek apakah user sudah klaim
func (s *RedEnvelopeService) HasClaimed(envelopeId *big.Int, user common.Address) (bool, error) {
	boundContract := bind.NewBoundContract(s.ContractAddress, s.ABI, s.Client, s.Client, s.Client)

	var result []interface{}
	err := boundContract.Call(&bind.CallOpts{}, &result, "hasUserClaimed", envelopeId, user)
	if err != nil {
		return false, fmt.Errorf("failed to check claim status: %v", err)
	}

	return result[0].(bool), nil
}

// RefundEnvelope refund envelope setelah expiry
func (s *RedEnvelopeService) RefundEnvelope(envelopeId *big.Int) (*types.Transaction, error) {
	nonce, err := s.Client.PendingNonceAt(context.Background(), s.Address)
	if err != nil {
		return nil, fmt.Errorf("failed to get nonce: %v", err)
	}

	gasPrice, err := s.Client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to get gas price: %v", err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(s.PrivateKey, s.ChainID)
	if err != nil {
		return nil, fmt.Errorf("failed to create transactor: %v", err)
	}
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)
	auth.GasLimit = uint64(300000)
	auth.GasPrice = gasPrice

	boundContract := bind.NewBoundContract(s.ContractAddress, s.ABI, s.Client, s.Client, s.Client)

	tx, err := boundContract.Transact(auth, "refundEnvelope", envelopeId)
	if err != nil {
		return nil, fmt.Errorf("failed to refund envelope: %v", err)
	}

	return tx, nil
}

// GetNextEnvelopeId mendapatkan next envelope ID
func (s *RedEnvelopeService) GetNextEnvelopeId() (*big.Int, error) {
	boundContract := bind.NewBoundContract(s.ContractAddress, s.ABI, s.Client, s.Client, s.Client)

	var result []interface{}
	err := boundContract.Call(&bind.CallOpts{}, &result, "nextEnvelopeId")
	if err != nil {
		return nil, fmt.Errorf("failed to get next envelope ID: %v", err)
	}

	return result[0].(*big.Int), nil
}

// GenerateRoomIdHash helper untuk generate room ID hash
func GenerateRoomIdHash(roomId string) [32]byte {
	hash := crypto.Keccak256Hash([]byte(roomId))
	return hash
}
