# Neovim

LazyVim-based config ([docs](https://lazyvim.github.io/installation)), installed via Nix — `~/.nix-profile/bin/nvim` is a wrapper around the nix-store binary. Custom plugins live in `lua/plugins/custom.lua`.

## Known issues

### macOS 26: SIGKILL (Code Signature Invalid) on startup

macOS 26 tightened runtime page validation for ad-hoc linker-signed binaries: the treesitter parser `.so` files in `~/.local/share/nvim/site/parser/` pass `codesign -vv` static verification but fail kernel-level page validation when `dlopen`ed by the Nix-store nvim. Re-sign after `:TSUpdate` or plugin updates:

```sh
for f in ~/.local/share/nvim/site/parser/*.so; do codesign -fs - "$f"; done
codesign -fs - ~/.local/share/nvim/lazy/telescope-fzf-native.nvim/build/libfzf.so
codesign -fs - ~/.local/share/nvim/lazy/blink.cmp/target/release/libblink_cmp_fuzzy.dylib
```

### LSP log bloat

terraform-ls, taplo, and helm_ls flood stderr, which lands as `[ERROR]` in `lsp.log`. `vim.lsp.set_log_level("OFF")` in `options.lua` silences it; raise to `"WARN"`/`"ERROR"` when debugging LSP issues.
