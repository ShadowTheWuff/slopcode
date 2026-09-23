# renamer

Renames every file inside one folder to the same name, numbered so they don't collide.
Extensions are kept, sub-folders are ignored, and you get a preview + confirmation first.

    photo.jpg, notes.txt, song.mp3  ->  vacation_1.jpg, vacation_2.txt, vacation_3.mp3

## Compile

Install Go (https://go.dev/dl/), then in this folder:

    go build            # makes ./renamer (or renamer.exe on Windows)

Or `./build.sh` to build Windows/macOS/Linux binaries into `dist/`,
or double-click `build.bat` on Windows.

## Use

    renamer <folder> <new-name>
    renamer "C:\Users\me\Pictures\trip" vacation

Or just double-click / run `renamer` with no arguments and it will ask you.
