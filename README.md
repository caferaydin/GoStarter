# Go Starter Auth API 🔐

Bu proje, Go ile JWT kimlik doğrulama ve katmanlı mimariyi temel alan, her projede kullanılabilecek **starter bir API iskeletidir**. Amaç; temiz kod, güvenli yapı, kolay geliştirme.

## Özellikler 🚀
- Katmanlı yapı: `handler`, `service`, `repository`, `model`, `middleware`, `config`
- JWT ile kullanıcı girişi ve korumalı endpoint
- Şifreleme: `bcrypt` ile password hashing
- Environment config: `.env` üzerinden güvenli yapılandırma
- PostgreSQL veritabanı bağlantısı (`sqlx`)

---

## Kullanım 🛠

### 1. .env Dosyası Oluştur
```env
DB_CONN=host=localhost port=5432 user=postgres password=1234 dbname=go_auth sslmode=disable
JWT_SECRET=SuperSecretKey
```

### 2. Veritabanı Yapılandırması
```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    password VARCHAR(100) NOT NULL
);

CREATE TABLE refresh_tokens (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    token TEXT NOT NULL,
    expires_at TIMESTAMP NOT NULL
);

```

### 3. Modülleri Yükle
```bash
go mod tidy
```

### 4. Uygulamayı Başlat
```bash
go run ./cmd/main.go
```

---

## API Endpoint'leri 🌐

| Yöntem | URL            | Açıklama             |
|--------|----------------|----------------------|
| POST   | /register      | Yeni kullanıcı kaydı |
| POST   | /login         | Giriş ve token alma  |
| GET    | /api/profile   | Token ile erişilen profil |
| POST   | /refresh-token | refresh-token        |

---

## Katmanlar Hakkında 📦
- `internal/config`: .env yönetimi, DB bağlantısı
- `internal/util`: JWT üretimi (optimize edilebilir)
- `internal/middleware`: Token doğrulama ve context taşıma
- `internal/service`: İş mantığı, şifre doğrulama/hash
- `internal/repository`: DB işlemleri (sqlx ile)

---

## İlham Veren Noktalar ✨
- Go'da Dependency Injection **manuel ama sade**
- Context ile kullanıcı verisi taşıma
- .env yönetimi ile güvenli yapılandırma
- Şifreleme ve JWT ile güvenlik temelleri

---

## Geliştirme Planları 🧩
- Token refresh sistemi
- Role-based yetkilendirme
- Swagger/OpenAPI dokümantasyonu
- Dockerfile & Compose entegrasyonu

---

## Lisans 📄
MIT

---

> "Temiz mimari, anlaşılır kod, kolay geliştirme. Her projede bir adım önde başla!"
