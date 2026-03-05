-- Options are automatically loaded before lazy.nvim startup
-- Default options that are always set: https://github.com/LazyVim/LazyVim/blob/main/lua/lazyvim/config/options.lua
-- Add any additional options here

-- Disable relative line numbers by default
vim.opt.relativenumber = false

-- Enable wrap by default
vim.opt.wrap = true

-- Suppress LSP log to prevent unbounded growth (terraform-ls, taplo, helm_ls flood stderr)
vim.lsp.set_log_level("OFF")
