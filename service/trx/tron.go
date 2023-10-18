package trx

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/aalexanderkevin/crypto-wallet/config"
	"github.com/aalexanderkevin/crypto-wallet/helper"
	"github.com/aalexanderkevin/crypto-wallet/model"
	"github.com/aalexanderkevin/crypto-wallet/service"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/fbsobreira/gotron-sdk/pkg/client"
	"github.com/fbsobreira/gotron-sdk/pkg/common"
	"github.com/fbsobreira/gotron-sdk/pkg/keys"
	"github.com/fbsobreira/gotron-sdk/pkg/proto/api"
	"github.com/fbsobreira/gotron-sdk/pkg/proto/core"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/proto"
)

type TronImpl struct {
	grpcClient *client.GrpcClient
	httpClient *http.Client
	config     config.Tron
}

func NewTronImpl(config config.Config) service.Tron {
	opts := make([]grpc.DialOption, 0)
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	conn := client.NewGrpcClient(config.Tron.NetgRPCUrl)
	if err := conn.Start(opts...); err != nil {
		_ = fmt.Errorf("error connecting GRPC Client: %v", err)
	}

	conn.SetAPIKey(config.Tron.ApiKey)

	return &TronImpl{
		grpcClient: conn,
		httpClient: &http.Client{
			Timeout: 5 * time.Second},
		config: config.Tron,
	}
}

func (t *TronImpl) Close() {
	t.grpcClient.Stop()
}

func (t *TronImpl) GetWallet(ctx context.Context, seedPhrase *string) *model.TrxHdWallet {
	privateKey, public := keys.FromMnemonicSeedAndPassphrase(*seedPhrase, "", 0)
	privateKeyECDSA := privateKey.ToECDSA()
	publicKeyECDSA := public.ToECDSA()

	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	address = "41" + address[2:]

	trxAddress := helper.ToTrxAddress(address)

	return &model.TrxHdWallet{
		PrivateKey: privateKeyECDSA,
		PublicKey:  publicKeyECDSA,
		Address:    trxAddress,
	}
}

func (t *TronImpl) SendTx(ctx context.Context, txOpts *model.TxOpts, wallet *model.TrxHdWallet) (transaction *api.TransactionExtention, err error) {
	logger := helper.GetLogger(ctx).WithField("method", "Service.Tron.SendTx")

	tx, err := t.grpcClient.Transfer(*wallet.Address, *txOpts.To, txOpts.Amount.Int64())
	if err != nil {
		logger.WithError(err).Warn("Failed to tranfer")
		return nil, err
	}

	rawData, err := proto.Marshal(tx.Transaction.GetRawData())
	if err != nil {
		logger.WithError(err).Warn("Failed to parse raw data")
		return nil, err
	}

	h256h := sha256.New()
	h256h.Write(rawData)
	hash := h256h.Sum(nil)

	signature, err := crypto.Sign(hash, wallet.PrivateKey)
	if err != nil {
		logger.WithError(err).Warn("Failed to sign transaction")
		return nil, err
	}

	tx.Transaction.Signature = append(tx.Transaction.Signature, signature)

	result, err := t.grpcClient.Broadcast(tx.Transaction)
	if err != nil {
		logger.WithError(err).Warn("Failed to broadcast message")
		return nil, err
	}

	if result.Code != api.Return_SUCCESS {
		err := errors.New(result.String())
		logger.WithError(err).Warn("broadcast transaction return not success")
		return nil, err
	}

	return tx, nil
}

func (t *TronImpl) GetBalance(ctx context.Context, address *string) (balance *int64, err error) {
	logger := helper.GetLogger(ctx).WithField("method", "Service.Tron.GetBalance")

	accDetailed, err := t.grpcClient.GetAccountDetailed(*address)
	if err != nil {
		logger.WithError(err).Warn("error get account detailed")
		return nil, err
	}

	return helper.Pointer(accDetailed.Balance), nil
}

func (t *TronImpl) GetCurrentBlock(ctx context.Context) (*int64, error) {
	logger := helper.GetLogger(ctx).WithField("method", "Service.Tron.GetCurrentBlock")

	block, err := t.grpcClient.GetNowBlock()
	if err != nil {
		logger.WithError(err).Warn("Failed GetNowBlock")
		return nil, err
	}

	return helper.Pointer(block.BlockHeader.RawData.Number), nil
}

func (t *TronImpl) GetTx(ctx context.Context, txhash string) (*core.TransactionInfo, error) {
	logger := helper.GetLogger(ctx).WithField("method", "Service.Tron.GetTx")

	tx, err := t.grpcClient.GetTransactionInfoByID(txhash)
	if err != nil {
		logger.WithError(err).Warn("Failed get tx")
		return nil, err
	}

	return tx, nil
}

