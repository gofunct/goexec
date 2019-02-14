package commands

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func init() {
	goBinCmd.Flags().StringSliceVarP(&pkgs, "packages", "p", []string{}, "path to binary program to download")
	getCmd.Flags().StringVarP(&src, "source", "s", "", "url path to target download")
	getCmd.Flags().StringVarP(&src, "dest", "d", ".", "path to download files to")
	LoadCmd.AddCommand(goBinCmd, getCmd)
}

var (
	pkgs []string
	src  string
	dest string
)

var LoadCmd = &cobra.Command{
	Use:   "load",
	Short: "download all of your executable dependencies",
}

var goBinCmd = &cobra.Command{
	Use:   "mod",
	Short: "download go mod executable dependencies",
	Run: func(cmd *cobra.Command, args []string) {
		if err := goBin(pkgs); err != nil {
			panic(errors.Wrap(err, "failed to download binaries"))
		}
	},
}
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "intelligently download folders/files from a remote url",
	Long:  "use // before the path to a file to download a single file",
	Run: func(cmd *cobra.Command, args []string) {
		load(src, dest)
	},
}
