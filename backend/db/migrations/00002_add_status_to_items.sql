-- +goose Up
ALTER TABLE basket_items
    ADD COLUMN status TEXT NOT NULL DEFAULT 'active'
        CHECK (status IN ('active', 'deleted'));

-- +goose Down
ALTER TABLE basket_items DROP COLUMN status;