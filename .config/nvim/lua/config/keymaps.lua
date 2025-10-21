-- Keymaps are automatically loaded on the VeryLazy event
-- Default keymaps that are always set: https://github.com/LazyVim/LazyVim/blob/main/lua/lazyvim/config/keymaps.lua
-- Add any additional keymaps here

-- Neo-tree: focus or open if closed
vim.keymap.set("n", "<leader>e", "<cmd>Neotree focus<cr>", { desc = "Focus Neo-tree" })

-- Close buffer with Ctrl+w (same as <leader>bd)
vim.keymap.set("n", "<C-w>", function()
  Snacks.bufdelete()
end, { desc = "Delete Buffer" })
