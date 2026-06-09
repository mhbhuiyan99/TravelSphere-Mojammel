# TravelSphere

A full-stack travel destination explorer built with **Beego v2** (Go).  
Browse countries, explore attractions, and manage a personal travel wishlist.

---

## Quick Start

### 1. Prerequisites

- Go 1.21+
- `bee` CLI:
```bash
  go install github.com/beego/bee/v2@latest
```
- An OpenTripMap API key — get one free at https://opentripmap.io/

---

### 2. Clone the project

```bash
git clone https://github.com/YOUR_USERNAME/TravelSphere-Mojammel.git
cd TravelSphere-Mojammel
```

---

### 3. Create required files

**Config file** (copy the example and fill in your API key):

```bash
cp conf/app.conf.example conf/app.conf
```

Open `conf/app.conf` and set your key:

```ini
opentripmap_api_key = YOUR_KEY_HERE
```

**Wishlist data directory:**

```bash
mkdir -p data
echo '{}' > data/wishlist.json
```

---

### 4. Install dependencies

```bash
go mod tidy
```

---

### 5. Run

```bash
bee run
```

Visit **http://localhost:8080**

---

### 6. Run Tests

**All packages:**
```bash
go test ./...
```

**All packages with coverage:**
```bash
go test ./... -coverprofile=coverage.out
go tool cover -func=coverage.out
```

**Single folder (by package path):**
```bash
# utils only
go test ./utils/...

# services only
go test ./services/...

# controllers/api only
go test ./controllers/api/...
```

**Single folder with coverage:**
```bash
# utils
go test ./utils/... -coverprofile=coverage.out && go tool cover -func=coverage.out

# services
go test ./services/... -coverprofile=coverage.out && go tool cover -func=coverage.out
```

**Verbose output (see each test name):**
```bash
go test ./utils/... -v
go test ./services/... -v
```

With coverage:

```bash
go test ./... -coverprofile=coverage.out
go tool cover -func=coverage.out
```

---

## Project Structure

# TravelSphere-Mojammel Project Structure

```text
TravelSphere-Mojammel/
│
├── conf/
│   ├── app.conf                 # App config (gitignored — copy from example)
│   └── app.conf.example         # Template config for new contributors
│
├── controllers/
│   ├── api/
│   │   ├── country_api.go       # GET /api/countries, GET /api/countries/:slug
│   │   └── wishlist_api.go      # CRUD /api/wishlist, GET /api/dashboard/summary
│   │
│   ├── base.go                  # BaseController — Prepare(), layout, session
│   ├── auth.go                  # GET/POST /login
│   ├── logout.go                # GET /logout
│   ├── home.go                  # GET /
│   ├── country.go               # GET /countries, GET /countries/:slug
│   ├── wishlist.go              # GET /wishlist (auth protected)
│   └── dashboard.go             # GET /dashboard (auth protected)
│
├── filters/
│   ├── auth.go                  # Auth filter — redirects to /login if no session
│   └── logging.go               # Logging filter — logs method, path, duration
│
├── models/
│   ├── country.go               # Country struct
│   ├── attraction.go            # Attraction struct
│   └── wishlist.go              # WishlistItem struct, AllowedStatuses
│
├── routers/
│   └── router.go                # All SSR and API route registrations
│
├── services/
│   ├── country_service.go       # GetAllCountries, GetCountryBySlug, GetCountriesBySlugs
│   ├── attraction_service.go    # GetAttractionsByCoords, GetPopularAttractions
│   └── wishlist_service.go      # Add, Get, Update, Delete wishlist items
│
├── utils/
│   ├── country_client.go        # REST Countries API client
│   ├── attraction_client.go     # OpenTripMap API client
│   ├── wishlist_store.go        # File-based JSON store (read/write)
│   └── config.go                # configOrDefault helper
│
├── views/
│   ├── layouts/
│   │   └── main.tpl             # Base HTML layout with {{.LayoutContent}}
│   │
│   ├── partials/
│   │   ├── header.tpl           # Navbar — session-aware login/logout
│   │   ├── footer.tpl           # Footer
│   │   └── country_cards.tpl    # Country card grid (SSR + AJAX target)
│   │
│   └── pages/
│       ├── home.tpl             # Featured countries, attractions, autocomplete
│       ├── countries.tpl        # Country Explorer with search/filter
│       ├── destination.tpl      # Country detail with attractions + wishlist button
│       ├── wishlist.tpl         # Wishlist table with inline edit/delete
│       ├── dashboard.tpl        # Trip stats summary
│       ├── login.tpl            # Name-only login form
│       └── 404.tpl              # Not found page
│
├── static/
│   ├── css/                     # (empty — using Tailwind CDN)
│   └── js/
│       ├── countries.js         # AJAX country search with debounce
│       ├── autocomplete.js      # Home page search autocomplete
│       └── wishlist.js          # Wishlist add, save, delete, dashboard refresh
│
├── data/
│   └── wishlist.json            # Persisted wishlist data (gitignored)
│
├── main.go                      # Entry point — registers template functions, starts server
├── go.mod
├── go.sum
└── README.md
```

---

## Pages

| Route | Type | Description |
|---|---|---|
| `/` | SSR | Home — featured countries, popular attractions, search autocomplete |
| `/countries` | SSR + AJAX | Country Explorer — search and region filter update results without reload |
| `/countries/:slug` | SSR | Country detail — flag, info, attractions, Add to Wishlist |
| `/wishlist` | SSR + AJAX | Travel wishlist — edit notes, update status, delete (auth required) |
| `/dashboard` | SSR + AJAX | Trip summary stats (auth required) |
| `/login` | SSR | Name-only login — sets session |
| `/logout` | Redirect | Clears session |

---

## API Endpoints

| Method | Route | Description |
|---|---|---|
| GET | `/api/countries` | List countries — supports `?search=` and `?region=` |
| GET | `/api/countries/:slug` | Single country detail as JSON |
| GET | `/api/wishlist` | Get current user's wishlist |
| POST | `/api/wishlist` | Add destination to wishlist |
| PUT | `/api/wishlist/:id` | Update note and status |
| DELETE | `/api/wishlist/:id` | Remove item |
| GET | `/api/dashboard/summary` | Returns total, planned, visited counts |

---

## Storage Approach

**File-based JSON persistence via API client.**

Wishlist data is stored in `data/wishlist.json`, keyed by username:

```json
{
  "mojammel": [
    {
      "ID": "1234567890",
      "CountryName": "Bangladesh",
      "Note": "Visit Dhaka 2026",
      "Status": "Visited",
      "CreatedAt": "2025-06-01T10:00:00Z"
    }
  ]
}
```

`utils.ReadStore` / `utils.WriteStore` act as an API adapter — consistent with the
REST Countries and OpenTripMap client pattern. No database or ORM is used.

**Login persistence:** user identity is a session holding the entered name. The username
is the key in the JSON store, so wishlist data survives logout and server restart as long
as the same name is used on re-login.

`data/wishlist.json` and `conf/app.conf` are excluded from version control via `.gitignore`.

---

## Environment / Config

All config lives in `conf/app.conf`. See `conf/app.conf.example` for the full template.

| Key | Description | Default |
|---|---|---|
| `opentripmap_api_key` | OpenTripMap API key (required) | — |
| `restcountries_base_url` | REST Countries base URL | `https://restcountries.com/v3.1` |
| `opentripmap_base_url` | OpenTripMap base URL | `https://api.opentripmap.com/0.1/en/places` |
| `wishlist_store_path` | Path to wishlist JSON file | `data/wishlist.json` |
