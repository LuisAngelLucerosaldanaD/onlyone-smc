package credentials

import (
	"context"
	"crypto/rsa"
	"github.com/form3tech-oss/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	grpcMetadata "google.golang.org/grpc/metadata"
	"io/ioutil"
	"net/http"
	"onlyone_smc/internal/env"
	"onlyone_smc/internal/grpc/accounting_proto"
	"onlyone_smc/internal/grpc/transactions_proto"
	"onlyone_smc/internal/grpc/wallet_proto"
	"onlyone_smc/internal/helpers"
	"onlyone_smc/internal/logger"
	"onlyone_smc/internal/msg"
	"strconv"
	"time"
)

var (
	signKey    *rsa.PrivateKey
	privateKey string
)

type handlerCredentials struct {
	DB   *sqlx.DB
	TxID string
}

func init() {
	c := env.NewConfiguration()
	privateKey = c.App.RSAPrivateKey
	signBytes, err := ioutil.ReadFile(privateKey)
	if err != nil {
		logger.Error.Printf("leyendo el archivo privado de firma: %s", err)
	}

	signKey, err = jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	if err != nil {
		logger.Error.Printf("realizando el parse en auth RSA private: %s", err)
	}
}

// create credential godoc
// @Summary OnlyOne Smart Contract
// @Description Create credential
// @Accept  json
// @Produce  json
// @Success 200 {object} requestCreateTransaction
// @Success 202 {object} responseCreateCredential
// @Router /api/v1/credentials/create [post]
func (h *handlerCredentials) createCredential(c *fiber.Ctx) error {
	res := responseCreateCredential{Error: true}
	m := requestCreateTransaction{}
	e := env.NewConfiguration()
	u := helpers.GetUserContext(c)
	err := c.BodyParser(&m)
	if err != nil {
		logger.Error.Printf("couldn't bind model create wallets: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(1, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	connTrx, err := grpc.Dial(e.TransactionsService.Port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Error.Printf("error conectando con el servicio trx de blockchain: %s", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(22, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}
	defer connTrx.Close()

	connAuth, err := grpc.Dial(e.AuthService.Port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Error.Printf("error conectando con el servicio trx de blockchain: %s", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(22, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}
	defer connAuth.Close()

	clientTrx := transactions_proto.NewTransactionsServicesClient(connTrx)
	clientWallet := wallet_proto.NewWalletServicesWalletClient(connAuth)
	clientAccount := accounting_proto.NewAccountingServicesAccountingClient(connAuth)

	bearer := c.Get("Authorization")
	tkn := bearer[7:]

	ctx := grpcMetadata.AppendToOutgoingContext(context.Background(), "authorization", tkn)

	if m.To == "" && m.Data.IdentityNumber == "" {
		logger.Error.Printf("La wallet de destino es requerido")
		res.Code, res.Type, res.Msg = msg.GetByCode(22, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	if m.To == "" {
		walletToByIN, err := clientWallet.GetWalletByIdentityNumber(ctx, &wallet_proto.RqGetByIdentityNumber{IdentityNumber: m.Data.IdentityNumber})
		if err != nil {
			logger.Error.Printf("couldn't get wallet by identity number: %v", err)
			res.Code, res.Type, res.Msg = msg.GetByCode(70, h.DB, h.TxID)
			return c.Status(http.StatusAccepted).JSON(res)
		}

		if walletToByIN == nil {
			logger.Error.Printf("couldn't get wallet by identity number")
			res.Code, res.Type, res.Msg = msg.GetByCode(70, h.DB, h.TxID)
			return c.Status(http.StatusAccepted).JSON(res)
		}

		if walletToByIN.Error {
			logger.Error.Printf(walletToByIN.Msg)
			res.Code, res.Type, res.Msg = msg.GetByCode(70, h.DB, h.TxID)
			return c.Status(http.StatusAccepted).JSON(res)
		}

		if walletToByIN == nil {
			wallet, err := clientWallet.CreateWalletBySystem(ctx, &wallet_proto.RqCreateWalletBySystem{})
			if err != nil {
				logger.Error.Printf("couldn't create wallet: %v", err)
				res.Code, res.Type, res.Msg = msg.GetByCode(70, h.DB, h.TxID)
				return c.Status(http.StatusAccepted).JSON(res)
			}

			if wallet == nil {
				logger.Error.Printf("couldn't create wallet")
				res.Code, res.Type, res.Msg = msg.GetByCode(70, h.DB, h.TxID)
				return c.Status(http.StatusAccepted).JSON(res)
			}

			if wallet.Error {
				logger.Error.Printf(wallet.Msg)
				res.Code, res.Type, res.Msg = msg.GetByCode(70, h.DB, h.TxID)
				return c.Status(http.StatusAccepted).JSON(res)
			}

			m.To = wallet.Data.Id

			resAccountTo, err := clientAccount.CreateAccounting(ctx, &accounting_proto.RequestCreateAccounting{
				Id:       uuid.New().String(),
				IdWallet: m.To,
				Amount:   0,
				IdUser:   u.ID,
			})
			if err != nil {
				logger.Error.Printf("couldn't create accounting to wallet: %v", err)
				res.Code, res.Type, res.Msg = msg.GetByCode(70, h.DB, h.TxID)
				return c.Status(http.StatusAccepted).JSON(res)
			}

			if resAccountTo == nil {
				logger.Error.Printf("couldn't create accounting to wallet: %v", err)
				res.Code, res.Type, res.Msg = msg.GetByCode(70, h.DB, h.TxID)
				return c.Status(http.StatusAccepted).JSON(res)
			}

			if resAccountTo.Error {
				logger.Error.Printf(resAccountTo.Msg)
				res.Code, res.Type, res.Msg = msg.GetByCode(70, h.DB, h.TxID)
				return c.Status(http.StatusAccepted).JSON(res)
			}
		} else {
			m.To = walletToByIN.Data.Id
		}
	}

	var files []*transactions_proto.File
	for _, file := range m.Data.Files {
		files = append(files, &transactions_proto.File{
			IdFile:     int32(file.FileID),
			Name:       file.Name,
			FileEncode: file.FileEncode,
		})
	}

	var identifiers []*transactions_proto.Identifier
	for _, identifier := range m.Data.Identifiers {
		var attributes []*transactions_proto.Attribute
		for _, attribute := range identifier.Attributes {
			attributes = append(attributes, &transactions_proto.Attribute{
				Id:    int32(attribute.Id),
				Name:  attribute.Name,
				Value: attribute.Value,
			})
		}
		identifiers = append(identifiers, &transactions_proto.Identifier{
			Name:       identifier.Name,
			Attributes: attributes,
		})
	}

	resCreateTrx, err := clientTrx.CreateTransaction(ctx, &transactions_proto.RequestCreateTransaction{
		From:   m.From,
		To:     m.To,
		Amount: m.Amount,
		TypeId: int32(m.TypeId),
		Data: &transactions_proto.Data{
			Category:       m.Data.Category,
			IdentityNumber: m.Data.IdentityNumber,
			Files:          files,
			Name:           m.Data.Name,
			Description:    m.Data.Description,
			Identifiers:    identifiers,
		},
	})

	if err != nil {
		logger.Error.Printf("No se pudo crear el usuario, error: %s", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(22, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	if resCreateTrx == nil {
		logger.Error.Printf("No se pudo crear el usuario")
		res.Code, res.Type, res.Msg = msg.GetByCode(22, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	if resCreateTrx.Error {
		logger.Error.Printf(resCreateTrx.Msg)
		res.Code, res.Type, res.Msg = msg.GetByCode(22, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	res.Data = "La credencial ha sido creada correctamente"
	res.Code, res.Type, res.Msg = msg.GetByCode(29, h.DB, h.TxID)
	res.Error = false
	return c.Status(http.StatusOK).JSON(res)
}

// get All Credentials godoc
// @Summary OnlyOne Smart Contract
// @Description Get All Credentials
// @Accept  json
// @Produce  json
// @Success 200 {object} requestCreateTransaction
// @Success 202 {object} responseCreateTransaction
// @Router /api/v1/credentials/all [get]
// @Authorization Bearer token
func (h *handlerCredentials) getAllCredentials(c *fiber.Ctx) error {
	res := responseAllCredentials{Error: true}
	e := env.NewConfiguration()
	block, err := strconv.Atoi(c.Params("block_id"))
	if err != nil {
		logger.Error.Printf("couldn't get block: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(1, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	limit, err := strconv.Atoi(c.Params("limit"))
	if err != nil {
		logger.Error.Printf("couldn't get limit: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(1, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	offset, err := strconv.Atoi(c.Params("offset"))
	if err != nil {
		logger.Error.Printf("couldn't get offset: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(1, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	connTrx, err := grpc.Dial(e.TransactionsService.Port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Error.Printf("error conectando con el servicio trx de blockchain: %s", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(22, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}
	defer connTrx.Close()

	clientTrx := transactions_proto.NewTransactionsServicesClient(connTrx)

	bearer := c.Get("Authorization")
	tkn := bearer[7:]

	ctx := grpcMetadata.AppendToOutgoingContext(context.Background(), "authorization", tkn)

	resTrx, err := clientTrx.GetAllTransactions(ctx, &transactions_proto.GetAllTransactionsRequest{
		Limit:   int64(limit),
		Offset:  int64(offset),
		BlockId: int64(block),
	})
	if err != nil {
		logger.Error.Printf("couldn't get all transactions: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(70, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	if resTrx == nil {
		logger.Error.Printf("couldn't get all transactions: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(70, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	if resTrx.Error {
		logger.Error.Printf(resTrx.Msg)
		res.Code, res.Type, res.Msg = msg.GetByCode(70, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	var credentials []*credential
	for _, trx := range resTrx.Data {
		credentials = append(credentials, &credential{
			From:   trx.From,
			To:     trx.To,
			Amount: trx.Amount,
			TypeId: int(trx.TypeId),
			Data:   trx.Data,
		})
	}

	res.Data = credentials
	res.Code, res.Type, res.Msg = msg.GetByCode(29, h.DB, h.TxID)
	res.Error = false
	return c.Status(http.StatusOK).JSON(res)
}

// get JWTTransaction godoc
// @Summary OnlyOne Smart Contract
// @Description Get JWTTransaction By ID
// @Accept  json
// @Produce  json
// @Success 200 {object} JwtTransactionRequest
// @Success 202 {object} JwtTransactionRequest
// @Router /api/v1/credentials/jwt [post]
// @Authorization Bearer token
func (h *handlerCredentials) getJWTTransaction(c *fiber.Ctx) error {
	res := JwtTransactionResponse{Error: true}

	m := JwtTransactionRequest{}
	err := c.BodyParser(&m)
	if err != nil {
		logger.Error.Printf("couldn't bind model create JwtTransactionRequest: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(1, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	tk := jwt.New(jwt.SigningMethodRS256)
	claims := tk.Claims.(jwt.MapClaims)
	claims["credential"] = m
	claims["exp"] = time.Now().Add(time.Minute * time.Duration(m.Ttl)).Unix()

	token, err := tk.SignedString(signKey)
	if err != nil {
		res.Code, res.Type, res.Msg = msg.GetByCode(1, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	res.Data = token
	res.Code, res.Type, res.Msg = msg.GetByCode(29, h.DB, h.TxID)
	res.Error = false
	return c.Status(http.StatusOK).JSON(res)
}
