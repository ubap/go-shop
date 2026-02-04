-- Baskets represent individual shopping lists.
-- Each basket is identified by a unique key.
CREATE TABLE IF NOT EXISTS baskets
(
    -- The public identifier (e.g., 'weekly-grocery').
    -- We make this the PRIMARY KEY so we can link other tables to it.
    key TEXT PRIMARY KEY NOT NULL,

    -- The UTC timestamp of when the basket was first created.
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
) WITHOUT ROWID;

CREATE TABLE IF NOT EXISTS basket_items (
                                              id INTEGER PRIMARY KEY AUTOINCREMENT,
                                              basket_key TEXT NOT NULL,
    -- The name of the item (e.g., 'Milk').
    -- COLLATE NOCASE ensures that 'Milk' and 'milk' are treated as the same.
                                              title TEXT NOT NULL,

    -- The 'checked' state of the item. 0 = pending, 1 = completed.
                                              completed BOOLEAN DEFAULT 0,

    -- Constraints:
    -- 1. If the parent basket is deleted, remove all its items.
                                              FOREIGN KEY (basket_key) REFERENCES baskets(key) ON DELETE CASCADE,

    -- 2. Prevent the same item name from being added twice to the SAME basket.
                                              UNIQUE(basket_key, title COLLATE NOCASE)
);