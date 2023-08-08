package main

import (
	"net"
	"strconv"
)

// DefaultPort appends given port to connection if not specified
func DefaultPort(conn string, port int) string {
	if _, _, err := net.SplitHostPort(conn); err != nil {
		conn = net.JoinHostPort(conn, strconv.Itoa(port))
	}

	return conn
}