# Lab Request Generator

Script ini digunakan untuk generate 1000 lab request dummy data untuk testing dan development.

## Struktur Data yang Digenerate

### 1. **Patients** (Pasien)

- Nama depan dan belakang Indonesia
- Tanggal lahir random (umur 1-90 tahun)
- Jenis kelamin (M/F)
- Nomor telepon Indonesia
- Alamat dan lokasi

### 2. **Work Orders** (Lab Request)

- Status: NEW
- Barcode unik
- Tanggal create random (30 hari terakhir)
- Terhubung dengan patient dan admin creator

### 3. **Specimens** (Sampel)

- 1-3 specimen per work order
- Tipe specimen: SER, URI, CSF, LIQ, PLM
- Barcode unik per specimen
- Tanggal collection dan received
- Source, condition, method, dan comments

### 4. **Observation Requests** (Test Request)

- 1-5 test per specimen
- Test types dari data seeder yang sudah ada
- Status: PENDING
- Tanggal request random
- Tidak ada duplicate test code dalam satu specimen

## Cara Penggunaan

### Option 1: Menggunakan Batch File (Windows)

```cmd
cd scripts
run_generator.bat
```

### Option 2: Manual dengan Go

```cmd
# Pastikan berada di root directory project
go run scripts/generate_lab_requests.go
```

## Prerequisites

1. **Database Sudah Disetup**: Pastikan database SQLite sudah ada di `./tmp/biosystem-lims.db`
2. **Seed Data**: Pastikan test types dan admin sudah ada (jalankan seeder terlebih dahulu)
3. **Go Environment**: Go version 1.24+ terinstall

## Output

Script akan menggenerate:

- ±1000 patients baru
- ±1000 work orders baru
- ±1500-3000 specimens baru (1-3 per work order)
- ±3000-15000 observation requests baru (1-5 per specimen)

## Progress Monitoring

Script akan menampilkan progress setiap 100 record:

```
Progress created=100 total=1000
Progress created=200 total=1000
...
✅ Successfully generated 1000 lab requests!
```

## Test Data Characteristics

### Nama Pasien

Menggunakan nama-nama Indonesia yang umum seperti:

- First names: Ahmad, Budi, Siti, Dewi, Andi, Rini, dll
- Last names: Pratama, Wijaya, Santoso, Utomo, dll

### Test Types

Menggunakan test types yang sudah ada di seeder:

- Hematology: RBC, HGB, HCT, WBC, PLT, dll
- Biochemistry: Glucose, Creatinine, Cholesterol, dll
- Dan test lainnya sesuai dengan data seeder

### Distribusi Data

- 70% specimen dalam kondisi "Good"
- Random distribution untuk specimen types
- Random test combinations per specimen
- Realistic date ranges (last 30 days for work orders, last 7 days for specimen collection)

## Troubleshooting

### Error: "no test types available"

Jalankan seeder terlebih dahulu:

```cmd
go run cmd/rest/main.go # untuk seed default data
```

### Error: "failed to connect database"

Pastikan aplikasi utama sudah pernah dijalankan untuk membuat database:

```cmd
go run cmd/rest/main.go
```

### Error: "failed to get admin"

Pastikan ada minimal 1 admin di database (dari seeder)

## Catatan

- Script menggunakan transaction untuk memastikan data consistency
- Barcode di-generate secara unique untuk mencegah collision
- Data random tapi realistis untuk keperluan testing
- Script bisa dijalankan berulang kali untuk menambah data
