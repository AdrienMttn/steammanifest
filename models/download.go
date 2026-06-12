package models

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"time"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)
var styleSuccess = lipgloss.NewStyle().Foreground(lipgloss.Color("#00ff00"))
var styleError = lipgloss.NewStyle().Foreground(lipgloss.Color("#ff0000"))
var styleInfo = lipgloss.NewStyle().Foreground(lipgloss.Color("#ffff00"))
var body = lipgloss.NewStyle().Padding(1, 2)


type download struct {
	ti       textinput.Model
	submited bool
	allMsg []downloadMsg
	pb progressBar
}

type progressBar struct {
	state int
	width int
	value int
}

func initialDownload() download {
	ti := textinput.New()
	ti.Placeholder = "Enter the Steam AppID"
	ti.Focus()
	ti.CharLimit = 156
	ti.SetWidth(20)
	return download{
		ti: ti,
		submited: false,
		allMsg: []downloadMsg{},
		pb: progressBar{
			state: 0,
			value: 0,
		},
	}
}

func (d download) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		d.pb.width = msg.Width
	case tea.KeyPressMsg:
		switch msg.String() {
			
		case "esc":
			return menu{
				choices: []string{"Download Game", "Setting", "Quit"},
			}, nil
		case "enter":
			d.submited = true
			gameID := d.ti.Value()
			d.allMsg = append(d.allMsg, initialDownloadMsg(styleInfo.Render("\nDownloading manifest for game ID: "+gameID+"\n")))
			d.pb.state = 0
      d.pb.value = 0
			return d, checkGameExistCmd(gameID)
		}
	case gameExistMsg:
		if !msg.ok {
			d.pb.value=10
			d.pb.state = 1
			d.allMsg = append(d.allMsg, initialDownloadMsg(styleError.Render("Game not found. Please check the AppID and try again.")))
			d.allMsg = append(d.allMsg, initialDownloadMsg(styleInfo.Render("\nPress ESC to return to the main menu.")))
			return d, nil
		}
		d.pb.value=10
		d.pb.state = 0
		d.allMsg = append(d.allMsg, initialDownloadMsg(styleSuccess.Render("Game found. Starting download ...")))
		return d, downloadManifestCmd(msg.appID)

	case downloadResultMsg:
		if msg.ok {
			d.pb.value=25
			d.pb.state = 0
			d.allMsg = append(d.allMsg, initialDownloadMsg(styleSuccess.Render("Download completed successfully.")))
		} else {
			d.pb.value=25
			d.pb.state = 1
			d.allMsg = append(d.allMsg, initialDownloadMsg(styleError.Render("Failed to download the manifest. Please try again later.")))
			d.allMsg = append(d.allMsg, initialDownloadMsg(styleInfo.Render("\nPress ESC to return to the main menu.")))
		}
		return d, unzipManifestCmd(msg.appID)
	case unzipResultMsg:
		if msg.ok {
			d.pb.value=75
			d.pb.state = 0
			d.allMsg = append(d.allMsg, initialDownloadMsg(styleSuccess.Render("Manifest unzipped successfully.")))
		} else {
			d.pb.value=75
			d.pb.state = 1
			d.allMsg = append(d.allMsg, initialDownloadMsg(styleError.Render("Failed to unzip the manifest. Please try again later.")))
			d.allMsg = append(d.allMsg, initialDownloadMsg(styleInfo.Render("\nPress ESC to return to the main menu.")))
		}
		return d, restartSteamCmd()
	case restartSteamMsg:
		if msg.ok {
			d.pb.value=100
			d.pb.state = 0
			d.allMsg = append(d.allMsg, initialDownloadMsg(styleSuccess.Render("Steam restarted successfully. Your new manifest should now be active.")))
			d.allMsg = append(d.allMsg, initialDownloadMsg("If you don't see the changes, make sure you have injected LumaCore.dll\n( You can see this in the 'Setting' menu ) and restart Steam again."))
			d.allMsg = append(d.allMsg, initialDownloadMsg(styleInfo.Render("\nPress ESC to return to the main menu.")))
		} else {
			d.pb.value=100
			d.pb.state = 0
			d.allMsg = append(d.allMsg, initialDownloadMsg(styleError.Render("Failed to restart Steam. Please restart it manually to apply the new manifest.")))
			d.allMsg = append(d.allMsg, initialDownloadMsg(styleInfo.Render("\nPress ESC to return to the main menu.")))
		}
	}
	if !d.submited {
		d.ti, cmd = d.ti.Update(msg)
		return d, cmd
	}
	return d, nil
}

// Messages for asynchronous steps
type gameExistMsg struct {
	ok bool
	appID  string
}

type downloadResultMsg struct {
	ok    bool
	appID string
}


