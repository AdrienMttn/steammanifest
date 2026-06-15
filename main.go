package main

import (
	"log"
	"os"

	"SteamManifest/config"
	. "SteamManifest/models"

	tea "charm.land/bubbletea/v2"
)

func main() {
    config.LoadConfig()
    p := tea.NewProgram(InitialMenu())
    if _, err := p.Run(); err != nil {
        config.WriteLog("Alas, there's been an error: "+ err.Error())
        log.Fatalf("Alas, there's been an error: %v", err)
        os.Exit(1)
    }
}