# Reference Range Implementation - Simplified Approach

## Overview
Implementasi reference range berdasarkan umur dan jenis kelamin menggunakan pendekatan **single table** dengan kolom JSON, bukan tabel terpisah.

## Database Schema

### Tabel: `test_types`
Ditambahkan 1 kolom baru:
- `specific_ref_ranges` (TEXT/JSON): Menyimpan array of reference ranges spesifik

### Structure JSON `specific_ref_ranges`:
```json
[
  {
    "gender": "M",          // "M", "F", or null
    "age_min": 0,           // minimum age (years), null = no minimum
    "age_max": 18,          // maximum age (years), null = no maximum  
    "low_ref_range": 3.5,   // numeric low value
    "high_ref_range": 5.0,  // numeric high value
    "normal_ref_string": null  // OR string value like "Negative"
  },
  {
    "gender": "F",
    "age_min": 18,
    "age_max": null,
    "low_ref_range": 4.0,
    "high_ref_range": 5.5,
    "normal_ref_string": null
  }
]
```

## Backend Implementation

### Entity (`internal/entity/test_type.go`)
```go
type TestType struct {
    // ... existing fields ...
    SpecificRefRanges    []SpecificReferenceRange  `json:"specific_ref_ranges,omitempty" gorm:"-"`
    SpecificRefRangesDB  string                    `json:"-" gorm:"column:specific_ref_ranges"`
}

type SpecificReferenceRange struct {
    Gender          *string  `json:"gender,omitempty"`
    AgeMin          *float64 `json:"age_min,omitempty"`
    AgeMax          *float64 `json:"age_max,omitempty"`
    LowRefRange     *float64 `json:"low_ref_range,omitempty"`
    HighRefRange    *float64 `json:"high_ref_range,omitempty"`
    NormalRefString *string  `json:"normal_ref_string,omitempty"`
}
```

### Key Methods
- `TestType.GetReferenceRangeForPatient(age, gender)` - Get reference range untuk pasien spesifik
- `SpecificReferenceRange.MatchesCriteria(age, gender)` - Check apakah criteria cocok
- `BeforeCreate/BeforeUpdate` - Serialize array ke JSON string
- `AfterFind` - Deserialize JSON string ke array

## Frontend Implementation

### Type Definition (`web/src/types/test_type.ts`)
```typescript
export interface TestType {
    // ... existing fields ...
    specific_ref_ranges?: SpecificReferenceRange[];
}

export interface SpecificReferenceRange {
    gender?: string | null;
    age_min?: number | null;
    age_max?: number | null;
    low_ref_range?: number | null;
    high_ref_range?: number | null;
    normal_ref_string?: string | null;
}
```

## Usage Example

### Creating Test Type with Specific Ranges via API
```json
POST /api/v1/test-type
{
    "name": "Hemoglobin",
    "code": "HGB",
    "unit": "g/dL",
    "low_ref_range": 12.0,
    "high_ref_range": 16.0,
    "specific_ref_ranges": [
        {
            "gender": "M",
            "age_min": 18,
            "age_max": null,
            "low_ref_range": 13.5,
            "high_ref_range": 17.5
        },
        {
            "gender": "F",
            "age_min": 18,
            "age_max": null,
            "low_ref_range": 12.0,
            "high_ref_range": 15.5
        },
        {
            "gender": null,
            "age_min": 0,
            "age_max": 18,
            "low_ref_range": 11.0,
            "high_ref_range": 15.0
        }
    ]
}
```

### Creating/Editing Test Type via Admin UI

1. **Navigate to Test Type**
   - Go to "Test Type" menu in admin panel
   - Click "Create" or "Edit" existing test type

2. **Fill Basic Information**
   - Name, Code, Category, etc.
   - Fill default/global reference ranges (Low Range, High Range, or Normal String)

3. **Add Specific Reference Ranges** (Optional)
   - Scroll to "Specific Reference Ranges" section
   - Click "Add" button to add new criteria
   - For each range, specify:
     - **Gender**: All Genders / Male / Female
     - **Min Age**: Minimum age in years (leave empty for no minimum)
     - **Max Age**: Maximum age in years (leave empty for no maximum)
     - **Low Range & High Range**: For numeric values
     - **OR Normal String**: For text values like "Negative", "Positive"

4. **Example Scenarios**:

   **Scenario 1: Age-based only**
   ```
   Range 1: All Genders, 0-18 years → 11.0 - 15.0
   Range 2: All Genders, 18+ years → 12.0 - 16.0
   ```

   **Scenario 2: Gender and age based**
   ```
   Range 1: Male, 18+ years → 13.5 - 17.5
   Range 2: Female, 18+ years → 12.0 - 15.5
   Range 3: All Genders, 0-18 years → 11.0 - 15.0
   ```

   **Scenario 3: Mixed numeric and string**
   ```
   Range 1: All Genders, All Ages → Normal String: "Negative"
   (For qualitative tests like pregnancy test)
   ```

5. **Validation**
   - System will automatically select the most specific range that matches patient criteria
   - If no specific range matches, default range will be used

### Reference Range Selection Logic
1. System mencari `specific_ref_ranges` yang match dengan criteria pasien (age + gender)
2. Jika ada yang match, gunakan range tersebut
3. Jika tidak ada yang match, fallback ke global range (`low_ref_range`, `high_ref_range`, `normal_ref_string`)

## Migration
File: `migrations/20251112223547_add_test_type_reference_ranges.up.sql`
```sql
ALTER TABLE `test_types` ADD COLUMN `specific_ref_ranges` text NULL;
```

## Advantages
✅ Sederhana - tidak perlu table & relationship terpisah  
✅ Atomic - semua data test type dalam 1 record  
✅ Flexible - mudah menambah/edit reference ranges  
✅ Performance - tidak perlu JOIN untuk query  
✅ Easy to maintain - less code complexity  

## Notes
- Gunakan `specific_ref_ranges` untuk criteria-based ranges
- Global ranges (`low_ref_range`, `high_ref_range`, `normal_ref_string`) tetap digunakan sebagai fallback
- Gender values: `"M"`, `"F"`, atau `null` (all genders)
- Age dalam tahun (years), bisa decimal (e.g., 0.5 = 6 bulan)
