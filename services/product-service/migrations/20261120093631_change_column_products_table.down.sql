ALTER TABLE products
    -- Restore NOT NULL
    ALTER COLUMN category_slug SET NOT NULL,
    ALTER COLUMN name SET NOT NULL,
    ALTER COLUMN image SET NOT NULL,

    -- Restore DEFAULT
    ALTER COLUMN regular_price SET DEFAULT 0,
    ALTER COLUMN sale_price SET DEFAULT 0,
    ALTER COLUMN unit SET DEFAULT 'gram',
    ALTER COLUMN weight SET DEFAULT 0,
    ALTER COLUMN variant SET DEFAULT 1,
    ALTER COLUMN status SET DEFAULT 'DRAFT';