package cmd

import (
	"fmt"
	"github.com/google/go-cmp/cmp"
	"github.com/nqd/flat"
	"github.com/r3labs/diff/v3"
	logger "github.com/sirupsen/logrus"
)

type entry struct {
	Filename string
	Settings map[string]any
}

type differ struct {
	settings []*entry
	log      *logger.Entry
}

func newDiffer(log *logger.Entry) *differ {
	return &differ{[]*entry{}, log}
}

func (d *differ) AddSettings(fName string, m map[string]any) {
	uf, _ := flat.Unflatten(m, nil)
	d.settings = append(d.settings, &entry{fName, uf})
}

func (d *differ) DebugDiff() {
	if d.log.Level <= logger.DebugLevel {
		for i := 1; i < len(d.settings); i++ {
			lEntry, rEntry := d.settings[i-1], d.settings[i]
			d.log.Debugf("Diffing [%s] and [%s]...", lEntry.Filename, rEntry.Filename)
			eDiff := cmp.Diff(lEntry.Settings, rEntry.Settings)
			d.log.Tracef("Result:\n[%s]", eDiff)
		}
	}
	return
}

func (d *differ) Diff() map[string]diff.Changelog {
	out := make(map[string]diff.Changelog)
	for i := 1; i < len(d.settings); i++ {
		lEntry, rEntry := d.settings[i-1], d.settings[i]
		changelog, _ := diff.Diff(lEntry.Settings, rEntry.Settings, diff.SliceOrdering(true))
		out[fmt.Sprintf("%s -> %s", lEntry.Filename, rEntry.Filename)] = changelog
	}
	return out
}
