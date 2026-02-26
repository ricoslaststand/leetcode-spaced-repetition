#!/bin/bash

OS="$(uname -s)"

case "$OS" in
  Darwin)
    osascript -e 'display notification "Task complete — ready for next instruction" with title "Claude Code" sound name "Glass"'
    ;;
  Linux)
    notify-send --urgency=normal --app-name="Claude Code" "Claude Code" "Task complete — ready for next instruction"
    ;;
esac
