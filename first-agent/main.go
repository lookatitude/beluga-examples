// Example: first-agent
//
// Demonstrates: building a single-LLM agent with a typed calculator tool
// that answers arithmetic questions using Beluga's agent + tool packages.
//
// Prerequisites:
//   - Go 1.23+
//   - OPENAI_API_KEY=your OpenAI API key (or swap provider — see comments)
//
// Run:
//
//	go run .
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/lookatitude/beluga-ai/agent"
	"github.com/lookatitude/beluga-ai/config"
	"github.com/lookatitude/beluga-ai/core"
	"github.com/lookatitude/beluga-ai/llm"
	"github.com/lookatitude/beluga-ai/tool"

	// Register the OpenAI provider via init(). Swap this import for any
	// other provider (anthropic, ollama, google, bedrock, etc.) and change
	// the llm.New call below to match.
	_ "github.com/lookatitude/beluga-ai/llm/providers/openai"
)

// CalculatorInput is the typed input schema for the calculator tool.
// Struct tags drive the generated JSON Schema that the LLM sees.
type CalculatorInput struct {
	A  float64 `json:"a" description:"Left operand" required:"true"`
	Op string  `json:"op" description:"Operator: one of + - * /" required:"true"`
	B  float64 `json:"b" description:"Right operand" required:"true"`
}

func newCalculatorTool() tool.Tool {
	return tool.NewFuncTool(
		"calculator",
		"Evaluate a single binary arithmetic operation (a op b).",
		func(ctx context.Context, in CalculatorInput) (*tool.Result, error) {
			var out float64
			switch strings.TrimSpace(in.Op) {
			case "+":
				out = in.A + in.B
			case "-":
				out = in.A - in.B
			case "*":
				out = in.A * in.B
			case "/":
				if in.B == 0 {
					return nil, core.Errorf(core.ErrInvalidInput, "calculator: divide by zero")
				}
				out = in.A / in.B
			default:
				return nil, core.Errorf(core.ErrInvalidInput, "calculator: unknown operator %q", in.Op)
			}
			return tool.TextResult(strconv.FormatFloat(out, 'f', -1, 64)), nil
		},
	)
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	model, err := llm.New("openai", config.ProviderConfig{
		Provider: "openai",
		APIKey:   os.Getenv("OPENAI_API_KEY"),
		Model:    "gpt-4o-mini",
	})
	if err != nil {
		log.Fatalf("llm.New: %v", err)
	}

	a := agent.New("math-tutor",
		agent.WithLLM(model),
		agent.WithPersona(agent.Persona{
			Role:      "patient math tutor",
			Goal:      "answer arithmetic questions clearly",
			Backstory: "You explain each step before giving the final answer.",
			Traits:    []string{"concise", "accurate"},
		}),
		agent.WithTools([]tool.Tool{newCalculatorTool()}),
	)

	answer, err := a.Invoke(ctx, "What is 17 times 42, minus 19?")
	if err != nil {
		log.Fatalf("agent.Invoke: %v", err)
	}

	fmt.Println(answer)
}
