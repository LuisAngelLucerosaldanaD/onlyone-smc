package credentials

import (
	"context"
	"crypto/rsa"
	"encoding/json"
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
	"onlyone_smc/internal/pwd"
	"onlyone_smc/pkg/auth"
	"onlyone_smc/pkg/cfg"
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

// createCredential godoc
// @Summary Create credential
// @Description Create credential
// @Tags Credentials
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Authorization" default(Bearer <Add access token here>)
// @Param Sign header string true "sign" default(<Add sign here>)
// @Param createCredential body requestCreateTransaction true "Request create transaction"
// @Success 200 {object} responseCreateCredential
// @Router /api/v1/credentials/create [post]
func (h *handlerCredentials) createCredential(c *fiber.Ctx) error {
	res := responseCreateCredential{Error: true}
	m := requestCreateTransaction{}
	e := env.NewConfiguration()
	sign := c.Get("sign")
	if sign == "" {
		logger.Error.Printf("couldn't get sign")
		res.Code, res.Type, res.Msg = 1, 2, "No se encontro la firma del mensaje"
		return c.Status(http.StatusAccepted).JSON(res)
	}

	u, err := helpers.GetUserContextV2(c)
	if err != nil {
		logger.Error.Printf("couldn't get user of token: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(1, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	srv := auth.NewServerAuth(h.DB, u, h.TxID)

	err = c.BodyParser(&m)
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

	token := c.Get("Authorization")[7:]

	ctx := grpcMetadata.AppendToOutgoingContext(context.Background(), "authorization", token)
	ctx = grpcMetadata.AppendToOutgoingContext(ctx, "sign", sign)

	if m.To == "" && m.IdentityNumber == "" {
		logger.Error.Printf("La wallet de destino es requerido")
		res.Code, res.Type, res.Msg = msg.GetByCode(22, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	if m.To == "" {
		walletToByIN, err := clientWallet.GetWalletByIdentityNumber(ctx, &wallet_proto.RqGetByIdentityNumber{IdentityNumber: m.IdentityNumber})
		if err != nil {
			logger.Error.Printf("couldn't get wallet by identity number: %v", err)
			res.Code, res.Type, res.Msg = msg.GetByCode(22, h.DB, h.TxID)
			return c.Status(http.StatusAccepted).JSON(res)
		}

		if walletToByIN == nil {
			logger.Error.Printf("couldn't get wallet by identity number")
			res.Code, res.Type, res.Msg = msg.GetByCode(22, h.DB, h.TxID)
			return c.Status(http.StatusAccepted).JSON(res)
		}

		if walletToByIN.Error {
			logger.Error.Printf(walletToByIN.Msg)
			res.Code, res.Type, res.Msg = msg.GetByCode(int(walletToByIN.Code), h.DB, h.TxID)
			return c.Status(http.StatusAccepted).JSON(res)
		}

		if walletToByIN == nil {
			wallet, err := clientWallet.CreateWalletBySystem(ctx, &wallet_proto.RqCreateWalletBySystem{})
			if err != nil {
				logger.Error.Printf("couldn't create wallet: %v", err)
				res.Code, res.Type, res.Msg = msg.GetByCode(3, h.DB, h.TxID)
				return c.Status(http.StatusAccepted).JSON(res)
			}

			if wallet == nil {
				logger.Error.Printf("couldn't create wallet")
				res.Code, res.Type, res.Msg = msg.GetByCode(3, h.DB, h.TxID)
				return c.Status(http.StatusAccepted).JSON(res)
			}

			if wallet.Error {
				logger.Error.Printf(wallet.Msg)
				res.Code, res.Type, res.Msg = msg.GetByCode(int(wallet.Code), h.DB, h.TxID)
				return c.Status(http.StatusAccepted).JSON(res)
			}

			_, code, err := srv.SrvUsersCredential.CreateUsersCredential(uuid.New().String(), wallet.Data.Key.Private, m.IdentityNumber, wallet.Data.Mnemonic)
			if err != nil {
				logger.Error.Printf("No se pudo registrar la key de la nueva wallet, error: %s", err.Error())
				res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
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
				res.Code, res.Type, res.Msg = msg.GetByCode(3, h.DB, h.TxID)
				return c.Status(http.StatusAccepted).JSON(res)
			}

			if resAccountTo == nil {
				logger.Error.Printf("couldn't create accounting to wallet: %v", err)
				res.Code, res.Type, res.Msg = msg.GetByCode(3, h.DB, h.TxID)
				return c.Status(http.StatusAccepted).JSON(res)
			}

			if resAccountTo.Error {
				logger.Error.Printf(resAccountTo.Msg)
				res.Code, res.Type, res.Msg = msg.GetByCode(int(resAccountTo.Code), h.DB, h.TxID)
				return c.Status(http.StatusAccepted).JSON(res)
			}
		} else {
			m.To = walletToByIN.Data.Id
		}
	}

	var files []*transactions_proto.File
	for _, file := range m.Files {
		files = append(files, &transactions_proto.File{
			IdFile:     int32(file.FileID),
			Name:       file.Name,
			FileEncode: file.FileEncode,
		})
	}

	dataReq := DataTrx{}
	err = json.Unmarshal([]byte(m.Data), &dataReq)
	if err != nil {
		logger.Error.Printf("La data no cumple con la estructura del modelo, error: %s", err.Error())
		res.Code, res.Type, res.Msg = msg.GetByCode(15, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	resCreateTrx, err := clientTrx.CreateTransaction(ctx, &transactions_proto.RequestCreateTransaction{
		From:   m.From,
		To:     m.To,
		Amount: m.Amount,
		TypeId: int32(m.TypeId),
		Data:   m.Data,
		Files:  files,
	})

	if err != nil {
		logger.Error.Printf("No se pudo crear la transaccion, error: %s", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(3, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	if resCreateTrx == nil {
		logger.Error.Printf("No se pudo crear el usuario")
		res.Code, res.Type, res.Msg = msg.GetByCode(3, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	if resCreateTrx.Error {
		logger.Error.Printf(resCreateTrx.Msg)
		res.Code, res.Type, res.Msg = msg.GetByCode(int(resCreateTrx.Code), h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	res.Data = resTrx{
		Id:        resCreateTrx.Data.Id,
		From:      resCreateTrx.Data.From,
		To:        resCreateTrx.Data.To,
		Amount:    resCreateTrx.Data.Amount,
		TypeId:    resCreateTrx.Data.TypeId,
		Data:      resCreateTrx.Data.Data,
		Block:     0,
		Files:     resCreateTrx.Data.Files,
		CreatedAt: resCreateTrx.Data.CreatedAt,
		UpdatedAt: resCreateTrx.Data.UpdatedAt,
	}
	res.Code, res.Type, res.Msg = msg.GetByCode(29, h.DB, h.TxID)
	res.Error = false
	return c.Status(http.StatusOK).JSON(res)
}

// getAllCredentials godoc
// @Summary Get all Credentials by block id
// @Description Get All Credentials ny block id
// @Tags Credentials
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Authorization" default(Bearer <Add access token here>)
// @Param block_id path int64 true "Block ID"
// @Param limit path int true "Limit of pagination"
// @Param offset path int true "Salt of pagination"
// @Success 200 {object} responseAllCredentials
// @Router /api/v1/credentials/all/{block_id}/{limit}/{offset} [get]
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
		res.Code, res.Type, res.Msg = msg.GetByCode(70, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}
	defer connTrx.Close()

	clientTrx := transactions_proto.NewTransactionsServicesClient(connTrx)

	tkn := c.Get("Authorization")[7:]

	ctx := grpcMetadata.AppendToOutgoingContext(context.Background(), "authorization", tkn)

	resTrx, err := clientTrx.GetAllTransactions(ctx, &transactions_proto.GetAllTransactionsRequest{
		Limit:   int64(limit),
		Offset:  int64(offset),
		BlockId: int64(block),
	})
	if err != nil {
		logger.Error.Printf("couldn't get all transactions: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(22, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	if resTrx == nil {
		logger.Error.Printf("couldn't get all transactions: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(22, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	if resTrx.Error {
		logger.Error.Printf(resTrx.Msg)
		res.Code, res.Type, res.Msg = msg.GetByCode(int(resTrx.Code), h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	var credentials []*credential
	for _, trx := range resTrx.Data {
		var files []File
		err := json.Unmarshal([]byte(trx.Files), &files)
		if err != nil {
			logger.Error.Printf("couldn't parsed trx files: %v", err)
			res.Code, res.Type, res.Msg = msg.GetByCode(22, h.DB, h.TxID)
			return c.Status(http.StatusAccepted).JSON(res)
		}
		credentials = append(credentials, &credential{
			Id:     trx.Id,
			From:   trx.From,
			To:     trx.To,
			Amount: trx.Amount,
			TypeId: int(trx.TypeId),
			Data:   trx.Data,
			Files:  files,
		})
	}

	res.Data = credentials
	res.Code, res.Type, res.Msg = msg.GetByCode(29, h.DB, h.TxID)
	res.Error = false
	return c.Status(http.StatusOK).JSON(res)
}

// getJWTTransaction godoc
// @Summary Get JWTTransaction By ID
// @Description Get JWTTransaction By ID
// @Tags Credentials
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Authorization" default(Bearer <Add access token here>)
// @Param getJWTTransaction body JwtTransactionRequest true "Generate JWT Request"
// @Success 200 {object} JwtTransactionResponse
// @Router /api/v1/credentials/jwt [post]
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

// getAllTransactionFiles godoc
// @Summary Get transaction files
// @Description Get transaction files
// @Tags Credentials
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Authorization" default(Bearer <Add access token here>)
// @Param trx path int true "transaction id"
// @Success 200 {object} ResGetFiles
// @Router /api/v1/credentials/files/{trx} [get]
func (h *handlerCredentials) getAllTransactionFiles(c *fiber.Ctx) error {
	res := ResGetFiles{Error: true}
	e := env.NewConfiguration()
	trxID := c.Params("trx")
	connTrx, err := grpc.Dial(e.TransactionsService.Port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Error.Printf("error conectando con el servicio trx de blockchain: %s", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(22, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}
	defer connTrx.Close()

	clientTrx := transactions_proto.NewTransactionsServicesClient(connTrx)

	tkn := c.Get("Authorization")[7:]

	ctx := grpcMetadata.AppendToOutgoingContext(context.Background(), "authorization", tkn)

	resFiles, err := clientTrx.GetFilesTransaction(ctx, &transactions_proto.GetFilesByTransactionRequest{TransactionId: trxID})
	if err != nil {
		logger.Error.Printf("couldn't get files: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(22, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	if resFiles == nil {
		logger.Error.Printf("couldn't get files")
		res.Code, res.Type, res.Msg = msg.GetByCode(22, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	if resFiles.Error {
		logger.Error.Printf(resFiles.Msg)
		res.Code, res.Type, res.Msg = msg.GetByCode(int(resFiles.Code), h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	var files []*File
	for _, file := range resFiles.Data {
		files = append(files, &File{
			FileID:     int(file.FileId),
			Name:       file.NameDocument,
			FileEncode: file.Encoding,
		})
	}

	res.Data = files
	res.Code, res.Type, res.Msg = msg.GetByCode(29, h.DB, h.TxID)
	res.Error = false
	return c.Status(http.StatusOK).JSON(res)
}

// sharedCredentials godoc
// @Summary Registra los datos de la credencial a compartir
// @Description Método para registrar los datos de la credencial que se va a compartir
// @Tags Credentials
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Authorization" default(Bearer <Add access token here>)
// @Param sharedCredentials body reqSharedCredentials true "Datos para crear la credencial a compartir"
// @Success 200 {object} ResAnny
// @Router /api/v1/credentials/shared [post]
func (h *handlerCredentials) sharedCredentials(c *fiber.Ctx) error {
	res := ResAnny{Error: true}
	req := reqSharedCredentials{}
	err := c.BodyParser(&req)
	if err != nil {
		logger.Error.Printf("couldn't bind model shared credentials: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(1, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	u, err := helpers.GetUserContextV2(c)
	if err != nil {
		logger.Error.Printf("couldn't get user of token: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(1, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	srvCfg := cfg.NewServerCfg(h.DB, u, h.TxID)

	resCredential, code, err := srvCfg.SrvSharedCredential.CreateSharedCredential(req.Data, u.ID, req.Password, req.ExpiredAt, req.MaxNumberQueries)
	if err != nil {
		logger.Error.Printf("No se pudo crear la credencial al compartir: %s", err.Error())
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	diff := req.ExpiredAt.Sub(time.Now())
	tk := jwt.New(jwt.SigningMethodRS256)
	claims := tk.Claims.(jwt.MapClaims)
	claims["id"] = resCredential.ID
	claims["exp"] = time.Now().Add(time.Hour * 24 * time.Duration(diff.Hours())).Unix()

	token, err := tk.SignedString(signKey)
	if err != nil {
		res.Code, res.Type, res.Msg = 22, 1, "No se pudo generar el token"
		return c.Status(http.StatusAccepted).JSON(res)
	}

	res.Data = token
	res.Code, res.Type, res.Msg = msg.GetByCode(29, h.DB, h.TxID)
	res.Error = false
	return c.Status(http.StatusOK).JSON(res)
}

// getSharedCredentials godoc
// @Summary Obtiene los datos de la credencial compartida
// @Description Método para obtener los datos de la credencial compartida
// @Tags Credentials
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Authorization" default(Bearer <Add access token here>)
// @Param sharedCredentials body reqGetSharedCredential true "Datos para obtener la credencial compartida"
// @Success 200 {object} ResAnny
// @Router /api/v1/credentials/get-shared [post]
func (h *handlerCredentials) getSharedCredentials(c *fiber.Ctx) error {
	res := ResAnny{Error: true}
	req := reqGetSharedCredential{}
	err := c.BodyParser(&req)
	if err != nil {
		logger.Error.Printf("couldn't bind model shared credentials: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(1, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	srvCfg := cfg.NewServerCfg(h.DB, nil, h.TxID)

	resCredential, code, err := srvCfg.SrvSharedCredential.GetSharedCredentialByID(req.Id)
	if err != nil {
		logger.Error.Printf("No se pudo obtener la credencial a compartir: %s", err.Error())
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	isValid := pwd.Compare(resCredential.UserId, resCredential.Password, req.Password)
	if !isValid {
		res.Code, res.Type, res.Msg = 22, 1, "La contraseña es incorrecta"
		return c.Status(http.StatusAccepted).JSON(res)
	}

	if resCredential.MaxNumberQueries == 0 {
		res.Code, res.Type, res.Msg = 22, 1, "Se ha superado el número máximo de consultas configuradas para esta credencial"
		return c.Status(http.StatusAccepted).JSON(res)
	}

	dateExpired := resCredential.ExpiredAt.Sub(time.Now())
	if dateExpired.Minutes() <= 0 {
		res.Code, res.Type, res.Msg = 22, 1, "La credencial ha caducado"
		return c.Status(http.StatusAccepted).JSON(res)
	}

	_, code, err = srvCfg.SrvSharedCredential.UpdateSharedCredential(resCredential.ID, resCredential.Data, resCredential.UserId, resCredential.Password, resCredential.ExpiredAt, resCredential.MaxNumberQueries-1)
	if err != nil {
		logger.Error.Printf("No se pudo actualizar la credencial a compartir: %s", err.Error())
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	res.Data = resCredential.Data
	res.Code, res.Type, res.Msg = msg.GetByCode(29, h.DB, h.TxID)
	res.Error = false
	return c.Status(http.StatusOK).JSON(res)
}
