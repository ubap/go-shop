-- Represents the global catalog of unique items.
-- This design prevents duplicate strings (e.g., "Milk" vs "Milk")
-- and allows for easy auto-suggestion across different baskets.
CREATE TABLE IF NOT EXISTS items
(
    -- Internal auto-incrementing ID for relational foreign keys
    id    INTEGER PRIMARY KEY AUTOINCREMENT,

    -- The display name of the item.
    -- UNIQUE ensures we don't store the same item name twice globally.
    title TEXT NOT NULL UNIQUE
);
-- Baskets represent individual shopping lists.
-- Each basket is identified by a unique key.
CREATE TABLE IF NOT EXISTS baskets
(
    id         INTEGER PRIMARY KEY AUTOINCREMENT,

    -- The public identifier for the basket (e.g., 'weekly-grocery' or 'party-list').
    -- Used in URLs: myapp.com/b/test-basket
    -- UNIQUE constraint prevents duplicate lists from being created.
    key        TEXT NOT NULL UNIQUE,

    -- The UTC timestamp of when the basket was first created.
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- basket_items is a junction table representing the many-to-many relationship
-- between baskets and items. It tracks which items are in which list and
-- whether they have been "checked off."
CREATE TABLE IF NOT EXISTS basket_items
(
    -- Reference to the parent basket.
    basket_id INTEGER NOT NULL,

    -- Reference to the master item catalog.
    item_id   INTEGER NOT NULL,

    -- Tracks the 'checked' state of an item within THIS specific basket.
    -- SQLite stores this as 0 (false) or 1 (true).
    completed BOOLEAN DEFAULT 0,

    -- The composite Primary Key ensures a specific item can only appear
    -- ONCE in a specific basket.
    PRIMARY KEY (basket_id, item_id),

    -- If a basket is deleted, automatically remove all its item links.
    FOREIGN KEY (basket_id) REFERENCES baskets (id) ON DELETE CASCADE,

    -- If a master item is deleted, automatically remove it from all baskets.
    FOREIGN KEY (item_id) REFERENCES items (id) ON DELETE CASCADE
);