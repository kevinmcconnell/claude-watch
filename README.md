# claude-watch

A Waybar module that shows your Claude Code usage level. Reads your existing Claude Code OAuth credentials and displays the weekly usage percentage with a robot icon.

## Install

```sh
make install
```

This builds the binary and installs it to `~/.local/bin/claude-watch`. To install elsewhere:

```sh
PREFIX=/usr/local make install
```

## Prerequisites

You need to be logged into Claude Code. The tool reads your OAuth token from `~/.claude/.credentials.json`, which Claude Code creates automatically.

## Waybar setup

Add the module to your waybar config (typically `~/.config/waybar/config.jsonc`):

```jsonc
// Add to modules-right (or wherever you prefer):
"modules-right": [
  // ...
  "custom/claude",
  "clock"
],

// Module configuration:
"custom/claude": {
  "exec": "$HOME/.local/bin/claude-watch",
  "return-type": "json",
  "interval": 600,
  "format": "{}",
  "states": {
    "warning": 80
  }
}
```

Add preferred styling to your waybar CSS (typically `~/.config/waybar/style.css`). For example:

```css
#custom-claude {
  min-width: 12px;
  margin: 0 7.5px;
}

#custom-claude.warning {
  color: #f7768e;
}
```

Reload waybar to activate:

```sh
killall -SIGUSR2 waybar
```

## What it shows

- **Bar**: `󰚩 41%` weekly (7-day) usage percentage
- **Tooltip**: session (5-hour) and weekly stats with reset times
- **Warning state**: turns red at 80%+ usage (driven by waybar's `states` config)

Errors (expired token, network issues) show `󰚩 –` with the error in the tooltip.
