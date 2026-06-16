# 🎮 SteamManifest

> A powerful tool to download game manifests directly from Steam. Perfect for developers, modders, and enthusiasts!

![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat-square&logo=go)
![Stars](https://img.shields.io/github/stars/AdrienMttn/steammanifest?style=flat-square&logo=github)

---

## 📋 Table of Contents

- [🎮 SteamManifest](#-steammanifest)
  - [📋 Table of Contents](#-table-of-contents)
  - [✨ Features](#-features)
  - [📦 Requirements](#-requirements)
  - [💾 Installation](#-installation)
    - [Download Pre-built Release](#download-pre-built-release)
    - [Build from Source](#build-from-source)
  - [🚀 Quick Start](#-quick-start)
  - [⚙️ Configuration](#️-configuration)
  - [📚 Usage Guide](#-usage-guide)
    - [Finding an App ID](#finding-an-app-id)
    - [Troubleshooting Failed Downloads](#troubleshooting-failed-downloads)
  - [🐛 Troubleshooting](#-troubleshooting)
  - [TODO ✅](#todo-)
  - [⭐ Star History](#-star-history)

---

## ✨ Features

- 🚀 **Fast & Efficient** - Quickly download manifests using Steam's API
- ⚙️ **Easy Configuration** - Intuitive settings panel for Steam path setup
- 🎯 **Simple Interface** - User-friendly menu for manifest downloads
- 📊 **Lightweight** - Minimal resource usage

---

## 📦 Requirements

- **Operating System**: Windows
- **Steam**: Must be installed and running on your computer
- **Go**: 1.21 or higher (only if building from source)

---

## 💾 Installation

### Download Pre-built Release

1. Go to the [releases page](https://github.com/AdrienMttn/steammanifest/releases/latest)
2. Download the latest version for your operating system
3. Extract the archive to your desired location
4. Run the executable

### Build from Source

```bash
git clone https://github.com/AdrienMttn/steammanifest.git
cd steammanifest
go build -o steammanifest
```

---

## 🚀 Quick Start

1. **Launch the Application**
   - Run the SteamManifest executable
   - You can have an pop-up window asking he protect your computer, click on "More info" and then "Run anyway" this is normal because the application is not signed

2. **Configure Steam Path**
   - Navigate to **Settings** ⚙️
   - Set the path to your Steam installation directory
   - Save the configuration by pressing **Enter**

3. **Enable LumaCore** (Optional)
   - Go to **Settings** → "Enable/Disable LumaCore"
   - Press **Enter** to **"Unable LumaCore"** or **"Disable LumaCore"** if needed

4. **Download a Manifest**
   - Select **"Download Game"** from the menu
   - Enter the App ID of the desired game
   - Find App IDs on [SteamDB](https://steamdb.info/)
   - Press **ENTER** and wait for completion

---

## ⚙️ Configuration

The application stores settings in `app.ini`. You can manually edit this file or use the built-in settings panel:

```ini
[general]
path = C:/Program Files (x86)/Steam
```

## 📚 Usage Guide

### Finding an App ID

1. Visit [SteamDB](https://steamdb.info/)
2. Search for your game
3. Copy the **App ID** from the URL or main page
4. Paste it into the SteamManifest application

### Troubleshooting Failed Downloads

- ✅ Ensure Steam is running
- ✅ Verify the App ID is correct
- ✅ Check that your Steam path is correctly configured
- ✅ Make sure you own the game on your Steam account or you have lumaCore enabled

---

## 🐛 Troubleshooting

**Application won't start**
- Verify that Steam is installed and running
- Check that the Steam path is correctly set in settings

**Download fails**
- Confirm the App ID is valid
- Ensure Steam is actively running
- Try restarting the application

**Settings not saving**
- Check file permissions for `app.ini`
- Ensure the application has write access to its directory

---

## TODO ✅

- [ ] Add support for Linux and macOS
- [ ] Implement denuvo bypass for manifest downloads
- [ ] Implement multiplayer fix using [online-fix.me](https://online-fix.me/)
- [ ] Add a feature to update manifests automatically

---

## ⭐ Star History

Show your support by starring this project! ⭐

[![Star History Chart](https://api.star-history.com/svg?repos=AdrienMttn/steammanifest&type=date&legend=top-left)](https://www.star-history.com/#AdrienMttn/steammanifest&type=date&legend=top-left)

---

**Made with ❤️ by [AdrienMttn](https://github.com/AdrienMttn)**