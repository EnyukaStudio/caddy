package status

import (
	"strconv"

	"github.com/mholt/caddy"
	"github.com/mholt/caddy/caddyhttp/httpserver"
)

// init registers Status plugin
func init() {
	caddy.RegisterPlugin("status", caddy.Plugin{
		ServerType: "http",
		Action:     setup,
	})
}

// setup configures new Status middleware instance.
func setup(c *caddy.Controller) error {
	rules, err := statusParse(c)
	if err != nil {
		return err
	}

	cfg := httpserver.GetConfig(c)
	mid := func(next httpserver.Handler) httpserver.Handler {
		return Status{Rules: rules, Next: next}
	}
	cfg.AddMiddleware(mid)

	return nil
}

// statusParse parses status directive
func statusParse(c *caddy.Controller) (map[string]int, error) {
	res := make(map[string]int)

	for c.Next() {
		hadBlock := false
		args := c.RemainingArgs()

		switch len(args) {
		case 1:
			statusCode, err := strconv.Atoi(args[0])
			if err != nil {
				return res, c.Errf("Expecting a numeric status code, got '%s'", args[0])
			}

			for c.NextBlock() {
				hadBlock = true
				path := c.Val()

				if _, exists := res[path]; exists {
					return res, c.Errf("Duplicate path: '%s'", path)
				}
				res[path] = statusCode

				if c.NextArg() {
					return res, c.ArgErr()
				}
			}

			if !hadBlock {
				return res, c.ArgErr()
			}
		case 2:
			statusCode, err := strconv.Atoi(args[0])
			if err != nil {
				return res, c.Errf("Expecting a numeric status code, got '%s'", args[0])
			}

			path := args[1]
			if _, exists := res[path]; exists {
				return res, c.Errf("Duplicate path: '%s'", path)
			}
			res[path] = statusCode
		default:
			return res, c.ArgErr()
		}
	}

	return res, nil
}
