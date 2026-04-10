-- Disky initial schema
-- Run once on first startup via docker-entrypoint-initdb.d

CREATE EXTENSION IF NOT EXISTS "pgcrypto";
CREATE EXTENSION IF NOT EXISTS "pg_trgm";

-- ============================================================
-- USERS
-- ============================================================
CREATE TABLE users (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email       VARCHAR(255) UNIQUE NOT NULL,
    name        VARCHAR(255),
    avatar_url  TEXT,
    google_id   VARCHAR(255) UNIQUE,
    microsoft_id VARCHAR(255) UNIQUE,
    apple_id    VARCHAR(255) UNIQUE,
    email_verified BOOLEAN DEFAULT FALSE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_users_email ON users(email);

-- ============================================================
-- LOCALES / COUNTRIES
-- ============================================================
CREATE TABLE locales (
    code            VARCHAR(10) PRIMARY KEY,   -- 'uk', 'us', 'de', etc.
    name            VARCHAR(100) NOT NULL,
    currency_code   VARCHAR(3) NOT NULL,       -- 'GBP', 'USD', etc.
    currency_symbol VARCHAR(5) NOT NULL,
    amazon_domain   VARCHAR(100) NOT NULL,     -- 'www.amazon.co.uk'
    amazon_region   VARCHAR(50) NOT NULL,      -- 'eu-west-1'
    partner_tag     VARCHAR(100) NOT NULL      -- affiliate tag
);

INSERT INTO locales VALUES
    ('uk', 'United Kingdom', 'GBP', '£', 'www.amazon.co.uk', 'eu-west-1', 'prbox-21');

-- ============================================================
-- CATEGORIES
-- ============================================================
CREATE TABLE categories (
    id          SERIAL PRIMARY KEY,
    name        VARCHAR(100) NOT NULL,
    slug        VARCHAR(100) UNIQUE NOT NULL,
    parent_id   INTEGER REFERENCES categories(id),
    product_type VARCHAR(20) NOT NULL CHECK (product_type IN ('storage','ram'))
);

INSERT INTO categories (name, slug, product_type) VALUES
    -- Storage categories
    ('Internal HDD',          'internal-hdd',    'storage'),
    ('Internal SSD',          'internal-ssd',    'storage'),
    ('M.2 NVMe SSD',          'm2-nvme',         'storage'),
    ('M.2 SATA SSD',          'm2-sata',         'storage'),
    ('External HDD',          'external-hdd',    'storage'),
    ('External SSD',          'external-ssd',    'storage'),
    ('NAS HDD',               'nas-hdd',         'storage'),
    ('USB Flash Drive',       'usb-flash',       'storage'),
    ('SD Card',               'sd-card',         'storage'),
    ('Enterprise SSD',        'enterprise-ssd',  'storage'),
    -- RAM categories
    ('Desktop RAM (DIMM)',     'desktop-ram',     'ram'),
    ('Laptop RAM (SO-DIMM)',   'laptop-ram',      'ram'),
    ('Server / ECC RAM',      'server-ram',      'ram'),
    ('DDR5 RAM',              'ddr5-ram',        'ram'),
    ('DDR4 RAM',              'ddr4-ram',        'ram');

-- ============================================================
-- PRODUCTS
-- ============================================================
CREATE TABLE products (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    asin            VARCHAR(20) NOT NULL,
    locale          VARCHAR(10) NOT NULL REFERENCES locales(code) DEFAULT 'uk',
    title           TEXT NOT NULL,
    brand           VARCHAR(150),
    model           VARCHAR(150),
    category_id     INTEGER REFERENCES categories(id),

    -- Common
    image_url       TEXT,
    amazon_url      TEXT,

    -- Storage fields
    capacity_gb         DECIMAL(12,3),
    form_factor         VARCHAR(50),   -- '3.5"', '2.5"', 'M.2 2280', 'mSATA', etc.
    storage_interface   VARCHAR(50),   -- 'SATA III', 'NVMe PCIe 4.0', 'PCIe 5.0', 'USB 3.2', etc.
    rpm                 INTEGER,       -- HDD rotation speed
    cache_mb            INTEGER,       -- HDD cache
    read_speed_mbps     INTEGER,
    write_speed_mbps    INTEGER,
    tbw                 INTEGER,       -- SSD terabytes written endurance

    -- RAM fields
    ram_type        VARCHAR(20),   -- 'DDR4', 'DDR5', 'DDR3', 'LPDDR5', etc.
    ram_speed_mhz   INTEGER,
    ram_capacity_gb DECIMAL(8,2),
    ram_modules     INTEGER,       -- kit size: 1, 2, 4
    ram_form_factor VARCHAR(20),   -- 'DIMM', 'SO-DIMM', 'ECC DIMM'
    is_ecc          BOOLEAN DEFAULT FALSE,
    buffer_type     VARCHAR(30),   -- 'Unbuffered', 'Registered', 'Fully Buffered', 'Load Reduced'
    cas_latency     INTEGER,       -- CL value e.g. 16
    timing          VARCHAR(30),   -- e.g. 'CL16-18-18-38'
    voltage         DECIMAL(4,2),  -- e.g. 1.35

    last_fetched_at TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE (asin, locale)
);

CREATE INDEX idx_products_asin    ON products(asin);
CREATE INDEX idx_products_locale  ON products(locale);
CREATE INDEX idx_products_category ON products(category_id);
CREATE INDEX idx_products_capacity ON products(capacity_gb);
CREATE INDEX idx_products_title_trgm ON products USING gin(title gin_trgm_ops);

-- ============================================================
-- PRICE HISTORY
-- ============================================================
CREATE TABLE price_history (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id  UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    price       DECIMAL(10,2) NOT NULL,
    currency    VARCHAR(3) NOT NULL DEFAULT 'GBP',
    availability VARCHAR(50),
    condition   VARCHAR(20) DEFAULT 'New',  -- 'New', 'Used', 'Collectible'
    merchant    VARCHAR(150),
    is_prime    BOOLEAN DEFAULT FALSE,
    recorded_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_price_history_product_id ON price_history(product_id);
CREATE INDEX idx_price_history_recorded_at ON price_history(recorded_at DESC);

-- Materialised view: latest price per product
CREATE MATERIALIZED VIEW latest_prices AS
SELECT DISTINCT ON (product_id)
    product_id,
    price,
    currency,
    availability,
    condition,
    merchant,
    is_prime,
    recorded_at
FROM price_history
ORDER BY product_id, recorded_at DESC;

CREATE UNIQUE INDEX idx_latest_prices_product ON latest_prices(product_id);

-- ============================================================
-- PRICE ALERTS
-- ============================================================
CREATE TABLE price_alerts (
    id                      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id                 UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    alert_type              VARCHAR(30) NOT NULL CHECK (
        alert_type IN ('product_price','category_price_per_gb','category_price_per_unit')
    ),

    -- For product-level alerts
    product_id              UUID REFERENCES products(id) ON DELETE CASCADE,

    -- For category-level alerts
    category_id             INTEGER REFERENCES categories(id) ON DELETE CASCADE,
    category_filters        JSONB,  -- e.g. {"form_factor":"M.2 2280","ram_type":"DDR5"}

    -- Thresholds (set at least one)
    threshold_price         DECIMAL(10,2),     -- absolute price
    threshold_price_per_gb  DECIMAL(10,4),     -- price per GB

    is_active               BOOLEAN NOT NULL DEFAULT TRUE,
    notify_email            BOOLEAN NOT NULL DEFAULT TRUE,
    last_triggered_at       TIMESTAMPTZ,
    created_at              TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CHECK (
        (product_id IS NOT NULL AND category_id IS NULL) OR
        (product_id IS NULL AND category_id IS NOT NULL)
    )
);

CREATE INDEX idx_alerts_user_id    ON price_alerts(user_id);
CREATE INDEX idx_alerts_product_id ON price_alerts(product_id);
CREATE INDEX idx_alerts_active     ON price_alerts(is_active) WHERE is_active = TRUE;

-- ============================================================
-- ALERT NOTIFICATIONS LOG
-- ============================================================
CREATE TABLE alert_notifications (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    alert_id    UUID NOT NULL REFERENCES price_alerts(id) ON DELETE CASCADE,
    product_id  UUID REFERENCES products(id),
    price       DECIMAL(10,2),
    price_per_gb DECIMAL(10,4),
    sent_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    channel     VARCHAR(20) NOT NULL DEFAULT 'email'
);

-- ============================================================
-- FUNCTIONS & TRIGGERS
-- ============================================================

-- Auto-update updated_at
CREATE OR REPLACE FUNCTION update_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();

CREATE TRIGGER trg_products_updated_at
    BEFORE UPDATE ON products
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();

-- Function to refresh latest_prices materialized view
CREATE OR REPLACE FUNCTION refresh_latest_prices()
RETURNS TRIGGER AS $$
BEGIN
    REFRESH MATERIALIZED VIEW CONCURRENTLY latest_prices;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

-- Trigger after price insert (deferred refresh via pg_notify instead for performance)
CREATE OR REPLACE FUNCTION notify_price_inserted()
RETURNS TRIGGER AS $$
BEGIN
    PERFORM pg_notify('price_inserted', NEW.product_id::text);
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_notify_price_inserted
    AFTER INSERT ON price_history
    FOR EACH ROW EXECUTE FUNCTION notify_price_inserted();
