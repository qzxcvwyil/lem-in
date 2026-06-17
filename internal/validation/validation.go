package validation

import (
	"errors"
	"fmt"

	"01.tomorrow-school.ai/git/msabyrga/lem-in.git/internal/model"
)

func Args(args []string) error {
	argsLength := len(args)

	if argsLength == 0 {
		return errors.New("argument list not provided")
	}

	if argsLength > 1 {
		return errors.New("argument list exceeded")
	}

	return nil
}

func Input(input *model.Input) error {
	// Check duplicated room names
	roomSet := make(map[string]struct{})

	// Adding start and end for set
	roomSet[input.StartRoom.Name] = struct{}{}
	if _, ok := roomSet[input.EndRoom.Name]; ok {
		return fmt.Errorf("duplicated room: %s", input.EndRoom.Name)
	}
	roomSet[input.EndRoom.Name] = struct{}{}

	for _, room := range input.Rooms {
		if _, ok := roomSet[room.Name]; ok {
			return fmt.Errorf("duplicated room: %s", room.Name)
		}
		roomSet[room.Name] = struct{}{}
	}

	for room, links := range input.Links {
		if _, ok := roomSet[room]; !ok {
			return fmt.Errorf("unknown link for room: %s", room)
		}

		for _, link := range links {
			if _, ok := roomSet[link]; !ok {
				return fmt.Errorf("unknown link for room: %s", link)
			}
		}
	}

	return nil
}