type unzipResultMsg struct {
	ok    bool
	appID string
}

type restartSteamMsg struct {
	ok bool
}

// Cmds to run the network operations asynchronously and return messages
func checkGameExistCmd(appID string) tea.Cmd {
	return func() tea.Msg {
		exists := (&download{}).GameExist(appID)
		return gameExistMsg{ok: exists, appID: appID}
	}
}

func downloadManifestCmd(appID string) tea.Cmd {
	return func() tea.Msg {
		ok := (&download{}).downloadManifest(appID)
		return downloadResultMsg{ok: ok, appID: appID}
	}
}

func unzipManifestCmd(appID string) tea.Cmd {
	return func() tea.Msg {
		ok := (&download{}).unzipManifest(appID)
		return unzipResultMsg{ok: ok, appID: appID}
	}
}

func restartSteamCmd() tea.Cmd {
	return func() tea.Msg {
		ok := (&download{}).restartSteam()
		return restartSteamMsg{ok: ok}
	}
}

func (d download) View() tea.View {
	var c *tea.Cursor
	if !d.ti.VirtualCursor() {
		c = d.ti.Cursor()
		c.Y += lipgloss.Height(art)
	}

	s := ""
	if d.submited {
		s = lipgloss.JoinVertical(lipgloss.Top, titleStyle.Render(art)+"\n\n")
		s += renderBar(d.pb.value, 20, d.pb.state)
		for _, msg := range d.allMsg {
			s += "\n" + msg.View()

		}
	}else {
		s = lipgloss.JoinVertical(lipgloss.Top, titleStyle.Render(art)+"\n\n",d.ti.View())
	}
	
	v := tea.NewView(s)
	v.Cursor = c
	return v

}

func renderBar(value int, width int, state int) string {
    if width <= 0 {
        width = 20
    }
    if value < 0 {
        value = 0
    }
    if value > 100 {
        value = 100
    }

    filled := (value * width) / 100
    bar := "["
    for i := 0; i < width; i++ {
        if i < filled {
						if state == 0 {
							bar += styleSuccess.Render("█")
						} else {
							bar += styleError.Render("█")
						}
        } else {
            bar += "░"
        }
    }
    bar += fmt.Sprintf("] %d%%", value)
    return bar
}

func (d download) GameExist(appID string) bool {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://walftech.com/proxy.php?id="+appID, nil)
	if err != nil {
		return false
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	req.Header.Set("Referer", "https://short.walftech.com/?id="+appID)
	req.Header.Set("cookie", "cf_clearance=")
	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == 200
}

func (d download) downloadManifest(appID string) (bool) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://walftech.com/proxy.php?id="+appID, nil)
	if err != nil {
		return false
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	req.Header.Set("Referer", "https://short.walftech.com/?id="+appID)
	req.Header.Set("cookie", "cf_clearance=")
	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
    return false
	}
	if err := os.WriteFile("temp/"+appID+".zip", data, 0644); err != nil {
    return false
	}
	fmt.Print(io.ReadAll(resp.Body))
	return true
}

func (d download) unzipManifest(appID string) bool {
	zipPath := "temp/" + appID + ".zip"
	zipReader, err := zip.OpenReader(zipPath)
	if err != nil {
		return false
	}
	defer zipReader.Close()
	targetName := appID + ".lua"
	destPath := "C:\\Program Files (x86)\\Steam\\config\\stplug-in/" + targetName
	for _, file := range zipReader.File {
		if file.Name == targetName {
			srcFile, err := file.Open()
			if err != nil {
				return false
			}
			defer srcFile.Close()
			destFile, err := os.Create(destPath)
			if err != nil {
				return false
			}
			defer destFile.Close()
			if _, err := io.Copy(destFile, srcFile); err != nil {
				return false
			}
			zipReader.Close()
			err = os.Remove(zipPath)
			if err != nil {
				return false
			}
			return true
		}
	}
	return false
}

func (d download) restartSteam() bool {
	defer time.Sleep(500 * time.Millisecond)
	// arrêter Steam
  switch runtime.GOOS {
  case "windows":
      exec.Command("taskkill", "/IM", "steam.exe", "/F").Run()
  case "darwin":
      exec.Command("pkill", "-f", "Steam").Run()
  default: // linux / autres
      exec.Command("pkill", "-f", "steam").Run()
  }
	// relancer Steam
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

func (d download) Init() tea.Cmd {
	return textinput.Blink
}


type downloadMsg struct {
	msg string
}

func initialDownloadMsg(msg string) downloadMsg {
	return downloadMsg{
		msg: msg,
	}
}

func (dMsg downloadMsg) String() string {
	return fmt.Sprint(dMsg.msg)
}

func (dMsg downloadMsg) View() string {
	return dMsg.String()
}

