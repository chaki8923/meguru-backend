-- image_urlカラムをimage_dataに戻す

-- 1. 新しいimage_dataカラムを追加
ALTER TABLE recipes ADD COLUMN image_data BYTEA;

-- 2. 既存のimage_urlカラムを削除
ALTER TABLE recipes DROP COLUMN image_url;

-- 3. コメントを更新
COMMENT ON COLUMN recipes.image_data IS 'レシピ画像データ（オプション）'; 