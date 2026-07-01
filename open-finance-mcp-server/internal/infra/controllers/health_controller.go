package controllers

import (
	"net/http"

	"github.com/BrunoPolaski/open-finance-mcp-server/internal/infra/controllers/dto"
	httphelper "github.com/BrunoPolaski/open-finance-mcp-server/internal/interfaces/http"
	"gorm.io/gorm"
)

type HealthController struct {
	database *gorm.DB
}

func NewHealthController(database *gorm.DB) *HealthController {
	return &HealthController{
		database: database,
	}
}

// @Summary Check API health
// @Description Checks the health of the API and its dependencies.
// @Tags /health
// @Produce json
// @Success 200 {object} dto.HealthResponse
// @Failure 503 {object} dto.HealthResponse
// @Router /health [get]
func (hc *HealthController) Check(w http.ResponseWriter, r *http.Request) {
	response := dto.HealthResponse{
		APIStatus:      "up",
		DatabaseStatus: "up",
	}

	database := hc.database.Exec("SELECT 1")
	if database.Error != nil {
		response.DatabaseStatus = "down"
	}

	if response.DatabaseStatus == "down" || response.CacheStatus == "down" {
		w.WriteHeader(http.StatusServiceUnavailable)
	}

	httphelper.JSON(response, w)
}
