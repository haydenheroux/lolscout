package export

import (
	"encoding/csv"
	"lolscout/data"
	"os"
)

func WriteMatches(name string, stats []data.MatchParticipantMetrics) error {
	file, err := os.Create(name)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	header := false

	for _, s := range stats {
		if !header {
			err := writer.Write(s.Header())

			if err != nil {
				return err
			}

			header = true
		}

		err := writer.Write(s.Row())
		if err != nil {
			return err
		}
	}

	return nil
}
