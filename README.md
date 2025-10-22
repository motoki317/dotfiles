# dotfiles

## Bootstrap

Prerequisites:
- `nix` installed

1. Clone this into home directory.
2. Setup home-manager:

```shell
nix run home-manager/master -- init --switch
```

(Running `home-manager switch` will be suffice for later package updates)

3. If shell aliases are correctly loaded, you can now launch Neovim with `n`.
