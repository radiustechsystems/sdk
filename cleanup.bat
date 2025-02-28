@echo off

echo Running cleanup for all Radius SDKs...

echo.
echo Running Go SDK cleanup...
cd go
call cleanup.bat
cd ..

echo.
echo Running TypeScript SDK cleanup...
cd typescript
call cleanup.bat
cd ..

echo.
echo All SDK checks completed successfully!
