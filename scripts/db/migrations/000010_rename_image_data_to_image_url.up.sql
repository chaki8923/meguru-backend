-- 既存のimage_dataカラムをimage_urlに変更

-- 1. 新しいimage_urlカラムを追加
ALTER TABLE recipes ADD COLUMN image_url TEXT;

-- 2. 既存のimage_dataカラムを削除
ALTER TABLE recipes DROP COLUMN image_data;

-- 3. コメントを更新
COMMENT ON COLUMN recipes.image_url IS 'レシピ画像URL（オプション）'; 