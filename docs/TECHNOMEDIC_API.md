# TechnoMedic API Integration

API bridging untuk integrasi dengan TechnoMedic telah berhasil dibuat. Berikut adalah dokumentasi endpoint yang tersedia:

## ⚙️ Configuration Requirements

Sebelum menggunakan API TechnoMedic, pastikan konfigurasi sudah diaktifkan:

### Aktivasi TechnoMedic Integration

1. **Login ke aplikasi** sebagai administrator
2. **Navigasi ke halaman Config** (Settings)
3. **Cari setting** `TechnoMedicIntegrationEnabled`
4. **Set value ke** `true`
5. **Save configuration**

**Atau via Database:**
```sql
UPDATE configs SET value = 'true' WHERE id = 'TechnoMedicIntegrationEnabled';
```

### Status Integration

- ✅ **Enabled (true)**: API TechnoMedic dapat diakses
- ❌ **Disabled (false)**: API TechnoMedic akan return `403 Forbidden`

**Error Response saat Disabled:**
```json
{
  "message": "TechnoMedic integration is not enabled"
}
```

### Cara Mengecek Status

```bash
# Test API endpoint
curl http://localhost:8080/api/v1/technomedic/test-types

# Jika disabled, akan return:
# HTTP 403 Forbidden
# {"message":"TechnoMedic integration is not enabled"}

# Jika enabled, akan return:
# HTTP 200 OK
# {"code":200,"status":"success","data":[...]}
```

## Base URL
```
http://localhost:{PORT}/api/v1/technomedic
```

## Endpoints

### 1. Get Test Types
Mendapatkan daftar semua jenis tes yang tersedia.

**Endpoint:** `GET /api/v1/technomedic/test-types`

**Response:**
```json
{
  "code": 200,
  "status": "success",
  "message": "Test types retrieved successfully",
  "data": [
    {
      "id": "1",
      "code": "HB",
      "name": "Hemoglobin",
      "category": "Hematologi",
      "sub_category": "Complete Blood Count",
      "specimen_type": "Whole Blood",
      "unit": "g/dL"
    }
  ]
}
```

### 2. Get Sub-Categories
Mendapatkan daftar semua sub-kategori tes dari tabel `sub_categories`.

**Endpoint:** `GET /api/v1/technomedic/sub-categories`

**Response:**
```json
{
  "code": 200,
  "status": "success",
  "message": "Sub-categories retrieved successfully",
  "data": [
    {
      "id": "1",
      "code": "CBC",
      "name": "Complete Blood Count",
      "category": "Hematologi",
      "description": "Pemeriksaan darah lengkap"
    }
  ]
}
```

### 3. Get Test Types by Sub-Category
Mendapatkan daftar test types berdasarkan sub-category ID.

**Endpoint:** `GET /api/v1/technomedic/sub-categories/:id/test-types`

**Path Parameters:**
- `id` (required): Sub-category ID

**Example:** `GET /api/v1/technomedic/sub-categories/1/test-types`

**Response:**
```json
{
  "code": 200,
  "status": "success",
  "message": "Test types retrieved successfully",
  "data": [
    {
      "id": "1",
      "code": "HB",
      "name": "Hemoglobin",
      "category": "Hematologi",
      "sub_category": "Complete Blood Count",
      "specimen_type": "Whole Blood",
      "unit": "g/dL"
    },
    {
      "id": "2",
      "code": "WBC",
      "name": "White Blood Cell",
      "category": "Hematologi",
      "sub_category": "Complete Blood Count",
      "specimen_type": "Whole Blood",
      "unit": "10^3/uL"
    }
  ]
}
```

### 4. Get Doctors
Mendapatkan daftar semua dokter.

**Endpoint:** `GET /api/v1/technomedic/doctors`

**Response:**
```json
{
  "code": 200,
  "status": "success",
  "message": "Doctors retrieved successfully",
  "data": [
    {
      "id": 1,
      "fullname": "Dr. John Doe",
      "username": "johndoe",
      "is_active": true
    }
  ]
}
```

### 5. Get Analysts
Mendapatkan daftar semua analis.

**Endpoint:** `GET /api/v1/technomedic/analysts`

