package btc

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"net/url"
	"regexp"
	"strconv"

	"github.com/aalexanderkevin/crypto-wallet/config"
	"github.com/aalexanderkevin/crypto-wallet/helper"
	"github.com/aalexanderkevin/crypto-wallet/model"
	"github.com/aalexanderkevin/crypto-wallet/service"

	"github.com/blockcypher/gobcy/v2"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/tyler-smith/go-bip39"
)

type BitcoinImpl struct {
	client *gobcy.API
	config config.Bitcoin
}

func NewBitcoinImpl(config config.Config) service.Bitcoin {
	//explicitly
	client := &gobcy.API{
		Token: config.Bitcoin.Token,
		Coin:  "btc",                //options: "btc","bcy","ltc","doge","eth"
		Chain: config.Bitcoin.Chain, //depending on coin: "main","test3","test"
	}

	return &BitcoinImpl{
		client: client,
		config: config.Bitcoin,
	}
}

func (b *BitcoinImpl) CheckAddress(address *string) bool {
	// Regular expression for a valid Bitcoin address
	pattern := "^(tb1|[mn2])[a-km-zA-HJ-NP-Z0-9]{25,39}$"
	chainConfig := &chaincfg.TestNet3Params
	if b.client.Chain == "main" {
		chainConfig = &chaincfg.MainNetParams
		pattern = "^bc1[ac-hj-np-z02-9]{25,39}$"
	}

	matched, _ := regexp.MatchString(pattern, *address)
	if !matched {
		return false
	}

	// Decode the address and perform checksum verification
	decoded, err := btcutil.DecodeAddress(*address, chainConfig)
	if err != nil {
		return false
	}

	return decoded.IsForNet(chainConfig)
}

func (b *BitcoinImpl) GetWallet(ctx context.Context, seedPhrase *string) (*model.BtcHdWallet, error) {
	logger := helper.GetLogger(ctx).WithField("method", "Service.Bitcoin.GetWallet")

	// Generate a seed from the mnemonic
	seed := bip39.NewSeed(*seedPhrase, "")

	// Create a master extended key from the seed
	masterKey, err := hdkeychain.NewMaster(seed, &chaincfg.TestNet3Params)
	if err != nil {
		logger.WithError(err).Warn("Failed create masker key")
		return nil, err
	}

	// Get the public key and Bitcoin address from the child key
	publicKey, err := masterKey.ECPubKey()
	if err != nil {
		logger.WithError(err).Warn("Failed get public key")
		return nil, err
	}

	// Get the private key in Wallet Import Format (WIF)
	privateKey, err := masterKey.ECPrivKey()
	if err != nil {
		logger.WithError(err).Warn("Failed get private key")
		return nil, err
	}

	// Convert the private key to WIF using the btcutil library
	wif, err := btcutil.NewWIF(privateKey, &chaincfg.TestNet3Params, true)
	if err != nil {
		logger.WithError(err).Warn("Failed convert private key")
		return nil, err
	}

	address, err := btcutil.NewAddressPubKey(publicKey.SerializeCompressed(), &chaincfg.TestNet3Params)
	if err != nil {
		logger.WithError(err).Warn("Failed generate new address")
		return nil, err
	}

	res := &model.BtcHdWallet{
		PublicKey:  publicKey,
		Address:    address,
		PrivateKey: privateKey,
		Wif:        wif,
	}
	return res, nil
}

func (b *BitcoinImpl) GetBalance(ctx context.Context, address string) (*big.Int, error) {
	logger := helper.GetLogger(ctx).WithField("method", "Service.Bitcoin.GetWallet")

	addr, err := b.client.GetAddr(address, nil)
	if err != nil {
		logger.WithError(err).Warn("Failed GetAddrBal")
		return nil, err
	}

	return &addr.Balance, nil
}

func (b *BitcoinImpl) SendTx(ctx context.Context, wallet *model.BtcHdWallet, txOpts *model.TxOpts) (*model.Transaction, error) {
	logger := helper.GetLogger(ctx).WithField("method", "Service.Bitcoin.GetWallet")

	tx := gobcy.TempNewTX(wallet.Address.EncodeAddress(), *txOpts.To, *txOpts.Amount)
	tx.Preference = "low"
	skels, err := b.client.NewTX(tx, false)
	if err != nil {
		logger.WithError(err).Warn("Failed create tx")
		return nil, err
	}

	prikHexs := []string{}
	for i := 0; i < len(skels.ToSign); i++ {
		prikHexs = append(prikHexs, fmt.Sprintf("%x", wallet.PrivateKey.Serialize()))
	}

	// Sign it locally
	err = skels.Sign(prikHexs)
	if err != nil {
		logger.WithError(err).Warn("Failed sign transaction")
		return nil, err
	}

	// Send TXSkeleton
	skels, err = b.client.SendTX(skels)
	if err != nil {
		logger.WithError(err).Warn("Failed send tx")
		return nil, err
	}

	return model.Transaction{}.FromModel(skels.Trans), nil
}

