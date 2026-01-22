#!/usr/bin/env python3
"""
Script untuk insert QC data CALCIUM ARSENAZO ke database
20 data untuk Level 1 dan 20 data untuk Level 2
Dengan increment waktu 1 hari per data
"""

import sqlite3
import random
import os
from datetime import datetime, timedelta

# Konfigurasi database - menggunakan path yang sama dengan aplikasi Go
LOCAL_APP_DATA = os.getenv("LOCALAPPDATA")
if LOCAL_APP_DATA:
    DB_PATH = os.path.join(LOCAL_APP_DATA, "TAMALabs", "database", "TAMALabs.db")
else:
    DB_PATH = "TAMALabs.db"  # Fallback jika LOCALAPPDATA tidak ada

print(f"Database path: {DB_PATH}")

# Konfigurasi QC Entry
# Anda perlu adjust ID ini sesuai dengan QC Entry yang sudah ada di database
# Atau uncomment bagian create QC entry di bawah
LEVEL1_QC_ENTRY_ID = 54  # Will be created or found
LEVEL2_QC_ENTRY_ID = 55  # Will be created or found

DEVICE_ID = 15  # Sesuaikan dengan device ID yang ada
TEST_TYPE_CODE = "CALCIUM ARSENAZO"

# Target values untuk CALCIUM ARSENAZO (contoh, sesuaikan dengan datasheet Anda)
LEVEL1_TARGET_MEAN = 9.70
LEVEL1_TARGET_SD = 0.39

LEVEL2_TARGET_MEAN = 13.90
LEVEL2_TARGET_SD = 0.60

# Start date untuk data (mundur 20 hari dari sekarang)
START_DATE = datetime.now() - timedelta(days=20)

def connect_db():
    """Connect to SQLite database"""
    conn = sqlite3.connect(DB_PATH)
    conn.row_factory = sqlite3.Row
    return conn

def get_or_create_test_type(conn, code):
    """Get test type ID by code"""
    cursor = conn.cursor()
    cursor.execute("SELECT id FROM test_types WHERE code = ?", (code,))
    row = cursor.fetchone()
    if row:
        return row[0]
    else:
        print(f"Test type '{code}' not found in database!")
        print("Please create the test type first or update TEST_TYPE_CODE")
        return None

def get_or_create_qc_entry(conn, device_id, test_type_id, qc_level, target_mean, target_sd):
    """Get or create QC entry"""
    cursor = conn.cursor()
    
    # Check if active QC entry exists
    cursor.execute("""
        SELECT id FROM qc_entries 
        WHERE device_id = ? AND test_type_id = ? AND qc_level = ? AND is_active = 1
    """, (device_id, test_type_id, qc_level))
    
    row = cursor.fetchone()
    if row:
        print(f"Found existing QC Entry Level {qc_level}: ID = {row[0]}")
        return row[0]
    
    # Create new QC entry
    lot_number = f"LOT-CA-{qc_level}-2026"
    ref_min = target_mean - (3 * target_sd)
    ref_max = target_mean + (3 * target_sd)
    
    cursor.execute("""
        INSERT INTO qc_entries (
            device_id, test_type_id, qc_level, lot_number,
            target_mean, target_sd, ref_min, ref_max,
            method, is_active, created_by, created_at, updated_at
        ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
    """, (
        device_id, test_type_id, qc_level, lot_number,
        target_mean, target_sd, ref_min, ref_max,
        'manual', 1, 'SYSTEM', datetime.now().isoformat(), datetime.now().isoformat()
    ))
    
    entry_id = cursor.lastrowid
    print(f"Created new QC Entry Level {qc_level}: ID = {entry_id}")
    return entry_id

def generate_qc_value(target_mean, target_sd, distribution='normal'):
    """
    Generate QC value with realistic distribution
    85% within ±2SD (In Control)
    10% within ±2SD to ±3SD (Warning)
    5% beyond ±3SD (Reject)
    """
    rand = random.random()
    
    if rand < 0.85:  # 85% In Control (±2SD)
        value = random.gauss(target_mean, target_sd * 0.6)
    elif rand < 0.95:  # 10% Warning (2SD to 3SD)
        if random.random() < 0.5:
            value = random.uniform(target_mean - 3*target_sd, target_mean - 2*target_sd)
        else:
            value = random.uniform(target_mean + 2*target_sd, target_mean + 3*target_sd)
    else:  # 5% Reject (beyond 3SD)
        if random.random() < 0.5:
            value = random.uniform(target_mean - 4*target_sd, target_mean - 3*target_sd)
        else:
            value = random.uniform(target_mean + 3*target_sd, target_mean + 4*target_sd)
    
    return round(value, 2)

