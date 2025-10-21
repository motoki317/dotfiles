return require("schema-companion").setup_client(
  require("schema-companion").adapters.yamlls.setup({
    sources = {
      require("schema-companion").sources.matchers.kubernetes.setup({
        version = "master",
      }),
      require("schema-companion").sources.lsp.setup(),
    },
  }),
  {
    -- Your normal yamlls configuration goes here
    filetypes = { "yaml" }, -- Exclude helm filetype
    settings = {
      yaml = {
        validate = true,
        format = { enable = true },
        hover = true,
        completion = true,
      },
      schemas = {
        ["https://json.schemastore.org/github-workflow"] = ".github/workflows/*",
        ["https://json.schemastore.org/github-action"] = "*action.{yml,yaml}",
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
  }
)
