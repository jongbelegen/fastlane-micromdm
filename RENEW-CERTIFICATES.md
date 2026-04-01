# MicroMDM Certificate Renewal Guide

Certificates expire yearly. Set a calendar reminder 9-10 months after renewal.

## Details

- **Server:** https://mdm.fastlane.events (Fly.io, app: `micromdm-app`)
- **Organization:** Pitstop Payments B.V. / Fastlane
- **Email:** dimitri@fastlane.events
- **Country:** NL
- **mdmctl version:** v1.13.0 (`/usr/local/bin/mdmctl`)
- **mdmctl config:** `~/.micromdm/servers.json` (profile: `production`)

## Renewal History

| Year | Directory              | Password            | Date       |
|------|------------------------|---------------------|------------|
| 2025 | `mdm-certificates/`   | (unknown/lost)      | 2025-03-20 |
| 2026 | `mdm-certificates-2026/` | `killianisthebest` | 2026-02-20 |

## Step-by-step renewal

All commands run from the project root. Replace `YEAR` and `PASSWORD` as needed.

### Step 1: Create new vendor certificate request

```bash
cd mdm-certificates-YEAR/
mdmctl mdmcert vendor -password=PASSWORD -country=NL -email=dimitri@fastlane.events
```

This creates `mdm-certificates/VendorCertificateRequest.csr` and `mdm-certificates/VendorPrivateKey.key`.

**Note:** `mdmctl` always creates an `mdm-certificates/` subdirectory relative to where you run it.

### Step 2: Get vendor cert from Apple Developer Portal

1. Go to https://developer.apple.com/account/resources/certificates/list
2. Click (+) next to Certificates
3. Select **MDM CSR** under Services, click Continue
4. Upload `mdm-certificates/VendorCertificateRequest.csr`, click Continue
5. Download `mdm.cer` into the `mdm-certificates/` subfolder

### Step 3: Create push certificate request

