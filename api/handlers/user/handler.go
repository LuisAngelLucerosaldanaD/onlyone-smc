package users

import (
	"context"
	"encoding/base64"
	"encoding/json"
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
	"onlyone_smc/internal/grpc/users_proto"
	"onlyone_smc/internal/grpc/wallet_proto"
	"onlyone_smc/internal/helpers"
	"onlyone_smc/internal/logger"
	"onlyone_smc/internal/models"
	"onlyone_smc/internal/msg"
	"onlyone_smc/internal/ocr"
	"onlyone_smc/internal/ws"
	"onlyone_smc/pkg/auth"
	"strconv"
	"strings"
)

type handlerUser struct {
	DB   *sqlx.DB
	TxID string
}

// createUser godoc
// @Summary Create User of OnlyOne - BLion
// @Description Create User
// @tags User
// @Accept  json
// @Produce  json
// @Param createUser body requestCreateUser true "Request create user"
// @Success 200 {object} responseCreateUser
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

	resCreateUser, err := clientUser.CreateUser(context.Background(), &users_proto.UserRequest{
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
		res.Code, res.Type, res.Msg = msg.GetByCode(3, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	if resCreateUser == nil {
		logger.Error.Printf("No se pudo crear el usuario")
		res.Code, res.Type, res.Msg = msg.GetByCode(3, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	if resCreateUser.Error {
		logger.Error.Printf(resCreateUser.Msg)
		res.Code, res.Type, res.Msg = msg.GetByCode(3, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	res.Data = "Usuario creado correctamente, se envió un correo de confirmación a su correo electrónico"
	res.Code, res.Type, res.Msg = msg.GetByCode(29, h.DB, h.TxID)
	res.Error = false
	return c.Status(http.StatusOK).JSON(res)
}

// validateEmail godoc
// @Summary validate user by email
// @Description validate user by email
// @tags User
// @Accept  json
// @Produce  json
// @Param email path string true "email of user"
// @Success 200 {object} responseUserValid
// @Router /api/v1/user/validate-email/{email} [get]
func (h *handlerUser) validateEmail(c *fiber.Ctx) error {
	res := responseUserValid{Error: true}
	emailStr := c.Params("email")

	e := env.NewConfiguration()
	connAuth, err := grpc.Dial(e.AuthService.Port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Error.Printf("error conectando con el servicio auth de blockchain: %s", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(70, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}
	defer connAuth.Close()

	clientUser := users_proto.NewAuthServicesUsersClient(connAuth)

	resValidate, err := clientUser.ValidateEmail(context.Background(), &users_proto.ValidateEmailRequest{Email: emailStr})
	if err != nil {
		logger.Error.Printf("couldn't get user by email: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(22, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	if resValidate == nil {
		logger.Error.Printf("couldn't get user by email: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(22, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	if resValidate.Error {
		logger.Error.Printf(resValidate.Msg)
		res.Data = resValidate.Data
		res.Code, res.Type, res.Msg = msg.GetByCode(int(resValidate.Code), h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	res.Data = resValidate.Data
	res.Code, res.Type, res.Msg = msg.GetByCode(29, h.DB, h.TxID)
	res.Error = false
	return c.Status(http.StatusOK).JSON(res)
}

// validateNickname godoc
// @Summary validate user by nickname
// @Description validate user by nickname
// @tags User
// @Accept  json
// @Produce  json
// @Param nickname path string true "username (nickname) of user"
// @Success 200 {object} responseUserValid
// @Router /api/v1/user/validate-nickname/{nickname} [get]
func (h *handlerUser) validateNickname(c *fiber.Ctx) error {
	res := responseUserValid{Error: true}
	nickname := c.Params("nickname")

	e := env.NewConfiguration()
	connAuth, err := grpc.Dial(e.AuthService.Port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Error.Printf("error conectando con el servicio auth de blockchain: %s", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(70, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}
	defer connAuth.Close()

	clientUser := users_proto.NewAuthServicesUsersClient(connAuth)

	resValidate, err := clientUser.ValidateNickname(context.Background(), &users_proto.ValidateNicknameRequest{Nickname: nickname})
	if err != nil {
		logger.Error.Printf("couldn't get user by nickname: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(22, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	if resValidate == nil {
		logger.Error.Printf("couldn't get user by nickname: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(22, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	if resValidate.Error {
		logger.Error.Printf(resValidate.Msg)
		res.Data = resValidate.Data
		res.Code, res.Type, res.Msg = msg.GetByCode(int(resValidate.Code), h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	res.Data = resValidate.Data
	res.Code, res.Type, res.Msg = msg.GetByCode(29, h.DB, h.TxID)
	res.Error = false
	return c.Status(http.StatusOK).JSON(res)
}

// getUserById godoc
// @Summary get user by ID
// @Description get user by ID
// @tags User
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Authorization" default(Bearer <Add access token here>)
// @Param id path string true "user ID"
// @Success 200 {object} responseUser
// @Router /api/v1/user/{id} [get]
func (h *handlerUser) getUserById(c *fiber.Ctx) error {
	res := responseUser{Error: true}
	usrId := c.Params("id")
	srvUser := auth.NewServerAuth(h.DB, nil, h.TxID)

	e := env.NewConfiguration()
	connAuth, err := grpc.Dial(e.AuthService.Port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Error.Printf("error conectando con el servicio auth de blockchain: %s", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(70, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}
	defer connAuth.Close()

	clientUser := users_proto.NewAuthServicesUsersClient(connAuth)

	tkn := c.Get("Authorization")[7:]

	ctx := grpcMetadata.AppendToOutgoingContext(context.Background(), "authorization", tkn)

	usr, err := clientUser.GetUserById(ctx, &users_proto.GetUserByIDRequest{Id: usrId})
	if err != nil {
		logger.Error.Printf("couldn't get User by ID: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(22, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	if usr == nil {
		logger.Error.Printf("couldn't get User by ID: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(22, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	if usr.Error {
		logger.Error.Printf(usr.Msg)
		res.Code, res.Type, res.Msg = msg.GetByCode(int(usr.Code), h.DB, h.TxID)
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
		RealIP:         user.RealIp,
	}
	res.Code, res.Type, res.Msg = msg.GetByCode(29, h.DB, h.TxID)
	res.Error = false
	return c.Status(http.StatusOK).JSON(res)
}

// validateIdentity godoc
// @Summary validity user identity and create wallet
// @Description validity user identity and create wallet
// @tags User
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Authorization" default(Bearer <Add access token here>)
// @Param validateIdentity body requestValidateIdentity true "request of validate user identity"
// @Success 200 {object} responseValidateUser
// @Router /api/v1/user/validate-identity [post]
func (h *handlerUser) validateIdentity(c *fiber.Ctx) error {
	res := responseValidateUser{Error: true}
	resPerson := responsePerson{}
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

	token := c.Get("Authorization")[7:]

	ctx := grpcMetadata.AppendToOutgoingContext(context.Background(), "authorization", token)

	identityBytes, err := base64.StdEncoding.DecodeString(req.DocumentFrontImg)
	if err != nil {
		logger.Error.Printf("couldn't decode identity: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(1, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	confirmBytes, err := base64.StdEncoding.DecodeString(req.SelfieImg)
	if err != nil {
		logger.Error.Printf("couldn't decode identity: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(1, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	var userFields *ocr.PersonInformation

	if req.Country == "PE" {
		userFields, err = ocr.GetDniInformation(e.Ocr.Url, "dni.png", identityBytes)
		if err != nil {
			logger.Error.Printf("couldn't get user fields values: %v", err)
			res.Code, res.Type, res.Msg = msg.GetByCode(22, h.DB, h.TxID)
			return c.Status(http.StatusAccepted).JSON(res)
		}
	} else {
		userFields, err = ocr.GetCedulaInformation(e.Ocr.Url, "dni.png", identityBytes)
		if err != nil {
			logger.Error.Printf("couldn't get user fields values: %v", err)
			res.Code, res.Type, res.Msg = msg.GetByCode(22, h.DB, h.TxID)
			return c.Status(http.StatusAccepted).JSON(res)
		}
	}

	if req.Country == "PE" && (userFields.Names == "" || userFields.Surnames == "" || userFields.IdentityNumber == "") {
		logger.Error.Printf("couldn't get user fields values: %v", err)
		res.Code, res.Type, res.Msg = 50, 1, "No se pudo obtener los datos del usuario, intente subir una foto del documento de identidad con mayor resolución"
		return c.Status(http.StatusAccepted).JSON(res)
	}

	if req.Country == "CO" && userFields.IdentityNumber != "" {

		resWs, code, err := ws.ConsumeWS(nil, e.App.UrlPersons+userFields.IdentityNumber, "GET", token)
		if err != nil || code != 200 {
			logger.Error.Printf("No se pudo obtener la persona por el número de identificación: %v", err)
			res.Code, res.Type, res.Msg = code, 1, "No se pudo obtener la persona por el número de identificación"
			return c.Status(http.StatusAccepted).JSON(res)
		}

		err = json.Unmarshal(resWs, &resPerson)
		if err != nil {
			logger.Error.Printf("No se pudo parsear la respuesta: %v", err)
			res.Code, res.Type, res.Msg = 22, 1, "Error al obtener los datos de la persona"
			return c.Status(http.StatusAccepted).JSON(res)
		}

		if resPerson.Error {
			logger.Error.Printf(resPerson.Msg)
			res.Code, res.Type, res.Msg = resPerson.Code, 1, resPerson.Msg
			return c.Status(http.StatusAccepted).JSON(res)
		}

		person := resPerson.Data

		userFields.Names = strings.TrimSpace(person.FirstName + person.SecondName)
		userFields.Surnames = strings.TrimSpace(person.Surname + " " + person.SecondSurname)

	} else {
		logger.Error.Printf("couldn't get user fields values: %v", err)
		res.Code, res.Type, res.Msg = 50, 1, "No se pudo obtener los datos del usuario, intente subir una foto del documento de identidad con mayor resolución"
		return c.Status(http.StatusAccepted).JSON(res)
	}

	resUser, err := clientUser.GetUserById(ctx, &users_proto.GetUserByIDRequest{Id: u.ID})
	if err != nil || resUser == nil {
		logger.Error.Printf("couldn't get user by id: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(22, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	if resUser.Error {
		logger.Error.Printf(resUser.Msg)
		res.Code, res.Type, res.Msg = msg.GetByCode(int(resUser.Code), h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	user := resUser.Data

	resUserWallet, err := clientUser.GetUserWalletByIdentityNumber(ctx, &users_proto.RqGetUserWalletByIdentityNumber{
		IdentityNumber: userFields.IdentityNumber,
		UserId:         user.Id,
	})
	if err != nil || resUserWallet == nil {
		logger.Error.Printf("couldn't get user wallet by user id and identity number: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(22, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	if resUserWallet.Error {
		logger.Error.Printf(resUserWallet.Msg)
		res.Code, res.Type, res.Msg = msg.GetByCode(int(resUserWallet.Code), h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	if resUserWallet.Data != nil {
		res.Data = nil
		res.Code, res.Type, res.Msg = 22, 1, "El usuario ya ha valido su identidad"
		res.Error = false
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

	idNumber, _ := strconv.ParseInt(userFields.IdentityNumber, 10, 64)

	f, err := srvAuth.SrvFiles.UploadFile(idNumber, userFields.IdentityNumber+".jpg", req.SelfieImg)
	if err != nil {
		logger.Error.Printf("couldn't upload file s3: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(3, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	resUpdateUser, err := clientUser.UpdateUserIdentity(ctx, &users_proto.RqUpdateUserIdentity{
		Id:       user.Id,
		Name:     userFields.Names,
		Lastname: userFields.Surnames,
		IdRole:   21,
	})
	if err != nil || resUpdateUser == nil {
		logger.Error.Printf("couldn't update user by id: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(22, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	if resUpdateUser.Error {
		logger.Error.Printf(resUpdateUser.Msg)
		res.Code, res.Type, res.Msg = msg.GetByCode(int(resUpdateUser.Code), h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	resWallet, err := clientWallet.GetWalletByIdentityNumber(ctx, &wallet_proto.RqGetByIdentityNumber{IdentityNumber: userFields.IdentityNumber})
	if err != nil || resWallet == nil {
		logger.Error.Printf("couldn't get wallet by identity number: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(22, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	if resWallet.Error {
		logger.Error.Printf(resWallet.Msg)
		res.Code, res.Type, res.Msg = msg.GetByCode(int(resWallet.Code), h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	wallet := &WalletIdentity{}

	if resWallet.Data != nil {
		infoWallet, code, err := srvAuth.SrvUsersCredential.GetUsersCredentialByIdentityNumber(userFields.IdentityNumber)
		if err != nil {
			logger.Error.Printf("No se pudo obtener la información de la wallet, error: %s", err.Error())
			res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
			return c.Status(http.StatusAccepted).JSON(res)
		}
		wallet = &WalletIdentity{
			ID:         resWallet.Data.Id,
			Mnemonic:   infoWallet.Mnemonic,
			RsaPublic:  resWallet.Data.Public,
			RsaPrivate: infoWallet.PrivateKey,
		}
	} else {
		newWallet, err := clientWallet.CreateWalletBySystem(ctx, &wallet_proto.RqCreateWalletBySystem{IdentityNumber: userFields.IdentityNumber})
		if err != nil || newWallet == nil {
			logger.Error.Printf("couldn't create wallet: %v", err)
			res.Code, res.Type, res.Msg = msg.GetByCode(5, h.DB, h.TxID)
			return c.Status(http.StatusAccepted).JSON(res)
		}

		if newWallet.Error {
			logger.Error.Printf(newWallet.Msg)
			res.Code, res.Type, res.Msg = msg.GetByCode(int(newWallet.Code), h.DB, h.TxID)
			return c.Status(http.StatusAccepted).JSON(res)
		}

		wallet = &WalletIdentity{
			ID:         newWallet.Data.Id,
			Mnemonic:   newWallet.Data.Mnemonic,
			RsaPublic:  newWallet.Data.Key.Public,
			RsaPrivate: newWallet.Data.Key.Private,
		}

		resCreateAccount, err := clientAccount.CreateAccounting(ctx, &accounting_proto.RequestCreateAccounting{
			Id:       uuid.New().String(),
			IdWallet: newWallet.Data.Id,
			Amount:   0,
			IdUser:   u.ID,
		})
		if err != nil || resCreateAccount == nil {
			logger.Error.Printf("couldn't create account to wallet: %v", err)
			res.Code, res.Type, res.Msg = msg.GetByCode(5, h.DB, h.TxID)
			return c.Status(http.StatusAccepted).JSON(res)
		}

		if resCreateAccount.Error {
			logger.Error.Printf(resCreateAccount.Msg)
			res.Code, res.Type, res.Msg = msg.GetByCode(int(resCreateAccount.Code), h.DB, h.TxID)
			return c.Status(http.StatusAccepted).JSON(res)
		}

	}

	resCreateUserWallet, err := clientUser.CreateUserWallet(ctx, &users_proto.RqCreateUserWallet{
		UserId:   u.ID,
		WalletId: wallet.ID,
	})
	if err != nil || resCreateUserWallet == nil {
		logger.Error.Printf("couldn't create users wallet: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(5, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	if resCreateUserWallet.Error {
		logger.Error.Printf(resCreateUserWallet.Msg)
		res.Code, res.Type, res.Msg = msg.GetByCode(5, h.DB, h.TxID)
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

	code, err = srvAuth.SrvUsersCredential.DeleteUsersCredentialByIdentityNumber(userFields.IdentityNumber)
	if err != nil {
		logger.Error.Printf("couldn't delete user key, error: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		res.Msg = err.Error()
		return c.Status(http.StatusAccepted).JSON(res)
	}

	res.Data = wallet
	res.Code, res.Type, res.Msg = msg.GetByCode(29, h.DB, h.TxID)
	res.Error = false
	return c.Status(http.StatusOK).JSON(res)
}

// getUserByIdentityNumber godoc
// @Summary get user by identity number
// @Description get user by identity number
// @tags User
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Authorization" default(Bearer <Add access token here>)
// @Param inumber path string true "user identity number"
// @Success 200 {object} responseUser
// @Router /api/v1/user/validate-identity-number/{inumber} [get]
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

	resUsr, err := clientUser.ValidIdentityNumber(context.Background(), &users_proto.RequestGetByIdentityNumber{IdentityNumber: usrIdentityNumber})
	if err != nil || resUsr == nil {
		logger.Error.Printf("couldn't get User by identity number: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(22, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	if resUsr.Error {
		logger.Error.Printf(resUsr.Msg)
		res.Data = resUsr.Msg
		res.Code, res.Type, res.Msg = msg.GetByCode(int(resUsr.Code), h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	res.Data = resUsr.Data
	res.Code, res.Type, res.Msg = msg.GetByCode(29, h.DB, h.TxID)
	res.Error = false
	return c.Status(http.StatusOK).JSON(res)
}

// getUserPictureProfile godoc
// @Summary get user by identity number
// @Description get profile picture of user
// @tags User
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Authorization" default(Bearer <Add access token here>)
// @Success 200 {object} responseUpdateUser
// @Router /api/v1/user/picture-profile [get]
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
	if err != nil || resUser == nil {
		logger.Error.Printf("Error trayendo el usuario por su id, error: %s", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(22, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	if resUser.Error {
		logger.Error.Printf(resUser.Msg)
		res.Code, res.Type, res.Msg = msg.GetByCode(22, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	user := resUser.Data

	profile, code, err := srvAuth.SrvUserProfile.GetUserProfileByUserID(user.Id)
	if err != nil {
		logger.Error.Printf("couldn't get profile picture: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	photo := strings.Split(user.FullPathPhoto, "-/")
	fileName := profile.Name
	filePath := profile.Name

	if user.FullPathPhoto != (profile.Path + "-/" + profile.Name) {
		_, code, err = srvAuth.SrvUserProfile.UpdateUserProfile(profile.ID, profile.UserId, photo[0], photo[1])
		if err != nil {
			logger.Error.Printf("couldn't update profile picture: %v", err)
			res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
			return c.Status(http.StatusAccepted).JSON(res)
		}
		fileName = photo[1]
		filePath = photo[0]
	}

	file, code, err := srvAuth.SrvFiles.GetFileByPath(filePath, fileName)
	if err != nil {
		logger.Error.Printf("couldn't get profile picture: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	res.Data = file.Encoding
	res.Code, res.Type, res.Msg = msg.GetByCode(29, h.DB, h.TxID)
	res.Error = false
	return c.Status(http.StatusOK).JSON(res)
}

// changePassword godoc
// @Summary change password of user
// @Description change password of user
// @tags User
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Authorization" default(Bearer <Add access token here>)
// @Param changePassword body requestUpdatePassword true "request of update password"
// @Success 200 {object} responseUpdateUser
// @Router /api/v1/user/update-password [post]
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

	token := c.Get("Authorization")[7:]

	ctx := grpcMetadata.AppendToOutgoingContext(context.Background(), "authorization", token)

	resChangePwd, err := clientUser.ChangePassword(ctx, &users_proto.RequestChangePwd{
		OldPassword:     m.OldPassword,
		NewPassword:     m.NewPassword,
		ConfirmPassword: m.ConfirmPassword,
	})
	if err != nil || resChangePwd == nil {
		logger.Error.Printf("No se pudo actualizar la contraseña, err: %s", err)
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
// @tags User
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Authorization" default(Bearer <Add access token here>)
// @Success 200 {object} responseGetWallets
// @Router /api/v1/user/wallets [get]
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

	tkn := c.Get("Authorization")[7:]
	ctx := grpcMetadata.AppendToOutgoingContext(context.Background(), "authorization", tkn)

	clientWallet := wallet_proto.NewWalletServicesWalletClient(connAuth)

	wt, err := clientWallet.GetWalletByUserId(ctx, &wallet_proto.RequestGetWalletByUserId{UserId: u.ID})
	if err != nil || wt == nil {
		logger.Error.Printf("couldn't get wallets by id: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(22, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	if wt.Error {
		logger.Error.Printf(wt.Msg)
		res.Code, res.Type, res.Msg = msg.GetByCode(95, h.DB, h.TxID)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	if wt.Data == nil {
		res.Code, res.Type, res.Msg = msg.GetByCode(10007, h.DB, h.TxID)
		res.Error = false
		return c.Status(http.StatusAccepted).JSON(res)
	}

	res.Data = &models.Wallet{
		ID:             wt.Data.Id,
		Mnemonic:       wt.Data.Mnemonic,
		RsaPublic:      wt.Data.Public,
		IpDevice:       wt.Data.IpDevice,
		StatusId:       int(wt.Data.StatusId),
		IdentityNumber: wt.Data.IdentityNumber,
	}
	res.Code, res.Type, res.Msg = msg.GetByCode(29, h.DB, h.TxID)
	res.Error = false
	return c.Status(http.StatusOK).JSON(res)
}
