package fingerprint

import (
	"fmt"
	"path/filepath"

	"github.com/debricked/cli/internal/file"
	"github.com/debricked/cli/internal/fingerprint"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var exclusions = file.DefaultExclusionsFingerprint()
var shouldFingerprintCompressedContent bool
var outputDir string
var minFingerprintContentLength int

const (
	ExclusionFlag                   = "exclusion"
	FingerprintCompressedContent    = "fingerprint-compressed-content"
	OutputDirFlag                   = "output-dir"
	MinFingerprintContentLengthFlag = "min-fingerprint-content-length"
)

func NewFingerprintCmd(fingerprinter fingerprint.IFingerprint) *cobra.Command {

	short := "Fingerprints files to match against the Debricked knowledge base."
	long := fmt.Sprintf("Fingerprint files for identification in a given path and writes it to %s.\nThis hashes all files to be used for matching against the Debricked knowledge base.", fingerprint.OutputFileNameFingerprints)
	cmd := &cobra.Command{
		Use:   "fingerprint [path]",
		Short: short,
		Long:  long,
		PreRun: func(cmd *cobra.Command, _ []string) {
			_ = viper.BindPFlags(cmd.Flags())
		},
		RunE: RunE(fingerprinter),
	}
	fileExclusionExample := filepath.Join("*", "**.pyc")
	dirExclusionExample := filepath.Join("**", "node_modules", "**")
	exampleFlags := fmt.Sprintf("-%s \"%s\" -%s \"%s\"", ExclusionFlag, fileExclusionExample, ExclusionFlag, dirExclusionExample)
	cmd.Flags().StringArrayVarP(&exclusions, ExclusionFlag, "", exclusions, `The following terms are supported to exclude paths:
Special Terms | Meaning
------------- | -------
"*"           | matches any sequence of non-Separator characters 
"/**/"        | matches zero or multiple directories
"?"           | matches any single non-Separator character
"[class]"     | matches any single non-Separator character against a class of characters ([see "character classes"])
"{alt1,...}"  | matches a sequence of characters if one of the comma-separated alternatives matches

Example: 
$ debricked files fingerprint . `+exampleFlags)
	cmd.Flags().BoolVar(&shouldFingerprintCompressedContent, FingerprintCompressedContent, false, `Fingerprint the contents of compressed files by unpacking them in memory, Supported files: `+fmt.Sprintf("%v", fingerprint.ZIP_FILE_ENDINGS))
	cmd.Flags().StringVar(&outputDir, OutputDirFlag, ".", "The directory to write the output file to")
	cmd.Flags().IntVar(&minFingerprintContentLength, MinFingerprintContentLengthFlag, 45, "Set minimum content length (in bytes) for files to fingerprint. Defaults to 45 bytes.")
	viper.MustBindEnv(ExclusionFlag)

	return cmd
}

func RunE(f fingerprint.IFingerprint) func(_ *cobra.Command, args []string) error {
	return func(_ *cobra.Command, args []string) error {
		path := ""
		if len(args) > 0 {
			path = args[0]
		}

		output, err := f.FingerprintFiles(path, exclusions, shouldFingerprintCompressedContent, minFingerprintContentLength)

		if err != nil {
			return err
		}
		outputFilePath := filepath.Join(outputDir, fingerprint.OutputFileNameFingerprints)
		err = output.ToFile(outputFilePath)
		if err != nil {
			return err
		}

		return nil
	}
}
