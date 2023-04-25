package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/ttani03/ouchi-ipam/gen/sqlc"
)

type subnetRequest struct {
	Name        string  `json:"name" example:"web"`                   // require
	Cidr        string  `json:"cidr" example:"192.168.0.0/24"`        // require
	Gateway     *string `json:"gateway" example:"192.168.0.1"`        // optional
	NameServer  *string `json:"name_server" example:"8.8.8.8"`        // optional
	Description *string `json:"description" example:"subnet for web"` // optional
}

type subnetResponse struct {
	ID          int32   `json:"id" example:"1"`                       // require
	Name        string  `json:"name" example:"web"`                   // require
	Cidr        string  `json:"cidr" example:"192.168.0.0/24"`        // require
	Gateway     *string `json:"gateway" example:"192.168.0.1"`        // optional
	NameServer  *string `json:"name_server" example:"8.8.8.8"`        // optional
	Description *string `json:"description" example:"subnet for web"` // optional
}

func subnetToSubnetResponse(subnet sqlc.Subnet) subnetResponse {
	res := subnetResponse{
		ID:          subnet.ID,
		Name:        subnet.Name,
		Description: subnet.Description,
	}

	res.Cidr = fmt.Sprintf("%s/%d", ipv4IntToString(subnet.NetworkAddress), subnet.MaskLength)

	if subnet.Gateway != nil {
		gw := ipv4IntToString(*subnet.Gateway)
		res.Gateway = &gw
	}
	if subnet.NameServer != nil {
		ns := ipv4IntToString(*subnet.NameServer)
		res.NameServer = &ns
	}

	return res
}

func subnetRequestToCreateSubnetParams(req subnetRequest) (*sqlc.CreateSubnetParams, error) {
	params := new(sqlc.CreateSubnetParams)

	if req.Name == "" {
		return nil, echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid name: %s", req.Name))
	}
	params.Name = req.Name

	if !isCIDR(req.Cidr) {
		return nil, echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid cidr: %s", req.Cidr))
	}
	params.NetworkAddress, params.MaskLength, _ = parseCIDR(req.Cidr)

	if req.Gateway != nil {
		gw := net.ParseIP(*req.Gateway)
		if gw == nil {
			return nil, echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid gateway: %s", *req.Gateway))
		}
		i := parseIPv4ToInt(gw)
		params.Gateway = &i
	}

	if req.NameServer != nil {
		ns := net.ParseIP(*req.NameServer)
		if ns == nil {
			return nil, echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid name_server: %s", *req.NameServer))
		}
		i := parseIPv4ToInt(ns)
		params.NameServer = &i
	}

	params.Description = req.Description

	return params, nil
}

// getSubnets godoc
//
//	@Summary		Get subnets
//	@Description	get all subnets
//	@Tags			subnets
//	@Produce		json
//	@Success		200	{object}	[]subnetResponse
//	@Router			/subnets [get]
func (s *Server) getSubnets(c echo.Context) error {
	subnets, err := s.db.GetSubnets(context.Background())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	res := make([]subnetResponse, len(subnets))
	for i, v := range subnets {
		res[i] = subnetToSubnetResponse(v)
	}

	return c.JSON(http.StatusOK, res)
}

// getSubnet godoc
//
//	@Summary		Get a subnet
//	@Description	get a subnet by ID
//	@Tags			subnets
//	@Produce		json
//	@Param			id	path		int	true	"subnet ID"
//	@Success		200	{object}	subnetResponse
//	@Router			/subnets/{id} [get]
func (s *Server) getSubnet(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	subnet, err := s.db.GetSubnet(context.Background(), int32(id))
	if err != nil {
		if subnet.ID == 0 {
			return echo.NewHTTPError(http.StatusNotFound)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, subnetToSubnetResponse(subnet))
}

// createSubnet godoc
//
//	@Summary		Create a subnet
//	@Description	Create a new subnet
//	@Tags			subnets
//	@Accept			json
//	@Produce		json
//	@Param			subnet		body		subnetRequest	true	"subnet info"
//	@Success		200			{object}	subnetResponse
//	@Router			/subnets	[post]
func (s *Server) createSubnet(c echo.Context) error {
	req := new(subnetRequest)

	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	params, err := subnetRequestToCreateSubnetParams(*req)
	if err != nil {
		return err
	}

	subnet, err := s.db.CreateSubnet(context.Background(), *params)
	if err != nil {
		s.echo.Logger.Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusCreated, subnetToSubnetResponse(subnet))
}

// deleteSubnet godoc
//
//	@Summary		Delete a subnet
//	@Description	Delete a subnet by ID
//	@Tags			subnets
//	@Param			id	path	int	true	"subnet ID"
//	@Success		204
//	@Router			/subnets/{id}	[delete]
func (s *Server) deleteSubnet(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	err = s.db.DeleteSubnet(context.Background(), int32(id))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.NoContent(http.StatusNoContent)
}
