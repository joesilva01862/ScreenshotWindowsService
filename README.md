## Screenshot Windows Service

This is a Windows service written in GO, to continually wake up at predefined intervals and invoke another program.
This other program is a program, also written in GO, to take a screenshot of your desktop and send it to
a WebDav server, such as Nextcloud.

## How to use it
Here, you'll find the souces for both the Service and the screenshot program, with install script and 
instructions on how to build both.

## Screenshots
<p align="center">
  <img align="center" width="600" height="250" src="https://github.com/joesilva01862/ScreenshotWindowsService/blob/master/WindowsService.jpg">
  <br>
  Fig 1: The Service panel shows Screenshot Service
</p>
  
  
<p align="center">
  <img align="center" width="600" height="250" src="https://github.com/joesilva01862/ScreenshotWindowsService/blob/master/RegistryEntries.jpg">
  <br>
  Fig 2: Registry variables required to run the service
</p>


<p align="center">
  <img align="center" width="600" height="250" src="https://github.com/joesilva01862/ScreenshotWindowsService/blob/master/Screenshot.jpg">
  <br>
  Fig 3: Sample screenshot taken by the program that the service invokes
</p>
