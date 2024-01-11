package main

import (
	"context"
	"fmt"
	"os"

	"github.com/gofireflyio/aiac/v3/libaiac"
)

func main() {
	client := libaiac.NewClient(os.Getenv("OPENAI_API_KEY"))
	ctx := context.TODO()

	// use the model-agnostic wrapper
	res, err := client.GenerateCode(
		ctx,
		libaiac.ModelGPT35Turbo0301,
		"generate terraform for ec2",
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed generating code: %s\n", err)
		os.Exit(1)
	}

	fmt.Fprintln(os.Stdout, res.Code)

	// use the completion API (for completion-only models)
	res, err = client.Complete(
		ctx,
		libaiac.ModelGPT35Turbo0301,
		"generate terraform for ec2",
	)

	// converse via a chat model
	chat := client.Chat(libaiac.ModelGPT35Turbo)
	res, err = chat.Send(ctx, "generate terraform for eks")
	res, err = chat.Send(ctx, "region must be eu-central-1")
}
