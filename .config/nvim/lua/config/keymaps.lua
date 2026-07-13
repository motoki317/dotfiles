-- Keymaps are automatically loaded on the VeryLazy event
-- Default keymaps that are always set: https://github.com/LazyVim/LazyVim/blob/main/lua/lazyvim/config/keymaps.lua
-- Add any additional keymaps here

-- Neo-tree: focus or open if closed
vim.keymap.set("n", "<leader>e", "<cmd>Neotree focus<cr>", { desc = "Focus Neo-tree" })

-- Hunk: review the working-tree diff in a floating terminal, lazygit-style.
-- `hunk diff` shows current changes (incl. untracked); run at the repo root so it
-- works from any buffer. Mirrors LazyVim's <leader>gg lazygit map. Bound to <leader>gv
-- ("git view"), a free key that leaves gitsigns' <leader>gh "hunks" group untouched.
vim.keymap.set("n", "<leader>gv", function()
  Snacks.terminal({ "hunk", "diff" }, { cwd = LazyVim.root.git() })
end, { desc = "Hunk (Root Dir)" })

-- Focus main buffer pane (move to right window)
vim.keymap.set("n", "<leader>.", "<cmd>wincmd l<cr>", { desc = "Focus main buffer" })

-- Close buffer with Ctrl+w (same as <leader>bd)
vim.keymap.set("n", "<C-w>", function()
  Snacks.bufdelete()
end, { desc = "Delete Buffer" })

-- Close all file buffers with Ctrl+Shift+w
vim.keymap.set("n", "<C-S-w>", function()
  local bufs = vim.api.nvim_list_bufs()
  for _, buf in ipairs(bufs) do
    if vim.api.nvim_buf_is_valid(buf) and vim.api.nvim_buf_get_option(buf, "buflisted") then
      local buftype = vim.api.nvim_buf_get_option(buf, "buftype")
      -- Only delete normal file buffers (not terminals, help, quickfix, etc.)
      if buftype == "" then
        vim.api.nvim_buf_delete(buf, { force = false })
      end
    end
  end
end, { desc = "Close All File Buffers" })
