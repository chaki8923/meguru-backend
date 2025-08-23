-- レシピ関連テーブル群の作成

-- 1. recipes テーブル
CREATE TABLE recipes (
    id BIGSERIAL PRIMARY KEY,
    recipe_id VARCHAR(255) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    author_comment TEXT,
    cook_time INTEGER, -- 分単位
    calories INTEGER,
    total_price INTEGER, -- 円単位
    cooking_point TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- 2. recipe_ingredients テーブル（レシピの材料）
CREATE TABLE recipe_ingredients (
    id BIGSERIAL PRIMARY KEY,
    recipe_ingredient_id VARCHAR(255) NOT NULL UNIQUE,
    recipe_id VARCHAR(255) NOT NULL REFERENCES recipes(recipe_id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    display_order INTEGER NOT NULL,
    amount_text VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- 3. recipe_seasonings テーブル（レシピの調味料）
CREATE TABLE recipe_seasonings (
    id BIGSERIAL PRIMARY KEY,
    recipe_seasoning_id VARCHAR(255) NOT NULL UNIQUE,
    recipe_id VARCHAR(255) NOT NULL REFERENCES recipes(recipe_id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    display_order INTEGER NOT NULL,
    amount_text VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- 4. recipe_steps テーブル（レシピの手順）
CREATE TABLE recipe_steps (
    id BIGSERIAL PRIMARY KEY,
    recipe_step_id VARCHAR(255) NOT NULL UNIQUE,
    recipe_id VARCHAR(255) NOT NULL REFERENCES recipes(recipe_id) ON DELETE CASCADE,
    instruction TEXT NOT NULL,
    step_number INTEGER NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- 5. saved_recipes テーブル（ユーザーが保存したレシピ）
CREATE TABLE saved_recipes (
    id BIGSERIAL PRIMARY KEY,
    saved_recipe_id VARCHAR(255) NOT NULL UNIQUE,
    recipe_id VARCHAR(255) NOT NULL REFERENCES recipes(recipe_id) ON DELETE CASCADE,
    user_id VARCHAR(255) NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    saved_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    
    -- 同じユーザーが同じレシピを複数回保存することを防ぐ
    UNIQUE(user_id, recipe_id)
);

-- インデックスの作成
CREATE INDEX idx_recipes_recipe_id ON recipes(recipe_id);
CREATE INDEX idx_recipes_name ON recipes(name);
CREATE INDEX idx_recipes_deleted_at ON recipes(deleted_at) WHERE deleted_at IS NULL;
CREATE INDEX idx_recipe_ingredients_recipe_id ON recipe_ingredients(recipe_id);
CREATE INDEX idx_recipe_ingredients_display_order ON recipe_ingredients(recipe_id, display_order);
CREATE INDEX idx_recipe_ingredients_deleted_at ON recipe_ingredients(deleted_at) WHERE deleted_at IS NULL;
CREATE INDEX idx_recipe_seasonings_recipe_id ON recipe_seasonings(recipe_id);
CREATE INDEX idx_recipe_seasonings_display_order ON recipe_seasonings(recipe_id, display_order);
CREATE INDEX idx_recipe_seasonings_deleted_at ON recipe_seasonings(deleted_at) WHERE deleted_at IS NULL;
CREATE INDEX idx_recipe_steps_recipe_id ON recipe_steps(recipe_id);
CREATE INDEX idx_recipe_steps_step_number ON recipe_steps(recipe_id, step_number);
CREATE INDEX idx_recipe_steps_deleted_at ON recipe_steps(deleted_at) WHERE deleted_at IS NULL;
CREATE INDEX idx_saved_recipes_user_id ON saved_recipes(user_id);
CREATE INDEX idx_saved_recipes_recipe_id ON saved_recipes(recipe_id);
CREATE INDEX idx_saved_recipes_saved_at ON saved_recipes(saved_at);
CREATE INDEX idx_saved_recipes_deleted_at ON saved_recipes(deleted_at) WHERE deleted_at IS NULL;

-- テーブルコメント
COMMENT ON TABLE recipes IS 'レシピ基本情報テーブル';
COMMENT ON COLUMN recipes.recipe_id IS 'レシピの一意識別子';
COMMENT ON COLUMN recipes.name IS 'レシピ名';
COMMENT ON COLUMN recipes.author_comment IS 'レシピ作成者コメント';
COMMENT ON COLUMN recipes.cook_time IS '調理時間（分）';
COMMENT ON COLUMN recipes.calories IS 'カロリー';
COMMENT ON COLUMN recipes.total_price IS 'レシピ合計金額（円）';
COMMENT ON COLUMN recipes.cooking_point IS '作る時のポイント';
COMMENT ON COLUMN recipes.deleted_at IS '論理削除日時（NULLの場合は削除されていない）';

COMMENT ON TABLE recipe_ingredients IS 'レシピ材料テーブル';
COMMENT ON COLUMN recipe_ingredients.recipe_ingredient_id IS 'レシピ材料の一意識別子';
COMMENT ON COLUMN recipe_ingredients.recipe_id IS 'レシピID（外部キー）';
COMMENT ON COLUMN recipe_ingredients.name IS '材料名';
COMMENT ON COLUMN recipe_ingredients.display_order IS '材料の表示順序';
COMMENT ON COLUMN recipe_ingredients.amount_text IS '分量・個数などのテキスト';
COMMENT ON COLUMN recipe_ingredients.deleted_at IS '論理削除日時（NULLの場合は削除されていない）';

COMMENT ON TABLE recipe_seasonings IS 'レシピ調味料テーブル';
COMMENT ON COLUMN recipe_seasonings.recipe_seasoning_id IS 'レシピ調味料の一意識別子';
COMMENT ON COLUMN recipe_seasonings.recipe_id IS 'レシピID（外部キー）';
COMMENT ON COLUMN recipe_seasonings.name IS '調味料名';
COMMENT ON COLUMN recipe_seasonings.display_order IS '調味料の表示順序';
COMMENT ON COLUMN recipe_seasonings.amount_text IS '分量・個数などのテキスト';
COMMENT ON COLUMN recipe_seasonings.deleted_at IS '論理削除日時（NULLの場合は削除されていない）';

COMMENT ON TABLE recipe_steps IS 'レシピ手順テーブル';
COMMENT ON COLUMN recipe_steps.recipe_step_id IS 'レシピ手順の一意識別子';
COMMENT ON COLUMN recipe_steps.recipe_id IS 'レシピID（外部キー）';
COMMENT ON COLUMN recipe_steps.instruction IS '手順の説明';
COMMENT ON COLUMN recipe_steps.step_number IS '手順の番号';
COMMENT ON COLUMN recipe_steps.deleted_at IS '論理削除日時（NULLの場合は削除されていない）';

COMMENT ON TABLE saved_recipes IS 'ユーザーが保存したレシピテーブル';
COMMENT ON COLUMN saved_recipes.saved_recipe_id IS '保存レシピの一意識別子';
COMMENT ON COLUMN saved_recipes.recipe_id IS 'レシピID（外部キー）';
COMMENT ON COLUMN saved_recipes.user_id IS 'ユーザーID（外部キー）';
COMMENT ON COLUMN saved_recipes.saved_at IS '保存日時';
COMMENT ON COLUMN saved_recipes.deleted_at IS '論理削除日時（NULLの場合は削除されていない）'; 