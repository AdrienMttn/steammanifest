package config

import (
	"log"
	"strings"

	"gopkg.in/ini.v1"
)

type Config struct {
	SteamPath string `ini:"STEAM_PATH"`
}

var AppConfig Config

func LoadConfig() {
  cfg, err := ini.Load("app.ini")
  if err != nil {
      log.Fatalf("Erreur de lecture du fichier INI: %v", err)
  }
  AppConfig.SteamPath= cfg.Section("General").Key("STEAM_PATH").String()
}

func SaveConfig() bool{
	cfg := ini.Empty()
	if (AppConfig.SteamPath[len(AppConfig.SteamPath)-1] != '\\' && AppConfig.SteamPath[len(AppConfig.SteamPath)-1] != '/') {
		AppConfig.SteamPath += "/"
		
	}
	AppConfig.SteamPath = strings.ReplaceAll(AppConfig.SteamPath, "\\", "/")
	cfg.Section("General").Key("STEAM_PATH").SetValue(AppConfig.SteamPath)
	err := cfg.SaveTo("app.ini")
	if err != nil {
		log.Fatalf("Erreur lors de l'enregistrement du fichier INI: %v", err)
		return false
	}
	return true
}