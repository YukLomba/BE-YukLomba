# BE-YukLomba
Backend for  YukLomba

## 🚀 Getting Started

### ✅ Prerequisites

- Go 1.21+
- A running PostgreSQL/MySQL database (or use SQLite)
- [Air](https://github.com/cosmtrek/air) for live reloading (optional but recommended)

### 📄 Documentation
Link to API Documentation
[Click here.](https://www.postman.com/yuk-lomba-1446/yuk-lomba-workspace/collection/to4rv27/yuklomba-api-documentation)

### 📥 Installation

```bash
git clone https://github.com/YukLomba/BE-YukLomba.git
cd BE-YukLomba
go mod tidy
```

### 🔧 Configuration
```bash
cp .env.example .env
```

### 💾 Migration
```bash
go run cmd/api/main.go db --migrate
```

migrate with fresh

```bash
go run cmd/api/main.go db --migrate --fresh
```

### 💾 Seeder
```bash
go run cmd/api/main.go db --seed
```

### 🚀 Running the Application
```bash
go run main.go
```

or with live reloading

```bash
air
```
