package format

import (
	"time"

	"k8s.io/apimachinery/pkg/util/duration"
)

// HumanDuration returns a human-readable approximation of d.
func HumanDuration(d time.Duration) string {
	return duration.HumanDuration(d)
}
