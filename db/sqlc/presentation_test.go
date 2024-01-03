package db

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func createRandomPresentationWithPolls(t *testing.T, len int) uuid.UUID {
	presentationID, err := testQueries.CreatePresentationAndPolls(context.Background(), []byte(presentationData[len]))
	require.NoError(t, err)
	require.NotEmpty(t, presentationID)

	return presentationID
}

func TestCreatePresentationAndPolls(t *testing.T) {
	createRandomPresentationWithPolls(t, 2)
}

func TestGetPresentation(t *testing.T) {
	createdPresenationID := createRandomPresentationWithPolls(t, 2)
	retrievedPresentation, err := testQueries.GetPresentation(context.Background(), createdPresenationID)
	require.NoError(t, err)
	require.NotEmpty(t, retrievedPresentation)

	require.Equal(t, createdPresenationID, retrievedPresentation.ID)
}

func TestGetPresentationCurrentPoll(t *testing.T) {
	createdPresenationID := createRandomPresentationWithPolls(t, 2)
	poll, err := testQueries.GetPresentationCurrentPoll(context.Background(), createdPresenationID)
	require.Error(t, err)
	require.Empty(t, poll)

	nextPoll, err := testQueries.GetNextPoll(context.Background(), createdPresenationID)
	require.NoError(t, err)
	require.NotEmpty(t, nextPoll)

	poll, err = testQueries.GetPresentationCurrentPoll(context.Background(), createdPresenationID)
	require.NoError(t, err)
	require.NotEmpty(t, poll)
	require.Equal(t, poll.ID, nextPoll.ID)
}

func TestListPolls(t *testing.T) {
	createdPresenationID := createRandomPresentationWithPolls(t, 8)
	polls, err := testQueries.ListPolls(context.Background(), createdPresenationID)
	require.NoError(t, err)
	require.NotEmpty(t, polls)
	require.Len(t, polls, 8)
}

func TestGetPresentationAndPoll(t *testing.T) {
	createdPresenationID := createRandomPresentationWithPolls(t, 8)
	polls, err := testQueries.ListPolls(context.Background(), createdPresenationID)
	require.NoError(t, err)
	require.NotEmpty(t, polls)
	require.Len(t, polls, 8)

	result, err := testQueries.GetPresentationAndPoll(context.Background(), GetPresentationAndPollParams{
		PresentationID: createdPresenationID,
		PollID:         polls[0].ID,
	})
	require.NoError(t, err)
	require.Equal(t, createdPresenationID, result.PresentationID)
	require.Equal(t, int32(0), result.Currentpollindex.Int32)
	require.Equal(t, polls[0].ID, result.PollID)

}

func TestGetNextPoll(t *testing.T) {
	createdPresenationID := createRandomPresentationWithPolls(t, 2)
	polls, err := testQueries.ListPolls(context.Background(), createdPresenationID)
	require.NoError(t, err)
	require.NotEmpty(t, polls)
	require.Len(t, polls, 2)

	poll, err := testQueries.GetPresentationCurrentPoll(context.Background(), createdPresenationID)
	require.Error(t, err)
	require.Empty(t, poll)

	for i := 0; i < 2; i++ {
		nextPoll, err := testQueries.GetNextPoll(context.Background(), createdPresenationID)
		require.NoError(t, err)
		require.NotEmpty(t, nextPoll)
		require.Equal(t, polls[i].ID.String(), nextPoll.ID.String())

		currPoll, err := testQueries.GetPresentationCurrentPoll(context.Background(), createdPresenationID)
		require.NoError(t, err)
		require.NotEmpty(t, currPoll)
		require.Equal(t, polls[i].ID.String(), currPoll.ID.String())
		require.Equal(t, nextPoll.ID.String(), currPoll.ID.String())
	}

	nextPoll, err := testQueries.GetNextPoll(context.Background(), createdPresenationID)
	require.Error(t, err)
	require.Empty(t, nextPoll)

	currPoll, err := testQueries.GetPresentationCurrentPoll(context.Background(), createdPresenationID)
	require.NoError(t, err)
	require.NotEmpty(t, currPoll)
	require.Equal(t, polls[1].ID.String(), currPoll.ID.String())

}

func TestGetPreviousPoll(t *testing.T) {
	createdPresenationID := createRandomPresentationWithPolls(t, 2)
	polls, err := testQueries.ListPolls(context.Background(), createdPresenationID)
	require.NoError(t, err)
	require.NotEmpty(t, polls)
	require.Len(t, polls, 2)

	// currently no poll is displayed
	poll, err := testQueries.GetPresentationCurrentPoll(context.Background(), createdPresenationID)
	require.Error(t, err)
	require.Empty(t, poll)

	// move forward twice to the second index
	for i := 0; i < 2; i++ {
		nextPoll, err := testQueries.GetNextPoll(context.Background(), createdPresenationID)
		require.NoError(t, err)
		require.NotEmpty(t, nextPoll)
		require.Equal(t, polls[i].ID.String(), nextPoll.ID.String())
	}

	// move backward to the first index
	prevPoll, err := testQueries.GetPreviousPoll(context.Background(), createdPresenationID)
	require.NoError(t, err)
	require.NotEmpty(t, prevPoll)
	require.Equal(t, polls[0].ID.String(), prevPoll.ID.String())

	//check the current displayed poll
	currPoll, err := testQueries.GetPresentationCurrentPoll(context.Background(), createdPresenationID)
	require.NoError(t, err)
	require.NotEmpty(t, currPoll)
	require.Equal(t, prevPoll.ID.String(), currPoll.ID.String())

	// move backward to the zero index  no poll is display
	prevPoll, err = testQueries.GetPreviousPoll(context.Background(), createdPresenationID)
	require.Error(t, err)
	require.Empty(t, prevPoll)
}
