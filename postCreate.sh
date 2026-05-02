git config --global user.email 'brotherlogic-automation@gmail.com'
git config --global user.name 'Brotherlogic Automation'

tic -x ghostty.terminfo

# Install tmux and emacs
sudo apt-get update && sudo apt-get install -y emacs tmux

# Install gh dash
gh extension install dlvhdr/gh-dash

# Auto-start tmux in zsh and bash
TMUX_BLOCK=$(cat << 'EOF'
if [ -z "$TMUX" ] && [ -n "$PS1" ]; then
  /workspaces/kubebrainz/start-tmux.sh && tmux attach-session -t prod
fi
EOF
)

grep -q "tmux attach-session" ~/.zshrc || echo "$TMUX_BLOCK" >> ~/.zshrc
grep -q "tmux attach-session" ~/.bashrc || echo "$TMUX_BLOCK" >> ~/.bashrc

# Ensure the session is created
/workspaces/kubebrainz/start-tmux.sh
