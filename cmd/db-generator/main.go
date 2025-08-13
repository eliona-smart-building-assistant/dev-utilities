package main

import (
	"flag"
	"github.com/jackc/pgtype"
	"log"

	"github.com/go-jet/jet/v2/generator/metadata"
	"github.com/go-jet/jet/v2/generator/postgres"
	"github.com/go-jet/jet/v2/generator/template"
	postgres2 "github.com/go-jet/jet/v2/postgres"
	"github.com/lib/pq"
)

func main() {
	dsn := flag.String("dsn", "", "Database DSN")
	schema := flag.String("schema", "", "Database schema")
	path := flag.String("path", "", "Destination directory")

	flag.Parse()

	if *dsn == "" || *schema == "" || *path == "" {
		log.Fatal("Missing required parameters. Please provide dsn, schema, and path.")
	}

	err := postgres.GenerateDSN(
		*dsn,
		*schema,
		*path,
		template.Default(postgres2.Dialect).
			UseSchema(func(s metadata.Schema) template.Schema {
				return template.DefaultSchema(s).
					UseModel(template.DefaultModel().
						UseTable(func(t metadata.Table) template.TableModel {
							return template.DefaultTableModel(t).
								UseField(func(c metadata.Column) template.TableModelField {
									f := template.DefaultTableModelField(c)

									// text[] -> pq.StringArray
									if c.DataType.Kind == metadata.ArrayType && c.DataType.Name == "text[]" {
										f.Type = template.NewType(pq.StringArray{})
										return f
									}

									// json / jsonb -> pgtype.JSONB
									if c.DataType.Kind == metadata.BaseType &&
										(c.DataType.Name == "json" || c.DataType.Name == "jsonb") {
										f.Type = template.NewType(pgtype.JSONB{})
										return f
									}

									return f
								})
						}),
					)
			}),
	)

	if err != nil {
		log.Fatalf("Failed to generate code: %v", err)
	}

	log.Println("Code generation completed successfully.")
}
