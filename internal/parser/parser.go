package parser

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"01.tomorrow-school.ai/git/msabyrga/lem-in.git/internal/model"
)

const (
	linkPartsLen = 2
	roomPartsLen = 3
)

type Parser struct {
	input     model.Input
	sc        *bufio.Scanner
	startSeen bool
	endSeen   bool
	linkSet   map[string]struct{}
}

func NewParser() *Parser {
	return &Parser{
		input: model.Input{
			Rooms:    make([]model.Room, 0),
			Comments: make([]string, 0),
			Links:    make(map[string][]string),
		},
		startSeen: false,
		linkSet:   make(map[string]struct{}),
	}
}

func (p *Parser) Parse(filePath string) (*model.Input, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("opening file: %w", err)
	}
	defer file.Close()

	p.sc = bufio.NewScanner(file)

	p.sc.Scan()
	ants := p.sc.Text()
	if err := p.handleAnts(ants); err != nil {
		return nil, fmt.Errorf("handling ants: %w", err)
	}

	for p.sc.Scan() {
		line := p.sc.Text()

		if line == "" {
			continue
		}

		if line == "##start" {
			if err := p.handleStart(); err != nil {
				return nil, fmt.Errorf("handling start room: %w", err)
			}
			continue
		}

		if line == "##end" {
			if err := p.handleEnd(); err != nil {
				return nil, fmt.Errorf("handling end room: %w", err)
			}
			continue
		}

		if strings.HasPrefix(line, "#") {
			p.input.Comments = append(p.input.Comments, line)
			continue
		}

		if strings.Contains(line, "-") {
			if err := p.handleLink(line); err != nil {
				return nil, fmt.Errorf("handling link: %w", err)
			}
			continue
		}

		room, err := parseRoom(line)
		if err != nil {
			return nil, fmt.Errorf("parsing room: %w", err)
		}
		p.input.Rooms = append(p.input.Rooms, *room)
	}

	if !p.startSeen {
		return nil, errors.New("start room not found")
	}

	if !p.endSeen {
		return nil, errors.New("end room not found")
	}

	return &p.input, nil
}

func (p *Parser) handleAnts(ants string) error {
	antsInt, err := strconv.Atoi(ants)
	if err != nil {
		return errors.New("ants must be int")
	}

	if antsInt <= 0 {
		return errors.New("ants must be more than 0")
	}

	p.input.Ants = antsInt

	return nil
}

func (p *Parser) handleStart() error {
	if p.startSeen {
		return errors.New("start marker duplicated")
	}

	if scanned := p.sc.Scan(); !scanned {
		return errors.New("file ended unexpectedly")
	}

	line := p.sc.Text()
	startRoom, err := parseRoom(line)
	if err != nil {
		return fmt.Errorf("parsing room: %w", err)
	}

	p.input.StartRoom = *startRoom
	p.startSeen = true

	return nil
}

func (p *Parser) handleEnd() error {
	if p.endSeen {
		return errors.New("end marker duplicated")
	}

	if scanned := p.sc.Scan(); !scanned {
		return errors.New("file ended unexpectedly")
	}

	line := p.sc.Text()
	endRoom, err := parseRoom(line)
	if err != nil {
		return fmt.Errorf("parsing room: %w", err)
	}

	p.input.EndRoom = *endRoom
	p.endSeen = true

	return nil
}

func (p *Parser) handleLink(line string) error {
	parts := strings.Split(line, "-")

	if len(parts) != linkPartsLen {
		return errors.New("link malformed")
	}

	var key string
	if parts[0] < parts[1] {
		key = parts[0] + "-" + parts[1]
	} else {
		key = parts[1] + "-" + parts[0]
	}

	if _, ok := p.linkSet[key]; ok {
		return fmt.Errorf("duplicated link: %s", key)
	}
	p.linkSet[key] = struct{}{}

	if parts[0] == parts[1] {
		return errors.New("self link not allowed")
	}

	// Bi-directional room links
	p.input.Links[parts[0]] = append(p.input.Links[parts[0]], parts[1])
	p.input.Links[parts[1]] = append(p.input.Links[parts[1]], parts[0])

	return nil
}

func parseRoom(line string) (*model.Room, error) {
	roomParts := strings.Split(line, " ")

	if len(roomParts) != roomPartsLen {
		return nil, errors.New("room malformed")
	}

	x, err := strconv.Atoi(roomParts[1])
	if err != nil {
		return nil, fmt.Errorf("x atoi conversion: %w", err)
	}

	y, err := strconv.Atoi(roomParts[2])
	if err != nil {
		return nil, fmt.Errorf("y atoi conversion: %w", err)
	}

	if strings.HasPrefix(roomParts[0], "#") || strings.HasPrefix(roomParts[0], "L") {
		return nil, errors.New("rooms cannot be started with # or L")
	}

	return &model.Room{
		Name: roomParts[0],
		X:    x,
		Y:    y,
	}, nil
}
