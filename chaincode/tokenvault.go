package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing a tokenized card vault
type SmartContract struct {
	contractapi.Contract
}

// CardToken represents a tokenized payment card stored in the vault
type CardToken struct {
	TokenID   string `json:"TokenID"`
	CardRef   string `json:"CardRef"`
	Holder    string `json:"Holder"`
	Status    string `json:"Status"`
	IssuedAt  string `json:"IssuedAt"`
	UpdatedAt string `json:"UpdatedAt"`
}

// HistoryEntry represents one revision of a token from the ledger history
type HistoryEntry struct {
	TxID      string     `json:"TxID"`
	Value     *CardToken `json:"Value"`
	Timestamp string     `json:"Timestamp"`
	IsDelete  bool       `json:"IsDelete"`
}

const (
	statusActive    = "ACTIVE"
	statusSuspended = "SUSPENDED"
)

// IssueToken creates a new card token in the vault with status "ACTIVE".
func (s *SmartContract) IssueToken(ctx contractapi.TransactionContextInterface, tokenID string, cardRef string, holder string) error {
	existing, err := ctx.GetStub().GetState(tokenID)
	if err != nil {
		return fmt.Errorf("failed to read from world state: %v", err)
	}
	if existing != nil {
		return fmt.Errorf("token %s already exists", tokenID)
	}

	now := txTimestamp(ctx)
	token := CardToken{
		TokenID:   tokenID,
		CardRef:   cardRef,
		Holder:    holder,
		Status:    statusActive,
		IssuedAt:  now,
		UpdatedAt: now,
	}

	return putToken(ctx, &token)
}

// GetToken returns the card token identified by tokenID.
func (s *SmartContract) GetToken(ctx contractapi.TransactionContextInterface, tokenID string) (*CardToken, error) {
	data, err := ctx.GetStub().GetState(tokenID)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if data == nil {
		return nil, fmt.Errorf("token %s does not exist", tokenID)
	}

	var token CardToken
	if err := json.Unmarshal(data, &token); err != nil {
		return nil, err
	}
	return &token, nil
}

// SuspendToken transitions a token from "ACTIVE" to "SUSPENDED".
func (s *SmartContract) SuspendToken(ctx contractapi.TransactionContextInterface, tokenID string) error {
	token, err := s.GetToken(ctx, tokenID)
	if err != nil {
		return err
	}
	if token.Status != statusActive {
		return fmt.Errorf("token %s is not ACTIVE (current status: %s)", tokenID, token.Status)
	}

	token.Status = statusSuspended
	token.UpdatedAt = txTimestamp(ctx)
	return putToken(ctx, token)
}

// ResumeToken transitions a token from "SUSPENDED" back to "ACTIVE".
func (s *SmartContract) ResumeToken(ctx contractapi.TransactionContextInterface, tokenID string) error {
	token, err := s.GetToken(ctx, tokenID)
	if err != nil {
		return err
	}
	if token.Status != statusSuspended {
		return fmt.Errorf("token %s is not SUSPENDED (current status: %s)", tokenID, token.Status)
	}

	token.Status = statusActive
	token.UpdatedAt = txTimestamp(ctx)
	return putToken(ctx, token)
}

// DeleteToken permanently removes a token from the vault.
func (s *SmartContract) DeleteToken(ctx contractapi.TransactionContextInterface, tokenID string) error {
	data, err := ctx.GetStub().GetState(tokenID)
	if err != nil {
		return fmt.Errorf("failed to read from world state: %v", err)
	}
	if data == nil {
		return fmt.Errorf("token %s does not exist", tokenID)
	}
	return ctx.GetStub().DelState(tokenID)
}

// GetTokenHistory returns the full revision history of a token, newest first.
func (s *SmartContract) GetTokenHistory(ctx contractapi.TransactionContextInterface, tokenID string) ([]HistoryEntry, error) {
	resultsIterator, err := ctx.GetStub().GetHistoryForKey(tokenID)
	if err != nil {
		return nil, fmt.Errorf("failed to get history for %s: %v", tokenID, err)
	}
	defer resultsIterator.Close()

	var history []HistoryEntry
	for resultsIterator.HasNext() {
		modification, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		entry := HistoryEntry{
			TxID:      modification.TxId,
			Timestamp: time.Unix(modification.Timestamp.Seconds, int64(modification.Timestamp.Nanos)).UTC().Format(time.RFC3339),
			IsDelete:  modification.IsDelete,
		}
		if !modification.IsDelete {
			var token CardToken
			if err := json.Unmarshal(modification.Value, &token); err != nil {
				return nil, err
			}
			entry.Value = &token
		}
		history = append(history, entry)
	}
	return history, nil
}

// --- helpers ---

func putToken(ctx contractapi.TransactionContextInterface, token *CardToken) error {
	bytes, err := json.Marshal(token)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState(token.TokenID, bytes)
}

// txTimestamp returns the deterministic transaction timestamp in RFC3339.
func txTimestamp(ctx contractapi.TransactionContextInterface) string {
	ts, err := ctx.GetStub().GetTxTimestamp()
	if err != nil {
		return ""
	}
	return time.Unix(ts.Seconds, int64(ts.Nanos)).UTC().Format(time.RFC3339)
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
