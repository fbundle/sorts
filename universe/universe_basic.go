package universe

import (
	"errors"

	"github.com/fbundle/sorts/sorts"
)

func Basic(universeHeader sorts.Name, initialHeader sorts.Name, terminalHeader sorts.Name) (Universe, error) {
	nameSet := make(map[sorts.Name]struct{})
	nameSet[universeHeader] = struct{}{}
	nameSet[initialHeader] = struct{}{}
	nameSet[terminalHeader] = struct{}{}
	if len(nameSet) != 3 {
		return nil, errors.New("universe, initial, terminal name must be distinct")
	}
	u := &universe{
		universeHeader: universeHeader,
		initialHeader:  initialHeader,
		terminalHeader: terminalHeader,
	}

	return u, nil
}
