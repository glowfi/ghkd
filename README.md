<div align="center">

# ⌨️ ghkd

<img src="./images/logo.png" width="420"/>

**A system-level hotkey daemon for Linux**

Wayland • X11 • TTY • Anywhere

![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go)
![Linux](https://img.shields.io/badge/platform-linux-success?logo=linux)
![Zero Dependencies](https://img.shields.io/badge/dependencies-none-lightgrey)
![License](https://img.shields.io/github/license/glowfi/ghkd)

</div>

---

## ✨ Overview

**ghkd** is a fast, kernel-level hotkey daemon for Linux.

Unlike compositor-bound solutions, ghkd reads input directly from **evdev**, allowing global hotkeys that work everywhere:

- Wayland compositors
- X11 environments
- Desktop environments
- TTY console sessions

No compositor configs. No display server coupling.

---

## 🚀 Features

- 🖥 **Display Server Agnostic** — Wayland, X11, or no GUI
- ⚡ **Kernel-Level Input** — reads directly from `/dev/input`
- 🧱 **Zero Dependencies** — pure Go binary
- 🔁 **Hot Reload** — update config without restarting
- 🧠 **Smart Device Detection** — ignores mice & peripherals
- 🔧 **Daemon Management** — built-in background control

### Execution Modes

| Mode       | Description                          |
| ---------- | ------------------------------------ |
| **Run**    | Execute commands directly            |
| **Script** | Inline Bash/Python/Node/Ruby scripts |
| **File**   | Execute external scripts             |

---

## 📦 Installation

### Option 1 — Build From Source

```bash
git clone https://github.com/glowfi/ghkd.git
cd ghkd
go build -o ghkd ./main.go
sudo mv ghkd /usr/local/bin/
```

---

### Option 2 — Download Release Binary

```bash
wget https://github.com/glowfi/ghkd/releases/download/v1.0.0/ghkd_linux_amd64
chmod +x ghkd_linux_amd64
sudo mv ghkd_linux_amd64 /usr/local/bin/ghkd
```

---

## 🔐 Permissions Setup

ghkd accesses hardware input devices.

Add your user to the `input` group:

```bash
sudo usermod -aG input $USER
```

Then **log out or reboot**.

---

## ⚙️ Configuration

Create:

```
~/.config/ghkd/config.yaml
```

---

## ⌨️ Keybinding Rules

### Core Rules

1. Exactly **one main key** per binding
2. Unlimited modifier keys allowed
3. Case-insensitive syntax
4. Keys joined using `+`

---

### Supported Keys

Full reference:
[https://raw.githubusercontent.com/glowfi/ghkd/main/internal/hotkey/keymap.go](https://raw.githubusercontent.com/glowfi/ghkd/main/internal/hotkey/keymap.go)

| Category   | Examples                           |
| ---------- | ---------------------------------- |
| Modifiers  | `ctrl`, `alt`, `shift`, `super`    |
| Standard   | `a-z`, `0-9`, `f1-f24`             |
| Navigation | `left`, `right`, `home`, `end`     |
| Special    | `enter`, `space`, `esc`, `tab`     |
| Media      | `volumeup`, `mute`, `brightnessup` |

---

## 🧠 Example Configuration

```yaml
keybindings:
    - name: Terminal
      keys: ctrl+alt+t
      run: alacritty

    - name: Volume Up
      keys: volumeup
      run: pactl set-sink-volume @DEFAULT_SINK@ +5%

    - name: System Info
      keys: super+i
      interpreter: python3
      script: |
          import platform
          print(platform.system())

    - name: Screenshot
      keys: meta+print
      interpreter: bash
      script: |
          file="$HOME/Pictures/screen-$(date +%s).png"
          grim "$file"
          notify-send "Screenshot taken"

    - name: Backup
      keys: super+b
      file: ~/scripts/backup.sh
```

---

## 🖥 CLI Usage

| Flag                 | Description              |
| -------------------- | ------------------------ |
| `-b`, `--background` | Run daemon in background |
| `-r`, `--reload`     | Reload configuration     |
| `-k`, `--kill`       | Stop running instance    |
| `-c`, `--config`     | Custom config path       |
| `-v`, `--version`    | Show version             |

### Quick Start

```bash
ghkd -b -c ~/.config/ghkd/config.yaml
```

---

## 🛠 Troubleshooting

### Permission denied

Ensure your user belongs to `input` group:

```bash
groups
```

Reboot after adding.

---

### No keyboards found

Check kernel device detection:

```bash
cat /proc/bus/input/devices
```

---

### Daemon already running

```bash
ghkd -k
```

(removes stale lock file)

---

## 🎯 Design Goals

- Universal hotkeys
- Minimal runtime overhead
- Display-server independence
- Predictable configuration
- System-level reliability

---

## 🤝 Contributing

Issues and PRs are welcome.

Small focused contributions preferred.

---

## 📄 License

GPL-3.0