**Response:**
```json
{
  "code": 200,
  "status": "success",
  "message": "Analysts retrieved successfully",
  "data": [
    {
      "id": 2,
      "fullname": "Jane Smith",
      "username": "janesmith",
      "is_active": true
    }
  ]
}
```

### 6. Create Order
Membuat order baru dari TechnoMedic.

**Endpoint:** `POST /api/v1/technomedic/order`

**Request Body:**
```json
{
  "no_order": "TM-2024-001",
  "patient": {
    "patient_id": "P001",
    "full_name": "John Doe",
    "sex": "M",
    "address": "Jl. Contoh No. 123",
    "birthdate": "1990-01-15",
    "medical_record_number": "MR001",
    "phone_number": "081234567890"
  },
  "test_type_ids": [1, 2, 3],
  "requested_by": "Dr. Smith",
  "requested_at": "2024-01-20 10:30:00"
}
```

**Opsi untuk Menentukan Test Types:**

Anda dapat menggunakan salah satu atau kombinasi dari opsi berikut:

1. **`test_type_ids`** (RECOMMENDED) - Array of test type IDs
   ```json
   "test_type_ids": [1, 2, 3, 4]
   ```

2. **`sub_category_ids`** (RECOMMENDED) - Array of sub-category IDs
   ```json
   "sub_category_ids": [1, 2]
   ```
   *Note: Akan mengambil semua test types dalam sub-category tersebut*

3. **`param_request`** - Array of test type codes
   ```json
   "param_request": ["HB", "WBC", "PLT"]
   ```

4. **`sub_category_request`** - Array of sub-category names
   ```json
   "sub_category_request": ["Complete Blood Count", "Liver Function"]
   ```

**Contoh Request dengan Sub-Category IDs:**
```json
{
  "no_order": "TM-2024-002",
  "patient": {
    "full_name": "Jane Doe",
    "sex": "F",
    "birthdate": "1985-05-20",
    "medical_record_number": "MR002"
  },
  "sub_category_ids": [1, 3]
}
```

**Contoh Request dengan Test Type IDs:**
```json
{
  "no_order": "TM-2024-003",
  "patient": {
    "full_name": "Bob Smith",
    "sex": "M",
    "birthdate": "1992-12-10",
    "medical_record_number": "MR003"
  },
  "test_type_ids": [1, 2, 5, 7]
}
```

**Contoh Request dengan Kombinasi:**
```json
{
  "no_order": "TM-2024-004",
  "patient": {
    "full_name": "Alice Johnson",
    "sex": "F",
    "birthdate": "1988-08-15",
    "medical_record_number": "MR004"
  },
  "sub_category_ids": [1],
  "test_type_ids": [10, 11]
}
```

**Notes:**
- **Minimal satu** dari: `test_type_ids`, `sub_category_ids`, `param_request`, atau `sub_category_request` harus diisi
- Jika menggunakan kombinasi, semua test types akan digabung (de-duplicated)
- `sex`: Must be "M" or "F"
- `birthdate`: Format YYYY-MM-DD
- `requested_at`: Format YYYY-MM-DD HH:MM:SS or YYYY-MM-DDTHH:MM:SSZ (optional, defaults to current time)

**Response:**
```json
{
  "code": 201,
  "status": "success",
  "message": "Order created successfully",
  "data": {
    "no_order": "TM-2024-001"
  }
}
```

### 7. Get Order
Mendapatkan detail order termasuk hasil pemeriksaan.

**Endpoint:** `GET /api/v1/technomedic/order/{no_order}`

**Response Structure:**

Response akan memiliki 2 bagian:
1. **`sub_categories`** - Test types yang memiliki sub-category (dikelompokkan)
2. **`parameters_result`** - Test types yang TIDAK memiliki sub-category (langsung di root level)

