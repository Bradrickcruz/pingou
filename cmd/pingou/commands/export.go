package commands

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

var exportCmd = &cobra.Command{
	Use:               "export-db",
	Short:             "Exporta banco de dados SQLite",
	Long:              "Exporta o banco de dados atual para um arquivo SQLite.\nSe --output nao for especificado, usa o nome do banco de PINGOU_DATABASE_URL ou 'pingou.db' no diretorio atual.",
	RunE:              runExportDB,
	PersistentPreRunE: requireKey,
}

var flagOutput string

func init() {
	exportCmd.Flags().StringVarP(&flagOutput, "output", "o", "", "caminho do arquivo de output (default: PINGOU_DATABASE_URL ou pingou.db)")
}

func runExportDB(cmd *cobra.Command, args []string) error {
	// Carregar .env primeiro
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "WARNING: failed to load .env: %v\n", err)
	}

	// Obter nome do banco via env diretamente
	srcPath := os.Getenv("PINGOU_DATABASE_URL")
	if srcPath == "" {
		srcPath = "pingou.db"
	}

	// Obter PWD
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("cannot get current directory: %w", err)
	}

	// Se nao for path absoluto, resolver contra PWD
	if !filepath.IsAbs(srcPath) {
		srcPath = filepath.Join(cwd, srcPath)
	}

	// Verificar se arquivo existe
	if _, err := os.Stat(srcPath); os.IsNotExist(err) {
		return fmt.Errorf("banco de dados nao encontrado: %s", srcPath)
	}

	// Determinar output path
	outputPath := flagOutput

	if outputPath == "" {
		// Sem --output: criar no PWD com prefixo "exported_"
		outputPath = filepath.Join(cwd, "exported_"+filepath.Base(srcPath))
	} else if !filepath.IsAbs(outputPath) {
		// Garantir output no PWD se nao for path absoluto
		outputPath = filepath.Join(cwd, outputPath)
	}

	// Copiar arquivo principal
	if err := copyFile(srcPath, outputPath); err != nil {
		return fmt.Errorf("falha ao copiar banco: %w", err)
	}

	// Copiar arquivos WAL e SHM se existirem
	for _, suffix := range []string{"-wal", "-shm"} {
		srcWal := srcPath + suffix
		dstWal := outputPath + suffix
		if _, err := os.Stat(srcWal); err == nil {
			if err := copyFile(srcWal, dstWal); err != nil {
				return fmt.Errorf("falha ao copiar %s: %w", suffix, err)
			}
		}
	}

	fmt.Printf("Banco exportado para: %s\n", outputPath)
	return nil
}

// copyFile copia um arquivo de src para dst
func copyFile(src, dst string) error {
	srcF, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcF.Close()

	dstF, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstF.Close()

	_, err = io.Copy(dstF, srcF)
	return err
}
