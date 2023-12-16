package db

import (
	"context"
	"testing"

	"github.com/aleale2121/interactive-presentation/util"
	"github.com/stretchr/testify/require"
)

func createRandomOption(t *testing.T) string {
	key := util.RandomOptionKey()
	presentation := createRandomPresentation(t)
	poll := createPresentationPoll(t, presentation)

	arg := CreateOptionParams{
		Pollid:      poll.ID,
		Optionkey:   key,
		Optionvalue: util.RandomOptionValue(),
	}

	err := testQueries.CreateOption(context.Background(), arg)
	require.NoError(t, err)
	return key
}

func createRandomPollOption(t *testing.T, poll Poll) string {
	key := util.RandomOptionKey()
	arg := CreateOptionParams{
		Pollid:      poll.ID,
		Optionkey:   key,
		Optionvalue: util.RandomOptionValue(),
	}

	err := testQueries.CreateOption(context.Background(), arg)
	
	require.NoError(t, err)
	return key
}

func TestCreateOption(t *testing.T) {
	createRandomOption(t)
}

func TestGetOption(t *testing.T) {
	optKey := createRandomOption(t)
	option2, err := testQueries.GetOptionByKey(context.Background(), optKey)
	require.NoError(t, err)
	require.NotEmpty(t, option2)
}

func TestGetPollOptions(t *testing.T) {
	pres := createRandomPresentation(t)
	poll := createPresentationPoll(t, pres)
	for i := 0; i < 10; i++ {
		createRandomPollOption(t, poll)
	}

	options, err := testQueries.GetPollOptions(context.Background(), poll.ID)
	require.NoError(t, err)
	require.NotEmpty(t, options)
	require.Len(t, options, 10)
	for _, option := range options {
		require.NotEmpty(t, option)
	}
}
