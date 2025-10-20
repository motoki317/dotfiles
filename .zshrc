# pure (zsh theme)
fpath+=($HOME/.zsh/pure)
autoload -U promptinit
promptinit
prompt pure

# Load common
if [ -f .profile.common ]; then
  . .profile.common
fi

# Load local
if [ -f ~/.zshrc.local ]; then
  . ~/.zshrc.local
fi

