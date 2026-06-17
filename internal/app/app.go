package app

import (
	"fmt"
	"os"

	"01.tomorrow-school.ai/git/msabyrga/lem-in.git/internal/parser"
	"01.tomorrow-school.ai/git/msabyrga/lem-in.git/internal/service/pathfinding"
	"01.tomorrow-school.ai/git/msabyrga/lem-in.git/internal/service/solver"
	"01.tomorrow-school.ai/git/msabyrga/lem-in.git/internal/validation"
)

func Run(args []string) error {
	if err := validation.Args(args); err != nil {
		return fmt.Errorf("args validation: %w", err)
	}

	pr := parser.NewParser()
	input, err := pr.Parse(args[0])
	if err != nil {
		return fmt.Errorf("parsing input: %w", err)
	}

	if err := validation.Input(input); err != nil {
		return fmt.Errorf("input validation: %w", err)
	}

	rightPaths, err := pathfinding.FindRightPaths(input)
	if err != nil {
		return fmt.Errorf("finding right paths: %w", err)
	}

	result, err := solver.Solve(rightPaths, input.Ants)
	if err != nil {
		return fmt.Errorf("solve function: %w", err)
	}

	if err := printResult(args[0], result); err != nil {
		return fmt.Errorf("reading file: %w", err)
	}

	return nil
}

func printResult(args string, result string) error {
	fileContent, err := os.ReadFile(args)
	if err != nil {
		return err
	}

	fmt.Println(string(fileContent))
	fmt.Println()
	fmt.Print(result)

	return nil
}
