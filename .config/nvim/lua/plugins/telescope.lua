return {
  {
    "nvim-telescope/telescope.nvim",
    opts = {
      pickers = {
        live_grep = {
          additional_args = { "--hidden", "--glob", "!**/.git/*" },
        },
        grep_string = {
          additional_args = { "--hidden", "--glob", "!**/.git/*" },
        },
      },
    },
  },
}
