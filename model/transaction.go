package model

import (
	"math/big"
	"time"

	"github.com/aalexanderkevin/crypto-wallet/helper"
	validation "github.com/go-ozzo/ozzo-validation/v4"

	"github.com/blockcypher/gobcy/v2"
)

type Transaction struct {
	Id              *string    `json:"id"`
	SenderAddress   []string   `json:"sender_address"`
	ReceiverAddress []string   `json:"receiver_address"`
	Amount          *int64     `json:"amount"`
	Fee             *int64     `json:"fee"`
	Block           *int64     `json:"block"`
	Confirmation    *int64     `json:"confirmation"`
	Status          *string    `json:"status"`
	ReceivedAt      *time.Time `json:"received_at"`
	CompletedAt     *time.Time `json:"completed_at"`
}

func (t Transaction) FromModel(data gobcy.TX) *Transaction {
	inputs := []string{}
	for _, input := range data.Inputs {
		inputs = append(inputs, input.Addresses...)
	}
	ouputs := []string{}
	for _, output := range data.Outputs {
		ouputs = append(ouputs, output.Addresses...)
	}

	var completedAt *time.Time
	status := "pending"
	if data.Confirmations == 6 {
		status = "success"
		completedAt = helper.Pointer(time.Now())
	}

	return &Transaction{
		Id:              &data.Hash,
		SenderAddress:   inputs,
		ReceiverAddress: ouputs,
		Amount:          helper.Pointer(data.Total.Int64()),
		Fee:             helper.Pointer(data.Fees.Int64()),
		Confirmation:    helper.Pointer(int64(data.Confirmations)),
		Status:          &status,
		ReceivedAt:      &data.Received,
		CompletedAt:     completedAt,
	}
}

type SendToken struct {
	Email           *string `json:"email"`
	ReceiverAddress *string `json:"receiver_address"`
	Amount          *int64  `json:"amount"`
	Token           *string `json:"token"`
}

func (s SendToken) Validate() error {
	return validation.ValidateStruct(
		&s,
		validation.Field(&s.Email, validation.Required),
		validation.Field(&s.ReceiverAddress, validation.Required),
		validation.Field(&s.Amount, validation.Required),
		validation.Field(&s.Token, validation.Required),
	)
}

type TxOpts struct {
	To          *string
	Amount      *big.Int
	AmountInt64 *int64
}
