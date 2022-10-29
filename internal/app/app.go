package app

// App
type App struct {
	Name        string   `json:"name"`         // app name
	ContextPath string   `json:"context_path"` // app context path
	Address     []string `json:"address"`      // app address lists
	Port        int      `json:"port"`         // app port
	Version     string   `json:"version"`      // app version
	Status      string   `json:"status"`       // app status
}

// GetApp get app by name
// name is the name of the app
// registry app registry endpoint
// Returns App object
func GetApp(name, registry string) (*App, error) {
	app := &App{
		Name:        name,
		ContextPath: "",
		Address:     []string{""},
		Port:        0,
		Version:     "",
		Status:      "",
	}
	return app, nil
}
