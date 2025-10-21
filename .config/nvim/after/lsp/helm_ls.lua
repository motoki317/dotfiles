return require("schema-companion").setup_client(
  require("schema-companion").adapters.helmls.setup({
    sources = {
      require("schema-companion").sources.matchers.kubernetes.setup({
        version = "master",
      }),
      require("schema-companion").sources.lsp.setup(),
    },
  }),
  {
    -- Your normal helm-ls configuration
    settings = {
      ["helm-ls"] = {
        logLevel = "info",
        valuesFiles = {
          mainValuesFile = "values.yaml",
          lintOverlayValuesFile = "values.lint.yaml",
          additionalValuesFilesGlobPattern = "values*.yaml",
        },
        yamlls = {
          enabled = true,
          diagnosticsLimit = 50,
          showDiagnosticsDirectly = false,
          path = "yaml-language-server",
        },
      },
    },
  }
)
