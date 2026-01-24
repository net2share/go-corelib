package tui

import "fmt"

// BannerConfig holds configuration for banner rendering.
type BannerConfig struct {
	AppName   string
	Version   string
	BuildTime string
	ASCII     string // Optional ASCII art
}

// PrintBanner displays an application banner with version info.
func PrintBanner(cfg BannerConfig) {
	fmt.Println()
	if cfg.ASCII != "" {
		fmt.Println(TitleStyle.Render(cfg.ASCII))
	}
	if cfg.BuildTime != "" {
		fmt.Printf("%s v%s (built %s)\n\n", cfg.AppName, cfg.Version, cfg.BuildTime)
	} else {
		fmt.Printf("%s v%s\n\n", cfg.AppName, cfg.Version)
	}
}

// PrintSimpleBanner displays a simple banner with just app name and version.
func PrintSimpleBanner(appName, version, buildTime string) {
	fmt.Println()
	if buildTime != "" {
		fmt.Printf("%s v%s (built %s)\n", appName, version, buildTime)
	} else {
		fmt.Printf("%s v%s\n", appName, version)
	}
	fmt.Println()
}
