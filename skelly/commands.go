package skelly

import (
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
)

var (
	helpSubCommand   = "help"
	addSubCommand    = "add"
	updateSubCommand = "update"
	deleteSubCommand = "delete"
	listSubCommand   = "list"
)

// HandleSlashCommand takes slack slash command configuration and executes it
func HandleSlashCommand(s *slack.SlashCommand) error {

	// parse slash command input
	args := strings.Fields(strings.TrimSpace(strings.ToLower(s.Text)))

	// join the command and the input
	command := strings.Join([]string{s.Command, strings.Join(args, " ")}, " ")

	logrus.Infof("handling command(%s)", command)

	// validate input
	if len(args) == 0 {

		// unsupported command, send help
		err := sendHelp(command, s.ResponseURL)
		if err != nil {
			err = errors.Wrap(err, "could not send help")
			return err
		}
		return nil
	}

	// execute the subcommand
	err := handleSubCommand(s, command, args)
	if err != nil {
		err = errors.Wrap(err, "could not send help")
		return err
	}
	return nil
}

// handleSubCommand takes slash command arguments and executes the appropriate subcommand
func handleSubCommand(s *slack.SlashCommand, command string, args []string) error {

	subcommand := args[0]

	// execute the command
	switch subcommand {

	// /skelly help
	case helpSubCommand:

		// send help message
		err := sendHelp(command, s.ResponseURL)
		if err != nil {
			err = errors.Wrap(err, "could not send help")
			return err
		}

		return nil

	// /skelly add :emoji: <@usergroup>
	case addSubCommand:

		// open add reaction modal
		err := openAddModal(s, command, args)
		if err != nil {
			err = errors.Wrap(err, "could not open add modal")
			return err
		}

		return nil

	// /skelly update :smile: <@usergroup>
	case updateSubCommand:

		// open update reaction modal
		err := openUpdateModal(s, command, args)
		if err != nil {
			err = errors.Wrap(err, "could not open update modal")
			return err
		}

		return nil

	// /skelly delete :smile: <@usergroup>
	case deleteSubCommand:

		// open delete reaction modal
		err := openDeleteModal(s, command, args)
		if err != nil {
			err = errors.Wrap(err, "could not open delete modal")
			return err
		}

		return nil

	// /skelly list
	case listSubCommand:

		// open delete reaction modal
		err := listReactions(s, command, args)
		if err != nil {
			err = errors.Wrap(err, "could not list reactions")
			return err
		}

		return nil

	// unsupported skelly subcommand
	default:

		// unsupported command, send help
		err := sendHelp(command, s.ResponseURL)
		if err != nil {
			err = errors.Wrap(err, "could not send help")
			return err
		}
		return nil
	}
}

// parseReactionSubCommandArgs takes input args and parses emoji and usergroup
// returns error if input is not valid
func parseReactionSubCommandArgs(args []string) (string, string, error) {

	// validate input
	if len(args) != 2 && len(args) != 3 {
		return "", "", errors.New("invalid number of args")
	}

	// extract emoji
	emoji := args[1]

	// default usergroup set to "none"
	usergroup := "none"

	// user specified usergroup
	if len(args) == 3 {
		usergroup = args[2]
	}
	return emoji, usergroup, nil
}

// parseListSubCommandArgs takes input args and checks for no args
// returns error if input is not valid
func parseListSubCommandArgs(args []string) error {

	// validate input
	if len(args) != 1 {
		return errors.New("invalid number of args")
	}

	return nil
}
