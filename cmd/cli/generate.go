/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"fmt"

	"github.com/popoffvg/go-mapstruct/generator"
	"github.com/popoffvg/go-mapstruct/templates"
	"github.com/spf13/cobra"
)

const (
	srcPkgPathFlagName = "srcPackagePath"
	dstPkgPathFlagName = "dstPackagePath"

	srcTypeNameFlagName = "srcTypeName"
	dstTypeNameFlagName = "dstTypeName"
)

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "generate code for type transformation",
	RunE: func(cmd *cobra.Command, _ []string) error {
		src, err := cmd.Flags().GetString(srcPkgPathFlagName)
		if err != nil {
			return fmt.Errorf("failed parse flag \"%s\" %w", srcPkgPathFlagName, err)
		}
		dst, err := cmd.Flags().GetString(dstPkgPathFlagName)
		if err != nil {
			return fmt.Errorf("failed parse flag \"%s\" %w", dstPkgPathFlagName, err)
		}

		if dst == "" {
			dst = src
		}

		srcTypeName, err := cmd.Flags().GetString(srcTypeNameFlagName)
		if err != nil {
			return fmt.Errorf("failed parse flag \"%s\" %w", srcTypeNameFlagName, err)
		}
		dstTypeName, err := cmd.Flags().GetString(dstTypeNameFlagName)
		if err != nil {
			return fmt.Errorf("failed parse flag \"%s\" %w", dstTypeNameFlagName, err)
		}

		g, err := generator.New(generator.Config{
			SrcTypeName: srcTypeName,
			DstTypeName: dstTypeName,
			DstPkg:  dst,
			SrcPkg:  src,
			Dir:         ".", // TODO: parse from args
		})

		if err != nil {
			return fmt.Errorf("transform failed: %w", err)
		}

		transformSettings, err := g.Run()
		if err != nil {
			return fmt.Errorf("transform failed: %w", err)
		}

		mng := templates.New()
		b, err := mng.Process(transformSettings)
		if err != nil {
			return fmt.Errorf("transform failed: %w", err)
		}
		cmd.Print(string(b))


		return nil
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)

	generateCmd.Flags().StringP(srcPkgPathFlagName, "s", "", "path to package with source type")
	generateCmd.MarkFlagRequired(srcPkgPathFlagName)

	generateCmd.Flags().StringP(dstPkgPathFlagName, "d", "", "path to package with destination type")

	generateCmd.Flags().StringP(srcTypeNameFlagName, "", "", "source type name")
	generateCmd.MarkFlagRequired(srcTypeNameFlagName)

	generateCmd.Flags().StringP(dstTypeNameFlagName, "", "", "destination type name")
	generateCmd.MarkFlagRequired(dstTypeNameFlagName)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// generateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// generateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