func (t *TronImpl) GetUnconfirmedTxAddress(ctx context.Context, address *string) (*service.GetTransactionResponse, error) {
	logger := helper.GetLogger(ctx).WithField("method", "Service.Tron.GetUnconfirmedTx")

	URL := fmt.Sprintf("%s/v1/accounts/%s/transactions?only_confirmed=false&only_unconfirmed=true", t.config.NetUrl, *address)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, URL, nil)
	if err != nil {
		logger.WithError(err).Warn("Failed create request")
		return nil, err
	}

	req.Header.Add("TRON-PRO-API-KEY", t.config.ApiKey)
	req.Header.Add("Content-Type", "application/json")

	resp, err := t.httpClient.Do(req)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		logger.
			WithField("header", resp.Header).
			WithField("body", string(b)).
			WithField("status", resp.Status).
			Error("Status is not OK", resp.Status)
		return nil, fmt.Errorf("status not OK: %s, body: %s", resp.Status, b)
	}

	var result service.GetTransactionResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (t *TronImpl) GetConfirmedTxAddress(ctx context.Context, address *string) (*service.GetTransactionResponse, error) {
	logger := helper.GetLogger(ctx).WithField("method", "Service.Tron.GetUnconfirmedTx")

	URL := fmt.Sprintf("%s/v1/accounts/%s/transactions?only_confirmed=true&only_unconfirmed=false", t.config.NetUrl, *address)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, URL, nil)
	if err != nil {
		logger.WithError(err).Warn("Failed create request")
		return nil, err
	}

	req.Header.Add("TRON-PRO-API-KEY", t.config.ApiKey)
	req.Header.Add("Content-Type", "application/json")

	resp, err := t.httpClient.Do(req)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		logger.
			WithField("header", resp.Header).
			WithField("body", string(b)).
			WithField("status", resp.Status).
			Error("Status is not OK", resp.Status)
		return nil, fmt.Errorf("status not OK: %s, body: %s", resp.Status, b)
	}

	var result service.GetTransactionResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (t *TronImpl) CheckAddress(address string) error {
	if _, err := common.DecodeCheck(address); err != nil {
		return err
	}
	return nil
}

func (t *TronImpl) GetTxByAccountAddress(ctx context.Context, address *string, filter *service.GetTxByAccountAddressFilter) (*service.GetTransactionResponse, error) {
	logger := helper.GetLogger(ctx).WithField("method", "Service.Tron.GetUnconfirmedTx")

	URL := fmt.Sprintf("%s/v1/accounts/%s/transactions", t.config.NetUrl, *address)

	if filter != nil {
		queryParams := []string{}
		if filter.OnlyUnconfirmed != nil {
			queryParams = append(queryParams, fmt.Sprintf("only_unconfirmed=%t", *filter.OnlyUnconfirmed))
		}

		if filter.OnlyConfirmed != nil {
			queryParams = append(queryParams, fmt.Sprintf("only_confirmed=%t", *filter.OnlyConfirmed))
		}

		if filter.OnlyTo != nil {
			queryParams = append(queryParams, fmt.Sprintf("only_to=%t", *filter.OnlyTo))
		}

		if filter.OnlyFrom != nil {
			queryParams = append(queryParams, fmt.Sprintf("only_from=%t", *filter.OnlyFrom))
		}

		if filter.Limit != nil {
			queryParams = append(queryParams, fmt.Sprintf("limit=%d", *filter.Limit))
		}

		if filter.Fingerprint != nil {
			queryParams = append(queryParams, fmt.Sprintf("fingerprint=%s", *filter.Fingerprint))
		}

		if filter.OrderBy != nil {
			queryParams = append(queryParams, fmt.Sprintf("order_by=%s", *filter.OrderBy))
		}

		if filter.MinTimestampMs != nil {
			queryParams = append(queryParams, fmt.Sprintf("min_timestamp=%d", *filter.MinTimestampMs))
		}

		if filter.MaxTimestampMs != nil {
			queryParams = append(queryParams, fmt.Sprintf("max_timestamp=%d", *filter.MaxTimestampMs))
		}

		if len(queryParams) > 0 {
			URL += "?" + strings.Join(queryParams, "&")
		}
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, URL, nil)
	if err != nil {
		logger.WithError(err).Warn("Failed create request")
		return nil, err
	}

	req.Header.Add("TRON-PRO-API-KEY", t.config.ApiKey)
	req.Header.Add("Content-Type", "application/json")

	resp, err := t.httpClient.Do(req)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		logger.
			WithField("header", resp.Header).
			WithField("body", string(b)).
			WithField("status", resp.Status).
			Error("Status is not OK", resp.Status)
		return nil, fmt.Errorf("status not OK: %s, body: %s", resp.Status, b)
	}

	var result service.GetTransactionResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
