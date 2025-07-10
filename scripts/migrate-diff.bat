@echo off
rem Script to run atlas migrate diff with a required argument

rem Check if an argument (the migration name) was provided.
if "%~1"=="" (
    rem Print error message to standard error
    echo Error: Missing migration name. ^>&2
    echo Usage: %~n0 ^<migration_name^> ^>&2
    exit /b 1
)

rem Get the migration name from the first argument
set "MIGRATION_NAME=%~1"

rem Echo the command that will be run (optional, but helpful for debugging)
echo Running: atlas migrate diff --env gorm %MIGRATION_NAME%

rem Execute the actual command
atlas migrate diff --env gorm %MIGRATION_NAME%
if %errorlevel% neq 0 (
    echo Error: Atlas command failed with exit code %errorlevel% ^>&2
    exit /b %errorlevel%
)

atlas migrate lint --env gorm
if %errorlevel% neq 0 (
    echo Error: Lint migration failed with exit code %errorlevel% ^>&2
    echo Please do the following: ^>&2
    echo 1. REMOVE the new migration file and revert the atlas.sum ^>&2
    echo 2. Modify the GORM entity to match the correct schema ^>&2
    echo 3. Run this script again ^>&2
    exit /b %errorlevel%
)

echo Atlas migrate diff completed.
exit /b 0
