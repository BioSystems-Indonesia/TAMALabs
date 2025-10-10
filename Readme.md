# LIMS HL Seven Client

## Overview

This project is a Laboratory Information Management System (LIMS) designed to handle messages from laboratories machine. It features a web-based frontend for user interaction and a Go backend to process and manage data.

---

## Tech Stack

### Backend

- **Language:** [Go](https://golang.org/)
- **Web Framework:** [Echo](https://echo.labstack.com/)
- **Database ORM:** [GORM](https://gorm.io/)
- **Database Migrations:** [Atlas](https://atlasgo.io/)
- **Linting:** [golangci-lint](https://golangci-lint.run/)
- **Live Reload:** [Air](https://github.com/cosmtrek/air)

### Frontend

- **Framework:** [React](https://reactjs.org/) with [Vite](https://vitejs.dev/)
- **UI Toolkit:** [React Admin](https://marmelab.com/react-admin/)

---

## Prerequisites

Before you begin, ensure you have the following installed:

- [Go](https://golang.org/doc/install) (version 1.21 or higher)
- [Node.js](https://nodejs.org/en/download/) (version 18 or higher)
- [Make](https://www.gnu.org/software/make/)
- [Atlas](https://atlasgo.io/cli/getting-started/)
- [Docker](https://www.docker.com/products/docker-desktop)

---

## Getting Started

Follow these steps to get your development environment running:

1.  **Clone the repository:**

    ```sh
    git clone https://github.com/BioSystems-Indonesia/TAMALabs
    cd lims-hl-seven
    ```

2.  **Install backend dependencies:**

    ```sh
    go mod tidy
    ```

3.  **Install frontend dependencies:**

    ```sh
    cd web
    npm install
    cd ..
    ```

4.  **Run the development servers:**
    - **Backend:**
      ```sh
      make dev-be
      ```
    - **Frontend:**
      ```sh
      make dev-fe
      ```

The backend will be running on `http://localhost:8322` and the frontend on `http://localhost:5132`.

---

## Database Migrations

The application uses SQLite as the database. The database file is located at `./tmp/biosystem-lims.db`.

This project uses [Atlas](https://atlasgo.io/) to manage database schemas, with migration files located in the `migrations` directory. The desired schema state is derived from the GORM entity files in `internal/entity`.

- **Create a new migration file:**
  **DO NOT MANUALLY CREATE MIGRATION FILE. YOU MUST MODIFY the internal/entity and then run the migrate-diff command**

  1.  Modify the GORM structs in the `internal/entity` directory to reflect your desired schema changes.
  2.  Run the following command to generate a new migration file. This command compares your GORM entities with the current state of the database and generates the necessary SQL.
      ```sh
      make migrate-diff desc="your_migration_name"
      ```
      _Replace `your_migration_name` with a descriptive name for your migration (e.g., `add_user_email_field`)._
  3.  If you get a checksum error, run the following command to rehash the migration files.
      ```sh
      make migrate-hash
      ```

- **Apply all pending migrations:**
  Simply run the app to apply all pending migrations. It will run all migration on startup.

---

## Linting

To ensure code quality, run the linter for the backend:

```sh
make lint
```

---

## Building for Production

To create a production-ready build that embeds the frontend into the Go binary:

```sh
make build
```

This command will:

1.  Build the frontend application into the `web/dist` directory.
2.  Embed the static assets into the Go binary.
3.  Compile the final executable to `bin/winapp.exe`.

---

## Create Windows Installer

To create a Windows installer for the application:

1.  Run the following command:

    ```sh
    make installer
    ```

    This will use Docker to build the installer, so ensure you have Docker
    installed and running on your system.

2.  The installer will be created in the `installer` directory.
