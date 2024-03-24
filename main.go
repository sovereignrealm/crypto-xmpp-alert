package main

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/sovereignrealm/crypto-xmpp-alert/internal/services/price"
	"github.com/sovereignrealm/crypto-xmpp-alert/internal/services/xmpp"
	"github.com/sovereignrealm/crypto-xmpp-alert/kit/config"
)

type CryptoTransaction struct {
	PurchasedDate  string
	PurchasedPrice float64
	CryptoAmount   float64
}

type CryptoAssetSS struct {
	Name     price.CrytpoAsset
	Boundary float64
}

type FileReader interface {
	ReadFile(cryptoAsset string) ([]byte, error)
}

type FileWriter interface {
	WriteFile(filename string, data []byte, perm fs.FileMode) error
}

type DefaultFileReaderWriter struct{}

func (r *DefaultFileReaderWriter) ReadFile(filename string) ([]byte, error) {
	return os.ReadFile(filename)
}

func (r *DefaultFileReaderWriter) WriteFile(filename string, data []byte, perm fs.FileMode) error {
	return os.WriteFile(filename, data, perm)
}

func main() {
	cfg := config.NewConfig(os.Getenv("ENV_FILE"))
	ps := price.NewPriceService()
	xs := xmpp.NewXmppService(xmpp.XmppServiceOptions{
		Username:     cfg.XmppUsername,
		Password:     cfg.XmppPassword,
		Server:       cfg.XmppServer,
		RecipientJID: cfg.XmppRecipientJid,
	})
	assets := []CryptoAssetSS{
		{Name: price.CrytpoAssetBitcoin, Boundary: 200.0},
		{Name: price.CrytpoAssetEthereum, Boundary: 100.0},
		{Name: price.CrytpoAssetCardano, Boundary: 100.0},
		{Name: price.CrytpoAssetPolkadot, Boundary: 50.0},
	}

	jsonData, err := os.ReadFile("./json/data.json")
	if err != nil {
		panic(fmt.Errorf("error reading JSON file: %w", err))
	}

	var transactions map[string][]CryptoTransaction
	if err := json.Unmarshal(jsonData, &transactions); err != nil {
		panic(fmt.Errorf("error unmarshalling JSON: %w", err))
	}

	for _, asset := range assets {
		transactions, ok := transactions[string(asset.Name)]
		if !ok {
			continue
		}
		totalAmountInvested, totalCryptoAsset := calculateTotalValues(transactions)
		currentPrice, err := ps.GetCurrentPrice(asset.Name)
		if err != nil {
			if err := xs.SendMessage(fmt.Sprintf("Error fetching current %s price", asset.Name)); err != nil {
				// TODO: Should panic???
				log.Println("error while sending xmpp message:", err)
			}
			continue
		}
		if currentPrice == 0 {
			if err := xs.SendMessage(fmt.Sprintf("Error fetching current %s price", asset.Name)); err != nil {
				// TODO: Should panic???
				log.Println("error while sending xmpp message:", err)
			}
			continue
		}
		fmt.Printf("-------- %s  --------\n", asset.Name)
		printSummary(asset.Name, totalAmountInvested, totalCryptoAsset, currentPrice)
		var fileReader FileReader = &DefaultFileReaderWriter{}
		readData, err := readFile(fileReader, asset.Name)
		if err != nil {
			fmt.Println("Error reading file: ", err)
			return
		}
		if isFirstTime(readData) {
			moneyNow := totalCryptoAsset * currentPrice
			totalGained := calculatePercentageChange(totalAmountInvested, moneyNow)
			if totalGained >= asset.Boundary {
				strVal := strconv.FormatFloat(totalGained, 'f', 2, 64)
				if err := xs.SendMessage(fmt.Sprintf("You have gained in %s: %s%%", asset.Name, strVal)); err != nil {
					// TODO: Should panic???
					log.Println("error while sending xmpp message:", err)
				}
			}
			var fileWriter FileWriter = &DefaultFileReaderWriter{}
			if err := writeFile(fileWriter, string(asset.Name)); err != nil {
				// TODO: Should panic???
				log.Println("error writing to file:", err)
			}
		}
		fmt.Printf("-------- END %s  --------\n", asset.Name)
	}
}

// TODO: Refactor all thes functions move to services package
func calculateTotalValues(transactions []CryptoTransaction) (float64, float64) {
	totalAmountInvested := 0.0
	totalCryptoAsset := 0.0

	for _, transaction := range transactions {
		totalAmountInvested += transaction.PurchasedPrice
		totalCryptoAsset += transaction.CryptoAmount
	}

	return totalAmountInvested, totalCryptoAsset
}

func printSummary(cryptoAsset price.CrytpoAsset, totalAmountInvested, totalCryptoAsset, currentPrice float64) {
	fmt.Printf("totalAmountInvested: %.2f\n", totalAmountInvested)
	fmt.Printf("total %s: %.2f\n", cryptoAsset, totalCryptoAsset)

	moneyNow := totalCryptoAsset * currentPrice
	fmt.Printf("moneyNow: %.2f\n", moneyNow)

	totalGained := calculatePercentageChange(totalAmountInvested, moneyNow)
	fmt.Printf("totalGained: %.2f%%\n", totalGained)
}

func isFirstTime(data []byte) bool {
	return strings.TrimSpace(string(data)) == ""
}

func readFile(fileRead FileReader, cryptoAsset price.CrytpoAsset) ([]byte, error) {
	filePath := "./input/" + strings.ToLower(string(cryptoAsset)) + ".txt"
	data, err := fileRead.ReadFile(filePath)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return nil, err
	}
	return data, nil
}

func writeFile(fileWriter FileWriter, cryptoAsset string) error {
	filePath := "./input/" + strings.ToLower(cryptoAsset) + ".txt"
	return fileWriter.WriteFile(filePath, []byte("false"), 0644)
}

func calculatePercentageChange(initialValue, finalValue float64) float64 {
	return ((finalValue - initialValue) / initialValue) * 100
}
