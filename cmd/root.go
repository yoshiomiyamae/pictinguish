/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/tajtiattila/metadata"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "pictinguish",
	Short: "",
	Long:  ``,
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		dest, err := cmd.Flags().GetString("dest")
		src := args[0]
		if err != nil {
			return err
		}

		if dest == "" {
			dest = src
		}

		fileInfoList, err := os.ReadDir(src)
		if err != nil {
			return err
		}

		for _, fileInfo := range fileInfoList {
			if fileInfo.IsDir() {
				continue
			}

			srcFileName := fileInfo.Name()
			srcFilePath := filepath.Join(src, srcFileName)

			fmt.Println(srcFileName)

			file, err := os.Open(srcFilePath)
			if err != nil {
				fmt.Println("=> Couldn't open", err)
				continue
			}

			defer file.Close()

			metadata, err := metadata.Parse(file)
			if err != nil {
				fmt.Println("=> Couldn't decode exif", err)
				continue
			}
			err = file.Close()
			if err != nil {
				fmt.Println("=> Couldn't close", err)
				continue
			}

			destDirPath := filepath.Join(dest, metadata.DateTimeCreated.Format("2006-01-02"))
			destDirInfo, err := os.Stat(destDirPath)
			if err != nil || !destDirInfo.IsDir() {
				err = os.Mkdir(destDirPath, 0755)
				if err != nil {
					fmt.Println(err)
				}
			}
			destFilePath := filepath.Join(
				destDirPath,
				srcFileName,
			)
			err = os.Rename(
				srcFilePath,
				destFilePath,
			)
			if err != nil {
				fmt.Println("=> Couldn't move", err)
				continue
			}
			fmt.Println("=> Skip")
		}

		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringP("dest", "d", "", "Destination")
}
