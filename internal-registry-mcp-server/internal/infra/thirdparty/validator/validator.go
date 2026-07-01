package validator

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/BrunoPolaski/go-rest-err/rest_err"
	valueobjects "github.com/BrunoPolaski/internal-registry-mcp-server/internal/core/entities/value_objects"
	"github.com/BrunoPolaski/internal-registry-mcp-server/internal/infra/thirdparty/logger"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translation "github.com/go-playground/validator/v10/translations/en"
	"go.uber.org/zap"
)

var (
	validate   *validator.Validate
	translator ut.Translator
)

func InitValidator() {
	validate = validator.New(validator.WithRequiredStructEnabled())

	enLocale := en.New()
	uni := ut.New(enLocale, enLocale)
	translator, _ = uni.GetTranslator("en")
	en_translation.RegisterDefaultTranslations(validate, translator)

	registerValidation("password", valueobjects.ValidatePassword, "{0} must be a valid password")
	registerValidation("user_type", valueobjects.ValidateUserType, "{0} must be a valid user type")
	registerValidation("phone", valueobjects.ValidatePhoneNumber, "{0} must be a valid phone number")
	registerValidation("document", valueobjects.ValidateDocument, "{0} must be a valid document")

	logger.Info("Custom validator initialized")
}

func registerValidation(tag string, fn func(string) bool, errorMsg string) {
	validate.RegisterValidation(tag, func(fl validator.FieldLevel) bool {
		return fn(fl.Field().String())
	})

	validate.RegisterTranslation(tag, translator,
		func(ut ut.Translator) error {
			return ut.Add(tag, errorMsg, true)
		},
		func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T(tag, fe.Field())
			return t
		},
	)
}

func ShouldBindJSON(r *http.Request, s any) *rest_err.RestErr {
	if err := json.NewDecoder(r.Body).Decode(s); err != nil {
		return rest_err.NewBadRequestError("%s", err.Error()).WithCause(err)
	}

	err := Validate(s)
	if err != nil {
		return err
	}

	return nil
}

func ShouldBindBytes(bytes []byte, s any) *rest_err.RestErr {
	if err := json.Unmarshal(bytes, s); err != nil {
		return rest_err.NewBadRequestError("%s", err.Error()).WithCause(err)
	}

	err := Validate(s)
	if err != nil {
		return err
	}

	return nil
}

func Validate(s any) *rest_err.RestErr {
	defer func() {
		if r := recover(); r != nil {
			var err error
			switch v := r.(type) {
			case string:
				err = errors.New(v)
			case error:
				err = v
			default:
				err = fmt.Errorf("%v", v)
			}

			logger.Error("Panic recovered in validator.Validate",
				err,
				zap.Any("struct", s),
			)
		}
	}()

	logger.Info("Validating struct", zap.Any("struct", s))

	err := validate.Struct(s)
	if err != nil {
		var jsonValidationError validator.ValidationErrors
		if ok := errors.As(err, &jsonValidationError); ok {
			errorCauses := make([]rest_err.Causes, len(jsonValidationError))
			for i, v := range jsonValidationError {
				errorCauses[i] = rest_err.Causes{
					Field:   v.Field(),
					Message: v.Translate(translator),
				}
			}
			return rest_err.NewRestErr(
				"Invalid request body",
				"bad_request",
				http.StatusBadRequest,
				errorCauses,
			)
		}

		return rest_err.NewBadRequestError("%s", err.Error()).WithCause(err)
	}
	return nil
}
