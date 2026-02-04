# TechnoMedic Integration Configuration

## 📋 Overview

TechnoMedic API integration memiliki sistem konfigurasi on/off yang terintegrasi dengan halaman Config aplikasi, sama seperti SIMRS dan Database Sharing bridging yang sudah ada.

## 🔧 Implementasi

### 1. Config Schema
**File**: `config/schema.go`

Added field:
```go
// TechnoMedic Configuration
TechnoMedicIntegrationEnabled string `validate:"-"`
```

### 2. Seed Config (Default Value)
**File**: `internal/app/seed_config.go`

Added default config:
```go
{
    ID:    "TechnoMedicIntegrationEnabled",
    Value: "false",  // Default: disabled
}
```

### 3. Integration Check Middleware
**File**: `internal/middleware/integration_check.go`

Added middleware method:
```go
func (m *IntegrationCheckMiddleware) CheckTechnoMedicEnabled() echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            ctx := c.Request().Context()
            
            // Check if TechnoMedic integration is enabled
            enabled, err := m.configGetter.Get(ctx, "TechnoMedicIntegrationEnabled")
            if err != nil || enabled != "true" {
                return echo.NewHTTPError(http.StatusForbidden, 
                    "TechnoMedic integration is not enabled")
            }
            
            return next(c)
        }
    }
}
```

### 4. Handler with Middleware
**File**: `internal/delivery/rest/technomedic.go`

Updated handler to accept and use middleware:
```go
type TechnoMedicHandler struct {
    usecase               *technomedicuc.Usecase
    integrationMiddleware *appMiddleware.IntegrationCheckMiddleware
}

func NewTechnoMedicHandler(
    usecase *technomedicuc.Usecase,
    integrationMiddleware *appMiddleware.IntegrationCheckMiddleware,
) *TechnoMedicHandler {
    return &TechnoMedicHandler{
        usecase:               usecase,
        integrationMiddleware: integrationMiddleware,
    }
}

func (h *TechnoMedicHandler) RegisterRoutes(router *echo.Group) {
    // Apply middleware to ALL TechnoMedic routes
    technomedic := router.Group("/technomedic", 
        h.integrationMiddleware.CheckTechnoMedicEnabled())
    
    // All routes protected by middleware
    technomedic.GET("/test-types", h.GetTestTypes)
    technomedic.GET("/sub-categories", h.GetSubCategories)
    technomedic.POST("/order", h.CreateOrder)
    // ... etc
}
```

### 5. Provider Update
**File**: `internal/app/provider.go`

Updated provider to inject middleware:
```go
func provideTechnoMedicHandler(
    technomedicUC *technomedicuc.Usecase,
    integrationCheckMiddleware *middleware.IntegrationCheckMiddleware,
) *rest.TechnoMedicHandler {
    return rest.NewTechnoMedicHandler(technomedicUC, integrationCheckMiddleware)
}
```

## 🎯 Cara Kerja

### Flow Diagram:
```
Request → Middleware Check → Handler
          ↓
    Config: TechnoMedicIntegrationEnabled
          ↓
    if "true" → Continue to Handler
    if "false" → Return 403 Forbidden
```

### Saat Config Disabled (Default):
```bash
curl http://localhost:8080/api/v1/technomedic/test-types

# Response:
HTTP/1.1 403 Forbidden
{
  "message": "TechnoMedic integration is not enabled"
}
```

### Saat Config Enabled:
```bash
# 1. Enable di database atau UI config
UPDATE configs SET value = 'true' WHERE id = 'TechnoMedicIntegrationEnabled';

# 2. Test endpoint
curl http://localhost:8080/api/v1/technomedic/test-types

# Response:
HTTP/1.1 200 OK
{
  "code": 200,
  "status": "success",
  "data": [...]
}
```

## 📊 Comparison dengan Integration Lain

### SIMRS Bridging:
```go
// Config: SimrsIntegrationEnabled
simrsExternalHandler.RegisterRoutes(unauthenticatedV1)
// Uses CheckSimrsEnabled() middleware
```

### Database Sharing (SIMGOS):
```go
// Config: SimgosIntegrationEnabled  
// Uses CheckSimgosEnabled() middleware
```

### TechnoMedic (NEW):
```go
// Config: TechnoMedicIntegrationEnabled
technomedicHandler.RegisterRoutes(unauthenticatedV1)
// Uses CheckTechnoMedicEnabled() middleware
```

