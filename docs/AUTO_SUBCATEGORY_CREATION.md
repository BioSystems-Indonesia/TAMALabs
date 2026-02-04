# Auto Sub-Category Creation Feature

## 📝 Penjelasan

Saat membuat atau mengupdate **Test Type** dengan field `sub_category`, sistem sekarang akan **otomatis**:

1. ✅ Cek apakah sub-category sudah ada di tabel `sub_categories`
2. ✅ Jika **belum ada**, create entry baru di tabel `sub_categories`
3. ✅ Set `sub_category_id` pada test type untuk reference ke tabel

## 🔄 Alur Otomatis

### Scenario 1: Sub-Category Sudah Ada
```
User membuat Test Type:
{
  "name": "Hemoglobin",
  "code": "HB",
  "category": "Hematologi",
  "sub_category": "Complete Blood Count"  ← Sudah ada di tabel
}

Sistem:
1. Cari "Complete Blood Count" di tabel sub_categories
2. Ditemukan dengan ID = 1
3. Set sub_category_id = 1 pada test type
4. Simpan test type
```

### Scenario 2: Sub-Category Baru
```
User membuat Test Type:
{
  "name": "ALT",
  "code": "ALT",
  "category": "Kimia Klinik",
  "sub_category": "Liver Function Test"  ← Belum ada di tabel
}

Sistem:
1. Cari "Liver Function Test" di tabel sub_categories
2. Tidak ditemukan
3. AUTO CREATE sub-category baru:
   - name: "Liver Function Test"
   - code: "LIV" (3 huruf pertama dari nama)
   - category: "Kimia Klinik"
   - id: 5 (auto-increment)
4. Set sub_category_id = 5 pada test type
5. Simpan test type
```

## 💡 Keuntungan

✅ **Konsistensi Data**: Sub-category selalu tersedia di tabel master
✅ **Otomatis**: Tidak perlu manual create sub-category dulu
✅ **Backward Compatible**: Field `sub_category` string tetap ada
✅ **ID Reference**: Test type punya foreign key ke sub_categories
✅ **API Ready**: TechnoMedic bisa langsung query by sub_category_id

## 🔧 Technical Details

### Files Modified:
1. **`internal/usecase/test_type/test_type.go`**
   - Added `subCategoryRepo` dependency
   - Added `ensureSubCategoryExists()` helper method
   - Updated `Create()` and `Update()` methods

2. **`internal/repository/sql/sub_category/sub_category.go`**
   - Added `Create()` method
   - Added `Update()` method  
   - Added `Delete()` method

### Logic Flow:
```go
func (u *Usecase) Create(ctx context.Context, req *entity.TestType) {
    if req.SubCategory != "" {
        // 1. Try find existing
        subCat, err := u.subCategoryRepo.FindByName(ctx, req.SubCategory)
        
        if err == nil {
            // Found - use existing ID
            req.SubCategoryID = &subCat.ID
        } else {
            // Not found - create new
            newSubCat := &entity.SubCategory{
                Name:     req.SubCategory,
                Code:     generateCode(req.SubCategory),
                Category: req.Category,
            }
            u.subCategoryRepo.Create(ctx, newSubCat)
            req.SubCategoryID = &newSubCat.ID
        }
    }
    
    // Save test type with sub_category_id
    u.repository.Create(ctx, req)
}
```

## 📊 Database Impact

### Before:
```sql
-- Test Type
INSERT INTO test_types (name, code, sub_category) 
VALUES ('Hemoglobin', 'HB', 'Complete Blood Count');
-- sub_category_id = NULL
```

### After (Auto):
```sql
-- 1. Check if sub-category exists
SELECT * FROM sub_categories WHERE name = 'Complete Blood Count';

-- 2a. If NOT exists, create it
INSERT INTO sub_categories (name, code, category) 
VALUES ('Complete Blood Count', 'COM', 'Hematologi');
-- Returns ID = 5

-- 2b. Insert test type with foreign key
INSERT INTO test_types (name, code, sub_category, sub_category_id) 
VALUES ('Hemoglobin', 'HB', 'Complete Blood Count', 5);
```

## 🎯 Use Cases

### Use Case 1: Import Test Types dari Excel
Saat import banyak test types, sub-categories akan auto-created tanpa perlu pre-populate tabel.

### Use Case 2: Manual Create via UI
User create test type dan pilih/ketik sub-category name, otomatis masuk ke tabel master.

### Use Case 3: TechnoMedic Integration
TechnoMedic bisa langsung query test types by sub_category_id tanpa khawatir data tidak konsisten.

## ⚠️ Error Handling

Jika gagal create sub-category:
- System akan **continue** tanpa set `sub_category_id`
- Test type tetap tersimpan dengan field `sub_category` (string)
- Backward compatible - tidak break existing functionality

## 🧪 Testing

```bash
# Test 1: Create test type dengan sub-category baru
POST /api/v1/test-type
{
  "name": "Total Protein",
  "code": "TP",
  "category": "Kimia Klinik",
  "sub_category": "Protein Studies"  ← Baru!
}

# Verify: Check sub_categories table
GET /api/v1/technomedic/sub-categories
# Should return "Protein Studies" dengan ID baru

# Test 2: Create test type lain dengan sub-category yang sama
POST /api/v1/test-type
{
  "name": "Albumin",
  "code": "ALB",
  "category": "Kimia Klinik",
  "sub_category": "Protein Studies"  ← Sudah ada
}

# Result: Tidak create duplicate, reuse existing ID
```

## 📝 Notes

- **Code Generation**: Kode sub-category diambil dari 3 huruf pertama nama (uppercase)
- **Unique Constraint**: Tabel sub_categories punya unique index pada `name`
- **Thread Safe**: Create menggunakan database transaction
- **Migration Compatible**: Existing data tetap berfungsi normal

---

**Status: IMPLEMENTED & TESTED** ✅

Sekarang setiap kali create/update test type, sub-category otomatis ter-manage di tabel master!