func (b *BitcoinImpl) GetTx(ctx context.Context, txhash string) (*gobcy.TX, error) {
	logger := helper.GetLogger(ctx).WithField("method", "Service.Bitcoin.GetTx")

	tx, err := b.client.GetTX(txhash, nil)
	if err != nil {
		logger.WithError(err).Warn("Failed get tx")
		return nil, err
	}

	return &tx, nil
}

func (b *BitcoinImpl) CreateWebhookConfirmedTx(ctx context.Context, address *string) (*gobcy.Hook, error) {
	logger := helper.GetLogger(ctx).WithField("method", "Service.Bitcoin.CreateWebhookConfirmedTx")
	hooks, err := b.client.ListHooks()
	for _, h := range hooks {
		err = b.client.DeleteHook(h.ID)
		fmt.Println(err)
	}

	hook, err := b.client.CreateHook(gobcy.Hook{
		Event: "tx-confirmation",
		// SignKey:       "preset",
		Address:       *address,
		URL:           b.config.WebhookURL + "/webhook/transaction",
		Confirmations: b.config.MinimalConfirmation,
	})
	if err != nil {
		logger.WithError(err).Warn("Failed create hook")
		return nil, err
	}

	return &hook, nil
}

func (b *BitcoinImpl) DeleteWebhook(ctx context.Context, id *string) error {
	logger := helper.GetLogger(ctx).WithField("method", "Service.Bitcoin.CreateWebhookConfirmedTx")

	if err := b.client.DeleteHook(*id); err != nil {
		logger.WithError(err).Warn("Failed delete hook")
		return err
	}

	return nil
}

type Hook struct {
	ID            string  `json:"id,omitempty"`
	Event         string  `json:"event"`
	SignKey       string  `json:"signkey"`
	Hash          string  `json:"hash,omitempty"`
	WalletName    string  `json:"wallet_name,omitempty"`
	Address       string  `json:"address,omitempty"`
	Confirmations int     `json:"confirmations,omitempty"`
	Confidence    float32 `json:"confidence,omitempty"`
	Script        string  `json:"script,omitempty"`
	URL           string  `json:"url,omitempty"`
	CallbackErrs  int     `json:"callback_errors,omitempty"`
}

func (b *BitcoinImpl) createHook(hook Hook) (result gobcy.Hook, err error) {
	u, err := b.buildURL("/hooks", nil)
	if err != nil {
		return
	}
	err = postResponse(u, &hook, &result)
	return
}

func (b *BitcoinImpl) buildURL(u string, params map[string]string) (target *url.URL, err error) {
	target, err = url.Parse("https://api.blockcypher.com/v1/" + b.client.Coin + "/" + b.client.Chain + u)
	if err != nil {
		return
	}
	values := target.Query()
	//Set parameters
	for k, v := range params {
		values.Set(k, v)
	}
	//add token to url, if present
	if b.client.Token != "" {
		values.Set("token", b.client.Token)
	}
	target.RawQuery = values.Encode()
	return
}

func postResponse(target *url.URL, encTarget interface{}, decTarget interface{}) (err error) {
	var data bytes.Buffer
	enc := json.NewEncoder(&data)
	if err = enc.Encode(encTarget); err != nil {
		return
	}
	resp, err := http.Post(target.String(), "application/json", &data)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		err = respErrorMaker(resp.StatusCode, resp.Body)
		return
	}
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(decTarget)
	return
}

func respErrorMaker(statusCode int, body io.Reader) (err error) {
	status := "HTTP " + strconv.Itoa(statusCode) + " " + http.StatusText(statusCode)
	if statusCode == 429 {
		err = errors.New(status)
		return
	}
	type errorJSON struct {
		Err    string `json:"error"`
		Errors []struct {
			Err string `json:"error"`
		} `json:"errors"`
	}
	var msg errorJSON
	dec := json.NewDecoder(body)
	err = dec.Decode(&msg)
	if err != nil {
		return err
	}
	var errtxt string
	errtxt += msg.Err
	for i, v := range msg.Errors {
		if i == len(msg.Errors)-1 {
			errtxt += v.Err
		} else {
			errtxt += v.Err + ", "
		}
	}
	if errtxt == "" {
		err = errors.New(status)
	} else {
		err = errors.New(status + ", Message(s): " + errtxt)
	}
	return
}
