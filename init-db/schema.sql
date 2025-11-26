-- photo_locations テーブル作成
CREATE TABLE IF NOT EXISTS photo_locations (
    location_id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    latitude DOUBLE PRECISION NOT NULL,
    longitude DOUBLE PRECISION NOT NULL,
    geohash VARCHAR(9) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- インデックスを付与
CREATE INDEX idx_photo_locations_geohash ON photo_locations (geohash);
CREATE INDEX IF NOT EXISTS idx_photo_locations_user_id ON photo_locations(user_id);