package db

import (
	"context"
	"testing"

	"github.com/aleale2121/interactive-presentation/util"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func createRandomPoll(t *testing.T) Poll {
	presentationID := createRandomPresentation(t)

	arg := CreatePollParams{
		Presentationid: presentationID,
		Question:       util.RandomQuestion(),
		Pollindex:      0,
	}

	poll, err := testQueries.CreatePoll(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, poll)

	require.Equal(t, arg.Presentationid, poll.Presentationid)
	require.Equal(t, arg.Question, poll.Question)
	require.Equal(t, arg.Pollindex, poll.Pollindex)
	require.NotZero(t, poll.ID)

	return poll
}

func createPresentationPoll(t *testing.T, presenationID uuid.UUID) Poll {

	arg := CreatePollParams{
		Presentationid: presenationID,
		Question:       util.RandomQuestion(),
	}

	poll, err := testQueries.CreatePoll(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, poll)

	require.Equal(t, arg.Presentationid, poll.Presentationid)
	require.Equal(t, arg.Question, poll.Question)
	require.Equal(t, arg.Pollindex, poll.Pollindex)
	require.NotZero(t, poll.ID)

	return poll
}

func TestCreatePoll(t *testing.T) {
	createRandomPoll(t)
}

func TestGetPoll(t *testing.T) {
	poll1 := createRandomPoll(t)
	poll2, err := testQueries.GetPoll(context.Background(), poll1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, poll2)

	require.Equal(t, poll1.ID, poll2.ID)
	require.Equal(t, poll1.Presentationid, poll2.Presentationid)
	require.Equal(t, poll1.Question, poll2.Question)
}

func TestGetPresentationPolls(t *testing.T) {
	presentationID := createRandomPresentation(t)

	for i := 0; i < 10; i++ {
		createPresentationPoll(t, presentationID)
	}

	polls, err := testQueries.GetPresentationPolls(context.Background(), presentationID)
	require.NoError(t, err)
	require.NotEmpty(t, polls)
	require.Len(t, polls, 10)
}