Run from the same `mdm-certificates-YEAR/` directory (NOT from inside the `mdm-certificates/` subfolder, or you'll get a nested directory):

```bash
mdmctl mdmcert push -password=PASSWORD -country=NL -email=dimitri@fastlane.events
```

This creates `PushCertificateRequest.csr` and `PushCertificatePrivateKey.key` in the `mdm-certificates/` subfolder.

If files end up in a nested `mdm-certificates/mdm-certificates/`, move them up to `mdm-certificates/`.

### Step 4: Sign the push request with vendor cert

Run from `mdm-certificates-YEAR/`:

```bash
mdmctl mdmcert vendor -sign -cert=./mdm-certificates/mdm.cer -password=PASSWORD
```

This creates `mdm-certificates/PushCertificateRequest.plist`.

### Step 5: Upload to Apple Push Portal (CRITICAL)

1. Go to https://identity.apple.com
2. Find your EXISTING push certificate (the one expiring soonest)
3. Click **RENEW** — do NOT create a new one!
4. Upload `mdm-certificates/PushCertificateRequest.plist`
5. Download the new `.pem` certificate into `mdm-certificates/`

**Using Renew preserves the APNS topic so devices do NOT need to re-enroll.**

To identify which cert is yours: `openssl x509 -in old-cert.pem -noout -dates` — match the expiry date.

### Step 6: Upload to MicroMDM

Run from `mdm-certificates-YEAR/`:

```bash
mdmctl mdmcert upload \
    -cert "./mdm-certificates/MDM_ Pitstop Payments B.V._Certificate.pem" \
    -private-key ./mdm-certificates/PushCertificatePrivateKey.key \
    -password=PASSWORD
```

### Step 7: Restart MicroMDM

```bash
fly apps restart micromdm-app
```

### Step 8: Verify

```bash
mdmctl get devices
```

## DEP Token Renewal (separate, also yearly)

```bash
mdmctl get dep-tokens -export-public-key /tmp/DEPPublicKey.pem
```

1. Go to Apple Business Manager (business.apple.com)
2. Go to Settings > your MDM server (MicroMDM)
3. **Upload the public key FIRST** — select `/tmp/DEPPublicKey.pem`
4. **AFTER** the upload succeeds, click Download Token
5. Import the token:

```bash
mdmctl apply dep-tokens -import /path/to/MicroMDM_Token_*_smime.p7m
```

6. Verify:

```bash
mdmctl get dep-account
```

**IMPORTANT:** You MUST upload the public key before downloading the token. If you download first, ABM encrypts it with the old key and you'll get a `pkcs7: no enveloped recipient` error.

## VPP Token Renewal (separate, also yearly)

The VPP (Volume Purchase Program) content token is used for app license management. It expires yearly.

### Renewal History

| Year | Expiry Date | Renewed On |
|------|-------------|------------|
| 2026 | 2026-03-15  | 2025-03-15 |
| 2027 | 2027-03-25  | 2026-03-25 |

### Step 1: Download new content token from Apple Business Manager

1. Go to https://business.apple.com
2. Sign in with your organization account
3. Go to **Settings** (bottom-left) > **Payments and Billing** > **Content Tokens**
4. Find the token for your MDM server / location
5. Click **Download** to get a new `sToken_for_Fastlane_B.V..vpptoken` file

### Step 2: Replace the token file

Copy the downloaded file into the repo root, replacing the existing one:

```bash
cp ~/Downloads/sToken_for_Fastlane_B.V..vpptoken /path/to/micromdm/sToken_for_Fastlane_B.V..vpptoken
```

### Step 3: Decode and verify the new token

The `.vpptoken` file is base64-encoded JSON. Decode it to verify the new expiry date:

```bash
cat sToken_for_Fastlane_B.V..vpptoken | base64 -d
```

You should see JSON like:
```json
{
  "expDate": "2027-03-25T03:47:21+0000",
  "token": "MHXEuDgd+Y1o6...",
  "orgName": "Fastlane B.V."
}
```

### Step 4: Update all token files

Three files need to be updated with the new token data:

1. **`vpp_token.json`** — Update with the decoded JSON (expDate, token, orgName)
2. **`token.vpptoken`** — Replace contents with the full base64 string from the downloaded file
3. **`config/vpp.json`** — Update the `sToken` field with the new raw token value (the `token` field from the decoded JSON, NOT the base64)

### Step 5: Update VPP_TOKEN in cashless-backend

The `VPP_TOKEN` environment variable in the cashless-backend `.env` file needs the **base64-encoded** string (the full contents of `sToken_for_Fastlane_B.V..vpptoken`):

```bash
cat sToken_for_Fastlane_B.V..vpptoken
```

Copy the output and update `VPP_TOKEN` in the cashless-backend `.env`.

### Step 6: Deploy / restart services

Restart any services that use the VPP token (MicroMDM, cashless-backend, etc.).

### Quick reference: file mapping

| File | Contents | Format |
|------|----------|--------|
| `sToken_for_Fastlane_B.V..vpptoken` | Downloaded from Apple | Base64-encoded JSON |
| `token.vpptoken` | Copy of the above | Base64-encoded JSON |
| `vpp_token.json` | Decoded token data | Plain JSON |
| `config/vpp.json` | Token for MicroMDM config | JSON with `sToken` field |
| cashless-backend `.env` `VPP_TOKEN` | For backend API | Base64-encoded JSON |

## SCEP Device Certificate Expiration

When a device enrolls, it receives a SCEP identity certificate used to sign all MDM communication.
If this certificate expires, the server rejects all messages from the device with:

```
pkcs7: signing time "..." is outside of certificate validity "..." to "..."
```

**There is no way to renew a SCEP certificate remotely.** The device must re-enroll (factory reset + DEP auto-enrollment).

### Custom patch: allow expired signing certificates (2026-04-01)

We run a **forked version** of MicroMDM with a patch in `micromdm/pkg/crypto/helpers.go` (`PKCS7Verifier.Verify`). The patch allows devices with expired SCEP certificates to keep communicating instead of getting blocked with a 500 error. This means devices can still receive commands (e.g. `EraseDevice`, `InstallProfile`) even after their cert expires.

**Why:** ~200 devices were enrolled with the default 1-year SCEP cert validity. Mass factory reset was not feasible. The patch lets them keep working until they naturally re-enroll (replacement, reset, etc.), at which point they get the new 50-year cert.

**Trade-off:** Slightly weaker security — a compromised device key would be accepted even after cert expiry. Acceptable for an internal MDM fleet.

**If upgrading MicroMDM upstream:** Re-apply this patch to `pkg/crypto/helpers.go`, or remove it once all devices have re-enrolled with 50-year certs. The Dockerfile builds from the local `micromdm/` submodule (not from GitHub upstream).

### Configuration

The SCEP client certificate validity is controlled by `MICROMDM_SCEP_CLIENT_VALIDITY` (in days) in `fly.toml`.

| Setting | Default | Current |
|---------|---------|---------|
| `MICROMDM_SCEP_CLIENT_VALIDITY` | 365 (1 year) | 18250 (50 years) |

Changed on 2026-03-27 after device `00008020-0005252A21F0402E` expired with the 365-day default.

### Recovery steps when a device certificate expires

1. Factory reset the device (Settings > General > Transfer or Reset > Erase All Content and Settings)
2. The device will auto-re-enroll via DEP/ABM during Setup Assistant
3. The new SCEP certificate will use the configured validity (currently 50 years)

Alternative: wipe remotely via Apple Business Manager (business.apple.com) — independent of MDM.

**Note:** You cannot renew a device's SCEP certificate remotely. The device holds the private key — replacing the cert in the database won't help because it wouldn't match. Re-enrollment is the only way.

### Verifying SCEP certificate validity after re-enrollment

After re-enrolling a device, verify the new certificate has 50-year validity:

```bash
# 1. Download the database from Fly
fly sftp get /var/db/micromdm/micromdm.db ./micromdm.db -a micromdm-app

# 2. Run the check script
cd scripts/
go run check-scep-certs.go ../micromdm.db
```

New certificates should show `(50.0 yrs)`. Pass `--all` to see every certificate:

```bash
go run check-scep-certs.go ../micromdm.db --all
```

## Tips

- The password only needs to be consistent within one renewal cycle (all steps use the same one)
- Back up the entire `mdm-certificates-YEAR/` directory to 1Password after renewal
- After upload to MicroMDM, the certs are stored server-side — the local copy is a backup