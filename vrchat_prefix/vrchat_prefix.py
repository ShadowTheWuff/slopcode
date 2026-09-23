"""Add the VRCHAT_ASSET_ prefix to every file in a folder.

    photo.png  ->  VRCHAT_ASSET_photo.png

Usage:
    python vrchat_prefix.py            # uses the folder this script is in
    python vrchat_prefix.py <folder>   # uses the given folder

Files that already start with the prefix, sub-folders and this script
itself are left alone. You get a preview and must confirm first.
"""

import os
import sys

PREFIX = "VRCHAT_ASSET_"


def main():
    if len(sys.argv) > 1:
        folder = sys.argv[1]
    else:
        # Folder of the script (or of the .exe when built with PyInstaller),
        # so double-clicking it works on the folder it's sitting in.
        me = sys.executable if getattr(sys, "frozen", False) else __file__
        folder = os.path.dirname(os.path.abspath(me))

    if not os.path.isdir(folder):
        print(f"Error: {folder!r} is not a folder.")
        return 1

    self_path = os.path.realpath(sys.executable if getattr(sys, "frozen", False) else __file__)

    changes = []
    for name in sorted(os.listdir(folder)):
        path = os.path.join(folder, name)
        if not os.path.isfile(path) or os.path.islink(path):
            continue
        # skip this program (and copies of it) and already-prefixed files
        if name == os.path.basename(self_path) or name.startswith(PREFIX):
            continue
        changes.append((name, PREFIX + name))

    if not changes:
        print(f"Nothing to rename in {folder}")
        return 0

    print(f"{len(changes)} file(s) in {folder} will be renamed:")
    for old, new in changes:
        print(f"  {old}  ->  {new}")
    if input("\nGo ahead? [y/N]: ").strip().lower() not in ("y", "yes"):
        print("Cancelled, nothing was changed.")
        return 0

    done = 0
    for old, new in changes:
        dst = os.path.join(folder, new)
        if os.path.exists(dst):
            print(f"  skipped {old}: {new} already exists")
            continue
        try:
            os.rename(os.path.join(folder, old), dst)
            done += 1
        except OSError as e:
            print(f"  failed {old}: {e}")

    print(f"Done! Renamed {done} file(s).")
    return 0


if __name__ == "__main__":
    code = main()
    if len(sys.argv) <= 1:
        input("Press Enter to exit...")  # keep the window open on double-click
    sys.exit(code)
