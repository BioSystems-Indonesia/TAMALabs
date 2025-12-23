#!/usr/bin/env python3
import socket
import time
import uuid
import random
from datetime import datetime

random.seed(42)
level1_values = []
for i in range(20):
    rand = random.random()
    if rand < 0.85:
        value = random.gauss(9.70, 0.39 * 0.6)
    elif rand < 0.95:
        if random.random() < 0.5:
            value = random.uniform(8.92, 9.31)
        else:
            value = random.uniform(10.09, 10.48)
    else:
        if random.random() < 0.5:
            value = random.uniform(8.50, 8.92)
        else:
            value = random.uniform(10.48, 10.90)
    level1_values.append(round(value, 2))

level2_values = []
for i in range(20):
    rand = random.random()
    if rand < 0.85:
        value = random.gauss(13.90, 0.60 * 0.6)
    elif rand < 0.95:
        if random.random() < 0.5:
            value = random.uniform(12.70, 13.10)
        else:
            value = random.uniform(14.70, 15.10)
    else:
        if random.random() < 0.5:
            value = random.uniform(12.10, 12.70)
        else:
            value = random.uniform(15.10, 15.70)
    level2_values.append(round(value, 2))

def create_hl7_message(level, lot_number, value, timestamp):
    message_id = str(uuid.uuid4())
    qc_sample_id = f"QC HUMAN {level} {lot_number}"
    
    hl7_msg = (
        f"MSH|^~\\&|BA400|Biosystems|Host|Host provider|{timestamp}||OUL^R22^OUL_R22|{message_id}|P|2.5.1|||ER|AL||UNICODE UTF-8|||LAB-29^IHE\r"
        f"SPM|1|{qc_sample_id}||NULL|||||||Q|||||||\r"
        f"INV|{qc_sample_id}|OK|CO|||||||||20260203022259||||{lot_number}XA\r"
        f"OBR||\"\"||CHOLESTEROL^CHOLESTEROL^A400|||||||||||||||||||||||||\r"
        f"ORC|OK||||CM||||{timestamp}\r"
        f"OBX|1|NM|CHOLESTEROL^CHOLESTEROL^A400||{value}|mg/dL^mg/dL^A400|8.5 - 10.8199997|NONE|||F|||||ADMIN||A400^Biosystems~834000240^Biosystems|{timestamp}\r"
    )
    return hl7_msg

def send_hl7_message(sock, hl7_message):
    start_block = b'\x0b'
    end_block = b'\x1c'
    carriage_return = b'\x0d'
    
    mllp_message = start_block + hl7_message.encode('utf-8') + end_block + carriage_return
    
    sock.sendall(mllp_message)
    
    response = sock.recv(4096)
    return response

TCP_HOST = 'localhost'
TCP_PORT = 5335

try:
    print("=" * 60)
    print(f"Sending QC Level I (Lot: 003980) - {len(level1_values)} measurements")
    print("Target: 9.70 ± 0.39 SD")
    print("=" * 60)
    for i, value in enumerate(level1_values, 1):
        sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        sock.connect((TCP_HOST, TCP_PORT))
        
        timestamp = datetime.now().strftime("%Y%m%d%H%M%S")
        hl7_msg = create_hl7_message("I", "003980", value, timestamp)
        
        response = send_hl7_message(sock, hl7_msg)
        
        error = (value - 9.70) / 0.39
        status = "Normal" if abs(error) <= 2 else ("Warning" if abs(error) <= 3 else "Error")
        print(f"Message {i:2d}/{len(level1_values)} - Value: {value:6.2f} - Error: {error:+5.2f} SD - {status}")
        
        sock.close()
        time.sleep(1)
    
    print("\n" + "=" * 60)
    print(f"Sending QC Level II (Lot: 004057) - {len(level2_values)} measurements")
    print("Target: 13.90 ± 0.60 SD")
    print("=" * 60)
    
    for i, value in enumerate(level2_values, 1):
        sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        sock.connect((TCP_HOST, TCP_PORT))
        
        timestamp = datetime.now().strftime("%Y%m%d%H%M%S")
        hl7_msg = create_hl7_message("II", "004057", value, timestamp)
        
        response = send_hl7_message(sock, hl7_msg)
        
        error = (value - 13.90) / 0.60
        status = "Normal" if abs(error) <= 2 else ("Warning" if abs(error) <= 3 else "Error")
        print(f"Message {i:2d}/{len(level2_values)} - Value: {value:6.2f} - Error: {error:+5.2f} SD - {status}")
        
        sock.close()
        time.sleep(1)
    
    print("\n" + "=" * 60)
    print("All messages sent successfully!")
    print("=" * 60)
    
except Exception as e:
    print(f"\nError: {e}")
