-- migrate:up
CREATE TABLE ip_addresses
(
    id SERIAL PRIMARY KEY,
    subnet_id INT REFERENCES subnets(id) NOT NULL,
    ip_address BIGINT NOT NULL,
    hostname VARCHAR(50) NOT NULL,
    UNIQUE (subnet_id, ip_address)
);

-- migrate:down
DROP TABLE ip_addresses;
