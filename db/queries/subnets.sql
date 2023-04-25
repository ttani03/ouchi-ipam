-- name: GetSubnets :many
SELECT * FROM subnets;

-- name: GetSubnet :one
SELECT * FROM subnets
WHERE id = $1 LIMIT 1;

-- name: CreateSubnet :one
INSERT INTO subnets (
    name, network_address, mask_length, gateway, name_server, description
) VALUES (
    $1, $2, $3, $4, $5, $6
)
RETURNING *;

-- name: DeleteSubnet :exec
DELETE FROM subnets
WHERE id = $1;
