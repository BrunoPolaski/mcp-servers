package httphelper

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/BrunoPolaski/go-rest-err/rest_err"
	"github.com/BrunoPolaski/registration-validation-mcp-server/internal/infra/thirdparty/logger"
	"go.uber.org/zap"
)

func ErrorResponse(err *rest_err.RestErr, w http.ResponseWriter) {
	if err.Code >= 500 {
		logger.Error(
			"Internal Server Error",
			err,
			zap.String("message", err.Message),
			zap.Int("status", err.Code),
			zap.Any("causes", err.Causes),
		)

		err.Message = http.StatusText(err.Code)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(err.Code)
	encoder := json.NewEncoder(w)

	encoder.Encode(err)
}

func JSON(data any, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	encoder.Encode(data)
}

func JSONWithStatus(data any, status int, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	encoder := json.NewEncoder(w)
	encoder.Encode(data)
}

type IntQueryParamOptions struct {
	Key          string
	DefaultValue int
	Required     bool
}

func IntQueryParam(r *http.Request, options IntQueryParamOptions) (int, *rest_err.RestErr) {
	values := r.URL.Query()[options.Key]
	if len(values) > 0 {
		if intValue, err := strconv.Atoi(values[0]); err == nil {
			return intValue, nil
		}

		return 0, rest_err.NewBadRequestError("invalid %s parameter", options.Key)
	}

	if options.Required {
		return 0, rest_err.NewBadRequestError("%s parameter is required", options.Key)
	}

	return options.DefaultValue, nil
}

func PathParam(r *http.Request, key string) (string, *rest_err.RestErr) {
	param := r.PathValue(key)

	if param == "" {
		return "", rest_err.NewBadRequestError("%s parameter is required", key)
	}

	return param, nil
}

func IntPathParam(r *http.Request, key string) (uint, *rest_err.RestErr) {
	param := r.PathValue(key)

	if param == "" {
		return 0, rest_err.NewBadRequestError("%s parameter is required", key)
	}

	intParam, err := strconv.Atoi(param)
	if err != nil {
		return 0, rest_err.NewBadRequestError("invalid %s parameter", key)
	}

	return uint(intParam), nil
}
