package rest

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/aleale2121/interactive-presentation/internal/constant/model"
	db "github.com/aleale2121/interactive-presentation/internal/storage/persistence"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type PresentationHandler interface {
	CreatePresentationHandler(c *gin.Context)
	GetPresentationHandler(c *gin.Context)
}

type presentationHandler struct {
	logger *logrus.Logger
	store  db.Store
}

func NewPresentationHandler(logger *logrus.Logger,
	store db.Store) PresentationHandler {
	return presentationHandler{
		logger: logger,
		store:  store,
	}
}

// CreatePresentationHandler handles the HTTP request for creating a new presentation.
func (server presentationHandler) CreatePresentationHandler(c *gin.Context) {
	var presenation model.CreatePresentionRequest
	if err := c.ShouldBindJSON(&presenation); err != nil {
		c.JSON(http.StatusBadRequest, "Mandatory body parameters missing or have incorrect type")
		return
	}

	if len(presenation.Polls) == 0 {
		c.JSON(http.StatusBadRequest, "no polls found")
		return
	}

	jsonb, err := json.Marshal(presenation.Polls)
	if err != nil {
		c.AbortWithStatus(400)
		return
	}

	presID, err := server.store.CreatePresentationAndPolls(context.Background(), []byte(jsonb))
	if err != nil {
		c.AbortWithStatus(400)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"presentation_id": presID,
	})
}

// GetPresentationHandler handles the HTTP request for retrieving the current poll of a presentation.
func (server presentationHandler) GetPresentationHandler(c *gin.Context) {
	presentationID, err := uuid.Parse(c.Param("presentation_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	presentation, err := server.store.GetPresentation(context.Background(), presentationID)
	if err != nil {
		c.JSON(http.StatusNotFound, "There is no presentation with the provided `presentation_id`")
		return
	}
	if presentation.Currentpollindex.Int32 == 0 {
		c.JSON(http.StatusConflict, "There are no polls currently displayed")
		return
	}

	polls, err := server.store.ListPolls(context.Background(), presentationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, struct {
		CurrentPollIndex int32             `json:"current_poll_index"`
		Polls            []db.ListPollsRow `json:"polls"`
	}{
		CurrentPollIndex: presentation.Currentpollindex.Int32,
		Polls:            polls,
	})
}
