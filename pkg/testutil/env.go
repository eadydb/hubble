package testutil

import "os"

// SetEnvs takes a map of key values to set using os.Setenv and returns
// a function that can be called to reset the envs to their previous values.
func (t *T) SetEnvs(envs map[string]string) {
	prevEnvs := map[string]string{}
	for key := range envs {
		prevEnvs[key] = os.Getenv(key)
	}

	t.Cleanup(func() { setEnvs(t, prevEnvs) })

	setEnvs(t, envs)
}

func setEnvs(t *T, envs map[string]string) {
	for key, value := range envs {
		if err := os.Setenv(key, value); err != nil {
			t.Error(err)
		}
	}
}
