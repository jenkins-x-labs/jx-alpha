package version

import (
	"fmt"

	"github.com/jenkins-x-labs/jxl/pkg/cmd/upgrade"
	"github.com/jenkins-x-labs/jxl/pkg/version"
	"github.com/jenkins-x/jx/v2/pkg/versionstream"

	"github.com/jenkins-x/jx/v2/pkg/cmd/helper"

	"github.com/jenkins-x/jx/v2/pkg/cmd/opts"
	"github.com/jenkins-x/jx-logging/pkg/log"
	"github.com/jenkins-x/jx/v2/pkg/util"
	"github.com/jenkins-x/jx/v2/pkg/util/system"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

type VersionOptions struct {
	*opts.CommonOptions

	Container      string
	Namespace      string
	HelmTLS        bool
	NoVersionCheck bool
	NoVerify       bool
}

func NewCmdVersion(commonOpts *opts.CommonOptions) *cobra.Command {
	options := &VersionOptions{
		CommonOptions: commonOpts,
	}

	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print the version information",
		Run: func(cmd *cobra.Command, args []string) {
			options.Cmd = cmd
			options.Args = args
			err := options.Run()
			helper.CheckErr(err)
		},
	}
	cmd.Flags().BoolVarP(&options.NoVersionCheck, "no-version-check", "n", false, "Disable checking of version upgrade checks")
	cmd.Flags().BoolVarP(&options.NoVerify, "no-verify", "", false, "Disable verification of package versions")
	cmd.Flags().StringVarP(&options.Namespace, "namespace", "", "", "The namespace to use to look for currently installed platform version")
	return cmd
}

func (o *VersionOptions) Run() error {
	ns := o.Namespace
	if ns == "" {
		var err error
		_, ns, err = o.JXClientAndDevNamespace()
		if err != nil {
			return err
		}
	}
	packages, table := o.GetPackageVersions(ns, o.HelmTLS)

	delete(packages, "jx")
	jxlVersion := version.GetVersion()
	packages["jxl"] = jxlVersion
	// replace the `jx` row
	table.Rows[1] = []string{"jxl", jxlVersion}

	// os version
	osVersion, err := o.GetOsVersion()
	if err != nil {
		log.Logger().Warnf("Failed to get OS version: %s", err)
	} else {
		table.AddRow("Operating System", util.ColorInfo(osVersion))
	}

	table.Render()
	if o.NoVerify {
		return nil
	}

	log.Logger().Info("\n\nverifying packages")
	versionResolver, err := o.GetVersionResolver()
	if err != nil {
		return err
	}

	if !o.NoVersionCheck {
		err = o.upgradeCliIfNeeded(versionResolver)
		if err != nil {
			return err
		}
		// lets not verify the jx version as the current executing binary
		// may have just been upgraded anyway and we've already warned the user if its old
		delete(packages, "jx")
	}

	// lets remove any non-package name before verifying
	delete(packages, "kubernetesCluster")

	return versionResolver.VerifyPackages(packages)
}

func (o *VersionOptions) upgradeCliIfNeeded(resolver *versionstream.VersionResolver) error {
	currentVersion, err := version.GetSemverVersion()
	if err != nil {
		return errors.Wrap(err, "getting current jx version")
	}

	newVersion, err := version.GetJXLVersion(resolver)
	if err != nil {
		return err
	}

	if currentVersion.LT(newVersion) {
		app := util.ColorInfo("jxl")
		log.Logger().Info("\n")
		log.Logger().Warnf("%s version %s is available in the version stream. You are using %s. We highly recommend you upgrade to it.", app, util.ColorInfo(newVersion.String()), util.ColorInfo(currentVersion.String()))
		if o.BatchMode {
			log.Logger().Warnf("To upgrade to this new version use: %s", util.ColorInfo("jxl upgrade cli"))
		} else {
			log.Logger().Info("\n")
			message := fmt.Sprintf("Would you like to upgrade to the %s version?", app)
			if answer, err := util.Confirm(message, true, "Please indicate if you would like to upgrade the binary version.", o.GetIOFileHandles()); err != nil {
				return err
			} else if answer {
				options := &upgrade.UpgradeCLIOptions{
					CommonOptions: o.CommonOptions,
				}
				options.Version = newVersion.String()
				options.NoBrew = true
				return options.Run()
			}
		}
	}
	return nil
}

// GetOsVersion returns a human friendly string of the current OS
// in the case of an error this still returns a valid string for the details that can be found.
func (o *VersionOptions) GetOsVersion() (string, error) {
	return system.GetOsVersion()
}
