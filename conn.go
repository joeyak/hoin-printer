package hoin

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"syscall"
)

type healingConn struct {
	addr string
	conn net.Conn
}

func newHealingConn(addr string) (io.ReadWriter, error) {
	conn := healingConn{addr: addr}
	err := conn.dial()
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func (c *healingConn) dial() error {
	if c.conn != nil {
		c.conn.Close()
		c.conn = nil
	}

	conn, err := net.Dial("tcp", c.addr)
	if err != nil {
		return fmt.Errorf("could not dial address: %w", err)
	}

	c.conn = conn

	return nil
}

func (c healingConn) Read(p []byte) (int, error) {
	n, err := c.conn.Read(p)
	if errors.Is(err, syscall.EPIPE) || errors.Is(err, context.DeadlineExceeded) {
		errRedial := c.dial()
		if errRedial != nil {
			return 0, errors.Join(fmt.Errorf("could not redial: %w", errRedial), err)
		}

		n, err = c.conn.Read(p)
	}

	return n, err
}

func (c healingConn) Write(p []byte) (int, error) {
	n, err := c.conn.Write(p)
	if errors.Is(err, syscall.EPIPE) || errors.Is(err, context.DeadlineExceeded) {
		errRedial := c.dial()
		if errRedial != nil {
			return 0, errors.Join(fmt.Errorf("could not redial: %w", errRedial), err)
		}

		n, err = c.conn.Write(p)
	}

	return n, err
}
