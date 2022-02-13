package main

import (
	"fmt"
	"os"
	"runtime"
	"syscall"

	"github.com/alecthomas/kong"
	"github.com/cli/safeexec"

	"github.com/k-jun/gcssurfer/pkg/c"
)

const version = "1.0.9"

var revision = "HEAD"

type buildInfo struct {
	Version  string
	Revision string
}

func (b buildInfo) String() string {
	return fmt.Sprintf(
		"gcssurfer %s (rev: %s/%s)",
		b.Version,
		b.Revision,
		runtime.Version(),
	)
}

// CLI ...
type CLI struct {
	Debug   string           `help:"Write debug log info file" short:"d" type:"path"`
	Version kong.VersionFlag `help:"Print version information and exit" short:"v"`

	Project string `help:"Project id search into" short:"p" arg:"required"`
	Bucket  string `help:"GCS bucket name" short:"b" arg:"optional"`
}

func init() {
	// https://github.com/rivo/tview/wiki/FAQ#why-do-my-borders-look-weird
	if os.Getenv("LC_CTYPE") != "en_US.UTF-8" && runtime.GOOS != "windows" {
		err := os.Setenv("LC_CTYPE", "en_US.UTF-8")
		if err != nil {
			panic(err)
		}
		env := os.Environ()
		argv0, err := safeexec.LookPath(os.Args[0])
		if err != nil {
			panic(err)
		}
		os.Args[0] = argv0
		/* #nosec G204 */
		if err := syscall.Exec(argv0, os.Args, env); err != nil {
			panic(err)
		}
	}
}

func main() {
	buildInfo := buildInfo{
		Version:  version,
		Revision: revision,
	}

	cli := CLI{}
	ctx := kong.Parse(&cli,
		kong.Name("gcssurfer"),
		kong.Description("gcssurfer is CLI tool for browsing gcs bucket and download objects.\nhttps://github.com/k-jun/gcssurfer"),
		kong.UsageOnError(),
		kong.Vars{"version": buildInfo.String()},
	)

	err := c.NewController(
		cli.Project,
		cli.Bucket,
		cli.Debug,
		buildInfo.String(),
	).Run()

	ctx.FatalIfErrorf(err)
}
