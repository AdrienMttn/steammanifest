package config

import (
	"log"
	"os"
	"strings"

	"gopkg.in/ini.v1"
)

type Config struct {
	SteamPath string `ini:"STEAM_PATH"`
}

var AppConfig Config

func LoadConfig() {
	checkIniFile() 
  cfg, err := ini.Load("app.ini")
  if err != nil {
			WriteLog("Erreur de lecture du fichier INI: " + err.Error())
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
		WriteLog("Erreur lors de l'enregistrement du fichier INI: " + err.Error())
		log.Fatalf("Erreur lors de l'enregistrement du fichier INI: %v", err)
		return false
	}
	return true
}

func checkIniFile() {
	os.Stat("app.ini")
	if _, err := os.Stat("app.ini"); os.IsNotExist(err) {
		// Le fichier n'existe pas, le créer avec les valeurs par défaut
		cfg := ini.Empty()
		cfg.Section("General").Key("STEAM_PATH").SetValue("C:/Program Files (x86)/Steam/")
		err := cfg.SaveTo("app.ini")
		if err != nil {
			WriteLog("Erreur lors de la création du fichier INI: " + err.Error())
			log.Fatalf("Erreur lors de la création du fichier INI: %v", err)
		}
	}
}

func WriteLog(message string) {
	logFile, err := os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Erreur lors de l'ouverture du fichier de log: %v", err)
	}
	defer logFile.Close()
}

// [General]
// STEAM_PATH = C:/Program Files (x86)/Steam/
