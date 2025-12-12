"""
Script to delete all QC results via API
"""
import requests

# Configuration
BASE_URL = "http://localhost:8080"
API_URL = f"{BASE_URL}/api/v1/quality-control/results"

def delete_all_qc_results():
    """Delete all QC results to start fresh"""
    try:
        # Get all results first
        response = requests.get(API_URL)
        if response.status_code == 200:
            data = response.json()
            results = data.get('data', {}).get('results', [])
            
            print(f"Found {len(results)} QC results")
            
            # Delete each result
            for result in results:
                result_id = result['id']
                delete_url = f"{API_URL}/{result_id}"
                del_response = requests.delete(delete_url)
                if del_response.status_code == 200:
                    print(f"✓ Deleted result ID {result_id}")
                else:
                    print(f"✗ Failed to delete result ID {result_id}: {del_response.status_code}")
            
            print(f"\nDeleted all QC results. Database is now clean.")
        else:
            print(f"Failed to get QC results: {response.status_code}")
    
    except Exception as e:
        print(f"Error: {e}")

if __name__ == "__main__":
    confirm = input("Delete ALL QC results? This cannot be undone! (yes/no): ")
    if confirm.lower() == 'yes':
        delete_all_qc_results()
    else:
        print("Cancelled.")
