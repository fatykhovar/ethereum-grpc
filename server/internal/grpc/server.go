package grpc

import (
	"context"
	ethereumv1 "ethereum-grpc/protos/gen/go/ethereum"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"google.golang.org/grpc"
)

var infura_url  = "https://mainnet.infura.io/v3/d46bdacab9a64306a97e95f419caebbd"

type serverAPI struct {
	ethereumv1.UnimplementedEthereumServiceServer
}

func Register(gRPCServer *grpc.Server) {
	ethereumv1.RegisterEthereumServiceServer(gRPCServer, &serverAPI{})
}

func (s *serverAPI) GetAccount(ctx context.Context, req *ethereumv1.GetAccountRequest) (res *ethereumv1.GetAccountResponse, err error) {
	isValidSignature := validateSignature(req.EthereumAddress, req.CryptoSignature)
	// fmt.Printf("is valid: %t", isValidSignature)

	if !isValidSignature {
		return nil, fmt.Errorf("invalid signature")
	}

	client, err := ethclient.Dial(infura_url)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client:%v", err)
	}
	defer client.Close()

	balance, err := client.BalanceAt(context.Background(), common.BytesToAddress(req.EthereumAddress), nil)
	if err != nil {
		return nil, fmt.Errorf("Failed to get balance:%v", err)
	}

	nonce, err := client.PendingNonceAt(context.Background(), common.BytesToAddress(req.EthereumAddress))
	if err != nil {
		return nil, fmt.Errorf("Failed to get nonce:%v", err)
	}

	return &ethereumv1.GetAccountResponse{
		GastokenBalance: balance.String(),
        WalletNonce:     nonce,
	}, nil
}

func validateSignature(address []byte, signature []byte) bool {
	publicKey, err := crypto.Ecrecover(address[:], signature)
	if err != nil {
		log.Fatal("Failed to return public key from signature: ", err)
	}

	verified := crypto.VerifySignature(publicKey, address, signature[:len(signature)-1])
	return verified
}
