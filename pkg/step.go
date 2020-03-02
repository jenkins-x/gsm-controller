package pkg

import (
	"fmt"
	"os/user"

	"github.com/pkg/errors"
)

// Run implements the command
func (o *options) Run() error {
	user, err := user.Current()
	if err != nil {
		return errors.Wrap(err, "failed to get current user")
	}
	message := hello(user.Username)
	fmt.Printf(message)

	return nil

}

func hello(name string) string {
	return fmt.Sprintf("Hello %s", name)
}
