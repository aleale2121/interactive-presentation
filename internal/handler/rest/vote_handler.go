package rest

import (
	"context"
	"net/http"

	"github.com/aleale2121/interactive-presentation/internal/constant/model"
	"github.com/aleale2121/interactive-presentation/internal/module/vote"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type VoteHandler interface {
	CreateVoteHandler(c *gin.Context)
	GetPollVotesHandler(c *gin.Context)
}

type voteHandler struct {
	logger  *logrus.Logger
	useCase vote.Usecase
}

func NewVoteHandler(logger *logrus.Logger,
	useCase vote.Usecase) VoteHandler {
	return voteHandler{
		logger:  logger,
		useCase: useCase,
	}
}

// CreateVoteHandler handles the HTTP request for creating a vote for the current poll.
func (server voteHandler) CreateVoteHandler(c *gin.Context) {
	presentationID, err := uuid.Parse(c.Param("presentation_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var voteParam model.CreateVoteRequestDTO
	if err := c.ShouldBindJSON(&voteParam); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = server.useCase.CreateVote(context.Background(), presentationID, voteParam)
	if err != nil {
		switch err {
		case model.ErrNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "Either `presentation_id` or `poll_id` not found"})
			return
		case model.ErrConflict:
			c.JSON(http.StatusBadRequest, gin.H{"error": "The `poll_id` in the request body doesn't match the current poll."})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
	}

	c.JSON(http.StatusNoContent, "")
}

// GetPollVotesHandler handles the HTTP request for retrieving votes for a specific poll.
func (server voteHandler) GetPollVotesHandler(c *gin.Context) {
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
	//TODO: error handle
	votes, err := server.useCase.GetPollVotes(context.Background(), presentationID, pollID)

	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	c.JSON(http.StatusOK, votes)
}
