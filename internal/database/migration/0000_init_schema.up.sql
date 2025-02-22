-- users
CREATE TABLE users (
    id SERIAL PRIMARY KEY ,
    email VARCHAR(255) UNIQUE NOT NULL ,
    password_hash TEXT NOT NULL,
    full_name VARCHAR(100),
    avatar_url TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- roles
CREATE TABLE roles (
                       id SERIAL PRIMARY KEY,
                       name VARCHAR(50) UNIQUE NOT NULL,
                       description TEXT
);

-- user-roles
CREATE TABLE user_roles (
                            user_id INT REFERENCES users(id) ON DELETE CASCADE,
                            role_id INT REFERENCES roles(id) ON DELETE CASCADE,
                            PRIMARY KEY (user_id, role_id)
);

-- sessions
CREATE TABLE sessions
(
    id            UUID PRIMARY KEY,
    user_id       INT         NOT NULL,
    refresh_token VARCHAR     NOT NULL,
    user_agent    VARCHAR     NOT NULL,
    client_ip     VARCHAR     NOT NULL,
    is_blocked    BOOLEAN     NOT NULL DEFAULT false,
    expires_at    TIMESTAMPTZ NOT NULL,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT (now()),
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

-- categories
CREATE TABLE categories (
                            id SERIAL PRIMARY KEY,
                            name VARCHAR(100) UNIQUE NOT NULL,
                            parent_id INT REFERENCES categories(id) ON DELETE SET NULL -- 支持二级分类
);

-- products
CREATE TABLE products (
                          id SERIAL PRIMARY KEY,
                          name VARCHAR(255) NOT NULL,
                          description TEXT,
                          price DECIMAL(10,2) NOT NULL,
                          stock INT NOT NULL CHECK (stock >= 0),
                          image_url TEXT,
                          category_id INT REFERENCES categories(id) ON DELETE SET NULL,
                          created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                          updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- inventory
CREATE TABLE inventory (
                           product_id INT PRIMARY KEY REFERENCES products(id) ON DELETE CASCADE,
                           stock INT NOT NULL CHECK (stock >= 0),
                           last_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- carts
CREATE TABLE carts (
                       id SERIAL PRIMARY KEY,
                       user_id INT REFERENCES users(id) ON DELETE CASCADE,
                       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- cart_times
CREATE TABLE cart_items (
                            id SERIAL PRIMARY KEY,
                            cart_id INT REFERENCES carts(id) ON DELETE CASCADE,
                            product_id INT REFERENCES products(id) ON DELETE CASCADE,
                            quantity INT NOT NULL CHECK (quantity > 0),
                            added_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- orders
CREATE TABLE orders (
                        id SERIAL PRIMARY KEY,
                        user_id INT REFERENCES users(id) ON DELETE CASCADE,
                        total_price DECIMAL(10,2) NOT NULL CHECK (total_price >= 0),
                        status VARCHAR(50) NOT NULL CHECK (status IN ('pending', 'paid', 'shipped', 'completed', 'cancelled')),
                        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- order_items
CREATE TABLE order_items (
                             id SERIAL PRIMARY KEY,
                             order_id INT REFERENCES orders(id) ON DELETE CASCADE,
                             product_id INT REFERENCES products(id) ON DELETE CASCADE,
                             quantity INT NOT NULL CHECK (quantity > 0),
                             price DECIMAL(10,2) NOT NULL CHECK (price >= 0),
                             subtotal DECIMAL(10,2) NOT NULL CHECK (subtotal >= 0)
);

-- payments
CREATE TABLE payments (
                          id SERIAL PRIMARY KEY,
                          order_id INT REFERENCES orders(id) ON DELETE CASCADE,
                          user_id INT REFERENCES users(id) ON DELETE CASCADE,
                          amount DECIMAL(10,2) NOT NULL CHECK (amount >= 0),
                          status VARCHAR(50) NOT NULL CHECK (status IN ('pending', 'success', 'failed', 'cancelled')),
                          transaction_id TEXT UNIQUE,
                          created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- index
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_sessions_user_id ON sessions (user_id);
CREATE INDEX idx_orders_user_id ON orders(user_id);
CREATE INDEX idx_payments_order_id ON payments(order_id);
CREATE INDEX idx_cart_items_cart_id ON cart_items(cart_id);
