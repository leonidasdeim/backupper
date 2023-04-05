# backupper

This app accept hot folder path and backup folder path, backing-up any file that is created or modified in the chosen
folder.

## Build with

App is written in Go. 

It uses one direct 3rd party library: [fsnotify](https://github.com/fsnotify/fsnotify). This library provides cross-platform filesystem notifications, it is widely used and maintained.

## How to use

Easiest way to build executable is to use provided Makefile:
```
make build
```

Run executable with required arguments for hot folder and backup folder:
```
./out/backupper --in={hot_folder_path} --out={backend_folder_path}
```

## Demo mode

Demo job is provided in Makefile, it will build and run executable in temporary folder:
```
make build_run_demo
```

Cleanup after running application in demo:
```
make clean
```

## Log view

On application start it will prompt for log viewer filter information.

First prompt will allow user to enable/disable log viewer.

Subsequent prompts will allow user to configure filter. 
On second application run it will ask to reuse previous filter.

Filter are built for date and filename. If user leaves empty prompt input for any of filters, it will disable appropriate filter.