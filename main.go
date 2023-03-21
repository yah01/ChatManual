package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/sashabaranov/go-openai"
	"github.com/urfave/cli/v2"
)

func main() {
	apiKey := os.Getenv("OPENAI_API_KEY")
	client := openai.NewClient(apiKey)

	app := &cli.App{
		Name:                   "cman",
		Usage:                  "man with ChatGPT",
		UseShortOptionHandling: true,

		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "short",
				Aliases: []string{"s"},
				Usage:   "show a short summary",
			},
			&cli.BoolFlag{
				Name:    "example",
				Aliases: []string{"e"},
				Usage:   "with example",
			},
			&cli.BoolFlag{
				Name:    "detail",
				Aliases: []string{"d"},
				Usage:   "with detail (signature and parameters' details)",
			},
			&cli.StringFlag{
				Name:    "format",
				Aliases: []string{"f"},
				Usage:   "output file format",
			},
			&cli.StringFlag{
				Name:    "language",
				Aliases: []string{"l"},
				Value:   "english",
				Usage:   "answer language",
			},
		},

		Action: func(ctx *cli.Context) error {
			query := strings.Join(ctx.Args().Slice(), " ")
			hint := "man"
			if ctx.Bool("short") {
				hint = "tl;dr"
			}
			if ctx.Bool("detail") {
				hint += " with code of function signature and description of parameters"
			}
			if ctx.Bool("example") {
				hint += " with a code example"
			}
			format := ctx.String("format")

			messages := []openai.ChatCompletionMessage{
				{Role: "system", Content: "You are a helpful programmer assistant."},
				{Role: "user", Content: fmt.Sprintf("%s: %s", hint, query)},
				// {Role: "user", Content: "with the function/class definition"},
				{Role: "user", Content: fmt.Sprintf("answer in %s language", ctx.String("language"))},
			}
			if len(format) > 0 {
				messages = append(messages, openai.ChatCompletionMessage{Role: "user", Content: fmt.Sprintf("output as %s format", format)})
			}

			resp, err := client.CreateChatCompletion(context.Background(), openai.ChatCompletionRequest{
				Model:       "gpt-3.5-turbo",
				Messages:    messages,
				Temperature: 0.3,
			})
			if err != nil {
				return err
			}

			fmt.Println(resp.Choices[0].Message.Content)
			return nil
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
