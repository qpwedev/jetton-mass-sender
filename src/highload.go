package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/log"

	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/liteclient"
	"github.com/xssnick/tonutils-go/tlb"
	"github.com/xssnick/tonutils-go/ton"
	"github.com/xssnick/tonutils-go/ton/jetton"
	"github.com/xssnick/tonutils-go/ton/wallet"
)

type MessageEntry struct {
	Amount  string `json:"amount"`
	Address string `json:"address"`
}

func massSender(seedPhrase string, jettonMasterAddress string, commentary string, messageEntryFilename string) {
	client := liteclient.NewConnectionPool()

	// Connect to mainnet lite server
	err := client.AddConnection(context.Background(), "135.181.140.212:13206", "K0t3+IWLOXHYMvMcrGZDPs+pn58a17LFbnXoQkKc2xw=")
	if err != nil {
		log.Error("connection err:", err.Error())
		return
	}

	ctx := client.StickyContext(context.Background())
	api := ton.NewAPIClient(client, ton.ProofCheckPolicyFast).WithRetry()

	words := strings.Split(seedPhrase, " ")

	log.Infof("Seed words: %s", words)

	// Initialize highload wallet
	w, err := wallet.FromSeed(api, words, wallet.HighloadV2R2)
	if err != nil {
		log.Error("FromSeed err:", err.Error())
		return
	}

	log.Infof("Wallet address: %s", w.WalletAddress())


	block, err := api.CurrentMasterchainInfo(context.Background())
	if err != nil {
		log.Error("CurrentMasterchainInfo err:", err.Error())
		return
	}

	balance, err := w.GetBalance(context.Background(), block)
	if err != nil {
		log.Error("GetBalance err:", err.Error())
		return
	}

	// Read message entries from file
	data, err := os.ReadFile(messageEntryFilename)
	if err != nil {
		log.Error("Error reading file:", err.Error())
		return
	}

	var messages []MessageEntry
	err = json.Unmarshal(data, &messages)
	if err != nil {
		log.Error("Error unmarshalling JSON:", err.Error())
		return
	}

	// if balance < len(messages) * 0.05 + their amount then exit
	sum := 0.0
	for _, msg := range messages {
		// from string to float
		amount, err := strconv.ParseFloat(msg.Amount, 64)
		if err != nil {
			log.Error("Error parsing amount:", err.Error())
			return
		}

		sum += amount
	}

	if float64(balance.Nano().Int64()) < (float64(len(messages)) * 5e7) + sum {
		log.Error("Not enough balance to send all messages")
		return
	}

	// Initialize token wallet
	token := jetton.NewJettonMasterClient(api, address.MustParseAddr(jettonMasterAddress))
	jettonWallet, err := token.GetJettonWallet(ctx, w.WalletAddress())
	if err != nil {
		log.Fatal(err)
	}

	jettonBalance, err := jettonWallet.GetBalance(ctx)
	if err != nil {
		log.Error("GetBalance err:", err.Error())
		return
	}

	log.Infof("Balance: %s jettons", jettonBalance)


	// Start sending messages in batches
	const batchSize = 100
	for i := 0; i < len(messages); i += batchSize {
		end := i + batchSize
		if end > len(messages) {
			end = len(messages)
		}

		batch := messages[i:end]
		var walletMessages []*wallet.Message
		for _, msg := range batch {
			log.Printf(msg.Address, msg.Amount)
			amountTokens := tlb.MustFromDecimal(msg.Amount, 9)
			comment, err := wallet.CreateCommentCell(commentary)
			if err != nil {
				log.Error("Error creating comment cell:", err.Error())
				log.Fatal(err)
			}
			to := address.MustParseAddr(msg.Address)
			transferPayload, err := jettonWallet.BuildTransferPayload(to, amountTokens, tlb.ZeroCoins, comment)
			if err != nil {
				log.Error("Error building transfer payload:", err.Error())
				log.Fatal(err)
			}
			walletMsg := wallet.SimpleMessage(jettonWallet.Address(), tlb.MustFromTON("0.05"), transferPayload)
			walletMessages = append(walletMessages, walletMsg)
		}

		log.Infof("Sending transaction and waiting for confirmation for batch starting from index %s...")

		txHash, err := w.SendManyWaitTxHash(ctx, walletMessages)
		if err != nil {
			log.Error("Transfer err:", err.Error())
			return
		}

		log.Infof("Batch transaction sent, hash: %s", base64.StdEncoding.EncodeToString(txHash))
		log.Infof("Explorer link: https://tonscan.org/tx/ %s", base64.URLEncoding.EncodeToString(txHash))

		log.Info("Sleeping for 30 seconds...")
		time.Sleep(30 * time.Second)
	}
}
