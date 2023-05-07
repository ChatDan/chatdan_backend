package utils

import (
	"ChatDanBackend/common/loggerx"
	"ChatDanBackend/common/utilsx"
	"ChatDanBackend/service/user/config"
	"ChatDanBackend/service/user/model"
	"fmt"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
	"time"
)

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
	var consumer ApisixConsumer

	// get request
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	req.SetRequestURI(config.CustomConfig.ApisixUrl + "/apisix/admin/consumers/" + consumerUsername)
	req.Header.Set("X-API-KEY", config.CustomConfig.ApisixAdminKey)
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
			loggerx.Logger.Error("get consumer failed",
				zap.Int("status_code", resp.StatusCode()),
				zap.ByteString("body", resp.Body()),
			)
			return nil, ErrGetConsumer
		}
	}

	data := resp.Body()
	err = json.Unmarshal(data, &consumer)
	if err != nil {
		loggerx.Logger.Error("unmarshal consumer failed",
			zap.ByteString("body", data),
		)
		return nil, err
	}

	return &consumer, nil
}

func CreateConsumer(userID int) (*ApisixConsumer, error) {
	consumerUsername := getConsumerUsername(userID)
	consumerSecret := utilsx.SecretGenerator(32)
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
	req.SetRequestURI(config.CustomConfig.ApisixUrl + "/apisix/admin/consumers")
	req.Header.Set("X-API-KEY", config.CustomConfig.ApisixAdminKey)
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
		loggerx.Logger.Error("create consumer failed",
			zap.String("consumer", consumerUsername),
			zap.Int("status_code", resp.StatusCode()),
			zap.ByteString("body", resp.Body()),
		)

		return nil, ErrCreateConsumer
	}
}

func CreateJWT(user *model.User) (string, error) {
	consumer, err := GetConsumer(user.ID)
	if err != nil {
		return "", err
	}
	claims := utilsx.UserClaims{
		UserID:  user.ID,
		IsAdmin: user.IsAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	}

	return utilsx.CreateJwtToken(claims, consumer.Plugins.JWTAuth.Secret)
}
