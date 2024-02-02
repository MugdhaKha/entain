package db

const (
	eventList = "list"
)

func getSportsQueries() map[string]string {
	return map[string]string{
		eventList: `
			SELECT 
				id,
				name,
				advertised_start_time,
				sport,
				home_team,
				away_team,
				competition_id,
				competition_name
			FROM sports
		`,
	}
}
