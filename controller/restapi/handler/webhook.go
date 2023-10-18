package handler

import (
	"net/http"

	"github.com/aalexanderkevin/crypto-wallet/container"
	"github.com/aalexanderkevin/crypto-wallet/helper"
	"github.com/aalexanderkevin/crypto-wallet/model"
	"github.com/aalexanderkevin/crypto-wallet/usecase"

	"github.com/blockcypher/gobcy/v2"
	"github.com/gin-gonic/gin"
)

type Webhook struct {
	appContainer *container.Container
}

func NewWebhook(appContainer *container.Container) *Webhook {
	return &Webhook{appContainer: appContainer}
}

// Crypto Wallet Webhook
// @Summary Crypto Wallet Webhook
// @Description Crypto Wallet Webhook
// @Tags Crypto Wallet
// @Accept json
// @Produce json
// @Param body body request.Webhook true " "
// @Success 200
// @Failure 401 {object} response.SendErrorResponse "When the auth token is missing or invalid"
// @Failure 422 {object} response.SendErrorResponse "When request validation failed"
// @Failure 500 {object} response.ErrorResponse "When server encountered unhandled error"
// @Security BearerAuth
// @Router /v1/public/candidate/update/:tenant-id [post]
func (w *Webhook) Transaction(c *gin.Context) {
	logger := helper.GetLogger(c).WithField("method", "Restapi.Handler.Transaction")

	// Validation
	var req gobcy.TX
	if err := c.ShouldBind(&req); err != nil {
		logger.WithError(err).Warning("bad request error")
		// response.SendErrorResponse(c, response.ErrBadRequest, "")
		return
	}

	trx := model.Transaction{}.FromModel(req)

	webhookUseCase := usecase.NewWebhook(w.appContainer)
	err := webhookUseCase.UpsertBitcoinTransaction(c, trx)
	if err != nil {
		logger.WithError(err).Warning("error UpsertBitcoinTransaction")
		c.JSON(http.StatusInternalServerError, nil)
		return
	}

	c.JSON(http.StatusOK, nil)

	// if err := req.Validate(); err != nil {
	// 	response.SendErrorResponse(c, response.ErrValidation, "")
	// 	return
	// }

	// // Get tenant from path
	// tenant := c.Param("tenant-id")

	// // Action
	// webhookUseCase := usecase.NewWebhook(w.appContainer)
	// err := webhookUseCase.CandidateUpdated(c, &tenant, req.CandidateID, req.EventDate, req.ActionType)
	// if err != nil {
	// 	response.Error(c, err, "")
	// 	return
	// }

	// response.Success(c)
}

// func extractSignature(header string) string {
// 	// You'll need to parse the "Signature" header to extract the signature.
// 	// This is a simplified example, and you may need to implement more robust parsing.
// 	const signaturePrefix = "signature=\""
// 	const signatureSuffix = "\""
// 	start := strings.Index(header, signaturePrefix)
// 	if start == -1 {
// 		return ""
// 	}
// 	start += len(signaturePrefix)
// 	end := strings.Index(header[start:], signatureSuffix)
// 	if end == -1 {
// 		return ""
// 	}
// 	return header[start : start+end]
// }

// func verifyWebhookSignature(webhookData, signature, publicKeyPEM string) (bool, error) {
// 	// Parse the PEM-encoded public key
// 	block, _ := pem.Decode([]byte(publicKeyPEM))
// 	if block == nil {
// 		return false, fmt.Errorf("Failed to decode public key")
// 	}
// 	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
// 	if err != nil {
// 		return false, err
// 	}

// 	// Compute the SHA-256 hash of the webhook data
// 	hashedData := sha256.Sum256([]byte(webhookData))

// 	// Verify the signature
// 	ecdsaPublicKey, ok := publicKey.(*ecdsa.PublicKey)
// 	if !ok {
// 		return false, fmt.Errorf("Public key is not of ECDSA type")
// 	}
// 	signatureBytes := []byte(signature)
// 	var r, s big.Int
// 	r.SetBytes(signatureBytes[:32])
// 	s.SetBytes(signatureBytes[32:])
// 	ok = ecdsa.Verify(ecdsaPublicKey, hashedData[:], &r, &s)
// 	if !ok {
// 		return false, errors.New("failed verify")
// 	}

// 	return true, nil
// }
