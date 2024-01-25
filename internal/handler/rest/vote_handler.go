package rest

import (
	"context"
	"net/http"

	"github.com/aleale2121/interactive-presentation/internal/constant/model"
	db "github.com/aleale2121/interactive-presentation/internal/storage/persistence"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type VoteHandler interface {
	CreateVoteHandler(c *gin.Context)
	GetPollVotesHandler(c *gin.Context)
}

type voteHandler struct {
	logger *logrus.Logger
	store  db.Store
}

func NewVoteHandler(logger *logrus.Logger,
	store db.Store) VoteHandler {
	return voteHandler{
		logger: logger,
		store:  store,
	}
}

// CreateVoteHandler handles the HTTP request for creating a vote for the current poll.
func (server voteHandler) CreateVoteHandler(c *gin.Context) {
	presentationID, err := uuid.Parse(c.Param("presentation_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var vote model.CreateVoteRequest
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
