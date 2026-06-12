package main

import (
	"fmt"
	"os"

	. "SteamManifest/models"

	tea "charm.land/bubbletea/v2"
)


var art = `┏┓┏┳┓┏┓┏┓┳┳┓  ┳┳┓┏┓┳┓┳┏┓┏┓┏┓┏┳┓
┗┓ ┃ ┣ ┣┫┃┃┃  ┃┃┃┣┫┃┃┃┣ ┣ ┗┓ ┃ 
┗┛ ┻ ┗┛┛┗┛ ┗  ┛ ┗┛┗┛┗┻┻ ┗┛┗┛ ┻ `






func main() {
    p := tea.NewProgram(InitialMenu())
    if _, err := p.Run(); err != nil {
        fmt.Printf("Alas, there's been an error: %v", err)
        os.Exit(1)
    }
}