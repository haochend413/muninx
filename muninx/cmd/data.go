package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/haochend413/muninx/config"
	"github.com/haochend413/muninx/sys"
	"github.com/spf13/cobra"
)

var ExportNoteCmd = &cobra.Command{
	Use:   "export",
	Short: "export",
	Long:  "export",
	Run: func(cmd *cobra.Command, args []string) {
		if err := globalDB.ExportNoteToJSON(globalCfg.DataFilePath + "/notes.json"); err != nil {
			sys.LogError(err)
			os.Exit(1)
		}
		fmt.Println("Exported notes to", globalCfg.DataFilePath+"/notes.json")
	},
}

var DataBackupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Backup muninx data",
	Long:  "Backup the muninx data folder to a specified destination, default to cwd.",
	Run: func(cmd *cobra.Command, args []string) {
		base, err := config.BasePathDefault()
		if err != nil {
			sys.LogError(err)
			return
		}

		// timestamp folder name
		dest := fmt.Sprintf("muninx_backup_%s", time.Now().Format("2006-01-02_15-04-05"))
		if len(args) > 0 {
			dest = args[0]
		}

		// absolute
		if !filepath.IsAbs(dest) {
			cwd, _ := os.Getwd()
			dest = filepath.Join(cwd, dest)
		}

		cpCmd := exec.Command("cp", "-r", base, dest)
		if output, err := cpCmd.CombinedOutput(); err != nil {
			sys.LogError(fmt.Errorf("Error backing up: %v\n%s", err, output))
			return
		}

		fmt.Printf("Backed up %s to %s\n", base, dest)
	},
}
