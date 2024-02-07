package main

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"

	"crypto-alerts/pkg/services"
	xmpp "crypto-alerts/pkg/xmpp"
)

type CryptoTransaction struct {
	PurchasedDate  string
	PurchasedPrice float64
	CryptoAmount   float64
}

type CryptoAsset struct {
	Name     string
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
	return ioutil.ReadFile(filename)
}

func (r *DefaultFileReaderWriter) WriteFile(filename string, data []byte, perm fs.FileMode) error {
	return ioutil.WriteFile(filename, data, perm)
}

func main() {
	assets := []CryptoAsset{
		{Name: "Bitcoin", Boundary: 200.0},
		{Name: "Ethereum", Boundary: 100.0},
		{Name: "Cardano", Boundary: 100.0},
		{Name: "Polkadot", Boundary: 50.0},
	}

	jsonData, err := ioutil.ReadFile("./json/data.json")
	if err != nil {
		fmt.Println("Error reading JSON file:", err)
		os.Exit(1)
	}

	var transactions map[string][]CryptoTransaction
	if err := json.Unmarshal(jsonData, &transactions); err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		os.Exit(1)
	}

	for _, asset := range assets {
		transactions, ok := transactions[asset.Name]
		if !ok {
			continue
		}
		totalAmountInvested, totalCryptoAsset := calculateTotalValues(transactions)
		client := http.Client{}
		currentPrice, err := services.GetCurrentPrice(&client, asset.Name)
		if err != nil {
			// fmt.Printf("Error fetching current %s price: %v\n", asset.Name, err)
			sendXmppErrorGettingPrice("Error fetching current " + asset.Name + " price")
			continue
		}
		if currentPrice == 0 {
			sendXmppErrorGettingPrice("Error fetching current " + asset.Name + " price")
			continue
		}
		fmt.Printf("-------- %s  --------\n", asset.Name)
		printSummary(asset.Name, totalAmountInvested, totalCryptoAsset, currentPrice, asset.Boundary)
		var fileReader FileReader = &DefaultFileReaderWriter{}
		readData, err := readFile(fileReader, asset.Name)
		if err != nil {
			fmt.Println("Error reading file: ", err)
			return
		}
		if isFirstTime(readData) {
			moneyNow := totalCryptoAsset * currentPrice
			sendXmppGainMsg(asset.Name, asset.Boundary, calculatePercentageChange(totalAmountInvested, moneyNow))
			var fileWriter FileWriter = &DefaultFileReaderWriter{}
			writeFile(fileWriter, asset.Name)
		}
		fmt.Printf("-------- END %s  --------\n", asset.Name)
	}
}

func calculateTotalValues(transactions []CryptoTransaction) (float64, float64) {
	totalAmountInvested := 0.0
	totalCryptoAsset := 0.0

	for _, transaction := range transactions {
		totalAmountInvested += transaction.PurchasedPrice
		totalCryptoAsset += transaction.CryptoAmount
	}

	return totalAmountInvested, totalCryptoAsset
}

func printSummary(cryptoAsset string, totalAmountInvested, totalCryptoAsset, currentPrice, boundary float64) {
	fmt.Printf("totalAmountInvested: %.2f\n", totalAmountInvested)
	fmt.Printf("total %s: %.2f\n", cryptoAsset, totalCryptoAsset)

	moneyNow := totalCryptoAsset * currentPrice
	fmt.Printf("moneyNow: %.2f\n", moneyNow)

	totalGained := calculatePercentageChange(totalAmountInvested, moneyNow)
	fmt.Printf("totalGained: %.2f%%\n", totalGained)
}

func isFirstTime(data []byte) bool {
	strValue := strings.TrimSpace(string(data))
	boolValue, err := strconv.ParseBool(strValue)
	if err != nil {
		fmt.Println("Error parsing into bool: ", err)
		return true
	}
	return boolValue
}

func readFile(fileRead FileReader, cryptoAsset string) ([]byte, error) {
	filePath := "./input/" + strings.ToLower(cryptoAsset) + ".txt"
	data, err := fileRead.ReadFile(filePath)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return nil, err
	}
	return data, nil
}

func writeFile(fileWriter FileWriter, cryptoAsset string) {
	filePath := "./input/" + strings.ToLower(cryptoAsset) + ".txt"
	err := fileWriter.WriteFile(filePath, []byte("false"), 0644)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}
}

func sendXmppErrorGettingPrice(message string) {
	xmpp.SendXMPP(message)
}

func sendXmppGainMsg(currency string, boundary float64, totalGained float64) {
	if totalGained >= boundary {
		strVal := strconv.FormatFloat(totalGained, 'f', 2, 64)
		xmpp.SendXMPP("You have gained in " + currency + ": " + strVal + "%")
	}
}

func calculatePercentageChange(initialValue, finalValue float64) float64 {
	return ((finalValue - initialValue) / initialValue) * 100
}
