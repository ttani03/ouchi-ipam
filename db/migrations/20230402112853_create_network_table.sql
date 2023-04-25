-- migrate:up
CREATE TABLE subnets
(
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) UNIQUE NOT NULL,
    network_address BIGINT NOT NULL,
    mask_length INT NOT NULL,
    gateway BIGINT,
    name_server BIGINT,
    description TEXT
);

-- migrate:down
DROP TABLE subnets;
