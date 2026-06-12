package models

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)


var disableStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#ff0000"))
type settings struct {
	choices []string // items on the to-do list
	cursor  int      // which to-do list item our cursor is pointing at
}

type lumacore struct {
	isEnabled bool
}

type lumacoreResponse struct {
	URL string `json:"url"`
	AssetsURL string `json:"assets_url"`
	UploadURL string `json:"upload_url"`
	HTMLURL string `json:"html_url"`
	ID int `json:"id"`
	Assets []struct {
		BrowserDownloadURL string `json:"browser_download_url"`
	}
}

func check() bool {
	destPath := "C:\\Program Files (x86)\\Steam/"
	_, err1 := os.Stat(destPath + "LumaCore.dll")
	_, err2 := os.Stat(destPath + "dwmapi.dll")
	return !os.IsNotExist(err1) && !os.IsNotExist(err2)
}

func initialLumaCore() lumacore {
	return lumacore{
		isEnabled: check(),
	}
}

func (l lumacore) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	// Is it a key press?
	case tea.KeyPressMsg:
		switch msg.String() {
		case "esc":
			return initialSettings(), nil
		case "enter", "space":
			switch {
				case l.isEnabled:
					if l.disable() {
						l.isEnabled = !l.isEnabled
					}
				default:	
					if l.enable(){
						l.isEnabled = !l.isEnabled
					}
			}
		}
	}
	return l, nil
}

func (l lumacore) View() tea.View {
	style := disableStyle
	checked := " "
	text := "Enable LumaCore"
	if l.isEnabled {
		style = activeStyle
		checked = "x"
		text = "Disable LumaCore"
	}
	str := titleStyle.Render(art) + "\n\n"
	str += fmt.Sprintf(" [%s] %s\n\n", style.Render(checked), style.Render(text))
	return tea.NewView(str)
}
func (l lumacore) enable() bool {
	// On ferme Steam pour éviter les conflits de fichiers
	switch runtime.GOOS {
  case "windows":
      exec.Command("taskkill", "/IM", "steam.exe", "/F").Run()
  case "darwin":
      exec.Command("pkill", "-f", "Steam").Run()
  default: // linux / autres
      exec.Command("pkill", "-f", "steam").Run()
  }
	rep, err := http.Get("https://api.github.com/repos/KoriaPolis/LumaCore/releases/latest")
	if err != nil {
		return false
	}
	defer rep.Body.Close()

	body, err := io.ReadAll(rep.Body)

	if rep.StatusCode != http.StatusOK {
		return false
	}
	var data lumacoreResponse
  if err := json.Unmarshal(body, &data); err != nil {
    return false
  }
	rep, err = http.Get(data.Assets[1].BrowserDownloadURL)
	if err != nil {
		return false
	}
	defer rep.Body.Close()
	file, err := io.ReadAll(rep.Body)
	if err != nil {
		return false
	}
	zipPath := "temp/LumaCore.zip"
	err = os.WriteFile(zipPath, file, 0644)
	zipReader, err := zip.OpenReader(zipPath)
	if err != nil {
		fmt.Printf("Error opening zip file: %v", err)
		return false
	}
	destPath := "C:\\Program Files (x86)\\Steam/"
	for _, file := range zipReader.File {
		srcFile, err := file.Open()
		if err != nil {
				fmt.Printf("Error opening file in zip: %v", err)
				return false
		}
		
		outFile, err := os.Create(destPath + file.Name)
		if err != nil {
				fmt.Printf("Error creating file: %v", err)
				return false
		}
		defer outFile.Close()
		if _, err := io.Copy(outFile, srcFile); err != nil {
			fmt.Printf("Error copying file: %v", err)
			return false
		}
		
	}
	zipReader.Close()
	err = os.Remove(zipPath)
	if err != nil {
		fmt.Printf("Error removing zip file: %v", err)
		return false
	}
	// Redémarrer Steam pour appliquer les changements
	var cmd *exec.Cmd
  switch runtime.GOOS {
  case "windows":
      steamPath := `C:\Program Files (x86)\Steam\steam.exe`
      cmd = exec.Command(steamPath)
  case "darwin":
      cmd = exec.Command("open", "-a", "Steam")
  default:
      cmd = exec.Command("steam")
  }
  if err := cmd.Start(); err != nil {
      fmt.Println("restartSteam: start error:", err)
      return false
  }
	return true
}

func (l lumacore) disable() bool {
	// On ferme Steam pour éviter les conflits de fichiers
	switch runtime.GOOS {
  case "windows":
      exec.Command("taskkill", "/IM", "steam.exe", "/F").Run()
  case "darwin":
      exec.Command("pkill", "-f", "Steam").Run()
  default: // linux / autres
      exec.Command("pkill", "-f", "steam").Run()
  }
	time.Sleep(1000 * time.Millisecond)
	dllToDelete := []string{"LumaCore.dll","dwmapi.dll"}
	destPath := "C:\\Program Files (x86)\\Steam/"
	for _, dll := range dllToDelete {
		err := os.Remove(destPath + dll)
		if err != nil {
			fmt.Printf("Error removing file: %v", err)
			return false
		}
	}
	// Redémarrer Steam pour appliquer les changements
	var cmd *exec.Cmd
	switch runtime.GOOS {
  case "windows":
      steamPath := `C:\Program Files (x86)\Steam\steam.exe`
      cmd = exec.Command(steamPath)
  case "darwin":
      cmd = exec.Command("open", "-a", "Steam")
  default:
      cmd = exec.Command("steam")
  }
  if err := cmd.Start(); err != nil {
      fmt.Println("restartSteam: start error:", err)
      return false
  }
	return true
}



func initialSettings() settings {
	return settings{
		choices: []string{"Edit steam PATH", "Enable/Disable LumaCore"},
	}
}


func (s settings) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	// Is it a key press?
	case tea.KeyPressMsg:
		// Cool, what was the actual key pressed?
		switch msg.String() {
		
		case "esc":
			return InitialMenu(), nil
		// The "up" and "k" keys move the cursor up
		case "up", "k":
			if s.cursor > 0 {
				s.cursor--
			}
		// The "down" and "j" keys move the cursor down
		case "down", "j":
			if s.cursor < len(s.choices)-1 {
				s.cursor++
			}

		// The "enter" key and the space bar toggle the selected state
		// for the item that the cursor is pointing at.
		case "enter", "space":
			switch s.cursor {
			case 0:
				return initialDownload(), nil
			case 1:
				return initialLumaCore(), nil
			}
		}
	}
	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return s, nil
}

func (s settings) View() tea.View {
	// The header
	str := titleStyle.Render(art) + "\n\n"

	// Iterate over our choices
	for i, choice := range s.choices {

		// Is the cursor pointing at this choice?
		cursor := " "  // no cursor
		checked := " " // not selected
		if s.cursor == i {
			cursor = activeStyle.Render(">")
			checked = activeStyle.Render("x")
			choice = activeStyle.Render(choice)
		}

		// Render the row
		str += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
	}
	// Send the UI for rendering
	return tea.NewView(str)
}

func (s settings) Init() tea.Cmd {
	return nil
}

func (l lumacore) Init() tea.Cmd {
	return nil
}