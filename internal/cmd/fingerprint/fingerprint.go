package fingerprint

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/debricked/cli/internal/fingerprint"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var exclusions = fingerprint.DefaultExclusionsFingerprint()
var inclusions []string
var shouldFingerprintCompressedContent bool
var outputDir string
var minFingerprintContentLength int
var shouldRegenerateFingerprintFile bool

const (
	ExclusionFlag                   = "exclusion"
	InclusionFlag                   = "inclusion"
	FingerprintCompressedContent    = "fingerprint-compressed-content"
	OutputDirFlag                   = "output-dir"
	MinFingerprintContentLengthFlag = "min-fingerprint-content-length"
	RegenerateFingerprintFile       = "regenerate"
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
	fileExclusionExample := filepath.Join("'*", "**.pyc'")
	dirExclusionExample := filepath.Join("'**", "node_modules", "**'")
	exampleFlags := fmt.Sprintf("-%s \"%s\" -%s \"%s\"", ExclusionFlag, fileExclusionExample, ExclusionFlag, dirExclusionExample)
	cmd.Flags().StringArrayVarP(&exclusions, ExclusionFlag, "e", exclusions, `The following terms are supported to exclude paths:
Special Terms | Meaning
------------- | -------
"*"           | matches any sequence of non-Separator characters 
"/**/"        | matches zero or multiple directories
"?"           | matches any single non-Separator character
"[class]"     | matches any single non-Separator character against a class of characters ([see "character classes"])
"{alt1,...}"  | matches a sequence of characters if one of the comma-separated alternatives matches

Example: 
$ debricked files fingerprint . `+exampleFlags)
	cmd.Flags().StringArrayVar(
		&inclusions,
		InclusionFlag,
		[]string{},
		`Forces inclusion of specified terms, see exclusion flag for more information on supported terms.
Examples: 
$ debricked scan . --include '**/node_modules/**'`)
	cmd.Flags().BoolVar(&shouldFingerprintCompressedContent, FingerprintCompressedContent, false, `Fingerprint the contents of compressed files by unpacking them in memory, Supported files: `+fmt.Sprintf("%v", fingerprint.ZIP_FILE_ENDINGS))
	cmd.Flags().StringVar(&outputDir, OutputDirFlag, ".", "The directory to write the output file to")
	cmd.Flags().IntVar(&minFingerprintContentLength, MinFingerprintContentLengthFlag, 45, "Set minimum content length (in bytes) for files to fingerprint. Defaults to 45 bytes.")
	cmd.Flags().BoolVar(&shouldRegenerateFingerprintFile, RegenerateFingerprintFile, true, `If generated fingerprint file should be overwritten on subequent scans. Defaults to true`)

	viper.MustBindEnv(ExclusionFlag)

	return cmd
}

func RunE(f fingerprint.IFingerprint) func(_ *cobra.Command, args []string) error {
	return func(_ *cobra.Command, args []string) error {
		path := ""
		if len(args) > 0 {
			path = args[0]
		}
		options := fingerprint.DebrickedOptions{
			Path:                         path,
			Exclusions:                   exclusions,
			Inclusions:                   inclusions,
			FingerprintCompressedContent: shouldFingerprintCompressedContent,
			MinFingerprintContentLength:  minFingerprintContentLength,
		}
		output, err := f.FingerprintFiles(options)

		if err != nil {
			return err
		}
		var outputFilePath string
		// TODO: decide if we should only create a date-suffix file if a file without date-suffix already exists.
		//       Depending on decision, check for existence of debricked.fingerprints.txt here.
		if shouldRegenerateFingerprintFile {
			outputFilePath = filepath.Join(outputDir, fingerprint.OutputFileNameFingerprints)
		} else {
			// TODO: decide on if date suffix is what we want or if an incrementing integer should be used.
			timestamp := time.Now().Format("2006-01-02_15-04-05") // reference date, do not change the date, only format
			ext := filepath.Ext(fingerprint.OutputFileNameFingerprints)
			ouputFileName := fmt.Sprintf("%s_%s%s", fingerprint.OutputFileNameFingerprints, timestamp, ext)
			outputFilePath = filepath.Join(outputDir, ouputFileName)
		}

		err = output.ToFile(outputFilePath)
		if err != nil {
			return err
		}

		return nil
	}
}
