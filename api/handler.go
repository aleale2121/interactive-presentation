package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	db "github.com/aleale2121/interactive-presentation/db/sqlc"
	"github.com/aleale2121/interactive-presentation/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Server serves HTTP requests for our banking service.
type Server struct {
	store  db.Store
	router *gin.Engine
}

// NewServer creates a new HTTP server and set up routing.
func NewServer(store db.Store) (*Server, error) {

	server := &Server{
		store: store,
	}

	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()

	router.POST("/presentations", server.CreatePresentationAndPolls)

	router.GET("/presentations/:presentation_id/polls/current", server.GetCurrentPoll)

	router.PUT("/presentations/:presentation_id/polls/current/forward", server.SlidePresentationPollToForwardHandler)
	router.PUT("/presentations/:presentation_id/polls/current/backward", server.SlidePresentationPollToPreviousHandler)

	router.POST("/presentations/:presentation_id/polls/current/votes", server.Vote)
	router.GET("/presentations/:presentation_id/polls/:poll_id/votes", server.PollVotes)

	router.GET("/ping", server.HealthCheck)

	server.router = router
}

// Start runs the HTTP server on a specific address.
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func (server *Server) CreatePresentation(c *gin.Context) {
	var presenation models.Presention
	if err := c.ShouldBindJSON(&presenation); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if len(presenation.Polls) == 0 {
		c.JSON(http.StatusBadRequest, "no polls found")
		return
	}
	presID, err := server.store.CreatePresentation(context.Background(), sql.NullInt32{
		Int32: 0,
		Valid: true,
	},
	)

	if err != nil {
		fmt.Println("1--> ", err)
		c.AbortWithStatus(400)
		return
	}

	for i, poll := range presenation.Polls {
		createdPoll, err := server.store.CreatePoll(context.Background(), db.CreatePollParams{
			Presentationid: presID,
			Question:       poll.Question,
			Pollindex:      int32(i),
		})
		if err != nil {
			c.AbortWithStatus(400)
			return
		}
		for _, opt := range poll.Options {
			err = server.store.CreateOption(context.Background(), db.CreateOptionParams{
				Pollid:      createdPoll.ID,
				Optionkey:   opt.Key,
				Optionvalue: opt.Value,
			})

			if err != nil {
				c.AbortWithStatus(400)
				return
			}
		}
	}

	c.JSON(http.StatusCreated, gin.H{
		"presentation_id": presID,
	})
}

func (server *Server) CreatePresentationAndPolls(c *gin.Context) {
	var presenation models.Presention
	if err := c.ShouldBindJSON(&presenation); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if len(presenation.Polls) == 0 {
		c.JSON(http.StatusBadRequest, "no polls found")
		return
	}
	jsonb, err := json.Marshal(presenation.Polls)
	// fmt.Println(json.RawMessage(jsonb))

	if err != nil {
		fmt.Println("1-->", err)
		c.AbortWithStatus(400)
		return
	}
	presID, err := server.store.CreatePresentationAndPolls(context.Background(), json.RawMessage(jsonb))

	if err != nil {
		fmt.Println("2-->", err)
		c.AbortWithStatus(400)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"presentation_id": presID,
	})
}

func (server *Server) SlidePresentationPollToForwardHandler(c *gin.Context) {
	presentationID, err := uuid.Parse(c.Param("presentation_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := server.store.UpdateCurrentPollToForwardTx(context.Background(), presentationID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var response models.PollResponse
	response.ID = result.ID
	response.Question = result.Question

	for _, opt := range result.Options {
		response.Options = append(response.Options, models.OptionResponse{
			OptionKey:   opt.Optionkey,
			OptionValue: opt.Optionvalue,
		})
	}
	c.JSON(http.StatusOK, response)
}

func (server *Server) SlidePresentationPollToPreviousHandler(c *gin.Context) {
	presentationID, err := uuid.Parse(c.Param("presentation_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := server.store.UpdateCurrentPollToPreviousTx(context.Background(), presentationID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var response models.PollResponse
	response.ID = result.ID
	response.Question = result.Question

	for _, opt := range result.Options {
		response.Options = append(response.Options, models.OptionResponse{
			OptionKey:   opt.Optionkey,
			OptionValue: opt.Optionvalue,
		})
	}
	c.JSON(http.StatusOK, response)
}

func (server *Server) GetCurrentPoll(c *gin.Context) {
	presentationID, err := uuid.Parse(c.Param("presentation_id"))
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	result, err := server.store.GetCurrentPoll(context.Background(), presentationID)

	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	var response models.PollResponse
	response.ID = result.ID
	response.Question = result.Question

	for _, opt := range result.Options {
		response.Options = append(response.Options, models.OptionResponse{
			OptionKey:   opt.Optionkey,
			OptionValue: opt.Optionvalue,
		})
	}
	c.JSON(http.StatusOK, response)
}

func (server *Server) Vote(c *gin.Context) {
	presentationID := c.Param("presentation_id")

	var vote models.CreateVoteRequest
	if err := c.ShouldBindJSON(&vote); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := server.store.VoteCurrentPollTx(context.Background(), db.VoteParams{
		PresentationID: presentationID,
		Pollid:         vote.PollID,
		Optionkey:      vote.Key,
		Clientid:       vote.ClientId,
	})

	if err != nil {
		c.AbortWithStatus(409)
		return
	}

	c.JSON(http.StatusNoContent, "")
}

func (server *Server) PollVotes(c *gin.Context) {
	presentationID := c.Param("presentation_id")
	pollID := c.Param("poll_id")
	polID, err := uuid.Parse(pollID)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	presID, err := uuid.Parse(presentationID)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	result, err := server.store.GetPollVotes(context.Background(), db.GetPollVotesParams{
		ID:   presID,
		ID_2: polID,
	})

	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	response := make([]models.VoteResponse, 0)
	for _, res := range result {
		response = append(response, models.VoteResponse{
			Key:      res.Optionkey,
			ClientId: res.Clientid,
		})
	}
	if err != nil {
		c.AbortWithStatus(409)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (server *Server) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"description": "The service is up and running",
	})
}