def insert_qc_result(conn, qc_entry_id, measured_value, created_at, operator="SYSTEM"):
    """Insert raw QC result to database"""
    cursor = conn.cursor()
    
    cursor.execute("""
        INSERT INTO qc_results (
            qc_entry_id, measured_value, 
            calculated_mean, calculated_sd, calculated_cv,
            error_sd, absolute_error, relative_error,
            sd_1, sd_2, sd_3, result, method,
            operator, created_by, created_at
        ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
    """, (
        qc_entry_id, measured_value,
        0, 0, 0,  # Will be calculated by backend
        0, 0, 0,  # Will be calculated by backend
        0, 0, 0,  # Will be calculated by backend
        'Raw', 'raw',  # Raw data, not yet calculated
        operator, operator, created_at
    ))
    
    return cursor.lastrowid

def main():
    print("=" * 60)
    print("QC CALCIUM ARSENAZO Data Insertion Script")
    print("=" * 60)
    
    random.seed(42)  # For reproducible results
    
    # Connect to database
    conn = connect_db()
    
    try:
        # Get test type ID
        test_type_id = get_or_create_test_type(conn, TEST_TYPE_CODE)
        if not test_type_id:
            return
        
        print(f"Test Type '{TEST_TYPE_CODE}' ID: {test_type_id}")
        
        # Get or create QC entries
        level1_entry_id = get_or_create_qc_entry(
            conn, DEVICE_ID, test_type_id, 1, 
            LEVEL1_TARGET_MEAN, LEVEL1_TARGET_SD
        )
        
        level2_entry_id = get_or_create_qc_entry(
            conn, DEVICE_ID, test_type_id, 2,
            LEVEL2_TARGET_MEAN, LEVEL2_TARGET_SD
        )
        
        print("\n" + "=" * 60)
        print("Inserting Level 1 QC Data (20 entries)")
        print("=" * 60)
        
        # Insert Level 1 data (20 entries, 1 day apart)
        for i in range(20):
            created_at = START_DATE + timedelta(days=i)
            value = generate_qc_value(LEVEL1_TARGET_MEAN, LEVEL1_TARGET_SD)
            
            result_id = insert_qc_result(
                conn, level1_entry_id, value, 
                created_at.isoformat(), "SYSTEM"
            )
            
            print(f"Day {i+1:2d} ({created_at.strftime('%Y-%m-%d')}): "
                  f"Value = {value:6.2f} | ID = {result_id}")
        
        print("\n" + "=" * 60)
        print("Inserting Level 2 QC Data (20 entries)")
        print("=" * 60)
        
        # Insert Level 2 data (20 entries, 1 day apart)
        for i in range(20):
            created_at = START_DATE + timedelta(days=i)
            value = generate_qc_value(LEVEL2_TARGET_MEAN, LEVEL2_TARGET_SD)
            
            result_id = insert_qc_result(
                conn, level2_entry_id, value,
                created_at.isoformat(), "SYSTEM"
            )
            
            print(f"Day {i+1:2d} ({created_at.strftime('%Y-%m-%d')}): "
                  f"Value = {value:6.2f} | ID = {result_id}")
        
        # Commit transaction
        conn.commit()
        
        print("\n" + "=" * 60)
        print("SUCCESS! Data inserted successfully")
        print("=" * 60)
        print(f"Total Level 1 entries: 20")
        print(f"Total Level 2 entries: 20")
        print(f"Date range: {START_DATE.strftime('%Y-%m-%d')} to "
              f"{(START_DATE + timedelta(days=19)).strftime('%Y-%m-%d')}")
        print("\nNote: Data is saved as RAW. Backend will calculate statistics")
        print("when you select Method filter (Manual/Statistic) in the UI.")
        
    except Exception as e:
        conn.rollback()
        print(f"\nERROR: {e}")
        import traceback
        traceback.print_exc()
    
    finally:
        conn.close()

if __name__ == "__main__":
    main()
