package db

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"git.neds.sh/matty/entain/racing/proto/racing"
)

func Test_ListRaces_WithVisibleFilter(t *testing.T) {
	query := getRaceQueries()[racesList]
	visible := true
	expectedQuery := query + " WHERE visible = ?"

	repo := racesRepo{}
	retQuery, retArgs := repo.applyFilter(query, &racing.ListRacesRequestFilter{Visible: &visible})

	// Check the correctness of the query
	assert.Equal(t, expectedQuery, retQuery)

	// check if returned arguments are not empty since we are filtering
	assert.NotEmpty(t, retArgs)

	// check the length of the arguments is 1 as only 1 filter is passed
	assert.Len(t, retArgs, 1)

	// check the argument value matches the filter value
	assert.Equal(t, retArgs[0], &visible)
}
