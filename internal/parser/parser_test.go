package parser

import (
	"os"
	"testing"
)

func TestParseValidInput(t *testing.T) {
	content := "4\n##start\n0 0 3\n2 2 5\n3 4 0\n##end\n1 8 3\n0-2\n2-3\n3-1\n"
	f, _ := os.CreateTemp("", "test*.txt")
	f.WriteString(content)
	f.Close()
	defer os.Remove(f.Name())

	p := NewParser()
	input, err := p.Parse(f.Name())
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if input.Ants != 4 {
		t.Errorf("expected 4 ants, got %d", input.Ants)
	}
	if input.StartRoom.Name != "0" {
		t.Errorf("expected start room '0', got '%s'", input.StartRoom.Name)
	}
	if input.EndRoom.Name != "1" {
		t.Errorf("expected end room '1', got '%s'", input.EndRoom.Name)
	}
}

func TestParseZeroAnts(t *testing.T) {
	content := "0\n##start\n0 0 3\n##end\n1 8 3\n0-1\n"
	f, _ := os.CreateTemp("", "test*.txt")
	f.WriteString(content)
	f.Close()
	defer os.Remove(f.Name())

	p := NewParser()
	_, err := p.Parse(f.Name())
	if err == nil {
		t.Fatal("expected error for zero ants, got nil")
	}
}

func TestParseMissingStart(t *testing.T) {
	content := "4\n0 0 3\n##end\n1 8 3\n0-1\n"
	f, _ := os.CreateTemp("", "test*.txt")
	f.WriteString(content)
	f.Close()
	defer os.Remove(f.Name())

	p := NewParser()
	_, err := p.Parse(f.Name())
	if err == nil {
		t.Fatal("expected error for missing start, got nil")
	}
}

func TestParseMissingEnd(t *testing.T) {
	content := "4\n##start\n0 0 3\n1 8 3\n0-1\n"
	f, _ := os.CreateTemp("", "test*.txt")
	f.WriteString(content)
	f.Close()
	defer os.Remove(f.Name())

	p := NewParser()
	_, err := p.Parse(f.Name())
	if err == nil {
		t.Fatal("expected error for missing end, got nil")
	}
}

func TestParseDuplicatedLink(t *testing.T) {
	content := "4\n##start\n0 0 3\n##end\n1 8 3\n0-1\n0-1\n"
	f, _ := os.CreateTemp("", "test*.txt")
	f.WriteString(content)
	f.Close()
	defer os.Remove(f.Name())

	p := NewParser()
	_, err := p.Parse(f.Name())
	if err == nil {
		t.Fatal("expected error for duplicated link, got nil")
	}
}
