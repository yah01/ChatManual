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
		Name:  "cman",
		Usage: "man with ChatGPT",
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
		},
		Action: func(ctx *cli.Context) error {
			query := strings.Join(ctx.Args().Slice(), " ")
			hint := "man"
			if ctx.Bool("short") {
				hint = "tl;dr"
			}
			if ctx.Bool("example") {
				hint += " with example"
			}
			resp, err := client.CreateChatCompletion(context.Background(), openai.ChatCompletionRequest{Model: "gpt-3.5-turbo", Messages: []openai.ChatCompletionMessage{
				{Role: "system", Content: "You are a helpful programmer assistant."},
				{Role: "user", Content: fmt.Sprintf("%s: %s", hint, query)},
			}})
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
