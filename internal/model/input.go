package model

type Input struct {
	Ants      int
	StartRoom Room
	Rooms     []Room
	Comments  []string
	EndRoom   Room
	Links     map[string][]string
}

type Room struct {
	Name string
	X    int
	Y    int
}
