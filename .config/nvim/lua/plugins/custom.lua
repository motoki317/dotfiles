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
      -- Always show neotree on startup
      vim.api.nvim_create_autocmd("VimEnter", {
        callback = function()
          vim.cmd("Neotree show")
        end,
      })
    end,
  },

  { "j-hui/fidget.nvim" },

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

  -- Helm files (ft helm)
  {
    "qvalentin/helm-ls.nvim",
    ft = "helm",
    opts = {
      conceal_templates = {
        -- Not working quite well
        enabled = false,
      },
    },
  },

  -- CRD files (ft yaml, does not work for helm files)
  {
    "diogo464/kubernetes.nvim",
    opts = {
      -- this can help with autocomplete. it sets the `additionalProperties` field on type definitions to false if it is not already present.
      schema_strict = true,
      -- true:  generate the schema every time the plugin starts
      -- false: only generate the schema if the files don't already exists. run `:KubernetesGenerateSchema` manually to generate the schema if needed.
      schema_generate_always = true,
      -- Patch yaml-language-server's validation.js file.
      patch = true,
      -- root path of the yamlls language server. by default it is assumed you are using mason but if not this option allows changing that path.
      yamlls_root = function()
        return vim.fs.joinpath(vim.fn.stdpath("data"), "/mason/packages/yaml-language-server/")
      end,
    },
  },
  {
    "neovim/nvim-lspconfig",
    dependencies = {
      "diogo464/kubernetes.nvim",
    },
    opts = function(_, opts)
      opts.servers = opts.servers or {}
      opts.servers.yamlls = {
        settings = {
          yaml = {
            schemas = {
              [require("kubernetes").yamlls_schema()] = "*.yaml",
              ["https://json.schemastore.org/github-workflow"] = ".github/workflows/*",
              ["https://json.schemastore.org/github-action"] = ".github/action.{yml,yaml}",
              ["https://json.schemastore.org/ansible-stable-2.9"] = "roles/tasks/*.{yml,yaml}",
              ["https://json.schemastore.org/prettierrc"] = ".prettierrc.{yml,yaml}",
              ["https://json.schemastore.org/kustomization"] = "kustomization.{yml,yaml}",
              ["https://json.schemastore.org/ansible-playbook"] = "*play*.{yml,yaml}",
              ["https://json.schemastore.org/chart"] = "Chart.{yml,yaml}",
              ["https://json.schemastore.org/dependabot-v2"] = ".github/dependabot.{yml,yaml}",
              ["https://json.schemastore.org/gitlab-ci"] = "*gitlab-ci*.{yml,yaml}",
              ["https://raw.githubusercontent.com/OAI/OpenAPI-Specification/main/schemas/v3.1/schema.json"] = "*api*.{yml,yaml}",
              ["https://raw.githubusercontent.com/compose-spec/compose-spec/master/schema/compose-spec.json"] = "*docker-compose*.{yml,yaml}",
              ["https://raw.githubusercontent.com/argoproj/argo-workflows/master/api/jsonschema/schema.json"] = "*flow*.{yml,yaml}",
            },
          },
        },
      }
      return opts
    end,
  },
}
