package main

import (
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing a tokenized card vault
type SmartContract struct {
	contractapi.Contract
}

// CardToken represents a tokenized payment card stored in the vault
type CardToken struct {
	TokenID   string `json:"TokenID"`   // unique token identifier, e.g. "tok1"
	CardRef   string `json:"CardRef"`   // masked card reference, e.g. "XXXX-XXXX-XXXX-1234"
	Holder    string `json:"Holder"`    // name of the card holder
	Status    string `json:"Status"`    // ACTIVE | SUSPENDED
	IssuedAt  string `json:"IssuedAt"`  // RFC3339 timestamp when the token was issued
	UpdatedAt string `json:"UpdatedAt"` // RFC3339 timestamp of the last status change
}

// HistoryEntry represents one revision of a token from the ledger history
type HistoryEntry struct {
	TxID      string     `json:"TxID"`
	Value     *CardToken `json:"Value"`
	Timestamp string     `json:"Timestamp"`
	IsDelete  bool       `json:"IsDelete"`
}

// IssueToken creates a new card token in the vault with status "ACTIVE".
// It must fail if a token with the same tokenID already exists.
func (s *SmartContract) IssueToken(ctx contractapi.TransactionContextInterface, tokenID string, cardRef string, holder string) error {

	return nil
}

// GetToken returns the card token identified by tokenID.
// It must fail if the token does not exist.
func (s *SmartContract) GetToken(ctx contractapi.TransactionContextInterface, tokenID string) (*CardToken, error) {

	return nil, nil
}

// SuspendToken transitions a token from "ACTIVE" to "SUSPENDED".
// It must fail if the token does not exist or is not currently ACTIVE.
func (s *SmartContract) SuspendToken(ctx contractapi.TransactionContextInterface, tokenID string) error {

	return nil
}

// ResumeToken transitions a token from "SUSPENDED" back to "ACTIVE".
// It must fail if the token does not exist or is not currently SUSPENDED.
func (s *SmartContract) ResumeToken(ctx contractapi.TransactionContextInterface, tokenID string) error {

	return nil
}

// DeleteToken permanently removes a token from the vault.
// It must fail if the token does not exist.
func (s *SmartContract) DeleteToken(ctx contractapi.TransactionContextInterface, tokenID string) error {

	return nil
}

// GetTokenHistory returns the full revision history of a token, newest first.
func (s *SmartContract) GetTokenHistory(ctx contractapi.TransactionContextInterface, tokenID string) ([]HistoryEntry, error) {

	return nil, nil
}

func main() {
	chaincode, err := contractapi.NewChaincode(&SmartContract{})
	if err != nil {
		panic("Error creating tokenvault chaincode: " + err.Error())
	}

	if err := chaincode.Start(); err != nil {
		panic("Error starting tokenvault chaincode: " + err.Error())
	}
}
