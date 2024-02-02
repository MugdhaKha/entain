package db

import (
	"fmt"
	"time"

	"syreclabs.com/go/faker"
)

func (s *sportsRepo) seed() error {
	statement, err := s.db.Prepare(`CREATE TABLE IF NOT EXISTS sports (id INTEGER PRIMARY KEY, name TEXT, advertised_start_time DATETIME, sport TEXT, home_team TEXT, away_team TEXT, competition_id INTEGER, competition_name TEXT)`)
	if err == nil {
		_, err = statement.Exec()
	}

	for i := 1; i <= 100; i++ {
		homeTeam := faker.Team().Name()
		awayTeam := faker.Team().Name()
		statement, err = s.db.Prepare(`INSERT OR IGNORE INTO sports(id, name, advertised_start_time, sport, home_team, away_team, competition_id, competition_name) VALUES (?,?,?,?,?,?,?,?)`)
		if err == nil {
			_, err = statement.Exec(
				i,
				fmt.Sprintf("%s vs %s", homeTeam, awayTeam),
				faker.Time().Between(time.Now().AddDate(0, 0, -1), time.Now().AddDate(0, 0, 2)).Format(time.RFC3339),
				faker.Name().Name(),
				homeTeam,
				awayTeam,
				faker.Number().NumberInt(5),
				faker.Company().CatchPhrase(),
			)
		}
	}

	return err
}
