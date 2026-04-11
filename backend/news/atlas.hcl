data "external_schema" "bun" {
  program = [
    "go",
    "run",
    "-mod=mod",
    "ariga.io/atlas-provider-bun",
    "load",
    "--path", "./internal/models",
    "--dialect", "postgres",
  ]
}

env "bun" {
  src = data.external_schema.bun.url
  dev = "docker://postgres/15/dev?search_path=public"
  migration {
    dir = "file://./internal/database/migrations/"
  }
}
