package main

import (
	"context"
	"log"

	ethereumv1 "ethereum-grpc/protos/gen/go/ethereum"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	infura_url = "https://mainnet.infura.io/v3/d46bdacab9a64306a97e95f419caebbd"
	hex_address = "d46bdacab9a64306a97e95f419caebbd"
)

func main() {
	// Connect to Ethereum node via Infura
	client, err := ethclient.Dial(infura_url)
	if err != nil {
		log.Fatal("Failed to create Ethereum client:%v", err)
	}
	defer client.Close()

	// Create gRPC client connection to the gRPC server
	conn, err := grpc.NewClient(":44044", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to create grpc client: %v", err)
	}
	defer conn.Close()

	clientGRPC := ethereumv1.NewEthereumServiceClient(conn)
	crypto_signature := createSignature([]byte(hex_address))
	address := []byte(hex_address)

	account, err := clientGRPC.GetAccount(context.Background(), &ethereumv1.GetAccountRequest{
		EthereumAddress: address,
		CryptoSignature: crypto_signature,
	})
	if err != nil {
		log.Fatalf("Failed to get account: %v", err)
	}
	log.Printf("Balance: %s, Nonce: %d", account.GastokenBalance, account.WalletNonce)

}

func createSignature(data []byte) []byte {
	hash := crypto.Keccak256Hash(data)
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		log.Fatal("Failed to generate private key:%v", err)
	}

	signature, err := crypto.Sign(hash.Bytes(), privateKey)
	if err != nil {
		log.Fatal("Failed to sign data:%v", err)
	}
	return signature
}
