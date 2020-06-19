## StatsAndroidApp

This is a Windows service written in GO, to continually wake up at predefined interval and invoke another program.
This other program is a program, also written in GO, to take a screenshot of your desktop and send it to
a WebDav server, such as Nextcloud.

## How to use it
Here, you'll find the souces for both the Service and the screenshot program, with install script and 
instructions on how to build both.

## Screenshots
<p align="center">
  <img align="center" src="https://github.com/joesilva01862/ScreenshotWindowsService/tree/master/service/pictures/WindowsService.jpg">
  <br>
  Fig 1: The Service panel shows Screenshot Service
</p>
  
<p align="center">
  <img align="center" src="https://github.com/joesilva01862/ScreenshotWindowsService/tree/master/service/pictures/RegistryEntries.png">
  <br>
  Fig 2: Registry variables required to run the service
</p>

<p align="center">
  <img align="center" src="https://github.com/joesilva01862/ScreenshotWindowsService/tree/master/service/pictures/Screenshot.jpg">
  <br>
  Fig 3: Sample screenshot taken by the program that the service invokes
</p>