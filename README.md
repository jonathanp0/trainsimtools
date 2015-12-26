# Jonathan's Train Simulation Tools

## Zusi 2 Browser & Journey Log
This tool, written in GO, is designed as a lightweight, local alternative to [Zusi Browser](http://zusi-info.steftones.de/).
Additionally, you can add the listed train services to a personal journey log, which is stored in JSON format.

### Download
Like all GO software it is compatible with Windows, MacOS and Linux, but a binary build is only available for Windows.
[Download latest version for Windows(released 24.12.15)](http://files.gu2.co.uk/zusi/zusibrowser.zip)

### Instructions
To run, simply double click the .exe file. It will scan your data files, then announce it is ready. Then go to http://localhost:8888/ to use the interface.
If you do not have Zusi installed in C:\\Program Files (x86)\\Zusi, you need to run it from the command line as 
> zusibrowser -zusi pathtozusi

## Trainsim-helper MIDI
A fork of the excellent trainsim-helper project, which consist of LUA scripts which import/export control values from DTG Train Simulator, which unbelievably has no public API usable with modern addons; plus an onscreen overlay, a joystick to Train Simulator connector and tools to patch LUA scripts.

This version has additional code to read simulator from MIDI control channels instead of joystick buttons and axes.

In separate repository [here](https://github.com/jonathanp0/trainsim-helper/tree/midi).
For user manual see [readme.txt](https://github.com/jonathanp0/trainsim-helper/blob/midi/main/readme.txt)


