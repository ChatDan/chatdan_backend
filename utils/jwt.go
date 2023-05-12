package utils

import (
	"ChatDanBackend/config"
	"fmt"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
	"time"
)

type UserClaims struct {
	UserID  int    `json:"user_id"`
	IsAdmin bool   `json:"is_admin"`
	Key     string `json:"key"`
	jwt.RegisteredClaims
}

type ApisixConsumerJwtAuth struct {
	Key    string `json:"key"`
	Secret string `json:"secret"`
}

type ApisixConsumerPlugins struct {
	JWTAuth ApisixConsumerJwtAuth `json:"jwt-auth"`
}

type ApisixConsumer struct {
	Username string                `json:"username"`
	Plugins  ApisixConsumerPlugins `json:"plugins"`
}

var (
	client = fasthttp.Client{}
)

var (
	ErrGetConsumer    = fmt.Errorf("get consumer failed")
	ErrCreateConsumer = fmt.Errorf("create consumer failed")
)

func getConsumerUsername(userID int) string {
	return fmt.Sprintf("chatdan_user_%d", userID)
}

func GetConsumer(userID int) (*ApisixConsumer, error) {
	consumerUsername := getConsumerUsername(userID)
	var responseStruct struct {
		Key   string         `json:"key"`
		Value ApisixConsumer `json:"value"`
	}

	// get request
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	req.SetRequestURI(config.Config.ApisixUrl + "/apisix/admin/consumers/" + consumerUsername)
	req.Header.Set("X-API-KEY", config.Config.ApisixAdminKey)
	req.Header.SetMethod(fasthttp.MethodGet)

	// get response
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	// execute request
	err := client.Do(req, resp)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != fiber.StatusOK {
		if resp.StatusCode() == fiber.StatusNotFound {
			return CreateConsumer(userID)
		} else {
			Logger.Error("get consumer failed",
				zap.Int("status_code", resp.StatusCode()),
				zap.ByteString("body", resp.Body()),
			)
			return nil, ErrGetConsumer
		}
	}

	data := resp.Body()
	err = json.Unmarshal(data, &responseStruct)
	if err != nil {
		Logger.Error("unmarshal consumer failed",
			zap.ByteString("body", data),
		)
		return nil, err
	}

	return &responseStruct.Value, nil
}

func CreateConsumer(userID int) (*ApisixConsumer, error) {
	consumerUsername := getConsumerUsername(userID)
	consumerSecret := SecretGenerator(32)
	consumer := ApisixConsumer{
		Username: consumerUsername,
		Plugins: ApisixConsumerPlugins{JWTAuth: ApisixConsumerJwtAuth{
			Key:    consumerUsername,
			Secret: consumerSecret,
		}},
	}

	// get request
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	req.SetRequestURI(config.Config.ApisixUrl + "/apisix/admin/consumers")
	req.Header.Set("X-API-KEY", config.Config.ApisixAdminKey)
	req.Header.SetContentType(fiber.MIMEApplicationJSON)
	req.Header.SetMethod(fasthttp.MethodPut)
	data, _ := json.Marshal(&consumer)
	req.SetBody(data)

	// get response
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	// execute request
	err := client.Do(req, resp)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() == fasthttp.StatusOK || resp.StatusCode() == fasthttp.StatusCreated {
		return &consumer, nil
	} else {
		Logger.Error("create consumer failed",
			zap.String("consumer", consumerUsername),
			zap.Int("status_code", resp.StatusCode()),
			zap.ByteString("body", resp.Body()),
		)

		return nil, ErrCreateConsumer
	}
}

func CreateJwtTokenFromApisix(claims UserClaims) (string, error) {
	consumer, err := GetConsumer(claims.UserID)
	if err != nil {
		return "", err
	}

	if claims.ExpiresAt == nil {
		claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(24 * time.Hour))
	}
	claims.Key = getConsumerUsername(claims.UserID)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(consumer.Plugins.JWTAuth.Secret))
}

func CreateJwtTokenStandalone(claims UserClaims, secret []byte) (string, error) {
	if claims.ExpiresAt == nil {
		claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(24 * time.Hour))
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}
