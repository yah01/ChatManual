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
		EnableBashCompletion:   true,
		UseShortOptionHandling: true,
		Suggest:                true,

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
			&cli.StringFlag{
				Name:    "tech",
				Aliases: []string{"t"},
				Usage:   "the tech you are talking about, e.g. C++, Rust, system calls...",
			},
		},

		Action: func(ctx *cli.Context) error {
			query := strings.Join(ctx.Args().Slice(), " ")
			hint := "manual"
			if ctx.Bool("short") {
				hint = "short summary"
			}

			messages := []openai.ChatCompletionMessage{
				{Role: "system", Content: "You are a helpful programmer assistant."},
				{Role: "user", Content: fmt.Sprintf("%s: %s", hint, query)},
				{Role: "user", Content: fmt.Sprintf("answer in %s language", ctx.String("language"))},
			}
			if ctx.Bool("example") {
				messages = append(messages, openai.ChatCompletionMessage{Role: "user", Content: "answer with a code example"})
			}
			if ctx.Bool("detail") {
				messages = append(messages, openai.ChatCompletionMessage{Role: "user", Content: "answer with function signature and description of parameters"})
			}
			format := ctx.String("format")
			tech := ctx.String("tech")
			if len(format) > 0 {
				messages = append(messages, openai.ChatCompletionMessage{Role: "user", Content: fmt.Sprintf("output as %s format", format)})
			}
			if len(tech) > 0 {
				messages = append(messages, openai.ChatCompletionMessage{Role: "user", Content: fmt.Sprintf("talk about %s", tech)})
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
