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
