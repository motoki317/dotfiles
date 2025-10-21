#!/usr/bin/env zsh

if [ -f ~/.zshrc ]; then
  . ~/.zshrc
fi

# Load common
if [ -f ~/.profile.common ]; then
  . ~/.profile.common
fi

# Load local
if [ -f ~/.zshenv.local ]; then
  . ~/.zshenv.local
fi

