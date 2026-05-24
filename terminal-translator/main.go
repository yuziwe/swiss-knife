package main

import (
    "flag"
	"fmt"
	tsize "github.com/kopoli/go-terminal-size"
	cache "github.com/yuziwe/swiss-knife/terminal-translator/cache"
    provider "github.com/yuziwe/swiss-knife/terminal-translator/provider"
    "os"
    "strings"
)

const (
	ColorRed    = "\x1b[0;31m"
	ColorGreen  = "\x1b[0;32m"
	ColorYellow = "\x1b[0;33m"
	ColorBlue   = "\x1b[0;34m"
	ColorReset  = "\x1b[0m"
)

const (
	DEFAULT_MODEL = "deepseek-v4-flash"
	API_URL_KEY   = "TERMINAL_TRANSLATOR_API_URL"
	API_KEY_KEY   = "TERMINAL_TRANSLATOR_API_KEY"
	SYSTEM_PROMPT = `You are a professional translation assistant that **exclusively performs precise text translation tasks** and strictly adheres to the following rules: 1. **Function Definition**  - Translate Chinese text input into English.  - Translate English text input into Chinese.  - Automatically detect the language of the input text.  2. **Output Rules**  - **Output only the translated text in the target language**, with **no** prefixes, suffixes, explanations, notes, punctuation clarifications, or formatting embellishments.  - Absolutely **do not** output lead-ins such as "Translation:", "Result:", or similar.  - Absolutely **do not** output any characters or line breaks that are not part of the translation itself.  3. **Examples**  - User Input: "你好，世界" → Your Output: "Hello, world"  - User Input: "How are you" → Your Output: "你好吗" Strictly follow these rules to ensure every response contains only the pure translation result.`
)

var (
    DebugMode bool
    Model string
)

type SystemCtx struct {
	Cache cache.CacheSchema
    Provider provider.LLMProvider
}

func ugly_separators() {
	ws, err := tsize.GetSize()
	if err != nil {
		fmt.Println("ERROR: get window size failed!: ", err)
		os.Exit(1)
	}

    fmt.Println(strings.Repeat("=", ws.Width))
}

func usage() {
    fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s <-d|-m> <English|Chinese>\n", os.Args[0])
    flag.PrintDefaults()
}

func prompt_banner() {
	ugly_separators()
	fmt.Println(ColorRed, "Current model: ", Model, ColorReset)
	ugly_separators()
}

func main() {
    flag.StringVar(&Model, "m", DEFAULT_MODEL, "Set backend model")
    flag.BoolVar(&DebugMode, "d", false, "Open debug mode")
    flag.Parse()

	if flag.NArg() < 1 {
        usage()
		os.Exit(1)
	}

    user_prompt := flag.Args()[0]

    prompt_banner()

	base_api_url := os.Getenv(API_URL_KEY)
	if base_api_url == "" {
		fmt.Println("ERROR: ", API_URL_KEY, " is empty!")
		os.Exit(1)
	}

	api_key := os.Getenv(API_KEY_KEY)
	if api_key == "" {
		fmt.Println("ERROR: ", API_KEY_KEY, " is empty!")
		os.Exit(1)
	}

	ctx := &SystemCtx{
		Cache: &cache.LocalCache{},
        Provider: &provider.OpenAI{
            BaseUrl: base_api_url,
            ApiKey: api_key,
        },
	}

	msgs := []provider.Message{
		{Role: "system", Content: SYSTEM_PROMPT},
		{Role: "user", Content: user_prompt},
	}

	// Init cache system
	if err := ctx.Cache.Init(DebugMode); err != nil {
		fmt.Println("ERROR: initialize cache system failed: ", err)
		os.Exit(1)
	}

	if ctx.Cache.Exist(user_prompt) {
		res, err := ctx.Cache.Rd(user_prompt)
		if err == nil {
			// Cache hit
			fmt.Println(ColorYellow, user_prompt, ColorReset)
			ugly_separators()
			fmt.Println(ColorGreen, res, ColorReset)
			ugly_separators()
			os.Exit(0)
		}
	}

	resp, err := ctx.Provider.Completions(Model, msgs)
	if err != nil {
		fmt.Println("ERROR: create completions failed!: ", err)
		os.Exit(1)
	}

	if len(resp.Choices) == 0 {
		fmt.Println("ERROR: got empty response!")
		os.Exit(1)
	}

	// Add to cache
	if err := ctx.Cache.Wt(user_prompt, resp.Choices[0].RMessage.Content); err != nil {
		fmt.Println("WARN: Write into cache failed, err: ", err)
	}

	// Output
	fmt.Println(ColorYellow, user_prompt, ColorReset)

	ugly_separators()

	fmt.Println(ColorGreen, resp.Choices[0].RMessage.Content, ColorReset)

	ugly_separators()
}
