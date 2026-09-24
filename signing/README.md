# Code-signing certificate

`TannerKnapp.cer` is the **public** certificate used to sign the `.exe` builds in
this repo. It contains no private key, so it's safe to share.

| | |
|---|---|
| Subject | CN=Tanner Knapp, O=Tanner Knapp, L=Lester Prairie, S=Minisota, C=US |
| Thumbprint (SHA-1) | `E3BEFA37B081A0F7C798629E972C014BE10377A3` |
| Valid | 2026-09-23 to 2031-09-23 |
| Usage | Code Signing |

It's self-signed, so Windows won't trust it until you install it.

## Trust it (Windows)

In this folder, run:

    powershell -ExecutionPolicy Bypass -File install-cert.ps1

Confirm the Windows security prompt after checking that the thumbprint matches the one above.

## Check a signed program

    Get-AuthenticodeSignature path\to\vrchat_prefix.exe

`Status` should be `Valid` and the signer thumbprint should match the one above.

## Remove it

Open `certmgr.msc` and delete "Tanner Knapp" from **Trusted Root Certification
Authorities** and **Trusted Publishers**.

## Automatic builds (GitHub Actions)

`.github/workflows/build-and-sign.yml` builds `vrchat_prefix.exe` and `renamer.exe`
on Windows, signs them, checks the signatures, and uploads them as the
**apps-signed** artifact on the workflow run. Pull-request builds are left unsigned.

### Releases

The signed builds are published as a GitHub Release (under **Releases** on the repo page) when:

- **You push to the default branch.** The release is named automatically: `v1.0.<run number>`.
- **You run it by hand.** Go to **Actions → Build and sign → Run workflow**, tick
  "Publish a GitHub Release", and optionally type a version like `v1.2.0`.
- **You push a tag.** For example, `git tag v2.0.0 && git push origin v2.0.0`.

Each release contains `vrchat_prefix.exe`, `renamer.exe`, `TannerKnapp.cer` and
`SHA256SUMS.txt`. Release builds skip the build cache so they start from scratch.

### One-time setup: add the signing key as secrets

1. On your PC, copy the `.pfx` to the clipboard as base64:

       $docs = [Environment]::GetFolderPath("MyDocuments")
       [Convert]::ToBase64String([IO.File]::ReadAllBytes("$docs\mycert.pfx")) | Set-Clipboard

2. On GitHub: repo **Settings → Secrets and variables → Actions → New repository secret**
   - `SIGNING_CERT_PFX_BASE64`: paste the clipboard
   - `SIGNING_CERT_PASSWORD`: the `.pfx` password

3. Clear your clipboard afterwards (copy something else).

How the workflow protects the key:
- The key only exists in the runner's temp folder while signing and is deleted right after.
- GitHub hides secret values in logs, and never gives them to pull requests from forks.
- The workflow refuses to sign unless the key's thumbprint matches `TannerKnapp.cer`.
- Third-party actions are pinned to exact commits, so a changed tag can't swap in new code.