**Response (Test dengan Sub-Category):**
```json
{
  "code": 200,
  "status": "success",
  "message": "Order retrieved successfully",
  "data": {
    "no_order": "TM-2024-001",
    "status": "SUCCESS",
    "patient": {
      "patient_id": "MR001",
      "full_name": "John Doe",
      "sex": "M",
      "address": "Jl. Contoh No. 123",
      "birthdate": "1990-01-15",
      "medical_record_number": "MR001",
      "phone_number": "081234567890"
    },
    "requested_by": "",
    "requested_at": "2024-01-20 10:30:00",
    "sub_categories": [
      {
        "id": "1",
        "name": "Complete Blood Count",
        "parameters_result": [
          {
            "id": "1",
            "code": "HB",
            "category_name": "Hemoglobin",
            "value": "14.5",
            "specimen_type": "Whole Blood",
            "unit": "g/dL",
            "ref": "12.0-16.0",
            "flag": ""
          }
        ]
      }
    ],
    "completed_at": "2024-01-20 11:00:00",
    "verified_at": "2024-01-20 11:15:00",
    "verified_by": "Dr. Smith"
  }
}
```

**Response (Test TANPA Sub-Category):**
```json
{
  "code": 200,
  "status": "success",
  "message": "Order retrieved successfully",
  "data": {
    "no_order": "TM-2024-002",
    "status": "SUCCESS",
    "patient": {...},
    "requested_by": "",
    "requested_at": "2024-01-20 10:30:00",
    "parameters_result": [
      {
        "id": "10",
        "code": "COVID19",
        "category_name": "COVID-19 Antigen Test",
        "value": "Negative",
        "specimen_type": "Nasopharyngeal Swab",
        "unit": "",
        "ref": "Negative",
        "flag": ""
      },
      {
        "id": "11",
        "code": "PREG",
        "category_name": "Pregnancy Test",
        "value": "Positive",
        "specimen_type": "Urine",
        "unit": "",
        "ref": "Negative",
        "flag": "H"
      }
    ]
  }
}
```

**Response (Kombinasi - Ada yang Dengan dan Tanpa Sub-Category):**
```json
{
  "code": 200,
  "status": "success",
  "data": {
    "no_order": "TM-2024-003",
    "status": "SUCCESS",
    "patient": {...},
    "sub_categories": [
      {
        "id": "1",
        "code": "CBC",
        "name": "Complete Blood Count",
        "parameters_result": [
          {
            "id": "1",
            "code": "HB",
            "category_name": "Hemoglobin",
            "value": "14.5",
            ...
          }
        ]
      }
    ],
    "parameters_result": [
      {
        "id": "10",
        "code": "COVID19",
        "category_name": "COVID-19 Antigen Test",
        "value": "Negative",
        ...
      }
    ]
  }
}
```

**Notes:**
- Jika order hanya berisi test tanpa sub-category, field `sub_categories` akan kosong/omitted
- Jika order hanya berisi test dengan sub-category, field `parameters_result` akan kosong/omitted
- Jika order berisi kombinasi, kedua field akan ada

## Error Responses

Semua endpoint menggunakan format error response yang sama:

```json
{
  "code": 400,
  "status": "error",
  "message": "Error message here"
}
```

**HTTP Status Codes:**
- `200 OK`: Request berhasil
- `201 Created`: Order berhasil dibuat
- `400 Bad Request`: Request tidak valid
- `404 Not Found`: Order tidak ditemukan
- `500 Internal Server Error`: Server error

## Database Changes

### Tabel Baru: `sub_categories`
Tabel baru untuk menyimpan sub-kategori tes:
- `id`: Primary key
- `name`: Nama sub-kategori (unique)
- `code`: Kode sub-kategori
- `category`: Kategori utama
- `description`: Deskripsi sub-kategori
- `created_at`: Timestamp pembuatan
- `updated_at`: Timestamp update

### Perubahan Tabel `test_types`
- `sub_category_id`: Foreign key ke tabel `sub_categories` (nullable)
- `sub_category`: Tetap ada untuk backward compatibility

### Perubahan Tabel `work_orders`
- `barcode_simrs`: Digunakan untuk menyimpan nomor order dari TechnoMedic (menggunakan field existing, bukan field baru)
- `source`: Sumber order (e.g., "technomedic")
- `verified_at`: Timestamp verifikasi
- `verified_by`: Nama verifikator
- `completed_at`: Timestamp penyelesaian

