package main

import (
	"context"
	"math"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/ttani03/ouchi-ipam/gen/sqlc"
)

type ipAddressResponse struct {
	ID        int32  `json:"id" example:"1"`                // require
	SubnetID  int32  `json:"subnet_id" example:"1"`         // require
	IpAddress string `json:"address" example:"192.168.0.2"` // require
	HostName  string `json:"hostname" example:"web01"`      // require
}

func ipAddrToIpAddrResponse(i sqlc.IpAddress) ipAddressResponse {
	return ipAddressResponse{
		ID:        i.ID,
		SubnetID:  i.SubnetID,
		IpAddress: ipv4IntToString(i.IpAddress),
		HostName:  i.Hostname,
	}
}

// getIpAddresses godoc
//
//	@Summary		Get IP Addresses
//	@Description	Get IP Addresses by subnet ID
//	@Tags			ip_addresses
//	@Produce		json
//	@Param			id	path		int	true	"Subnet ID"
//	@Success		200	{object}	[]ipAddressResponse
//	@Router			/subnets/{id}/ip [get]
func (s *Server) getIpAddresses(c echo.Context) error {
	subnetId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	addresses, err := s.db.GetIPAddresses(context.Background(), int32(subnetId))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	res := make([]ipAddressResponse, len(addresses))
	for i, v := range addresses {
		res[i] = ipAddrToIpAddrResponse(v)
	}

	return c.JSON(http.StatusOK, res)
}

// reserveIpAddress godoc
//
//	@Summary		Reserve an IP Address
//	@Description	Reserve an IP Address from the specified subnet
//	@Tags			ip_addresses
//	@Produce		json
//	@Param			id	path		int	true	"Subnet ID"
//	@Success		200	{object}	ipAddressResponse
//	@Router			/subnets/{id}/ip [post]
func (s *Server) reserveIpAddress(c echo.Context) error {
	subnetId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	subnet, err := s.db.GetSubnet(context.Background(), int32(subnetId))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	reserved, err := s.db.GetIPAddresses(context.Background(), int32(subnetId))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	if len(reserved) == int(math.Pow(2, float64(32-subnet.MaskLength))-2) {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to reserve a new IP because all IP addresses are already reserved.")
	}

	// Find the first IP address that is not reserved
	address := subnet.NetworkAddress + 1
	for _, v := range reserved {
		if v.IpAddress-address > 0 {
			break
		}
		address++
	}

	ip, err := s.db.ReserveIPAddress(context.Background(), sqlc.ReserveIPAddressParams{SubnetID: int32(subnetId), IpAddress: address, Hostname: c.QueryParam("hostname")})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, ipAddrToIpAddrResponse(ip))
}

// reserveSpecifiedIpAddress godoc
//
//	@Summary		Reserve a specified IP Address
//	@Description	Reserve a specified IP Address from the specified subnet
//	@Tags			ip_addresses
//	@Produce		json
//	@Param			id			path		int		true	"Subnet ID"
//	@Param			ip			path		string	true	"IP Address"
//	@Param			hostname	query		string	false	"Hostname"
//	@Success		200			{object}	ipAddressResponse
//	@Router			/subnets/{id}/ip/{ip} [post]
func (s *Server) reserveSpecifiedIpAddress(c echo.Context) error {
	subnetId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	subnet, err := s.db.GetSubnet(context.Background(), int32(subnetId))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	ip, err := stringToIPv4Int(c.Param("ip"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	if ip <= subnet.NetworkAddress || ip >= subnet.NetworkAddress+int64(math.Pow(2, float64(32-subnet.MaskLength))-1) {
		return echo.NewHTTPError(http.StatusBadRequest, "Specified IP address is not in the subnet.")
	}

	reserved, err := s.db.ReserveIPAddress(context.Background(), sqlc.ReserveIPAddressParams{SubnetID: int32(subnetId), IpAddress: ip, Hostname: c.QueryParam("hostname")})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, ipAddrToIpAddrResponse(reserved))
}

// freeIpAddress godoc
//
//	@Summary		Free a specified IP Address
//	@Description	Free a specified IP Address from the specified subnet
//	@Tags			ip_addresses
//	@Pram			id										path			int	true	"Subnet ID"
//	@Param			ip	path	string	true	"IP Address"
//	@Success		204
//	@Router			/subnets/{id}/ip/{ip} [delete]
func (s *Server) freeIpAddress(c echo.Context) error {
	subnetId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	ip, err := stringToIPv4Int(c.Param("ip"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	err = s.db.FreeIPAddress(context.Background(), sqlc.FreeIPAddressParams{SubnetID: int32(subnetId), IpAddress: ip})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.NoContent(http.StatusNoContent)
}
