# VPP License Manager

This script automates the process of assigning VPP (Volume Purchase Program) licenses to devices in your MDM environment.

## Setup

1. Install the required dependencies:
```bash
pip install -r requirements.txt
```

2. Make the script executable:
```bash
chmod +x vpp_license_manager.py
```

## Usage

The script requires three parameters:
- `--token`: Your VPP Bearer token from Apple School/Business Manager
- `--apps`: Comma-separated list of App Store IDs
- `--serials`: Comma-separated list of device serial numbers

Example:
```bash
./vpp_license_manager.py \
  --token "your_vpp_bearer_token" \
  --apps "123456789,987654321" \
  --serials "DMQZ1234XXXX,DMQZ5678XXXX"
```

## Response

The script will output the JSON response from the VPP API. If there's an error, it will display both the error message and the API response for debugging.

## Security Note

Store your VPP Bearer token securely and never commit it to version control. Consider using environment variables or a secure secrets management system in production environments. 