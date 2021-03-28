package main

import (
	"embed"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/aarondl/oa3/elm"
	"github.com/aarondl/oa3/generator"
	"github.com/aarondl/oa3/goserver"
	"github.com/aarondl/oa3/openapi3spec"
	"github.com/aarondl/oa3/tsclient"
	"github.com/spf13/cobra"

	_ "embed"
)

//go:embed templates
var templates embed.FS

var (
	wd      string
	version string = "unknown"
)

var (
	flagParams      []string
	flagDebug       bool
	flagWipe        bool
	flagOutputDir   string
	flagTemplateDir string
)

var rootCmd = &cobra.Command{
	Use:   "oa3 [flags] <generator> <openapifile>",
	Short: "Generate a language library for an openapi3 spec file",

	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		openapi3spec.DebugOutput = flagDebug
		if len(flagOutputDir) == 0 {
			flagOutputDir = filepath.Join(wd, "out", args[0])
		}
	},

	RunE:          rootCmdRun,
	Args:          cobra.ExactArgs(2),
	SilenceErrors: true,
	SilenceUsage:  true,
}

func main() {
	for _, a := range os.Args {
		if a == "--version" {
			fmt.Println("oa3 version " + version)
			return
		}
	}

	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println("failed to determine current working directory")
		os.Exit(1)
	}
	wd = cwd

	rootCmd.PersistentFlags().StringSliceVarP(&flagParams, "param", "p", nil, "key=value params to the generator")
	rootCmd.PersistentFlags().BoolVarP(&flagDebug, "debug", "", false, "debug output")
	rootCmd.PersistentFlags().BoolVarP(&flagWipe, "wipe", "w", false, "rm output directory before generation")
	rootCmd.PersistentFlags().StringVarP(&flagOutputDir, "output", "o", "", "output directory (default: "+filepath.Join(wd, "out", "<generator>")+")")
	rootCmd.PersistentFlags().StringVarP(&flagTemplateDir, "templates", "t", "", "template directory (default: embedded templates)")

	// ignored, only for docs
	rootCmd.PersistentFlags().BoolP("version", "", false, "output version and exit")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func rootCmdRun(cmd *cobra.Command, args []string) error {
	params := make(map[string]string, len(flagParams))
	for i, p := range flagParams {
		splits := strings.SplitN(p, "=", 2)
		if len(splits) != 2 || len(splits[0]) == 0 || len(splits[1]) == 0 {
			return fmt.Errorf("--param[%d] invalid: must be key=value pair, got: %s", i, p)
		}

		params[splits[0]] = splits[1]
	}

	files, err := generate(args[0], args[1], params)
	if err != nil {
		return err
	}

	if flagWipe {
		_ = os.RemoveAll(flagOutputDir)
	}

	if err = os.MkdirAll(flagOutputDir, 0775); err != nil {
		return err
	}

	for _, f := range files {
		if err := ioutil.WriteFile(filepath.Join(flagOutputDir, f.Name), f.Contents, 0640); err != nil {
			return err
		}
	}

	return nil
}

func generate(generatorID string, openapiFile string, params map[string]string) ([]generator.File, error) {
	var gen generator.Interface
	switch generatorID {
	case "go":
		gen = goserver.New()
	case "elm":
		gen = elm.New()
	case "ts":
		gen = tsclient.New()
	default:
		return nil, fmt.Errorf("unknown generator: %s", generatorID)
	}

	oa, err := openapi3spec.LoadYAML(openapiFile, true)
	if err != nil {
		return nil, err
	}

	var templateFS fs.FS = templates
	if len(flagTemplateDir) != 0 {
		templateFS = os.DirFS(flagTemplateDir)
	} else {
		templateFS, err = fs.Sub(templates, generatorID)
		if err != nil {
			return nil, fmt.Errorf("failed to root fs for generator: %w", err)
		}
	}

	if err := gen.Load(templateFS); err != nil {
		return nil, err
	}

	return gen.Do(oa, params)
}
