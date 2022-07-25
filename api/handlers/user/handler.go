package users

import (
	"context"
	"encoding/base64"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	grpcMetadata "google.golang.org/grpc/metadata"
	"net/http"
	"onlyone_smc/internal/aws_ia"
	"onlyone_smc/internal/env"
	"onlyone_smc/internal/grpc/accounting_proto"
	"onlyone_smc/internal/grpc/auth_proto"
	"onlyone_smc/internal/grpc/users_proto"
	"onlyone_smc/internal/grpc/wallet_proto"
	"onlyone_smc/internal/helpers"
	"onlyone_smc/internal/logger"
	"onlyone_smc/internal/models"
	"onlyone_smc/internal/msg"
	"onlyone_smc/pkg/auth"
	"strconv"
)

type handlerUser struct {
	DB   *sqlx.DB
	TxID string
}

// create user godoc
// @Summary OnlyOne Smart Contract
// @Description Create User
// @Accept  json
// @Produce  json
// @Success 200 {object} requestCreateUser
// @Success 202 {object} responseCreateUser
// @Router /api/v1/user/create [post]
func (h *handlerUser) createUser(c *fiber.Ctx) error {
	res := responseCreateUser{Error: true}
	m := requestCreateUser{}
	e := env.NewConfiguration()
	err := c.BodyParser(&m)
	if err != nil {
		logger.Error.Printf("couldn't bind model create wallets: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(1, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	if m.Password != m.ConfirmPassword {
		logger.Error.Printf("this password is not equal to confirm_password")
		res.Code, res.Type, res.Msg = msg.GetByCode(10005, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	connAuth, err := grpc.Dial(e.AuthService.Port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Error.Printf("error conectando con el servicio auth de blockchain: %s", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(22, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}
	defer connAuth.Close()

	clientUser := users_proto.NewAuthServicesUsersClient(connAuth)

	clientAuth := auth_proto.NewAuthServicesUsersClient(connAuth)

	resLogin, err := clientAuth.Login(context.Background(), &auth_proto.LoginRequest{
		Email:    nil,
		Nickname: &e.App.UserLogin,
		Password: e.App.UserPassword,
	})
	if err != nil {
		logger.Error.Printf("No se pudo obtener el token de autenticacion: %s", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(22, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	if resLogin == nil {
		logger.Error.Printf("No se pudo obtener el token de autenticacion")
		res.Code, res.Type, res.Msg = msg.GetByCode(22, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	if resLogin.Error {
		logger.Error.Printf(resLogin.Msg)
		res.Code, res.Type, res.Msg = msg.GetByCode(22, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	ctx := grpcMetadata.AppendToOutgoingContext(context.Background(), "authorization", resLogin.Data.AccessToken)

	resCreateUser, err := clientUser.CreateUser(ctx, &users_proto.UserRequest{
		Id:              "",
		Nickname:        m.Nickname,
		Email:           m.Email,
		Password:        m.Password,
		ConfirmPassword: m.ConfirmPassword,
		Name:            m.Name,
		Lastname:        m.Lastname,
		IdType:          int32(m.IdType),
		IdNumber:        m.IdNumber,
		Cellphone:       m.Cellphone,
		BirthDate:       m.BirthDate.String(),
	})
	if err != nil {
		logger.Error.Printf("No se pudo crear el usuario, error: %s", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(22, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	if resCreateUser == nil {
		logger.Error.Printf("No se pudo crear el usuario")
		res.Code, res.Type, res.Msg = msg.GetByCode(22, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	if resCreateUser.Error {
		logger.Error.Printf(resCreateUser.Msg)
		res.Code, res.Type, res.Msg = msg.GetByCode(22, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	res.Data = "Usuario creado correctamente, se envió un correo de confirmación a su correo electrónico"
	res.Code, res.Type, res.Msg = msg.GetByCode(29, h.DB, h.TxID)
	res.Error = false
	return c.Status(http.StatusOK).JSON(res)
}

// validate user by email
// @Summary OnlyOne Smart Contract
// @Description validate user by email
// @Accept  json
// @Produce  json
// @Success 202 {object} responseUserValid
// @Router /validate-email/:email [get]
func (h *handlerUser) validateEmail(c *fiber.Ctx) error {
	res := responseUserValid{Error: true}
	emailStr := c.Params("email")

	e := env.NewConfiguration()
	connAuth, err := grpc.Dial(e.AuthService.Port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Error.Printf("error conectando con el servicio auth de blockchain: %s", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(22, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}
	defer connAuth.Close()

	clientUser := users_proto.NewAuthServicesUsersClient(connAuth)

	bearer := c.Get("Authorization")
	tkn := bearer[7:]

	ctx := grpcMetadata.AppendToOutgoingContext(context.Background(), "authorization", tkn)

	resValidate, err := clientUser.ValidateEmail(ctx, &users_proto.ValidateEmailRequest{Email: emailStr})
	if err != nil {
		logger.Error.Printf("couldn't get user by email: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(70, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	if resValidate == nil {
		logger.Error.Printf("couldn't get user by email: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(70, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	if resValidate.Error {
		logger.Error.Printf(resValidate.Msg)
		res.Code, res.Type, res.Msg = msg.GetByCode(70, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	res.Data = resValidate.Data
	res.Code, res.Type, res.Msg = msg.GetByCode(29, h.DB, h.TxID)
	res.Error = false
	return c.Status(http.StatusOK).JSON(res)
}

// validate user by nickname
// @Summary OnlyOne Smart Contract
// @Description validate user by nickname
// @Accept  json
// @Produce  json
// @Success 202 {object} responseUserValid
// @Router /validate-nickname/:nickname [get]
func (h *handlerUser) validateNickname(c *fiber.Ctx) error {
	res := responseUserValid{Error: true}
	nickname := c.Params("nickname")

	e := env.NewConfiguration()
	connAuth, err := grpc.Dial(e.AuthService.Port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Error.Printf("error conectando con el servicio auth de blockchain: %s", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(22, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}
	defer connAuth.Close()

	clientUser := users_proto.NewAuthServicesUsersClient(connAuth)

	bearer := c.Get("Authorization")
	tkn := bearer[7:]

	ctx := grpcMetadata.AppendToOutgoingContext(context.Background(), "authorization", tkn)

	resValidate, err := clientUser.ValidateNickname(ctx, &users_proto.ValidateNicknameRequest{Nickname: nickname})
	if err != nil {
		logger.Error.Printf("couldn't get user by nickname: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(70, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	if resValidate == nil {
		logger.Error.Printf("couldn't get user by nickname: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(70, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	if resValidate.Error {
		logger.Error.Printf(resValidate.Msg)
		res.Code, res.Type, res.Msg = msg.GetByCode(70, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	res.Data = resValidate.Data
	res.Code, res.Type, res.Msg = msg.GetByCode(29, h.DB, h.TxID)
	res.Error = false
	return c.Status(http.StatusOK).JSON(res)
}

// get user by ID
// @Summary OnlyOne Smart Contract
// @Description get user by ID
// @Accept  json
// @Produce  json
// @Success 202 {object} responseUser
// @Router /api/v1/user/:id [get]
func (h *handlerUser) getUserById(c *fiber.Ctx) error {
	res := responseUser{Error: true}
	usrId := c.Params("id")
	srvUser := auth.NewServerAuth(h.DB, nil, h.TxID)

	e := env.NewConfiguration()
	connAuth, err := grpc.Dial(e.AuthService.Port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Error.Printf("error conectando con el servicio auth de blockchain: %s", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(22, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}
	defer connAuth.Close()

	clientUser := users_proto.NewAuthServicesUsersClient(connAuth)

	bearer := c.Get("Authorization")
	tkn := bearer[7:]

	ctx := grpcMetadata.AppendToOutgoingContext(context.Background(), "authorization", tkn)

	usr, err := clientUser.GetUserById(ctx, &users_proto.GetUserByIDRequest{Id: usrId})
	if err != nil {
		logger.Error.Printf("couldn't get User by ID: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(70, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	if usr == nil {
		logger.Error.Printf("couldn't get User by ID: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(70, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	if usr.Error {
		logger.Error.Printf(usr.Msg)
		res.Code, res.Type, res.Msg = msg.GetByCode(70, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	user := usr.Data

	if user.FullPathPhoto != "" {

		profile, code, err := srvUser.SrvUserProfile.GetUserProfileByUserID(user.Id)
		if err != nil {
			logger.Error.Printf("couldn't get user photo by ID: %v", err)
			res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
			return c.Status(http.StatusAccepted).JSON(res)
		}

		file, _, err := srvUser.SrvFiles.GetFileByPath(profile.Path, profile.Name)
		if err != nil {
			logger.Error.Printf("couldn't get profile picture: %v", err)
		}

		if file != nil {
			user.FullPathPhoto = file.Encoding
		}
	}

	res.Data = &models.User{
		ID:             user.Id,
		Nickname:       user.Nickname,
		Email:          user.Email,
		Name:           user.Name,
		Lastname:       user.Lastname,
		IdType:         int(user.IdType),
		IdNumber:       user.IdNumber,
		Cellphone:      user.Cellphone,
		StatusId:       int(user.StatusId),
		FailedAttempts: int(user.FailedAttempts),
		IdRole:         int(user.IdRole),
		FullPathPhoto:  user.FullPathPhoto,
		RsaPublic:      user.RsaPublic,
		RealIP:         user.RealIp,
	}
	res.Code, res.Type, res.Msg = msg.GetByCode(29, h.DB, h.TxID)
	res.Error = false
	return c.Status(http.StatusOK).JSON(res)
}

// validate user by identity number
// @Summary blockchain
// @Description validate user by identity number
// @Accept  json
// @Produce  json
// @Success 202 {object} responseValidateUser
// @Router /api/v1/user/validate-identity [post]
func (h *handlerUser) validateIdentity(c *fiber.Ctx) error {
	res := responseValidateUser{Error: true}
	req := requestValidateIdentity{}
	e := env.NewConfiguration()
	u, err := helpers.GetUserContextV2(c)
	if err != nil {
		logger.Error.Printf("couldn't get user token: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(1, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}
	srvAuth := auth.NewServerAuth(h.DB, nil, h.TxID)

	err = c.BodyParser(&req)
	if err != nil {
		logger.Error.Printf("couldn't bind model create wallets: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(1, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	connAuth, err := grpc.Dial(e.AuthService.Port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Error.Printf("error conectando con el servicio auth de blockchain: %s", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(22, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}
	defer connAuth.Close()

	clientUser := users_proto.NewAuthServicesUsersClient(connAuth)
	clientWallet := wallet_proto.NewWalletServicesWalletClient(connAuth)
	clientAccount := accounting_proto.NewAccountingServicesAccountingClient(connAuth)

	bearer := c.Get("Authorization")
	tkn := bearer[7:]

	ctx := grpcMetadata.AppendToOutgoingContext(context.Background(), "authorization", tkn)

	identityBytes, err := base64.StdEncoding.DecodeString(req.IdentityEncode)
	if err != nil {
		logger.Error.Printf("couldn't decode identity: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(1, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	confirmBytes, err := base64.StdEncoding.DecodeString(req.ConfirmEncode)
	if err != nil {
		logger.Error.Printf("couldn't decode identity: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(1, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	userFields, err := aws_ia.GetUserFields(identityBytes)
	if err != nil {
		logger.Error.Printf("couldn't get user fields values: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(22, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	if userFields.Names == "" || userFields.SecondSurname == "" || userFields.Surname == "" || userFields.IdentityNumber == "" {
		logger.Error.Printf("couldn't get user fields values: %v", err)
		res.Code, res.Type, res.Msg = 50, 1, "No se pudo obtener los datos del usuario, intente subir una foto del documento de identidad con mayor resolución"
		return c.Status(http.StatusAccepted).JSON(res)
	}

	resp, err := aws_ia.CompareFaces(identityBytes, confirmBytes)
	if err != nil {
		logger.Error.Printf("couldn't decode identity: %v", err)
		res.Code, res.Type, res.Msg = 22, 1, "No se pudo comparar las dos fotos de identificación"
		return c.Status(http.StatusAccepted).JSON(res)
	}

	if !resp {
		logger.Error.Printf("couldn't decode identity: %v", err)
		res.Code, res.Type, res.Msg = 22, 1, "La foto de identidad no coincide"
		return c.Status(http.StatusAccepted).JSON(res)
	}

	idNumber, _ := strconv.ParseInt(req.IdentityNumber, 10, 64)

	f, err := srvAuth.SrvFiles.UploadFile(idNumber, req.IdentityNumber+".jpg", req.ConfirmEncode)
	if err != nil {
		logger.Error.Printf("couldn't upload file s3: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(3, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	resUser, err := clientUser.GetUserById(ctx, &users_proto.GetUserByIDRequest{Id: u.ID})
	if err != nil {
		logger.Error.Printf("couldn't get user by id: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(70, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	if resUser == nil {
		logger.Error.Printf("couldn't get user by id: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(70, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	if resUser.Error {
		logger.Error.Printf(resUser.Msg)
		res.Code, res.Type, res.Msg = msg.GetByCode(70, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	user := resUser.Data

	resUserWallet, err := clientUser.GetUserWalletByIdentityNumber(ctx, &users_proto.RqGetUserWalletByIdentityNumber{
		IdentityNumber: req.IdentityNumber,
		UserId:         user.Id,
	})
	if err != nil {
		logger.Error.Printf("couldn't get user wallet by user id and identity number: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(70, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	if resUserWallet == nil {
		logger.Error.Printf("couldn't get user wallet by user id and identity number: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(70, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	if resUserWallet.Error {
		logger.Error.Printf(resUserWallet.Msg)
		res.Code, res.Type, res.Msg = msg.GetByCode(70, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	if resUserWallet.Data != nil {
		res.Data = nil
		res.Code, res.Type, res.Msg = 22, 1, "El usuario ya ha valido su identidad"
		res.Error = false
		return c.Status(http.StatusAccepted).JSON(res)
	}

	resUpdateUser, err := clientUser.UpdateUser(ctx, &users_proto.RqUpdateUser{
		Id:            user.Id,
		Nickname:      user.Nickname,
		Email:         user.Email,
		Password:      user.Password,
		Name:          userFields.Names,
		Lastname:      userFields.Surname + " " + userFields.SecondSurname,
		IdType:        user.IdType,
		IdNumber:      userFields.IdentityNumber,
		Cellphone:     user.Cellphone,
		BirthDate:     user.BirthDate,
		VerifiedCode:  user.VerifiedCode,
		VerifiedAt:    user.VerifiedAt,
		FullPathPhoto: user.FullPathPhoto,
		RsaPrivate:    user.RsaPrivate,
		RsaPublic:     user.RsaPublic,
		IdRole:        21,
	})
	if err != nil {
		logger.Error.Printf("couldn't update user by id: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(70, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	if resUpdateUser == nil {
		logger.Error.Printf("couldn't update user by id: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(70, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	if resUpdateUser.Error {
		logger.Error.Printf(resUpdateUser.Msg)
		res.Code, res.Type, res.Msg = msg.GetByCode(70, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	resWallet, err := clientWallet.GetWalletByIdentityNumber(ctx, &wallet_proto.RqGetByIdentityNumber{IdentityNumber: req.IdentityNumber})
	if err != nil {
		logger.Error.Printf("couldn't get wallet by identity number: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(70, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	if resWallet == nil {
		logger.Error.Printf("couldn't get wallet by identity number")
		res.Code, res.Type, res.Msg = msg.GetByCode(70, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	if resWallet.Error {
		logger.Error.Printf(resWallet.Msg)
		res.Code, res.Type, res.Msg = msg.GetByCode(70, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	wallet := resWallet.Data

	var walletID string

	if wallet == nil {
		newWallet, err := clientWallet.CreateWalletBySystem(ctx, &wallet_proto.RqCreateWalletBySystem{IdentityNumber: req.IdentityNumber})
		if err != nil {
			logger.Error.Printf("couldn't create wallet: %v", err)
			res.Code, res.Type, res.Msg = msg.GetByCode(70, h.DB, h.TxID)
			return c.Status(http.StatusAccepted).JSON(res)
		}

		if newWallet == nil {
			logger.Error.Printf("couldn't create wallet: %v", err)
			res.Code, res.Type, res.Msg = msg.GetByCode(70, h.DB, h.TxID)
			return c.Status(http.StatusAccepted).JSON(res)
		}

		if newWallet.Error {
			logger.Error.Printf(newWallet.Msg)
			res.Code, res.Type, res.Msg = msg.GetByCode(70, h.DB, h.TxID)
			return c.Status(http.StatusAccepted).JSON(res)
		}

		wallet = &wallet_proto.Wallet{
			Id:       newWallet.Data.Id,
			Mnemonic: newWallet.Data.Mnemonic,
		}

		walletID = newWallet.Data.Id

		resCreateAccount, err := clientAccount.CreateAccounting(ctx, &accounting_proto.RequestCreateAccounting{
			Id:       uuid.New().String(),
			IdWallet: newWallet.Data.Id,
			Amount:   0,
			IdUser:   u.ID,
		})
		if err != nil {
			logger.Error.Printf("couldn't create account to wallet: %v", err)
			res.Code, res.Type, res.Msg = msg.GetByCode(70, h.DB, h.TxID)
			return c.Status(http.StatusAccepted).JSON(res)
		}

		if resCreateAccount == nil {
			logger.Error.Printf("couldn't create account to wallet: %v", err)
			res.Code, res.Type, res.Msg = msg.GetByCode(70, h.DB, h.TxID)
			return c.Status(http.StatusAccepted).JSON(res)
		}

		if resCreateAccount.Error {
			logger.Error.Printf(resCreateAccount.Msg)
			res.Code, res.Type, res.Msg = msg.GetByCode(70, h.DB, h.TxID)
			return c.Status(http.StatusAccepted).JSON(res)
		}

	} else {
		walletID = wallet.Id
		resUpdateWallet, err := clientWallet.UpdateWallet(ctx, &wallet_proto.RqUpdateWallet{
			Id:               wallet.Id,
			RsaPublic:        wallet.RsaPublic,
			RsaPrivate:       wallet.RsaPrivate,
			RsaPublicDevice:  wallet.RsaPublicDevice,
			RsaPrivateDevice: wallet.RsaPrivateDevice,
			IpDevice:         wallet.IpDevice,
			IdentityNumber:   wallet.IdentityNumber,
			StatusId:         wallet.StatusId,
		})
		if err != nil {
			logger.Error.Printf("couldn't update wallet: %v", err)
			res.Code, res.Type, res.Msg = msg.GetByCode(70, h.DB, h.TxID)
			return c.Status(http.StatusAccepted).JSON(res)
		}

		if resUserWallet == nil {
			logger.Error.Printf("couldn't update wallet: %v", err)
			res.Code, res.Type, res.Msg = msg.GetByCode(70, h.DB, h.TxID)
			return c.Status(http.StatusAccepted).JSON(res)
		}

		if resUpdateWallet.Error {
			logger.Error.Printf(resUpdateWallet.Msg)
			res.Code, res.Type, res.Msg = msg.GetByCode(70, h.DB, h.TxID)
			return c.Status(http.StatusAccepted).JSON(res)
		}

		wallet = &wallet_proto.Wallet{
			Id:       resUpdateWallet.Data.Id,
			Mnemonic: resUpdateWallet.Data.Mnemonic,
		}
	}

	resCreateUserWallet, err := clientUser.CreateUserWallet(ctx, &users_proto.RqCreateUserWallet{
		UserId:   u.ID,
		WalletId: walletID,
	})
	if err != nil {
		logger.Error.Printf("couldn't create users wallet: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(70, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	if resCreateUserWallet == nil {
		logger.Error.Printf("couldn't create users wallet: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(70, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	if resCreateUserWallet.Error {
		logger.Error.Printf(resCreateUserWallet.Msg)
		res.Code, res.Type, res.Msg = msg.GetByCode(70, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	user.FullPathPhoto = f.Path + "-/" + f.FileName

	_, code, err := srvAuth.SrvUserProfile.CreateUserProfile(user.IdUser, f.Path, f.FileName)
	if err != nil {
		logger.Error.Printf("couldn't update profile picture: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		res.Msg = err.Error()
		return c.Status(http.StatusAccepted).JSON(res)
	}

	res.Data = &models.Wallet{
		ID:       wallet.Id,
		Mnemonic: wallet.Mnemonic,
	}
	res.Code, res.Type, res.Msg = msg.GetByCode(29, h.DB, h.TxID)
	res.Error = false
	return c.Status(http.StatusOK).JSON(res)
}

// getUserByIdentityNumber
// @Summary OnlyOne Smart Contract
// @Description get user by identity number
// @Accept  json
// @Produce  json
// @Success 202 {object} responseUser
// @Router /api/v1/user/:inumber [get]
func (h *handlerUser) getUserByIdentityNumber(c *fiber.Ctx) error {
	res := responseUserValid{Error: true}
	usrIdentityNumber := c.Params("inumber")
	e := env.NewConfiguration()
	connAuth, err := grpc.Dial(e.AuthService.Port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Error.Printf("error conectando con el servicio auth de blockchain: %s", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(22, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}
	defer connAuth.Close()

	clientUser := users_proto.NewAuthServicesUsersClient(connAuth)

	bearer := c.Get("Authorization")
	tkn := bearer[7:]

	ctx := grpcMetadata.AppendToOutgoingContext(context.Background(), "authorization", tkn)

	resUsr, err := clientUser.ValidIdentityNumber(ctx, &users_proto.RequestGetByIdentityNumber{IdentityNumber: usrIdentityNumber})
	if err != nil {
		logger.Error.Printf("couldn't get User by identity number: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(70, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	if resUsr == nil {
		logger.Error.Printf("couldn't get User by identity number: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(70, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	if resUsr.Error {
		logger.Error.Printf(resUsr.Msg)
		res.Code, res.Type, res.Msg = msg.GetByCode(70, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	res.Data = resUsr.Data
	res.Code, res.Type, res.Msg = msg.GetByCode(29, h.DB, h.TxID)
	res.Error = false
	return c.Status(http.StatusOK).JSON(res)
}

// getUserPictureProfile godoc
// @Summary OnlyOne Smart Contract
// @Description getUserPictureProfile
// @Accept  json
// @Produce  json
// @Success 202 {object} responseUpdateUser
// @Router /api/v1/user/picture-profile [GET]
// @Authorization Bearer token
func (h *handlerUser) getUserPictureProfile(c *fiber.Ctx) error {
	res := responseUpdateUser{Error: true, Data: ""}
	u, err := helpers.GetUserContextV2(c)
	if err != nil {
		logger.Error.Printf("couldn't get token user: %s", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(1, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	e := env.NewConfiguration()

	srvAuth := auth.NewServerAuth(h.DB, nil, h.TxID)

	connAuth, err := grpc.Dial(e.AuthService.Port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Error.Printf("error conectando con el servicio auth de blockchain: %s", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(22, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}
	defer connAuth.Close()

	clientUser := users_proto.NewAuthServicesUsersClient(connAuth)

	bearer := c.Get("Authorization")
	tkn := bearer[7:]

	ctx := grpcMetadata.AppendToOutgoingContext(context.Background(), "authorization", tkn)

	resUser, err := clientUser.GetUserById(ctx, &users_proto.GetUserByIDRequest{Id: u.ID})
	if err != nil {
		logger.Error.Printf("Error trayendo el usuario por su id, error: %s", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(22, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	if resUser == nil {
		logger.Error.Printf("Error trayendo el usuario por su id")
		res.Code, res.Type, res.Msg = msg.GetByCode(22, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	if resUser.Error {
		logger.Error.Printf(resUser.Msg)
		res.Code, res.Type, res.Msg = msg.GetByCode(22, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	user := resUser.Data

	if user.FullPathPhoto != "" {

		profile, code, err := srvAuth.SrvUserProfile.GetUserProfileByUserID(user.Id)
		if err != nil {
			logger.Error.Printf("couldn't get profile picture: %v", err)
			res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
			return c.Status(http.StatusAccepted).JSON(res)
		}

		file, code, err := srvAuth.SrvFiles.GetFileByPath(profile.Path, profile.Name)
		if err != nil {
			logger.Error.Printf("couldn't get profile picture: %v", err)
			res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
			return c.Status(http.StatusAccepted).JSON(res)
		}

		res.Data = file.Encoding
	}

	res.Code, res.Type, res.Msg = msg.GetByCode(29, h.DB, h.TxID)
	res.Error = false
	return c.Status(http.StatusOK).JSON(res)
}

// changePassword godoc
// @Summary OnlyOne Smart Contract
// @Description changePassword
// @Accept  json
// @Produce  json
// @Success 202 {object} responseUpdateUser
// @Router /api/v1/user/update-password [GET]
// @Authorization Bearer token
func (h *handlerUser) changePassword(c *fiber.Ctx) error {
	res := responseUpdateUser{Error: true, Data: ""}
	e := env.NewConfiguration()
	m := requestUpdatePassword{}
	err := c.BodyParser(&m)
	if err != nil {
		logger.Error.Printf("couldn't bind model login: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(1, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	connAuth, err := grpc.Dial(e.AuthService.Port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Error.Printf("error conectando con el servicio auth de blockchain: %s", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(22, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}
	defer connAuth.Close()

	clientUser := users_proto.NewAuthServicesUsersClient(connAuth)

	bearer := c.Get("Authorization")
	tkn := bearer[7:]

	ctx := grpcMetadata.AppendToOutgoingContext(context.Background(), "authorization", tkn)

	resChangePwd, err := clientUser.ChangePassword(ctx, &users_proto.RequestChangePwd{
		OldPassword:     m.OldPassword,
		NewPassword:     m.NewPassword,
		ConfirmPassword: m.ConfirmPassword,
	})
	if err != nil {
		logger.Error.Printf("No se pudo actualizar la contraseña, err: %s", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(22, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	if resChangePwd == nil {
		logger.Error.Printf("No se pudo actualizar la contraseña")
		res.Code, res.Type, res.Msg = msg.GetByCode(22, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	if resChangePwd.Error {
		logger.Error.Printf(resChangePwd.Msg)
		res.Code, res.Type, res.Msg = msg.GetByCode(22, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	res.Data = "contraseña actualizada correctamente"
	res.Code, res.Type, res.Msg = msg.GetByCode(29, h.DB, h.TxID)
	res.Error = false
	return c.Status(http.StatusOK).JSON(res)
}

// getWalletByUserId godoc
// @Summary OnlyOne Smart Contract
// @Description Get Wallet by user id
// @Accept  json
// @Produce  json
// @Success 200 {object} responseGetWallets
// @Router /api/v1/wallet/user [get]
func (h *handlerUser) getWalletByUserId(c *fiber.Ctx) error {
	e := env.NewConfiguration()
	res := responseGetWallets{Error: true}
	u, err := helpers.GetUserContextV2(c)
	if err != nil {
		logger.Error.Printf("couldn't get token user, error: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(1, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	connAuth, err := grpc.Dial(e.AuthService.Port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Error.Printf("error conectando con el servicio auth de blockchain: %s", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(22, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}
	defer connAuth.Close()

	bearer := c.Get("Authorization")
	tkn := bearer[7:]
	ctx := grpcMetadata.AppendToOutgoingContext(context.Background(), "authorization", tkn)

	clientWallet := wallet_proto.NewWalletServicesWalletClient(connAuth)

	wt, err := clientWallet.GetWalletByUserId(ctx, &wallet_proto.RequestGetWalletByUserId{UserId: u.ID})
	if err != nil {
		logger.Error.Printf("couldn't get wallets by id: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(70, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	if wt == nil {
		logger.Error.Printf("couldn't get wallets by id: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(70, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	if wt.Error {
		logger.Error.Printf(wt.Msg)
		res.Code, res.Type, res.Msg = msg.GetByCode(70, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	var wallets []*models.Wallet

	for _, wallet := range wt.Data {
		wallets = append(wallets, &models.Wallet{
			ID:               wallet.Id,
			Mnemonic:         wallet.Mnemonic,
			RsaPublic:        wallet.RsaPublic,
			RsaPrivate:       wallet.RsaPrivate,
			RsaPublicDevice:  wallet.RsaPublicDevice,
			RsaPrivateDevice: wallet.RsaPrivateDevice,
			IpDevice:         wallet.IpDevice,
			StatusId:         int(wallet.StatusId),
			IdentityNumber:   wallet.IdentityNumber,
		})
	}

	res.Data = wallets
	res.Code, res.Type, res.Msg = msg.GetByCode(29, h.DB, h.TxID)
	res.Error = false
	return c.Status(http.StatusOK).JSON(res)
}
