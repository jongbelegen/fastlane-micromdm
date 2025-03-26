#!/usr/bin/env python3

import argparse
import json
import requests
from typing import List, Dict

class VPPLicenseManager:
    def __init__(self, bearer_token: str):
        self.bearer_token = bearer_token
        self.base_url = "https://vpp.itunes.apple.com/mdm/v2"
        self.headers = {
            "Authorization": f"Bearer {bearer_token}",
            "Content-Type": "application/json"
        }

    def assign_licenses(self, app_ids: List[str], serial_numbers: List[str]) -> Dict:
        """
        Assign VPP licenses for given apps to specified device serial numbers
        
        Args:
            app_ids: List of App Store IDs
            serial_numbers: List of device serial numbers
        
        Returns:
            API response as dictionary
        """
        endpoint = f"{self.base_url}/assets/associate"
        
        assets = [{"adamId": app_id, "pricingParam": "STDQ"} for app_id in app_ids]
        
        payload = {
            "assets": assets,
            "serialNumbers": serial_numbers
        }
        
        response = requests.post(
            endpoint,
            headers=self.headers,
            json=payload
        )
        
        response.raise_for_status()
        return response.json()

def main():
    parser = argparse.ArgumentParser(description='Manage VPP license assignments')
    parser.add_argument('--token', required=True, help='VPP Bearer token')
    parser.add_argument('--apps', required=True, help='Comma-separated list of App Store IDs')
    parser.add_argument('--serials', required=True, help='Comma-separated list of device serial numbers')
    
    args = parser.parse_args()
    
    # Split comma-separated inputs into lists
    app_ids = [app.strip() for app in args.apps.split(',')]
    serial_numbers = [serial.strip() for serial in args.serials.split(',')]
    
    try:
        manager = VPPLicenseManager(args.token)
        result = manager.assign_licenses(app_ids, serial_numbers)
        print(json.dumps(result, indent=2))
    except requests.exceptions.RequestException as e:
        print(f"Error assigning licenses: {e}")
        if hasattr(e, 'response') and e.response is not None:
            print(f"Response: {e.response.text}")

if __name__ == "__main__":
    main() 