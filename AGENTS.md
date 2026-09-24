# AGENTS.md

Guidance for AI coding agents working in this repo.

## What's here

Two small standalone Windows-friendly file-renaming tools, plus code-signing material.

| Path | What it is |
|---|---|
| `renamer/` | Go CLI: renames every file in a folder to `<name>_<n>.<ext>` (numbered, extensions kept, sub-folders skipped). Single file, `main.go`, standard library only. |
| `vrchat_prefix/` | Python 3 script: prefixes every file in a folder with `VRCHAT_ASSET_`. Single file, `vrchat_prefix.py`, standard library only. Built into an .exe with PyInstaller. |
| `signing/` | Public code-signing certificate (`TannerKnapp.cer`), an install script for it, and signing docs. |

## Build and check

```sh
# renamer
cd renamer && go vet ./... && go build -o /tmp/renamer .
./build.sh                      # cross-compiles into renamer/dist/

# vrchat_prefix
python3 -m py_compile vrchat_prefix/vrchat_prefix.py
pip install -r vrchat_prefix/requirements-build.txt   # pinned PyInstaller
pyinstaller --onefile vrchat_prefix/vrchat_prefix.py
```

There is no test suite. To check a change, run the tool against a scratch folder
(never the repo itself). Both tools show a preview and ask `y/N`, so pipe the answer in:

```sh
mkdir -p /tmp/t && touch /tmp/t/a.png /tmp/t/b.txt
(cd renamer && echo y | go run . /tmp/t test)
echo y | python3 vrchat_prefix/vrchat_prefix.py /tmp/t
```

## Behaviour to preserve

Both tools:
- Rename **files only**. They skip folders and symlinks, and never recurse.
- **Never rename themselves** if they sit in the target folder. `renamer` skips its own
  running executable; `vrchat_prefix` skips any file with its own name, so copies are skipped too.
- Show a preview and **require confirmation** before changing anything.
- When run with no arguments (a double-click on Windows), prompt for input and
  wait for Enter before exiting, so the window stays open.

`renamer` specifically:
- Renames in two passes via temporary names and rolls back on error, so it never overwrites a file.
- Sorts files by name and zero-pads numbers (`01`, `02`, … `10`).
- Rejects names containing `/ \ : * ? " < > |`.

`vrchat_prefix` specifically:
- Keeps the prefix in the `PREFIX` constant at the top of the file.
- Skips files that already start with the prefix, so running it twice is harmless.
- Skips (doesn't overwrite) a file whose new name already exists.

## Conventions

- Keep each tool a single file with no third-party runtime dependencies.
- Match the surrounding style: short comments that explain *why*, plain names.
- Go: `gofmt`. Python: standard library, no type-checker config.
- Users are non-developers on Windows, so keep messages plain, avoid jargon,
  and make sure paths work with spaces and backslashes.

## Security: don't break these

- **Never commit private keys.** `.pfx`/`.p12` are git-ignored; keep it that way.
  `signing/TannerKnapp.cer` is the *public* certificate and is meant to be committed.
- Signing certificate thumbprint: `E3BEFA37B081A0F7C798629E972C014BE10377A3`.
  It appears in `signing/install-cert.ps1`, `signing/README.md` and the CI config;
  if the certificate is replaced, update all three.
- Don't log, echo or write secrets (signing key, its password, GitHub tokens) anywhere.
