CREATE TABLE IF NOT EXISTS contents (
    id            VARCHAR(255)     NOT NULL,
    source        VARCHAR(100)     NOT NULL,
    title         VARCHAR(500)     NOT NULL,
    type          VARCHAR(10)      NOT NULL,
    url           TEXT,
    published_at  TIMESTAMPTZ      NOT NULL,
    views         INT              NOT NULL DEFAULT 0,
    likes         INT              NOT NULL DEFAULT 0,
    reading_time  INT              NOT NULL DEFAULT 0,
    reactions     INT              NOT NULL DEFAULT 0,
    score         DOUBLE PRECISION NOT NULL DEFAULT 0,
    search_vector TSVECTOR         GENERATED ALWAYS AS (to_tsvector('simple', title)) STORED,
    created_at    TIMESTAMPTZ      NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ      NOT NULL DEFAULT NOW(),
    PRIMARY KEY (id, source)
);

CREATE INDEX IF NOT EXISTS idx_contents_type ON contents (type);
CREATE INDEX IF NOT EXISTS idx_contents_score ON contents (score DESC);
CREATE INDEX IF NOT EXISTS idx_contents_published_at ON contents (published_at DESC);
CREATE INDEX IF NOT EXISTS idx_contents_search ON contents USING gin(search_vector);
