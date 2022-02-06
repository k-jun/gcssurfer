package c

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gdamore/tcell/v2"
	"github.com/k-jun/gcssurfer/pkg/m"
	"github.com/k-jun/gcssurfer/pkg/v"
	"github.com/rivo/tview"
)

type Controller struct {
	dfp *os.File
	v   v.View
	m   m.GCSManager
}

// NewController ...
func NewController(
	project string,
	bucket string,
	debug string,
	version string,
) Controller {

	var dfp *os.File
	if debug != "" {
		var err error
		if dfp, err = os.OpenFile(filepath.Clean(debug), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600); err != nil {
			panic(err)
		}
	}

	fmt.Printf("fetch available buckets...\n")
	m := m.NewGCSManager(project)
	if bucket != "" {
		err := m.SetBucket(bucket)
		if err != nil {
			panic(err)
		}
	}

	v := v.NewView()
	v.Frame.AddText(version, true, tview.AlignCenter, tcell.ColorWhite)

	c := Controller{
		dfp,
		v,
		m,
	}

	return c
}
