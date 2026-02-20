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

## Tips

- The password only needs to be consistent within one renewal cycle (all steps use the same one)
- Back up the entire `mdm-certificates-YEAR/` directory to 1Password after renewal
- After upload to MicroMDM, the certs are stored server-side — the local copy is a backup