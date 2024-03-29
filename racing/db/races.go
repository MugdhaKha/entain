package db

import (
	"database/sql"
	"strings"
	"sync"
	"time"

	"github.com/golang/protobuf/ptypes"
	_ "github.com/mattn/go-sqlite3"

	"git.neds.sh/matty/entain/racing/proto/racing"
)

// RacesRepo provides repository access to races.
type RacesRepo interface {
	// Init will initialise our races repository.
	Init() error

	// List will return a list of races.
	List(filter *racing.ListRacesRequestFilter, orderBy *string) ([]*racing.Race, error)

	// Get will return single race based on race ID.
	Get(id int64) (*racing.Race, error)
}

type racesRepo struct {
	db   *sql.DB
	init sync.Once
}

// NewRacesRepo creates a new races repository.
func NewRacesRepo(db *sql.DB) RacesRepo {
	return &racesRepo{db: db}
}

// Init prepares the race repository dummy data.
func (r *racesRepo) Init() error {
	var err error

	r.init.Do(func() {
		// For test/example purposes, we seed the DB with some dummy races.
		err = r.seed()
	})

	return err
}

func (r *racesRepo) List(filter *racing.ListRacesRequestFilter, orderBy *string) ([]*racing.Race, error) {
	var (
		err   error
		query string
		args  []interface{}
	)

	query = getRaceQueries()[racesList]

	query, args = r.applyFilter(query, filter)

	query = r.applyOrder(query, orderBy)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	return r.scanRaces(rows)
}

func (r *racesRepo) Get(id int64) (*racing.Race, error) {
	var (
		err   error
		query string
	)

	query = getRaceQueries()[raceGet]

	rows, err := r.db.Query(query, id)
	if err != nil {
		return nil, err
	}

	// re-using the funcs
	races, err := r.scanRaces(rows)
	if err != nil {
		return nil, err
	}

	if len(races) == 1 {
		return races[0], nil
	}

	// if more than 1 or no race is returned
	return nil, nil
}

func (r *racesRepo) applyFilter(query string, filter *racing.ListRacesRequestFilter) (string, []interface{}) {
	var (
		clauses []string
		args    []interface{}
	)

	if filter == nil {
		return query, args
	}

	if len(filter.MeetingIds) > 0 {
		clauses = append(clauses, "meeting_id IN ("+strings.Repeat("?,", len(filter.MeetingIds)-1)+"?)")

		for _, meetingID := range filter.MeetingIds {
			args = append(args, meetingID)
		}
	}

	/*
		verify if visible filter is provided and if so add it to the db query clause
	*/
	if filter.Visible != nil {
		clauses = append(clauses, "visible = ?")

		args = append(args, filter.Visible)
	}

	if len(clauses) != 0 {
		query += " WHERE " + strings.Join(clauses, " AND ")
	}

	return query, args
}

func (r *racesRepo) applyOrder(query string, order *string) string {
	// adding order by clause to sort by default
	query += " ORDER BY advertised_start_time"

	// NB - assumption that asc or desc will be sent in order-by clause
	if order != nil && len(*order) > 0 {
		query += " " + strings.ToUpper(*order)
	}

	return query
}

func (r *racesRepo) scanRaces(
	rows *sql.Rows,
) ([]*racing.Race, error) {
	var races []*racing.Race

	currentTimeUTC := time.Now().UTC()

	for rows.Next() {
		var race racing.Race
		var advertisedStart time.Time

		if err := rows.Scan(&race.Id, &race.MeetingId, &race.Name, &race.Number, &race.Visible, &advertisedStart); err != nil {
			if err == sql.ErrNoRows {
				return nil, nil
			}

			return nil, err
		}

		ts, err := ptypes.TimestampProto(advertisedStart)
		if err != nil {
			return nil, err
		}

		race.AdvertisedStartTime = ts

		/*
			Shifted the race status calculation here from the racing service. Status calculation is an improvement of Q3.
			Reason:
				- Less complexity: this is already iterating through the races individually
				- Optimisation: do not have to explicitly call any function from multiple places to calculate the status
		*/
		race.Status = "OPEN"
		if race.AdvertisedStartTime.AsTime().UTC().Before(currentTimeUTC) {
			race.Status = "CLOSED"
		}

		races = append(races, &race)
	}

	return races, nil
}
