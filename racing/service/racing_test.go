package service

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"

	"git.neds.sh/matty/entain/racing/proto/racing"
)

func Test_ListRaces_CalculateStatus(t *testing.T) {
	srv := &racingService{}

	races := []*racing.Race{
		{
			Id:                  1,
			MeetingId:           1,
			Name:                "Test 1",
			Number:              1,
			Visible:             false,
			AdvertisedStartTime: timestamppb.New(time.Now().UTC().Add(-1 * time.Hour)),
		},
		{
			Id:                  2,
			MeetingId:           2,
			Name:                "Test 2",
			Number:              2,
			Visible:             true,
			AdvertisedStartTime: timestamppb.New(time.Now().UTC().Add(1 * time.Hour)),
		},
	}

	srv.calculateStatus(races)

	// check if returned races are not empty
	assert.NotNil(t, races)

	// check the length of the races is 2
	assert.Len(t, races, 2)

	// check that race 1 has CLOSED status
	assert.Equal(t, races[0].Status, "CLOSED")

	// check that race 2 has OPEN status
	assert.Equal(t, races[1].Status, "OPEN")
}
