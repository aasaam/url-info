package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/urfave/cli/v2"
)

func runServer(c *cli.Context) error {
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
		Prefork:               false,
		UnescapePath:          true,
		CaseSensitive:         true,
		StrictRouting:         true,
		BodyLimit:             c.Int("body-limit-size"),

		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError

			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}

			ctx.Set(fiber.HeaderContentType, fiber.MIMETextPlainCharsetUTF8)
			ctx.Status(code).SendString(err.Error())

			return nil
		},
	})

	app.Post("/info", func(c *fiber.Ctx) error {
		var processInfo urlInfoProcess
		if err := c.BodyParser(&processInfo); err != nil {
			return err
		}

		cacheID := processInfo.cache()
		cacheFile := cachePath(cacheID, "json")

		if fileExist(cacheFile) {
			return c.SendFile(cacheFile)
		}

		u, uE := newURL(processInfo.URL)

		if uE != nil {
			return uE
		}

		u.process(processInfo)

		b, bErr := json.Marshal(u)

		if bErr != nil {
			return bErr
		}

		writeErr := os.WriteFile(cacheFile, b, 0666)
		if writeErr != nil {
			return writeErr
		}

		c.Type("application/json", "utf8")
		return c.Send(b)
	})

	return app.Listen(c.String("listen"))
}

func main() {

	app := cli.NewApp()
	app.Usage = "URL Information"
	app.EnableBashCompletion = true
	app.Commands = []*cli.Command{
		{
			Name:   "run",
			Usage:  "Run server",
			Action: runServer,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "listen",
					Usage:    "Application listen http ip:port address",
					Value:    "0.0.0.0:4000",
					Required: false,
					EnvVars:  []string{"ASM_URL_INFO_LISTEN_ADDRESS"},
				},
				&cli.IntFlag{
					Name:     "body-limit-size",
					Usage:    "Request limit size",
					Value:    2 * 1024 * 1024,
					Required: false,
					EnvVars:  []string{"ASM_URL_INFO_BODY_LIMIT_SIZE"},
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
