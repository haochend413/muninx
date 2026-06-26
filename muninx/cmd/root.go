package cmd

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	"github.com/haochend413/muninx/config"
	"github.com/haochend413/muninx/internal/app"
	"github.com/haochend413/muninx/internal/clients"
	"github.com/haochend413/muninx/internal/db"
	"github.com/haochend413/muninx/internal/ui"
	"github.com/haochend413/muninx/state"
	"github.com/haochend413/muninx/sys"
	"github.com/spf13/cobra"
)

var globalCfg *config.Config
var globalDB *db.DB
var globalApp *app.App
var globalModel *ui.Model
var globalEmbedClient *clients.EmbedClient

var flagEmbed bool

var rootCmd = &cobra.Command{
	Use:   "muninx",
	Short: "muninx",
	Long:  "muninx",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Load app config
		fmt.Println(`                                         G  
                                        IG  
                                       GGGH 
                                     IGGGGH 
               t                   LGGFGGGX 
          GGGGGGGGGGG            OGFGGGGGG  
        iGF   LGGGGGGGG       0 FGGGGGGGF k 
    NGGGGGFGFGGGGGGGGGGGF0  HGKGFFFI GGGGFG 
  GKTXWPJGGGGGGGGGGGGGGI  GGGGGFFFn GFFGGF  
           5GGGGGGGGGF  KFGGGFFFG 9FFFGGFS  
             GGGGGGGF  GGFFGGFG  GFFFFFFg   
             GGGGGGG zGGFGFGGN wFFGFFFL     
             GGGGGFUzFGGGGFG  IGGGFI   Y    
            GGGGGGG FGGGFFT rGGGGILGFFE     
           lFGGGG7GGGGGFF  GGGGGGGGFFG      
           ZGGFF FGGGGGH  FFGGGGGGFF        
           3GFO GFGGGGP  FFFGGGGFFz         
            GF QFGGGFQ  EGGGGFGR            
            4  FFGFFX uGGGGGL sJF0          
              KGFFFO  GGGGGGGGGG            
              GGFFM  GGGGGGFFFh             
              GGGG  GGGGGGFG                
             OGFG  GGGGGH                   
              GFJ GGGGd4GG                  
              GG  GGGGGL                    
              DG GGGG                       
                 GI                         
                qw                          
               3G                           
              GGe                           
            eGfGG                           
            I6 HG                           
              Hm                            
                                            
                                            `)
		cfg := config.LoadOrCreateConfig()
		globalCfg = &cfg

		// Start and connect to the embedding server only when -E is given.
		if flagEmbed {
			if err := launchEmbedServer(cfg.VenvDir, cfg.EmbedServerDir); err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
			globalEmbedClient = clients.NewEmbedClient(embedBaseURL)
		}

		var err error
		globalDB, err = db.NewDB(cfg.DataFilePath+"/notes_dev.db", globalEmbedClient) // TODO: change this back in official version!
		if err != nil {
			sys.LogError(fmt.Errorf("Failed to connect to database: %v", err))
			os.Exit(1)
		}
	},
	Run: func(cmd *cobra.Command, argså []string) {
		// Get state (can be nil if first run)
		s, err := state.LoadState(globalCfg.StateFilePath)
		if err != nil {
			// Use default state if load fails
			s = state.DefaultState()
		}

		// Initialize application with AppState
		globalApp = app.NewApp(globalDB, &s.App, globalEmbedClient)

		// Initialize UI model with full state
		model := ui.NewModel(globalApp, globalCfg, s)
		globalModel = &model

		// Run Bubble Tea program
		p := tea.NewProgram(model)
		if _, err := p.Run(); err != nil {
			sys.LogError(err)
			os.Exit(1)
		}
	},
}

func Execute() {
	defer func() {
		if globalDB != nil {
			globalDB.Close()
		}
	}()

	if err := rootCmd.Execute(); err != nil {
		sys.LogError(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolVarP(&flagEmbed, "embed", "E", false, "start the embedding server before launching")
	rootCmd.AddCommand(ExportNoteCmd)
	rootCmd.AddCommand(LaunchGUICmd)
	rootCmd.AddCommand(DataBackupCmd)
	rootCmd.AddCommand(RelatedNotesCmd)
}
