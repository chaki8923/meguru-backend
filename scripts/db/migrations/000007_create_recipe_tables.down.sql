-- レシピ関連テーブル群の削除（ロールバック）

-- 外部キー制約があるため、依存関係を考慮して削除順序を設定
-- 1. saved_recipes テーブル（recipes と users に依存）
DROP TABLE IF EXISTS saved_recipes;

-- 2. recipe_steps テーブル（recipes に依存）
DROP TABLE IF EXISTS recipe_steps;

-- 3. recipe_seasonings テーブル（recipes に依存）
DROP TABLE IF EXISTS recipe_seasonings;

-- 4. recipe_ingredients テーブル（recipes に依存）
DROP TABLE IF EXISTS recipe_ingredients;

-- 5. recipes テーブル（他のテーブルに依存されていない）
DROP TABLE IF EXISTS recipes; 