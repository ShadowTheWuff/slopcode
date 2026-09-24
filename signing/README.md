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

## Automatic builds (CircleCI)

`.circleci/config.yml` runs three jobs:

1. **build** (Windows): builds `vrchat_prefix.exe` and `renamer.exe`. Runs for every branch and `v*` tag.
2. **sign** (Windows): signs both, checks the signatures, and saves them as build
   artifacts (`apps-signed`). Runs only for `main`, `claude/*` branches and `v*` tags.
3. **deploy** (Linux): publishes a GitHub Release, tracked as a CircleCI deploy
   (environment `github-releases`, component `vrchat-tools`). Runs for pushes to
   the default branch and `v*` tags.

Release names:
- Pushes to the default branch: `v1.0.<pipeline number>`
- A pushed tag such as `v2.0.0`: the tag name
- A manual run (**Trigger Pipeline** in CircleCI): add a string parameter `version`, e.g. `v1.2.0`

Each release contains `vrchat_prefix.exe`, `renamer.exe`, `TannerKnapp.cer` and
`SHA256SUMS.txt`. Tag builds skip the build cache so they start from scratch.

### One-time setup

1. **Connect the repo.** Sign in at circleci.com with GitHub and set up the
   `slopcode` project using the existing `.circleci/config.yml`.

2. **Create the `code-signing` context** (Organization Settings → Contexts) with:
   - `SIGNING_CERT_PFX_BASE64`: your `.pfx` as base64. On your PC:

         $docs = [Environment]::GetFolderPath("MyDocuments")
         [Convert]::ToBase64String([IO.File]::ReadAllBytes("$docs\mycert.pfx")) | Set-Clipboard

   - `SIGNING_CERT_PASSWORD`: the `.pfx` password

3. **Create the `github-release` context** with `GH_TOKEN`: a GitHub
   [fine-grained token](https://github.com/settings/personal-access-tokens/new)
   limited to this repo, with **Contents: Read and write**.

4. **Create the deploy environment.** In CircleCI go to **Deploys → Environments**,
   create an environment integration named `github-releases`. The component
   `vrchat-tools` shows up after the first deploy.

5. **Remove the old GitHub Actions secrets** (`SIGNING_CERT_PFX_BASE64`,
   `SIGNING_CERT_PASSWORD`) from GitHub. The workflow that used them is gone.

How the key is protected:
- Only the **sign** job gets the `code-signing` context, and only for this repo's
  own branches and tags. CircleCI never gives contexts to builds from forks.
- The key only exists in the VM's temp folder while signing and is deleted right after.
- The job refuses to sign unless the key's thumbprint matches `TannerKnapp.cer`.
- Optionally, restrict both contexts to yourself with a security group or
  project restriction in the context settings.
