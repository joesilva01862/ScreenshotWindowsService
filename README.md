## Screenshot Windows Service

This is a Windows service written in GO, to continually wake up at predefined intervals and invoke another program.
This other program, also written in GO, takes a screenshot of your desktop and sends it to
a WebDav server, such as Nextcloud.

## How to use it
Here you'll find the sources for both the Service and the screenshot program, along with install/uninstall scripts and 
instructions on how to build both.

## Screenshots
<p align="center">
  <img align="center" width="600" height="250" src="https://github.com/joesilva01862/ScreenshotWindowsService/blob/master/WindowsService.jpg"/>
  <br>
  Fig 1: The Service panel shows Screenshot Service
</p>
  

<p>

<p align="center">
  <img align="center" width="600" height="250" src="https://github.com/joesilva01862/ScreenshotWindowsService/blob/master/Screenshot.jpg"/>
  <br>
  Fig 2: Sample screenshot taken by the program that the service invokes
</p>

<p>
  
<p align="center">
  <img align="center" width="600" height="250" src="https://github.com/joesilva01862/ScreenshotWindowsService/blob/master/RegistryEntries.jpg"/>
  <br>
  Fig 3: Registry variables required to run the service
</p>
