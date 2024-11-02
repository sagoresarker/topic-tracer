package cmd

import (
	"fmt"
	"os"

	"github.com/sagoresarker/topic-tracer/internal/indexer"
	"github.com/sagoresarker/topic-tracer/internal/search"

	"github.com/spf13/cobra"
)

var (
	directory string
	query     string
	rootCmd   = &cobra.Command{
		Use:   "topic-tracer",
		Short: "Search through PDF files in a directory",
		Long:  `A fast PDF search tool that helps you find content across multiple PDF files.`,
	}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	indexCmd := &cobra.Command{
		Use:   "index",
		Short: "Index PDF files in a directory",
		Run:   runIndex,
	}
	indexCmd.Flags().StringVarP(&directory, "directory", "d", "", "Directory containing PDF files")
	indexCmd.MarkFlagRequired("directory")

	searchCmd := &cobra.Command{
		Use:   "search",
		Short: "Search indexed PDF files",
		Run:   runSearch,
	}
	searchCmd.Flags().StringVarP(&query, "query", "q", "", "Search query")
	searchCmd.MarkFlagRequired("query")

	rootCmd.AddCommand(indexCmd, searchCmd)
}

func runIndex(cmd *cobra.Command, args []string) {
	idx := indexer.New()
	if err := idx.IndexDirectory(directory); err != nil {
		fmt.Printf("Error indexing directory: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Indexing completed successfully")
}

func runSearch(cmd *cobra.Command, args []string) {
	searcher, err := search.New()
	if err != nil {
		fmt.Printf("Error initializing searcher: %v\n", err)
		os.Exit(1)
	}

	results, err := searcher.Search(query)
	if err != nil {
		fmt.Printf("Error searching: %v\n", err)
		os.Exit(1)
	}

	for _, result := range results {
		fmt.Printf("\nFile: %s\nPage: %d\nContext: %s\n",
			result.Filename, result.Page, result.Context)
	}
}
