# This block defines how Atlas loads the desired schema state.
# It uses the atlas-provider-gorm to inspect your Go structs.
data "external_schema" "gorm" {
  program = [
    "go",
    "run",
    "-mod=mod",
    "./migrations/loader",
  ]
}

# Defines the Atlas environment named "gorm"
env "gorm" {
  # Defines the "desired" state by referencing the external schema loader above.
  src = data.external_schema.gorm.url

  # Defines the "current" state database URL Atlas connects to for comparison.
  dev = "sqlite://tmp/temp_migration_purpose_only.db?cache=shared&_fk=1"

  # URL of the ACTUAL database to apply migrations TO.
  url = "sqlite://tmp/biosystem-lims.db?cache=shared&_fk=1"

  # Defines where migration files are stored (e.g., ./migrations directory).
  migration {
    dir = "file://migrations"
    format = "golang-migrate"
  }

  format {
    migrate {
      diff = "{{ sql . \"  \" }}"
    }
  }
}