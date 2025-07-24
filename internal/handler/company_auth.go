package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go-code-runner-microservice/api-gateway/internal/model"
	"go-code-runner-microservice/api-gateway/internal/service/grpc/company_auth"
)

type CompanyHandler struct {
	client *company_auth.Client
}

func NewCompanyHandler(client *company_auth.Client) *CompanyHandler {
	return &CompanyHandler{
		client: client,
	}
}

func (h *CompanyHandler) Register(c *gin.Context) {
	var req model.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.RegisterResponse{
			Success: false,
			Error:   "Invalid request: " + err.Error(),
		})
		return
	}

	resp, err := h.client.Register(c.Request.Context(), req.Name, req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.RegisterResponse{
			Success: false,
			Error:   "Failed to register company: " + err.Error(),
		})
		return
	}

	if !resp.Success {
		errorMsg := ""
		if resp.Error != nil {
			errorMsg = *resp.Error
		}
		c.JSON(http.StatusBadRequest, model.RegisterResponse{
			Success: false,
			Error:   errorMsg,
		})
		return
	}

	var company *model.Company
	if resp.Company != nil {
		company = &model.Company{
			ID:    int(resp.Company.Id),
			Name:  resp.Company.Name,
			Email: resp.Company.Email,
		}
		if resp.Company.ApiKey != nil {
			company.APIKey = resp.Company.ApiKey
		}
		if resp.Company.ClientId != nil {
			company.ClientID = resp.Company.ClientId
		}
	}

	c.JSON(http.StatusOK, model.RegisterResponse{
		Success: true,
		Company: company,
	})
}

func (h *CompanyHandler) Login(c *gin.Context) {
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.LoginResponse{
			Success: false,
			Error:   "Invalid request: " + err.Error(),
		})
		return
	}

	resp, err := h.client.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.LoginResponse{
			Success: false,
			Error:   "Failed to login: " + err.Error(),
		})
		return
	}

	if !resp.Success {
		errorMsg := ""
		if resp.Error != nil {
			errorMsg = *resp.Error
		}
		c.JSON(http.StatusUnauthorized, model.LoginResponse{
			Success: false,
			Error:   errorMsg,
		})
		return
	}

	var company *model.Company
	if resp.Company != nil {
		company = &model.Company{
			ID:    int(resp.Company.Id),
			Name:  resp.Company.Name,
			Email: resp.Company.Email,
		}
		if resp.Company.ApiKey != nil {
			company.APIKey = resp.Company.ApiKey
		}
		if resp.Company.ClientId != nil {
			company.ClientID = resp.Company.ClientId
		}
	}

	token := ""
	if resp.Token != nil {
		token = *resp.Token
	}

	c.JSON(http.StatusOK, model.LoginResponse{
		Success: true,
		Company: company,
		Token:   token,
	})
}

func (h *CompanyHandler) GenerateAPIKey(c *gin.Context) {
	var req model.GenerateAPIKeyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.GenerateAPIKeyResponse{
			Success: false,
			Error:   "Invalid request: " + err.Error(),
		})
		return
	}

	resp, err := h.client.GenerateAPIKey(c.Request.Context(), int32(req.CompanyID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.GenerateAPIKeyResponse{
			Success: false,
			Error:   "Failed to generate API key: " + err.Error(),
		})
		return
	}

	if !resp.Success {
		errorMsg := ""
		if resp.Error != nil {
			errorMsg = *resp.Error
		}
		c.JSON(http.StatusBadRequest, model.GenerateAPIKeyResponse{
			Success: false,
			Error:   errorMsg,
		})
		return
	}

	apiKey := ""
	if resp.ApiKey != nil {
		apiKey = *resp.ApiKey
	}

	c.JSON(http.StatusOK, model.GenerateAPIKeyResponse{
		Success: true,
		APIKey:  apiKey,
	})
}

func (h *CompanyHandler) GenerateClientID(c *gin.Context) {
	var req model.GenerateClientIDRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.GenerateClientIDResponse{
			Success: false,
			Error:   "Invalid request: " + err.Error(),
		})
		return
	}

	resp, err := h.client.GenerateClientID(c.Request.Context(), int32(req.CompanyID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.GenerateClientIDResponse{
			Success: false,
			Error:   "Failed to generate client ID: " + err.Error(),
		})
		return
	}

	if !resp.Success {
		errorMsg := ""
		if resp.Error != nil {
			errorMsg = *resp.Error
		}
		c.JSON(http.StatusBadRequest, model.GenerateClientIDResponse{
			Success: false,
			Error:   errorMsg,
		})
		return
	}

	clientID := ""
	if resp.ClientId != nil {
		clientID = *resp.ClientId
	}

	c.JSON(http.StatusOK, model.GenerateClientIDResponse{
		Success:  true,
		ClientID: clientID,
	})
}
