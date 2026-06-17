package validation

import (
	"testing"

	"01.tomorrow-school.ai/git/msabyrga/lem-in.git/internal/model"
)

func TestArgsEmpty(t *testing.T) {
	if err := Args([]string{}); err == nil {
		t.Fatal("expected error for empty args")
	}
}

func TestArgsTooMany(t *testing.T) {
	if err := Args([]string{"a", "b"}); err == nil {
		t.Fatal("expected error for too many args")
	}
}

func TestArgsValid(t *testing.T) {
	if err := Args([]string{"file.txt"}); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestInputDuplicatedRoom(t *testing.T) {
	input := &model.Input{
		StartRoom: model.Room{Name: "start"},
		EndRoom:   model.Room{Name: "end"},
		Rooms:     []model.Room{{Name: "A"}, {Name: "A"}},
		Links:     map[string][]string{},
	}
	if err := Input(input); err == nil {
		t.Fatal("expected error for duplicated room")
	}
}

func TestInputUnknownLink(t *testing.T) {
	input := &model.Input{
		StartRoom: model.Room{Name: "start"},
		EndRoom:   model.Room{Name: "end"},
		Rooms:     []model.Room{},
		Links:     map[string][]string{"ghost": {"end"}},
	}
	if err := Input(input); err == nil {
		t.Fatal("expected error for unknown room in links")
	}
}
