UPDATE categories
SET icon = ''
WHERE icon IS NULL;

ALTER TABLE categories
    ALTER COLUMN icon SET NOT NULL;