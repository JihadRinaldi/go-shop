CREATE UNIQUE INDEX IF NOT EXISTS uniq_active_cart_product
ON cart_items (cart_id, product_id)
WHERE deleted_at IS NULL;