**Note**: Kami menggunakan field `barcode_simrs` yang sudah ada untuk menyimpan nomor order dari external system seperti TechnoMedic, sehingga tidak perlu menambah kolom baru ke database.
## Implementation Files

1. **Entity**: 
   - `internal/entity/external/technomedic.go`
   - `internal/entity/sub_category.go`
2. **Repository**: 
   - `internal/repository/sql/sub_category/sub_category.go`
3. **Usecase**: 
   - `internal/usecase/external/technomedic/technomedic.go`
   - `internal/usecase/sub_category/sub_category.go`
4. **Handler**: `internal/delivery/rest/technomedic.go`
5. **Migration**: 
   - `migrations/20260202111405_add_technomedic_fields_to_work_orders.up.sql`
   - `migrations/20260202112047_create_sub_categories_table.up.sql`

## Testing

Untuk testing API, gunakan tools seperti Postman atau curl:

```bash
# Get all test types
curl http://localhost:8080/api/v1/technomedic/test-types

# Get all sub-categories
curl http://localhost:8080/api/v1/technomedic/sub-categories

# Get test types by sub-category ID
curl http://localhost:8080/api/v1/technomedic/sub-categories/1/test-types

# Get all doctors
curl http://localhost:8080/api/v1/technomedic/doctors

# Get all analysts
curl http://localhost:8080/api/v1/technomedic/analysts

# Create order with test_type_ids (RECOMMENDED)
curl -X POST http://localhost:8080/api/v1/technomedic/order \
  -H "Content-Type: application/json" \
  -d '{
    "no_order": "TM-2024-001",
    "patient": {
      "full_name": "John Doe",
      "sex": "M",
      "birthdate": "1990-01-15",
      "medical_record_number": "MR001"
    },
    "test_type_ids": [1, 2, 3]
  }'

# Create order with sub_category_ids (RECOMMENDED)
curl -X POST http://localhost:8080/api/v1/technomedic/order \
  -H "Content-Type: application/json" \
  -d '{
    "no_order": "TM-2024-002",
    "patient": {
      "full_name": "Jane Smith",
      "sex": "F",
      "birthdate": "1985-05-20",
      "medical_record_number": "MR002"
    },
    "sub_category_ids": [1]
  }'

# Create order with test codes (Alternative)
curl -X POST http://localhost:8080/api/v1/technomedic/order \
  -H "Content-Type: application/json" \
  -d '{
    "no_order": "TM-2024-003",
    "patient": {
      "full_name": "Bob Johnson",
      "sex": "M",
      "birthdate": "1992-03-10",
      "medical_record_number": "MR003"
    },
    "param_request": ["HB", "WBC"]
  }'

# Get order details
curl http://localhost:8080/api/v1/technomedic/order/TM-2024-001
```

## Alur Penggunaan

**Opsi 1: Menggunakan IDs (RECOMMENDED)**
```
1. GET /sub-categories 
   → Dapatkan list sub-kategori dengan ID

2. GET /sub-categories/{id}/test-types
   → Lihat test types dalam sub-kategori (dapat ID test types)

3. POST /order dengan test_type_ids ATAU sub_category_ids
   → Kirim array of IDs langsung

4. GET /order/{no_order}
   → Cek hasil pemeriksaan
```

**Opsi 2: Menggunakan Codes/Names (Alternative)**
```
1. GET /test-types
   → Dapatkan test codes

2. POST /order dengan param_request
   → Kirim array of test codes

3. GET /order/{no_order}
   → Cek hasil pemeriksaan
```

## Catatan Penting

- **4 Cara kirim test request**: `test_type_ids` (ID langsung), `sub_category_ids` (ID sub-kategori), `param_request` (kode test), atau `sub_category_request` (nama sub-kategori)
- **Rekomendasi**: Gunakan `test_type_ids` atau `sub_category_ids` untuk performa lebih baik
- Tabel `sub_categories` sekarang menyimpan master data sub-kategori
- Setiap test type dapat dikaitkan dengan sub-category melalui `sub_category_id`
- Field `sub_category` di tabel `test_types` masih ada untuk backward compatibility
- Migration akan otomatis memigrate data existing ke tabel `sub_categories`
