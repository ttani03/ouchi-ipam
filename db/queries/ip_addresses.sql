-- name: GetIPAddresses :many
SELECT * FROM ip_addresses
WHERE subnet_id = $1
ORDER BY ip_address ASC;

-- name: ReserveIPAddress :one
INSERT INTO ip_addresses (
    subnet_id, ip_address, hostname
) VALUES (
    $1, $2, $3
)
RETURNING *;

-- name: FreeIPAddress :exec
DELETE FROM ip_addresses
WHERE subnet_id = $1 AND ip_address = $2;
