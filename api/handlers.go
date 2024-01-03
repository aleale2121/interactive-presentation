package api

import (
	"context"
	"encoding/json"
	"net/http"

	db "github.com/aleale2121/interactive-presentation/db/sqlc"
	"github.com/aleale2121/interactive-presentation/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CreatePresentationHandler handles the HTTP request for creating a new presentation.
func (server *Server) CreatePresentationHandler(c *gin.Context) {
	var presenation models.CreatePresentionRequest
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

// UpdateCurrentPollHandler handles the HTTP request for updating the current poll index to move forward.
func (server *Server) UpdateCurrentPollHandler(c *gin.Context) {
	presentationID, err := uuid.Parse(c.Param("presentation_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	_, err = server.store.GetPresentation(context.Background(), presentationID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No presentation found"})
		return
	}
	currentPoll, err := server.store.UpdateCurrentPollToForwardTx(context.Background(), presentationID)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "The presentation ran out of polls."})
		return
	}

	c.JSON(http.StatusOK, currentPoll)
}

// GetCurrentPollHandler handles the HTTP request for retrieving the current poll of a presentation.
func (server *Server) GetCurrentPollHandler(c *gin.Context) {
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

	currentPoll, err := server.store.GetCurrentPoll(context.Background(), presentationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, currentPoll)
}

// GetPresentationHandler handles the HTTP request for retrieving the current poll of a presentation.
func (server *Server) GetPresentationHandler(c *gin.Context) {
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
		CurrentPollIndex int32                            `json:"current_poll_index"`
		Polls            []db.ListPollsRow `json:"polls"`
	}{
		CurrentPollIndex: presentation.Currentpollindex.Int32,
		Polls:            polls,
	})
}

// CreateVoteHandler handles the HTTP request for creating a vote for the current poll.
func (server *Server) CreateVoteHandler(c *gin.Context) {
	presentationID, err := uuid.Parse(c.Param("presentation_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var vote models.CreateVoteRequest
	if err := c.ShouldBindJSON(&vote); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = server.store.VoteCurrentPollTx(context.Background(), db.VoteParams{
		PresentationID: presentationID,
		Pollid:         vote.PollID,
		Optionkey:      vote.Key,
		Clientid:       vote.ClientId,
	})
	if err != nil {
		switch err {
		case models.ErrNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "Either `presentation_id` or `poll_id` not found"})
			return
		case models.ErrConflict:
			c.JSON(http.StatusBadRequest, gin.H{"error": "The `poll_id` in the request body doesn't match the current poll."})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
	}

	c.JSON(http.StatusNoContent, "")
}

// GetPollVotesHandler handles the HTTP request for retrieving votes for a specific poll.
func (server *Server) GetPollVotesHandler(c *gin.Context) {
	presentationID, err := uuid.Parse(c.Param("presentation_id"))
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	pollID, err := uuid.Parse(c.Param("poll_id"))
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	result, err := server.store.GetPresentationAndPoll(context.Background(), db.GetPresentationAndPollParams{
		PresentationID: presentationID,
		PollID:         pollID,
	})

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Either `presentation_id` or `poll_id` not found"})
		return
	}
	if result.Currentpollindex.Int32 != result.Pollindex {
		c.JSON(http.StatusBadRequest, gin.H{"error": "The `poll_id` in the request body doesn't match the current poll."})
		return
	}
	votes, err := server.store.GetVotes(context.Background(), db.GetVotesParams{
		PresentationID: presentationID,
		PollID:         pollID,
	})

	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	c.JSON(http.StatusOK, votes)
}

// HealthCheck handles the HTTP request for health checking the service.
func (server *Server) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"description": "The service is up and running",
	})
}
