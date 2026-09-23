@echo off
rem Builds renamer.exe for Windows
cd /d "%~dp0"
go build -o renamer.exe .
echo Built renamer.exe
pause