## 🔐 Security Benefits

1. ✅ **Prevent Unauthorized Access**: Jika tidak dikonfigurasi, API tidak bisa diakses
2. ✅ **Admin Control**: Hanya admin yang bisa enable/disable via config page
3. ✅ **Consistent Pattern**: Sama seperti SIMRS/SIMGOS bridging
4. ✅ **Runtime Toggle**: Bisa enable/disable tanpa restart aplikasi
5. ✅ **Clear Error Messages**: User langsung tahu integration belum diaktifkan

## 📝 Configuration Management

### Via UI (Recommended):

**Steps:**
1. Login sebagai **Admin**
2. Navigasi ke **Config** page (menu Settings)
3. Cari section **"SIMRS Bridging"**
4. Enable toggle **"SIMRS Bridging"**
5. Pada dropdown **"SIMRS"**, pilih **"TechnoMedic (API)"**
6. Documentation accordion akan muncul dengan endpoint details
7. Click **"Save"** button
8. Config akan tersimpan dan API endpoints akan aktif

**Screenshot Flow:**
```
[x] SIMRS Bridging
    ┌─────────────────────────────┐
    │ SIMRS: TechnoMedic (API) ▼  │
    └─────────────────────────────┘
    
    ▼ TechnoMedic API Documentation
      • GET /api/v1/technomedic/test-types
      • GET /api/v1/technomedic/sub-categories
      • POST /api/v1/technomedic/order
      • GET /api/v1/technomedic/order/:no_order
      
    [Save] button
```

### Via Database:
```sql
-- Enable
UPDATE configs SET value = 'true' WHERE id = 'TechnoMedicIntegrationEnabled';

-- Disable  
UPDATE configs SET value = 'false' WHERE id = 'TechnoMedicIntegrationEnabled';

-- Check current status
SELECT id, value FROM configs WHERE id = 'TechnoMedicIntegrationEnabled';
```

### Via API (jika ada endpoint config):
```bash
# Enable
curl -X PUT http://localhost:8080/api/v1/config/TechnoMedicIntegrationEnabled \
  -H "Content-Type: application/json" \
  -d '{"value": "true"}'

# Disable
curl -X PUT http://localhost:8080/api/v1/config/TechnoMedicIntegrationEnabled \
  -H "Content-Type: application/json" \
  -d '{"value": "false"}'
```

## 🧪 Testing Integration Config

### Test 1: Config Disabled (Default)
```bash
# Semua endpoint harus return 403
curl http://localhost:8080/api/v1/technomedic/test-types
curl http://localhost:8080/api/v1/technomedic/sub-categories
curl http://localhost:8080/api/v1/technomedic/doctors

# Expected: HTTP 403 Forbidden
```

### Test 2: Enable Config
```sql
UPDATE configs SET value = 'true' WHERE id = 'TechnoMedicIntegrationEnabled';
```

### Test 3: Config Enabled
```bash
# Semua endpoint harus return 200 OK
curl http://localhost:8080/api/v1/technomedic/test-types
curl http://localhost:8080/api/v1/technomedic/sub-categories
curl http://localhost:8080/api/v1/technomedic/doctors

# Expected: HTTP 200 OK dengan data
```

## 📁 Files Modified

1. ✅ `config/schema.go` - Added TechnoMedicIntegrationEnabled field
2. ✅ `internal/app/seed_config.go` - Added default config
3. ✅ `internal/middleware/integration_check.go` - Added CheckTechnoMedicEnabled()
4. ✅ `internal/delivery/rest/technomedic.go` - Apply middleware to routes
5. ✅ `internal/app/provider.go` - Inject middleware to handler
6. ✅ `internal/app/wire_gen.go` - Auto-regenerated by Wire
7. ✅ `docs/TECHNOMEDIC_API.md` - Updated with config requirements

## ✅ Benefits

- ✅ **Security**: API tidak bisa diakses tanpa aktivasi
- ✅ **Control**: Admin control via config page
- ✅ **Consistency**: Pattern sama dengan SIMRS/SIMGOS
- ✅ **User-Friendly**: Clear error messages
- ✅ **Production-Ready**: Safe default (disabled)

---

**Status: IMPLEMENTED** ✅

TechnoMedic API sekarang memiliki kontrol akses via konfigurasi, sama seperti SIMRS dan Database Sharing bridging!
