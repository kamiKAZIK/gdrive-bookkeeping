package cmd

import (
	"fmt"
	"time"
    "context"
    "encoding/json"
    "log"
    "net/http"
    "os"

    "golang.org/x/oauth2"
    "golang.org/x/oauth2/google"
    "google.golang.org/api/drive/v3"
    "google.golang.org/api/option"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var year int

var createFoldersCmd = &cobra.Command{
	Use:   "createFolders",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("createFolders called")

        ctx := context.Background()
        b, err := os.ReadFile("credentials.json")
        if err != nil {
            log.Fatalf("Unable to read client secret file: %v", err)
        }

        config, err := google.ConfigFromJSON(b, drive.DriveMetadataReadonlyScope)
        if err != nil {
            log.Fatalf("Unable to parse client secret file to config: %v", err)
        }
        client := getClient(config)

        srv, err := drive.NewService(ctx, option.WithHTTPClient(client))
        if err != nil {
            log.Fatalf("Unable to retrieve Drive client: %v", err)
        }

        r, err := srv.Files.List().Q("mimeType = 'application/vnd.google-apps.folder' and name = ''").PageSize(10).Fields("nextPageToken, files(id, name)").Do()
        if err != nil {
            log.Fatalf("Unable to retrieve files: %v", err)
        }
        fmt.Println("Files:")
        if len(r.Files) == 0 {
            fmt.Println("No files found.")
        } else {
            for _, i := range r.Files {
                fmt.Printf("%s (%s)\n", i.Name, i.Id)
            }
        }
	},
}

func init() {
	rootCmd.AddCommand(createFoldersCmd)

	createFoldersCmd.Flags().IntVarP(&year, "year", "y", time.Now().Year(), "Year for which the folders need to be created")
	createFoldersCmd.MarkFlagRequired("year")

	viper.BindPFlag("year", createFoldersCmd.Flags().Lookup("year"))
}

func getClient(config *oauth2.Config) *http.Client {
    tokFile := "token.json"
    tok, err := tokenFromFile(tokFile)
    if err != nil {
        tok = getTokenFromWeb(config)
        saveToken(tokFile, tok)
    }

    return config.Client(context.Background(), tok)
}

func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
    authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
    fmt.Printf("Go to the following link in your browser then type the authorization code: \n%v\n", authURL)

    var authCode string
    if _, err := fmt.Scan(&authCode); err != nil {
        log.Fatalf("Unable to read authorization code %v", err)
    }

    tok, err := config.Exchange(context.TODO(), authCode)
    if err != nil {
        log.Fatalf("Unable to retrieve token from web %v", err)
    }

    return tok
}

func tokenFromFile(file string) (*oauth2.Token, error) {
    f, err := os.Open(file)
    if err != nil {
        return nil, err
    }
    defer f.Close()

    tok := &oauth2.Token{}
    err = json.NewDecoder(f).Decode(tok)

    return tok, err
}

func saveToken(path string, token *oauth2.Token) {
    fmt.Printf("Saving credential file to: %s\n", path)
    f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
    if err != nil {
        log.Fatalf("Unable to cache oauth token: %v", err)
    }
    defer f.Close()
    json.NewEncoder(f).Encode(token)
}
