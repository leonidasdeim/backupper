# backupper

This app accepts hot folder and backup folder paths, backing-up any file that is created or modified in the chosen
folder.

## Build with

App is written in Go. 

It uses one direct 3rd party library: [fsnotify](https://github.com/fsnotify/fsnotify). This library provides cross-platform filesystem notifications, it is widely used and maintained.

## How to use

Easiest way to build executable is to use provided Makefile job:
```
make build
```

Then run executable with arguments for hot and backup folders:
```
./out/backupper --in={hot_folder_path} --out={backup_folder_path}
```

Application will save state of folders paths and log view filter. On subsequent start, arguments are not required - previous state would be used.

### Run in demo mode

Demo job is provided in the Makefile, it will build and run executable in 'tmp' folder and provide paths to 'hot' and 'backup' folders inside it:
```
make build_run_demo
```

To demonstrate states persistence, there is another job in the Makefile, it is supposed to be run after 'build_run_demo' is stopped. It does not provide folder paths, so application load previously used paths:
```
make second_run_demo
```

### Other makefile jobs

Cleanup after building/running application:
```
make clean
```

Run tests:
```
make test
```

### Log view

On application start it will prompt for log viewer filter information.

First prompt will allow user to enable/disable log viewer.

Subsequent prompts will allow user to configure filter. 
On second application run it will ask to reuse previous state.

Filter is built for date and filename. If user leaves empty prompt input, it will disable appropriate filter.
