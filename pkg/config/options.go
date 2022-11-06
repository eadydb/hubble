package config

import (
	"time"
)

// WaitForDeletions configures the wait for pending deletions.
type WaitForDeletions struct {
	Max     time.Duration
	Delay   time.Duration
	Enabled bool
}

// HubbleOptions are options that are set by command line arguments not included in the config file itself
type HubbleOptions struct {
	CheckClusterNodePlatforms bool
	GlobalConfig              string
	EventLogFile              string
	RenderOutput              string
	User                      string
	CustomTag                 string
	Namespace                 string
	CacheFile                 string
	Trigger                   string
	KubeContext               string
	KubeConfig                string
	Command                   string
	StatusCheck               BoolOrUndefined
}

type RunMode string

var RunModes = struct {
	Build    RunMode
	Dev      RunMode
	Debug    RunMode
	Run      RunMode
	Deploy   RunMode
	Render   RunMode
	Delete   RunMode
	Diagnose RunMode
}{
	Build:    "build",
	Dev:      "dev",
	Debug:    "debug",
	Run:      "run",
	Deploy:   "deploy",
	Render:   "render",
	Delete:   "delete",
	Diagnose: "diagnose",
}

func (opts *HubbleOptions) Mode() RunMode {
	return RunMode(opts.Command)
}
