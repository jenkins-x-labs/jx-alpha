package get

import (
	jxget "github.com/jenkins-x/jx/v2/pkg/cmd/get"
	"github.com/jenkins-x/jx/v2/pkg/cmd/helper"

	"github.com/spf13/cobra"

	"github.com/jenkins-x/jx/v2/pkg/cmd/opts"
	"github.com/jenkins-x/jx/v2/pkg/cmd/templates"
)

// Options is the start of the data required to perform the operation.  As new fields are added, add them here instead of
// referencing the cmd.Flags()
type Options struct {
	*opts.CommonOptions

	Output string
}

const (
	validResources = `Valid resource types include:

    * environments (aka 'env')
    * pipelines (aka 'pipe')
    * urls (aka 'url')
    `
)

var (
	getLong = templates.LongDesc(`
		Display one or more resources.

		` + validResources + `

`)

	getExample = templates.Examples(`
		# List all pipelines
		jx get pipeline

		# List all URLs for services in the current namespace
		jx get url
	`)
)

// NewCmdGet creates a command object for the generic "get" action, which
// retrieves one or more resources from a server.
func NewCmdGet(commonOpts *opts.CommonOptions) *cobra.Command {
	options := &Options{
		CommonOptions: commonOpts,
	}

	cmd := &cobra.Command{
		Use:     "get TYPE [flags]",
		Short:   "Display one or more resources",
		Long:    getLong,
		Example: getExample,
		Run: func(cmd *cobra.Command, args []string) {
			options.Cmd = cmd
			options.Args = args
			err := options.Run()
			helper.CheckErr(err)
		},
		SuggestFor: []string{"list", "ps"},
	}

	cmd.AddCommand(jxget.NewCmdGetActivity(commonOpts))
	cmd.AddCommand(jxget.NewCmdGetApplications(commonOpts))
	cmd.AddCommand(jxget.NewCmdGetBuild(commonOpts))
	cmd.AddCommand(jxget.NewCmdGetBuildPack(commonOpts))
	cmd.AddCommand(jxget.NewCmdGetDevPod(commonOpts))
	cmd.AddCommand(jxget.NewCmdGetEnv(commonOpts))
	cmd.AddCommand(jxget.NewCmdGetLang(commonOpts))
	cmd.AddCommand(jxget.NewCmdGetPipeline(commonOpts))
	cmd.AddCommand(jxget.NewCmdGetPreview(commonOpts))
	cmd.AddCommand(jxget.NewCmdGetQuickstartLocation(commonOpts))
	cmd.AddCommand(jxget.NewCmdGetQuickstarts(commonOpts))
	cmd.AddCommand(jxget.NewCmdGetRelease(commonOpts))
	cmd.AddCommand(jxget.NewCmdGetStorage(commonOpts))
	cmd.AddCommand(jxget.NewCmdGetURL(commonOpts))
	cmd.AddCommand(jxget.NewCmdGetStream(commonOpts))
	return cmd
}

// Run implements this command
func (o *Options) Run() error {
	return o.Cmd.Help()
}

// AddGetFlags adds flags for get
func (o *Options) AddGetFlags(cmd *cobra.Command) {
	o.Cmd = cmd
	cmd.Flags().StringVarP(&o.Output, "output", "o", "", "The output format such as 'yaml'")
}
