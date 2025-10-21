#!/usr/bin/env zsh

# pure (zsh theme)
fpath+=("$HOME/.zsh/pure")
autoload -U promptinit
promptinit
prompt pure
