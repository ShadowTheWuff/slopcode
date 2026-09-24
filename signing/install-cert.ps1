# Trusts Tanner Knapp's code-signing certificate on this PC, so programs
# signed with it (like vrchat_prefix.exe) show as signed and trusted.
#
# Run in PowerShell:
#   powershell -ExecutionPolicy Bypass -File install-cert.ps1
#
# It installs for your Windows user only, so admin isn't needed. Windows will
# show a security warning asking you to confirm. Check that the thumbprint
# matches the one below before clicking Yes.

$ExpectedThumbprint = "E3BEFA37B081A0F7C798629E972C014BE10377A3"
$certPath = Join-Path $PSScriptRoot "TannerKnapp.cer"

$cert = New-Object System.Security.Cryptography.X509Certificates.X509Certificate2 $certPath
if ($cert.Thumbprint -ne $ExpectedThumbprint) {
    Write-Error "Certificate thumbprint doesn't match. Expected $ExpectedThumbprint, got $($cert.Thumbprint). Not installing."
    exit 1
}

Write-Host "Installing certificate:"
Write-Host "  Subject:    $($cert.Subject)"
Write-Host "  Thumbprint: $($cert.Thumbprint)"
Write-Host "  Expires:    $($cert.NotAfter)"

Import-Certificate -FilePath $certPath -CertStoreLocation Cert:\CurrentUser\Root | Out-Null
Import-Certificate -FilePath $certPath -CertStoreLocation Cert:\CurrentUser\TrustedPublisher | Out-Null

Write-Host "Done. Programs signed by Tanner Knapp are now trusted for your user."
