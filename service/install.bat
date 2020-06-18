@echo off
rem run this script as admin

if not exist service.exe (
    echo Build the service program before installing by running the build.bat file.
    goto :exit
)

sc create screenshot-service binpath= "%CD%\service.exe" start= auto DisplayName= "Screenshot Service"
sc description screenshot-service "Takes a screenshot of the desktop at interval of time specified in the Registry"
sc start screenshot-service
sc query screenshot-service

echo Check the Windoes event viewer

:exit
