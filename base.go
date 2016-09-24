package cut

import (
	"github.com/spf13/cobra"
	"reflect"
	"strings"
)

// Binder can bound some options to cobra.
type Binder interface {
	Bind(cc *cobra.Command)
}

type Command interface {

	// Run command with rest of args
	Run(args ...string) (err error)
}

// RunNone
type RunNone struct {}
func (r RunNone) Run(args ...string) (err error) { return }


// Attach command. If command implements Binder it will be also evaluated
func Attach(c Command, binders []Binder, cmds ...*cobra.Command) (cc *cobra.Command) {
	cc = &cobra.Command{}
	for _, binder := range binders {
		binder.Bind(cc)
	}
	if c1, ok := c.(Binder); ok {
		c1.Bind(cc)
	} else {
		// TypeOf().Name() return empty string on pointer types
		cmd := reflect.TypeOf(c).String()
		chunks := strings.Split(cmd, ".")
		cc.Use = strings.ToLower(strings.TrimPrefix(chunks[len(chunks)-1], "*"))
	}
	if len(cmds) == 0 {
		cc.RunE = func(cc *cobra.Command, args []string) (err error) {
			return c.Run(args...)
		}
	}
	cc.AddCommand(cmds...)
	return
}
