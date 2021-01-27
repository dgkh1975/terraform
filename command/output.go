package command

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform/command/views"
	"github.com/hashicorp/terraform/states"
	"github.com/hashicorp/terraform/tfdiags"
)

// OutputCommand is a Command implementation that reads an output
// from a Terraform state and prints it.
type OutputCommand struct {
	Meta
}

func (c *OutputCommand) View() views.View {
	return views.NewView(c.Ui, false)
}

func (c *OutputCommand) Run(args []string) int {
	// Parse and validate args
	args = c.Meta.process(args)
	var statePath string
	var jsonOutput, rawOutput bool
	cmdFlags := c.Meta.defaultFlagSet("output")
	cmdFlags.BoolVar(&jsonOutput, "json", false, "json")
	cmdFlags.BoolVar(&rawOutput, "raw", false, "raw")
	cmdFlags.StringVar(&statePath, "state", "", "path")
	cmdFlags.Usage = func() { c.Ui.Error(c.Help()) }
	if err := cmdFlags.Parse(args); err != nil {
		c.Ui.Error(fmt.Sprintf("Error parsing command-line flags: %s\n", err.Error()))
		return 1
	}

	args = cmdFlags.Args()
	if len(args) > 1 {
		c.Ui.Error(
			"The output command expects exactly one argument with the name\n" +
				"of an output variable or no arguments to show all outputs.\n")
		cmdFlags.Usage()
		return 1
	}

	if jsonOutput && rawOutput {
		c.Ui.Error("The -raw and -json options are mutually-exclusive.\n")
		cmdFlags.Usage()
		return 1
	}

	if rawOutput && len(args) == 0 {
		c.Ui.Error("You must give the name of a single output value when using the -raw option.\n")
		cmdFlags.Usage()
		return 1
	}

	name := ""
	if len(args) > 0 {
		name = args[0]
	}

	if statePath != "" {
		c.Meta.statePath = statePath
	}

	// Prepare view for configured UI mode
	var view views.Output
	if jsonOutput {
		view = &views.OutputJSON{View: c.View()}
	} else if rawOutput {
		view = &views.OutputRaw{}
	} else {
		view = &views.OutputText{View: c.View()}
	}

	// Fetch data from state
	outputs, diags := c.Outputs()
	if diags.HasErrors() {
		c.showDiagnostics(diags)
		return 1
	}

	// Render the view
	viewDiags := view.Output(name, outputs)
	diags = diags.Append(viewDiags)

	c.showDiagnostics(diags)

	if diags.HasErrors() {
		return 1
	}

	return 0
}

func (c *OutputCommand) Outputs() (map[string]*states.OutputValue, tfdiags.Diagnostics) {
	var diags tfdiags.Diagnostics

	// Load the backend
	b, backendDiags := c.Backend(nil)
	diags = diags.Append(backendDiags)
	if diags.HasErrors() {
		return nil, diags
	}

	// This is a read-only command
	c.ignoreRemoteBackendVersionConflict(b)

	env, err := c.Workspace()
	if err != nil {
		diags = diags.Append(fmt.Errorf("Error selecting workspace: %s", err))
		return nil, diags
	}

	// Get the state
	stateStore, err := b.StateMgr(env)
	if err != nil {
		diags = diags.Append(fmt.Errorf("Failed to load state: %s", err))
		return nil, diags
	}

	if err := stateStore.RefreshState(); err != nil {
		diags = diags.Append(fmt.Errorf("Failed to load state: %s", err))
		return nil, diags
	}

	state := stateStore.State()
	if state == nil {
		state = states.NewState()
	}

	return state.RootModule().OutputValues, nil
}

func (c *OutputCommand) Help() string {
	helpText := `
Usage: terraform output [options] [NAME]

  Reads an output variable from a Terraform state file and prints
  the value. With no additional arguments, output will display all
  the outputs for the root module.  If NAME is not specified, all
  outputs are printed.

Options:

  -state=path      Path to the state file to read. Defaults to
                   "terraform.tfstate".

  -no-color        If specified, output won't contain any color.

  -json            If specified, machine readable output will be
                   printed in JSON format.

  -raw             For value types that can be automatically
                   converted to a string, will print the raw
                   string directly, rather than a human-oriented
                   representation of the value.
`
	return strings.TrimSpace(helpText)
}

func (c *OutputCommand) Synopsis() string {
	return "Show output values from your root module"
}
