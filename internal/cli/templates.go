package cli

import (
	"fmt"
	"strings"

	"github.com/Hermithic/aiask/internal/config"
	"github.com/Hermithic/aiask/internal/llm"
	"github.com/Hermithic/aiask/internal/shell"
	"github.com/Hermithic/aiask/internal/templates"
	"github.com/Hermithic/aiask/internal/ui"
	"github.com/spf13/cobra"
)

var (
	templateDesc string
)

var templatesCmd = &cobra.Command{
	Use:     "templates",
	Aliases: []string{"template", "tmpl"},
	Short:   "Manage saved prompt templates",
	Long: `Manage saved prompt templates for frequently used requests.

Examples:
  aiask templates                    # List all templates
  aiask templates list               # List all templates
  aiask save git-log "show recent commits"  # Save a template
  aiask run git-log                  # Run a saved template
  aiask templates remove git-log     # Remove a template`,
	Run: runTemplatesList,
}

var templatesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all saved templates",
	Run:   runTemplatesList,
}

var templatesSaveCmd = &cobra.Command{
	Use:   "save <name> <prompt>",
	Short: "Save a new template",
	Long: `Save a prompt as a named template for later use.

Examples:
  aiask save git-log "show recent commits"
  aiask save find-large "find files larger than 100MB" -d "Find large files"`,
	Args: cobra.MinimumNArgs(2),
	Run:  runTemplatesSave,
}

var templatesRunCmd = &cobra.Command{
	Use:     "run <name>",
	Aliases: []string{"use", "exec"},
	Short:   "Run a saved template",
	Args:    cobra.ExactArgs(1),
	Run:     runTemplatesRun,
}

var templatesRemoveCmd = &cobra.Command{
	Use:     "remove <name>",
	Aliases: []string{"rm", "delete"},
	Short:   "Remove a saved template",
	Args:    cobra.ExactArgs(1),
	Run:     runTemplatesRemove,
}

func init() {
	templatesSaveCmd.Flags().StringVarP(&templateDesc, "description", "d", "", "Description for the template")

	templatesCmd.AddCommand(templatesListCmd)
	templatesCmd.AddCommand(templatesRemoveCmd)

	// Add save and run as both subcommands of templates and as root commands
	rootCmd.AddCommand(templatesSaveCmd)
	rootCmd.AddCommand(templatesRunCmd)
}

func runTemplatesList(cmd *cobra.Command, args []string) {
	t, err := templates.Load()
	if err != nil {
		ui.ShowError(fmt.Errorf("failed to load templates: %w", err))
		return
	}

	if len(t.Items) == 0 {
		fmt.Println("No templates saved yet.")
		fmt.Println()
		fmt.Println("Save a template with:")
		fmt.Printf("  %saiask save <name> \"<prompt>\"%s\n", ui.ColorCyan, ui.ColorReset)
		return
	}

	fmt.Printf("%s╔══════════════════════════════════════════╗%s\n", ui.ColorCyan, ui.ColorReset)
	fmt.Printf("%s║           Saved Templates                ║%s\n", ui.ColorCyan, ui.ColorReset)
	fmt.Printf("%s╚══════════════════════════════════════════╝%s\n", ui.ColorCyan, ui.ColorReset)
	fmt.Println()

	for _, tmpl := range t.List() {
		fmt.Printf("%s%s%s", ui.ColorBold, tmpl.Name, ui.ColorReset)
		if tmpl.UsageCount > 0 {
			fmt.Printf(" %s(used %d times)%s", ui.ColorDim, tmpl.UsageCount, ui.ColorReset)
		}
		fmt.Println()

		if tmpl.Description != "" {
			fmt.Printf("  %s%s%s\n", ui.ColorDim, tmpl.Description, ui.ColorReset)
		}
		fmt.Printf("  %sPrompt:%s %s\n", ui.ColorDim, ui.ColorReset, tmpl.Prompt)
		fmt.Println()
	}

	fmt.Printf("%sRun a template with: aiask run <name>%s\n", ui.ColorDim, ui.ColorReset)
}

func runTemplatesSave(cmd *cobra.Command, args []string) {
	name := args[0]
	prompt := strings.Join(args[1:], " ")

	t, err := templates.Load()
	if err != nil {
		ui.ShowError(fmt.Errorf("failed to load templates: %w", err))
		return
	}

	if err := t.Add(name, prompt, templateDesc); err != nil {
		ui.ShowError(err)
		return
	}

	if err := t.Save(); err != nil {
		ui.ShowError(fmt.Errorf("failed to save templates: %w", err))
		return
	}

	fmt.Printf("%s✓ Template '%s' saved successfully!%s\n", ui.ColorGreen, name, ui.ColorReset)
	fmt.Printf("  Run it with: %saiask run %s%s\n", ui.ColorCyan, name, ui.ColorReset)
}

func runTemplatesRun(cmd *cobra.Command, args []string) {
	name := args[0]

	t, err := templates.Load()
	if err != nil {
		ui.ShowError(fmt.Errorf("failed to load templates: %w", err))
		return
	}

	tmpl, err := t.Get(name)
	if err != nil {
		ui.ShowError(err)
		fmt.Println()
		fmt.Println("Available templates:")
		for _, item := range t.List() {
			fmt.Printf("  - %s\n", item.Name)
		}
		return
	}

	// Increment usage count
	t.IncrementUsage(name)
	_ = t.Save()

	// Load config and run the prompt
	cfg, err := config.Load()
	if err != nil {
		ui.ShowError(fmt.Errorf("configuration error: %w", err))
		return
	}

	shellInfo := shell.Detect()

	provider, err := llm.NewProvider(cfg)
	if err != nil {
		ui.ShowError(fmt.Errorf("failed to create LLM provider: %w", err))
		return
	}
	defer llm.CloseProvider(provider)

	fmt.Printf("%sRunning template '%s': %s%s\n", ui.ColorDim, name, tmpl.Prompt, ui.ColorReset)

	// Run the interaction loop with the template prompt
	runInteractionLoop(provider, tmpl.Prompt, shellInfo, cfg)
}

func runTemplatesRemove(cmd *cobra.Command, args []string) {
	name := args[0]

	t, err := templates.Load()
	if err != nil {
		ui.ShowError(fmt.Errorf("failed to load templates: %w", err))
		return
	}

	if err := t.Remove(name); err != nil {
		ui.ShowError(err)
		return
	}

	if err := t.Save(); err != nil {
		ui.ShowError(fmt.Errorf("failed to save templates: %w", err))
		return
	}

	fmt.Printf("%s✓ Template '%s' removed%s\n", ui.ColorGreen, name, ui.ColorReset)
}

