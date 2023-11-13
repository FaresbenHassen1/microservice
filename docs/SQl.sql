CREATE TABLE users (
    id_user uuid PRIMARY KEY default uuid_generate_v4(),
    name character(255)
);
CREATE TABLE wallet (
    id_wallet uuid PRIMARY key default uuid_generate_v4(),
    created_date timestamp default NOW(),
    balance float NOT NULL,
    currency character(3),
    users_id uuid,
    FOREIGN KEY (users_id) REFERENCES users (id_user)
);
CREATE TABLE transaction (
    id_transaction uuid PRIMARY key default uuid_generate_v4(),
    type character (255),
    amount float NOT NULL,
    date timestamp default NOW(),
    wallet_id uuid,
    FOREIGN KEY (wallet_id) REFERENCES wallet (id_wallet)
);
