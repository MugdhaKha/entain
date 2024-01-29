package service

import (
	"time"

	"golang.org/x/net/context"

	"git.neds.sh/matty/entain/racing/db"
	"git.neds.sh/matty/entain/racing/proto/racing"
)

type Racing interface {
	// ListRaces will return a collection of races.
	ListRaces(ctx context.Context, in *racing.ListRacesRequest) (*racing.ListRacesResponse, error)
}

// racingService implements the Racing interface.
type racingService struct {
	racesRepo db.RacesRepo
}

// NewRacingService instantiates and returns a new racingService.
func NewRacingService(racesRepo db.RacesRepo) Racing {
	return &racingService{racesRepo}
}

func (s *racingService) ListRaces(ctx context.Context, in *racing.ListRacesRequest) (*racing.ListRacesResponse, error) {
	races, err := s.racesRepo.List(in.Filter, in.OrderBy)
	if err != nil {
		return nil, err
	}

	s.calculateStatus(races)

	return &racing.ListRacesResponse{Races: races}, nil
}

func (s *racingService) calculateStatus(races []*racing.Race) {
	currentTimeUTC := time.Now().UTC()

	for _, race := range races {
		race.Status = "OPEN"

		// check if race advertised_start_time (in UTC) is before current time (in UTC)
		if race.AdvertisedStartTime.AsTime().UTC().Before(currentTimeUTC) {
			race.Status = "CLOSED"
		}
	}
}
