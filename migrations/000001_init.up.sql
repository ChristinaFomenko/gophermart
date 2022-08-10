CREATE TABLE users
(
    id       SERIAL PRIMARY KEY,
    login    TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL
);

CREATE TABLE orders
(
    id        BIGSERIAL PRIMARY KEY,
    order_num BIGINT UNIQUE,
    user_id   INT NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users (id)
);

CREATE TABLE accruals
(
    order_num   BIGINT PRIMARY KEY,
    user_id     INT  NOT NULL,
    status      TEXT NOT NULL            DEFAULT 'NEW',
    amount      REAL                     DEFAULT 0,
    uploaded_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    FOREIGN KEY (user_id) REFERENCES users (id),
    FOREIGN KEY (order_num) REFERENCES orders (order_num)
);

CREATE TABLE withdrawals
(
    order_num    BIGINT PRIMARY KEY,
    user_id      INT NOT NULL,
    amount       REAL                     DEFAULT 0,
    processed_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    FOREIGN KEY (user_id) REFERENCES users (id),
    FOREIGN KEY (order_num) REFERENCES orders (order_num)
);