return {
  {
    "nvim-neo-tree/neo-tree.nvim",
    branch = "v3.x",
    dependencies = {
      "nvim-lua/plenary.nvim",
      "MunifTanjim/nui.nvim",
      "nvim-tree/nvim-web-devicons", -- optional, but recommended
    },
    lazy = false, -- neo-tree will lazily load itself
    opts = {
      filesystem = {
        filtered_items = {
          visible = true,
          hide_dotfiles = false,
          hide_gitignored = true,
        },
        window = {
          mappings = {
            ["O"] = "expand_all_subnodes",
          },
        },
      },
    },
    config = function(_, opts)
      require("neo-tree").setup(opts)
      -- Show Neo-tree on startup, but not when launched as another tool's $EDITOR.
      -- Detected from argv shape rather than an env flag, since some launchers
      -- (e.g. Claude Code's Ctrl+G) tokenize $EDITOR naively and lose embedded args.
      local function launched_as_editor()
        if vim.fn.argc() ~= 1 then
          return false
        end
        local arg = vim.fn.fnamemodify(vim.fn.argv(0), ":p")
        if
          arg:match("COMMIT_EDITMSG$")
          or arg:match("MERGE_MSG$")
          or arg:match("TAG_EDITMSG$")
          or arg:match("git%-rebase%-todo$")
          or arg:match("addp%-hunk%-edit%.diff$")
        then
          return true
        end
        -- Strip macOS /private symlink prefix so /private/var/folders/... → /var/folders/...
        local normalized = arg:gsub("^/private/", "/")
        local tmpdir = vim.env.TMPDIR
        if tmpdir and tmpdir ~= "" then
          local stripped = tmpdir:gsub("^/private/", "/")
          if normalized:sub(1, #stripped) == stripped then
            return true
          end
        end
        return normalized:match("^/tmp/") ~= nil or normalized:match("^/var/folders/") ~= nil
      end
      vim.api.nvim_create_autocmd("VimEnter", {
        callback = function()
          if not launched_as_editor() then
            vim.cmd("Neotree show")
          end
        end,
      })
    end,
  },

  -- Lualine config
  {
    "nvim-lualine/lualine.nvim",
    opts = {
      sections = {
        lualine_c = {
          {
            "filename",
            path = 1,
            shorting_target = 10,
          },
        },
      },
    },
  },

  -- Nice LSP progress message
  { "j-hui/fidget.nvim" },

  -- Shellchecks, terraform support
  {
    "mfussenegger/nvim-lint",
    opts = {
      linters_by_ft = {
        sh = { "shellcheck" },
        bash = { "shellcheck" },
        zsh = { "shellcheck" },
        terraform = { "terraform_validate" },
        tf = { "terraform_validate" },
      },
    },
  },
  {
    "stevearc/conform.nvim",
    optional = true,
    opts = {
      formatters_by_ft = {
        hcl = { "packer_fmt" },
        terraform = { "terraform_fmt" },
        tf = { "terraform_fmt" },
        ["terraform-vars"] = { "terraform_fmt" },
      },
    },
  },

  -- k8s yaml / helm
  {
    "cenk1cenk2/schema-companion.nvim",
    dependencies = {
      { "nvim-lua/plenary.nvim" },
    },
    config = function()
      require("schema-companion").setup({
        log_level = vim.log.levels.INFO,
      })
    end,
    keys = {
      {
        "<leader>ys",
        function()
          require("schema-companion").select_schema()
        end,
        desc = "Select YAML Schema",
      },
      {
        "<leader>ym",
        function()
          require("schema-companion").select_from_matching_schema()
        end,
        desc = "Select from Matching Schemas",
      },
      {
        "<leader>yr",
        function()
          require("schema-companion").match()
        end,
        desc = "Re-match Schema",
      },
    },
  },

  -- helm-ls fix
  {
    "qvalentin/helm-ls.nvim",
    ft = "helm",
    opts = {
      conceal_templates = {
        -- Not working quite well
        enabled = false,
      },
      action_highlight = {
        -- Disable due to invalid tree-sitter query
        enabled = false,
      },
    },
  },

  -- Code action preview
  {
    "rachartier/tiny-code-action.nvim",
    dependencies = {
      { "nvim-lua/plenary.nvim" },
      {
        "folke/snacks.nvim",
        opts = {
          terminal = {},
        },
      },
    },
    event = "LspAttach",
    opts = {
      picker = "snacks",
    },
  },
  {
    "neovim/nvim-lspconfig",
    opts = {
      servers = {
        ["*"] = {
          keys = {
            { "<leader>ca", false },
            {
              "<leader>ca",
              function()
                require("tiny-code-action").code_action({})
              end,
              desc = "Code Action (tiny-code-action)",
              mode = { "n", "x" },
            },
          },
        },
      },
    },
  },
}
