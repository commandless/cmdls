package cmd

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	cobra.OnInitialize(initConfig)
}

// A single command
func Execute() {
	root := newRootCmd(nil)
	err := root.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

type rootCmd struct {
	client client

	apiURL string
	debug  bool
}

func newRootCmd(c client) *cobra.Command {
	inst := rootCmd{client: c}

	cmd := &cobra.Command{
		Use:   "rigs",
		Short: "short",
		Long:  `long`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return inst.run()
		},
	}

	f := cmd.Flags()
	f.BoolVarP(&inst.debug, "debug", "d", false, "enable debug logging")
	f.StringVarP(&inst.apiURL, "url", "u", "http://host.docker.internal:8000", "url of api")
	return cmd
}

func (c *rootCmd) run() error {
	c.client = ensureClient(c.client, c.apiURL)

	res := Response{}
	err := c.client.Get("/api.json", &res)
	if err != nil {
		return err
	}

	// choose recipe
	for i, result := range res.Recipes {
		idx := i + 1
		fmt.Printf("[%d] %s\n", idx, result.Description)
	}
	fmt.Printf("\nSelect your recipe:\n> ")
	input := bufio.NewScanner(os.Stdin)
	input.Scan()
	choice := input.Text()
	index, err := strconv.ParseInt(choice, 10, 64)
	if err != nil {
		return fmt.Errorf("%s is not a valid choice ", choice)
	}
	index = index - 1
	if index < 0 || index > int64(len(res.Recipes)) {
		return fmt.Errorf("'%d' is not a valid choice", index)
	}
	recipe := res.Recipes[index]

	// check if the binary is installed
	command := ""
	exists := c.checkInstalled("foobar")
	if !exists {
		fmt.Printf("'%s' does not exist on your machine. Do you want to install it? [Y/N]\n> ", res.Resolution.Bin)
		input := bufio.NewScanner(os.Stdin)
		input.Scan()
		choice = input.Text()
		if choice == "Y" || choice == "y" {
			// Run command
			fmt.Printf("Installing '%s'\n", res.Resolution.Bin)
		} else {
			command = "npx"
		}
	}

	// get values for all required arguments
	commandArgs := []string{}
	fmt.Println("Provide values for required arguments")
	for _, arg := range recipe.Inputs.Args {
		fmt.Printf("(%s) %s\n", arg.Type, arg.Description)
		fmt.Printf("> ")
		input := bufio.NewScanner(os.Stdin)
		input.Scan()
		value := input.Text()
		commandArgs = append(commandArgs, value)
	}

	// get values for all flags
	fmt.Println("Provide values for flags")
	for _, flag := range recipe.Inputs.Flags {
		fmt.Printf("(%s) --%s\n", flag.Type, flag.Long)
		fmt.Printf("> ")
		input := bufio.NewScanner(os.Stdin)
		input.Scan()
		value := input.Text()
		if value == "" {
			if flag.Default == "" {
				continue
			} else {
				value = flag.Default
			}
		}
		commandArgs = append(commandArgs, flag.Long)
		commandArgs = append(commandArgs, "--"+value)
	}

	if command == "" {
		return execute(res.Resolution.Bin, commandArgs...)
	} else {
		args := append([]string{res.Resolution.Bin}, commandArgs...)
		return execute(command, args...)
	}

	return nil
}

func (c rootCmd) checkInstalled(bin string) bool {
	cmd := exec.Command("which", bin)
	cmd.Stdin = strings.NewReader("some input")
	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err == nil {
		return true
	}
	return false
}
