# vrchat_prefix

Adds `VRCHAT_ASSET_` to the front of every file in a folder:

    avatar.fbx  ->  VRCHAT_ASSET_avatar.fbx

Sub-folders, files that already have the prefix, and the script itself are skipped.
You see a preview and confirm before anything changes.

## Run

Needs Python 3 (https://www.python.org/downloads/).

- Copy `vrchat_prefix.py` into the folder and double-click it, **or**
- `python vrchat_prefix.py "C:\path\to\folder"`

## Compile to an .exe (optional)

    pip install pyinstaller
    pyinstaller --onefile vrchat_prefix.py

The program ends up in `dist/`. Drop it into a folder and double-click it.

To use a different prefix, change `PREFIX` at the top of `vrchat_prefix.py`.
