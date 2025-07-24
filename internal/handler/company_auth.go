package handler

import (
	"go.uber.org/zap"
	"net/http"

	"github.com/gin-gonic/gin"
	"go-code-runner-microservice/api-gateway/internal/logger"
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

	log := logger.WithContext(c.Request.Context())

	var req model.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warn("invalid register request",
			zap.Error(err),
			zap.String("client_ip", c.ClientIP()),
		)
		c.JSON(http.StatusBadRequest, model.RegisterResponse{
			Success: false,
			Error:   "Invalid request format",
		})
		return
	}

	// Log registration attempt (without password)
	log.Info("company registration attempt",
		zap.String("company_name", req.Name),
		zap.String("email", req.Email),
	)

	resp, err := h.client.Register(c.Request.Context(), req.Name, req.Email, req.Password)
	if err != nil {
		log.Error("failed to register company",
			zap.Error(err),
			zap.String("company_name", req.Name),
			zap.String("email", req.Email),
		)
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
		log.Warn("company registration failed",
			zap.String("reason", errorMsg),
			zap.String("email", req.Email),
		)
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

		log.Info("company registered successfully",
			zap.Int("company_id", company.ID),
			zap.String("company_name", company.Name),
			zap.String("email", company.Email),
		)
	}

	c.JSON(http.StatusOK, model.RegisterResponse{
		Success: true,
		Company: company,
	})
}

func (h *CompanyHandler) Login(c *gin.Context) {
	log := logger.WithContext(c.Request.Context())

	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warn("invalid login request",
			zap.Error(err),
			zap.String("client_ip", c.ClientIP()),
		)
		c.JSON(http.StatusBadRequest, model.LoginResponse{
			Success: false,
			Error:   "Invalid request: " + err.Error(),
		})
		return
	}

	log.Info("company login attempt",
		zap.String("email", req.Email),
		zap.String("client_ip", c.ClientIP()),
	)

	resp, err := h.client.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		log.Error("failed to process login",
			zap.Error(err),
			zap.String("email", req.Email),
		)
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
		log.Warn("company login failed",
			zap.String("reason", errorMsg),
			zap.String("email", req.Email),
			zap.String("client_ip", c.ClientIP()),
		)
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
		log.Info("company login successful",
			zap.Int("company_id", company.ID),
			zap.String("company_name", company.Name),
			zap.String("email", company.Email),
		)
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
	log := logger.WithContext(c.Request.Context())

	var req model.GenerateAPIKeyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warn("invalid generate API key request", zap.Error(err))

		c.JSON(http.StatusBadRequest, model.GenerateAPIKeyResponse{
			Success: false,
			Error:   "Invalid request: " + err.Error(),
		})
		return
	}

	log.Info("generating API key",
		zap.Int("company_id", req.CompanyID),
	)

	resp, err := h.client.GenerateAPIKey(c.Request.Context(), int32(req.CompanyID))
	if err != nil {
		log.Error("failed to generate API key",
			zap.Error(err),
			zap.Int("company_id", req.CompanyID),
		)
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
		log.Warn("API key generation failed",
			zap.String("reason", errorMsg),
			zap.Int("company_id", req.CompanyID),
		)
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

	log.Info("API key generated successfully",
		zap.Int("company_id", req.CompanyID),
	)

	c.JSON(http.StatusOK, model.GenerateAPIKeyResponse{
		Success: true,
		APIKey:  apiKey,
	})
}

func (h *CompanyHandler) GenerateClientID(c *gin.Context) {
	log := logger.WithContext(c.Request.Context())

	var req model.GenerateClientIDRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warn("invalid generate client ID request", zap.Error(err))

		c.JSON(http.StatusBadRequest, model.GenerateClientIDResponse{
			Success: false,
			Error:   "Invalid request: " + err.Error(),
		})
		return
	}

	log.Info("generating client ID",
		zap.Int("company_id", req.CompanyID),
	)

	resp, err := h.client.GenerateClientID(c.Request.Context(), int32(req.CompanyID))
	if err != nil {
		log.Error("failed to generate client ID",
			zap.Error(err),
			zap.Int("company_id", req.CompanyID),
		)
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
		log.Warn("client ID generation failed",
			zap.String("reason", errorMsg),
			zap.Int("company_id", req.CompanyID),
		)
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
	log.Info("client ID generated successfully",
		zap.Int("company_id", req.CompanyID),
	)

	c.JSON(http.StatusOK, model.GenerateClientIDResponse{
		Success:  true,
		ClientID: clientID,
	})
}
