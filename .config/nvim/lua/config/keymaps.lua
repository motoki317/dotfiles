-- Keymaps are automatically loaded on the VeryLazy event
-- Default keymaps that are always set: https://github.com/LazyVim/LazyVim/blob/main/lua/lazyvim/config/keymaps.lua
-- Add any additional keymaps here

-- Neo-tree: focus or open if closed
vim.keymap.set("n", "<leader>e", "<cmd>Neotree focus<cr>", { desc = "Focus Neo-tree" })

-- Focus main buffer pane (move to right window)
vim.keymap.set("n", "<leader>.", "<cmd>wincmd l<cr>", { desc = "Focus main buffer" })

-- Close buffer with Ctrl+w (same as <leader>bd)
vim.keymap.set("n", "<C-w>", function()
  Snacks.bufdelete()
end, { desc = "Delete Buffer" })
