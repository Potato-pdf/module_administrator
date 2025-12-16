package validations

import (
	"errors"
	"luminaMO-orq-modAI/src/domain/valueobjects"
)

func ValidateCommand(command string) error {
	if command == "" {
		return errors.New("command is required")
	}
	if !IsValidCommand(command) {
		return errors.New("command is invalid")
	}
	return nil
}

func IsValidCommand(command string) bool {
	for _, validCmd := range valueobjects.ValidCommands {
		if string(validCmd) == command {
			return true
		}
	}
	return false
}

func GetAllCommands() []string {
	commands := make([]string, len(valueobjects.ValidCommands))
	for i, cmd := range valueobjects.ValidCommands {
		commands[i] = string(cmd)
	}
	return commands
}