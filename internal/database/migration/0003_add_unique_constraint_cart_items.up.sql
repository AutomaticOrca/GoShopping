ALTER TABLE cart_items ADD CONSTRAINT unique_cart_product UNIQUE(cart_id, product_id);
