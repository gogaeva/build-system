package gomodule

import (
  "bytes"
  "strings"
  "testing"

  "github.com/google/blueprint"
  "github.com/roman-mazur/bood"
)

func TestDocFactory(t *testing.T) {
  ctx := blueprint.NewContext()

  ctx.MockFileSystem(map[string][]byte{
    "Blueprints": []byte(`
      go_doc {
        name: "test-doc",
        pkg: ".",
        srcs: ["test-src.go"],
      }  
    `),
    "test-src.go": nil,
  })

  ctx.RegisterModuleType("go_doc", DocFactory)

  cfg := bood.NewConfig()

  _, errs := ctx.ParseBlueprintsFiles(".", cfg)
  if len(errs) != 0 {
    t.Errorf("Syntax errors in the test blueprint file: %s", errs)
  }

  _, errs = ctx.PrepareBuildActions(cfg)
  if len(errs) != 0 {
    t.Errorf("Unexpected errors while preparing build actions: %s", errs)
  }

  buffer := new(bytes.Buffer)
  if err := ctx.WriteBuildFile(buffer); err != nil {
    t.Errorf("Error writing ninja file: %s", err)
  } else {
    text := buffer.String()
    t.Logf("Generated ninja build file:\n%s", text)
    if !strings.Contains(text, "out/docs/test-doc.txt: ") {
      t.Errorf("Generated ninja file does not have build statement for documentation generation")
    }
    if !strings.Contains(text, "test-src.go") {
      t.Errorf("Generated ninja file does not have source dependency")
    }
  }
}