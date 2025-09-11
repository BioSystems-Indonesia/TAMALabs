@echo off
echo Generating 1000 lab requests...
echo.

cd /d "%~dp0.."
go run scripts/generate_lab_requests.go

if %ERRORLEVEL% EQU 0 (
    echo.
    echo ✅ Successfully generated 1000 lab requests!
    echo.
    echo You can now:
    echo - Check the database at ./tmp/biosystem-lims.db
    echo - View the data via the web interface
    echo - Use the REST API to access the generated data
) else (
    echo.
    echo ❌ Failed to generate lab requests. Error code: %ERRORLEVEL%
    echo Please check the error messages above.
)

echo.
pause
