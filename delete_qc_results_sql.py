"""
Delete QC results directly via SQL
"""
import sqlite3
import os

db_path = "TAMALabs.db"

if not os.path.exists(db_path):
    print(f"Database not found: {db_path}")
    exit(1)

conn = sqlite3.connect(db_path)
cursor = conn.cursor()

# Count current results
cursor.execute("SELECT COUNT(*) FROM qc_results")
count = cursor.fetchone()[0]
print(f"Found {count} QC results")

if count > 0:
    confirm = input(f"Delete ALL {count} QC results? (yes/no): ")
    if confirm.lower() == 'yes':
        cursor.execute("DELETE FROM qc_results")
        conn.commit()
        print(f"✓ Deleted all {count} QC results")
        print("\nDatabase is now clean. You can test with fresh data.")
    else:
        print("Cancelled.")
else:
    print("No QC results to delete.")

conn.close()
