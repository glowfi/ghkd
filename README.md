# 🎹 ghkd - Go Hotkey Daemon

**ghkd** is a blazing fast, system-level hotkey daemon for Linux. 🚀

It reads input directly from the kernel (`evdev`), which means it works **everywhere**: Wayland, X11, and even the TTY console. No more fighting with compositor-specific config files!

## ✨ Features

- **🖥️ Display Server Agnostic:** Works perfectly on Hyprland, Sway, Gnome, KDE, X11, or no GUI at all.
- **📦 Zero Dependencies:** Written in pure Go. No bloat, no X11 libraries required.
- **⚡ 3 Execution Modes:**
    1.  **Run:** Execute simple commands.
    2.  **Script:** Write inline Bash/Python/Node scripts directly in your config.
    3.  **File:** Execute external scripts.
- **🔄 Hot Reload:** Update your config on the fly without restarting.
- **🛡️ Smart Detection:** Automatically detects keyboards and ignores mice/peripherals.
- **👻 Background Mode:** Built-in daemon management (start, stop, reload).

## 🛠️ Installation

### Prerequisites

- Linux Kernel
- Go 1.21+ (to build)

### Build from Source

```bash
git clone https://github.com/glowfi/ghkd.git
cd ghkd
go build -o ghkd ./main.go
sudo mv ghkd /usr/local/bin/
cd ..
rm -rf ghkd
```

## 🔐 Permissions Setup

Since **ghkd** reads hardware input directly, it needs access to `/dev/input/`. You do **not** need root if you add your user to the `input` group.

1.  **Add user to group:**
    ```bash
    sudo usermod -aG input $USER
    ```
2.  **Reboot** (or log out & log in) for changes to take effect. ⚠️

## ⚙️ Configuration

Create your config at `~/.config/ghkd/config.yaml`.

### 🔑 Syntax

- **Modifiers:** `ctrl`, `alt`, `shift`, `super` (meta/win).
- **Keys:** `a-z`, `0-9`, `f1-f12`, `print`, `space`, `enter`, etc.
- **Media:** `volumeup`, `mute`, `playpause`, `brightnessup`, etc.

### 📝 Example Config

```yaml
settings:
    log_level: info

keybindings:
    # 🚀 MODE 1: Simple Command
    - name: Terminal
      keys: ctrl+alt+t
      run: alacritty

    - name: Volume Up
      keys: volumeup
      run: pactl set-sink-volume @DEFAULT_SINK@ +5%

    # 🐍 MODE 2: Inline Script
    # Great for logic involving variables or pipes!
    - name: System Info
      keys: super+i
      interpreter: python3
      script: |
          import platform
          print(f"OS: {platform.system()}")

    - name: Screenshot
      keys: print
      interpreter: bash
      script: |
          file="$HOME/Pictures/screen-$(date +%s).png"
          grim "$file"
          notify-send "📸 Screenshot taken"

    # 📂 MODE 3: External File
    - name: Backup
      keys: super+b
      file: ~/scripts/backup.sh
```

## 💻 CLI Usage

Manage the daemon easily with flags.

| Flag                | Description                               |
| :------------------ | :---------------------------------------- |
| `-b` `--background` | 👻 Run ghkd in the background.            |
| `-r` `--reload`     | 🔄 Reload config of the running instance. |
| `-k` `--kill`       | 💀 Gracefully kill the running instance.  |
| `-c` `--config`     | 📂 Use a custom config path.              |
| `-v` `--version`    | ℹ️ Show version.                          |

### ⚡ Quick Workflow

1.  **Start the daemon:**
    ```bash
    ghkd -b
    ```
2.  **Edit your config file.**
3.  **Apply changes:**
    ```bash
    ghkd -r
    ```
4.  **Stop it:**
    ```bash
    ghkd -k
    ```

## ❓ Troubleshooting

- **⛔ "Permission denied":**
  Run `groups` in your terminal. If you don't see `input`, run the permission setup command above and **reboot**.

- **⌨️ "No keyboards found":**
  ghkd filters out mice/power buttons strictly. Ensure your kernel sees your device as a keyboard via `cat /proc/bus/input/devices`.

- **⚠️ "Daemon already running":**
  ghkd uses a lock file at `/tmp/ghkd.pid`. If it crashed hard, run `ghkd -k` to clean it up.
