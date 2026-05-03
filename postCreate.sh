git config --global user.email 'brotherlogic-automation@gmail.com'
git config --global user.name 'Brotherlogic Automation'

tic -x ghostty.terminfo

# Install tmux and emacs
sudo apt-get update && sudo apt-get install -y emacs tmux

# Install gh dash
gh extension install dlvhdr/gh-dash

# Setup tmux for Ghostty and graphics support
cat << 'EOF' > "$HOME/.tmux.conf"
# Set proper default terminal for better TUI rendering
set -g default-terminal "tmux-256color"

# Standard tmux best practice for modern TUIs
set -s escape-time 0

# Allow programs to use the terminal's graphics capabilities
set -g allow-passthrough on

# Support Ghostty terminal capabilities
set -as terminal-overrides ',xterm-ghostty:Sync:Tc,ghostty:Sync:Tc'
EOF


# Auto-start tmux in zsh and bash
TMUX_BLOCK=$(cat << 'EOF'
if [ -z "$TMUX" ] && [ -n "$PS1" ]; then
  export LANG=C.UTF-8
  export LC_ALL=C.UTF-8
  /workspaces/kubebrainz/start-tmux.sh && tmux -u attach-session -t prod
fi
EOF
)

grep -q "tmux attach-session" ~/.zshrc || echo "$TMUX_BLOCK" >> ~/.zshrc
grep -q "tmux attach-session" ~/.bashrc || echo "$TMUX_BLOCK" >> ~/.bashrc

# Ensure the session is created
/workspaces/kubebrainz/start-tmux.sh
