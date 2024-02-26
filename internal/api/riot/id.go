package riot

import (
	"errors"
	"fmt"
	"strings"
)

func Split(id string) (string, string, error) {
	fields := strings.Split(id, "#")

	if len(fields) != 2 {
		return "", "", errors.New("bad riot id")
	}

	gameName := fields[0]
	tagLine := fields[1]

	return gameName, tagLine, nil
}

func Join(gameName, tagLine string) string {
	return fmt.Sprintf("%s#%s", gameName, tagLine)
}
