#!/bin/bash

# Ensure UTF-8 locale
export LANG=C.UTF-8
export LC_ALL=C.UTF-8

# Ensure the 'prod' session exists
if ! tmux has-session -t prod 2>/dev/null; then
  # Create a new session named 'prod', detached
  tmux -u new-session -A -d -s prod
  
  # Split the window horizontally (-h)
  # The left pane will remain a terminal
  # The right pane will run 'gh dash'
  tmux split-window -h -t prod "gh dash"
fi
