# ghkd - Go Hotkey Daemon

<p align="center">
<img src="./images/logo.png" alt="Project Logo" width=400/>
</p>

**ghkd** is a fast, system-level hotkey daemon for Linux.

It reads input directly from the kernel (`evdev`), which means it works **everywhere**: Wayland, X11, and even the TTY console. No more fighting with compositor-specific config files.

## Features

- **Display Server Agnostic:** Works perfectly on Hyprland, Sway, Gnome, KDE, X11, or no GUI at all.
- **Zero Dependencies:** Written in pure Go. No bloat, no X11 libraries required.
- **3 Execution Modes:**
    1.  **Run:** Execute simple commands.
    2.  **Script:** Write inline Bash/Python/Node/Ruby scripts directly in your config.
    3.  **File:** Execute external scripts.
- **Hot Reload:** Update your config on the fly without restarting.
- **Smart Detection:** Automatically detects keyboards and ignores mice/peripherals.
- **Background Mode:** Built-in daemon management (start, stop, reload).

## Installation

### Prerequisites

- Linux Kernel
- Go 1.21+ (only if building from source)

### Build from Source

```bash
git clone https://github.com/glowfi/ghkd.git
cd ghkd
go build -o ghkd ./main.go
sudo mv ghkd /usr/local/bin/
cd ..
rm -rf ghkd
```

### Install from Releases

```bash
wget "https://github.com/glowfi/ghkd/releases/download/v1.0.0/ghkd_linux_amd64" -O "ghkd_linux_amd64"
chmod +x ghkd_linux_amd64
sudo mv ghkd_linux_amd64 /usr/local/bin/
```

## Permissions Setup

Since **ghkd** reads hardware input directly, it needs access to `/dev/input/`. You do **not** need root if you add your user to the `input` group.

1.  **Add user to group:**
    ```bash
    sudo usermod -aG input $USER
    ```
2.  **Reboot** (or log out & log in) for changes to take effect.

## Configuration

Create your config at `~/.config/ghkd/config.yaml`.

## Key Syntax & Rules

Defining hotkeys is case-insensitive. Keys are combined using the `+` symbol.

### Important Rules

1.  **Exactly One Main Key:** Your binding must have exactly **one** non-modifier key (e.g., `t`, `enter`, `space`). You cannot combine two main keys like `a+b`.
2.  **Modifiers:** You can use as many modifiers as you like (`ctrl`, `alt`, `shift`, `super`).

### Supported Keys

[See full keymap reference](https://raw.githubusercontent.com/glowfi/ghkd/refs/heads/main/internal/hotkey/keymap.go)

| Category       | Available Keys                                                           |
| :------------- | :----------------------------------------------------------------------- |
| **Modifiers**  | `super`, `win`, `ctrl`, `alt`, `shift`                                   |
| **Standard**   | `a-z`, `0-9`, `f1-f24`                                                   |
| **Navigation** | `left`, `right`, `up`, `down`, `home`, `end`                             |
| **Special**    | `space`, `enter`, `tab`, `esc`, `backspace`, `print`, `insert`, `delete` |
| **Media**      | `volumeup`, `volumedown`, `mute`, `playpause`, `brightnessup`            |

### Examples

| Status      | Combo               | Reason                              |
| :---------- | :------------------ | :---------------------------------- |
| **Valid**   | `ctrl+alt+t`        | Modifiers + 1 Main Key.             |
| **Valid**   | `super+shift+enter` | Multiple modifiers are allowed.     |
| **Valid**   | `volumeup`          | Special keys can work alone.        |
| **Invalid** | `ctrl+a+b`          | Error: Two main keys (`a` and `b`). |
| **Invalid** | `ctrl+alt`          | Error: No main key specified.       |

### Example Config

```yaml
keybindings:
    # MODE 1: Simple Command
    - name: Terminal
      keys: ctrl+alt+t
      run: alacritty

    - name: Volume Up
      keys: volumeup
      run: pactl set-sink-volume @DEFAULT_SINK@ +5%

    # MODE 2: Inline Script
    - name: System Info
      keys: super+i
      interpreter: python3
      script: |
          import platform
          print(f"OS: {platform.system()}")

    - name: Screenshot
      keys: meta+print
      interpreter: bash
      script: |
          file="$HOME/Pictures/screen-$(date +%s).png"
          grim "$file"
          notify-send "Screenshot taken"

    # MODE 3: External File
    - name: Backup
      keys: super+b
      file: ~/scripts/backup.sh
```

## CLI Usage

Manage the daemon with these flags:

| Flag                | Description                            |
| :------------------ | :------------------------------------- |
| `-b` `--background` | Run ghkd in the background.            |
| `-r` `--reload`     | Reload config of the running instance. |
| `-k` `--kill`       | Gracefully kill the running instance.  |
| `-c` `--config`     | Use a custom config path.              |
| `-v` `--version`    | Show version.                          |

### Quick Start

```sh
ghkd -b -c ~/.config/ghkd/config.yaml
```

## Troubleshooting

- **"Permission denied":**
  Run `groups` in your terminal. If you don't see `input`, run the permission setup command above and **reboot**.

- **"No keyboards found":**
  ghkd filters out mice/power buttons. Ensure your kernel sees your device as a keyboard via `cat /proc/bus/input/devices`.

- **"Daemon already running":**
  ghkd uses a lock file at `/tmp/ghkd.pid`. If it crashed, run `ghkd -k` to clean it up.
