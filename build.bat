@echo off

if not exist ".\build\daemon" mkdir ".\build\daemon"
if not exist ".\build\client" mkdir ".\build\client"

go build -o .\build\ttrackerclient.exe .\timetracker\src\client\client.go
go build -o .\build\ttrackerd.exe .\timetracker\src\daemon\daemon.go

pause