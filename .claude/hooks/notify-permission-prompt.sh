#!/bin/bash

OS="$(uname -s)"

case "$OS" in
  Darwin)
    osascript -e 'display notification "Claude Code needs your input" with title "Claude Code" sound name "Ping"'
    ;;
  Linux)
    notify-send --urgency=normal --app-name="Claude Code" "Claude Code" "Claude Code needs your input"
    ;;
esac
