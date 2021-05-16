package gomodule

import (
  "fmt"
  "path"

  "github.com/google/blueprint"
  "github.com/roman-mazur/bood"
)

var goDoc = pctx.StaticRule("doc", blueprint.RuleParams{
  Command:     "cd $workDir && go doc -all -u $pkg > $outputPath",
  Description: "generate documentation for $pkg",
}, "workDir", "outputPath", "pkg")

type goDocModule struct {
  blueprint.SimpleName

  properties struct {
    Pkg  string
    Srcs []string
  }
}

func (gd *goDocModule) GenerateBuildActions(ctx blueprint.ModuleContext) {
  name := ctx.ModuleName()
  config := bood.ExtractConfig(ctx)
  config.Debug.Printf("Adding build actions for go documentation module %s", name)

  var inputs []string
  outputPath := path.Join(config.BaseOutputDir, "docs", fmt.Sprintf("%s.txt", name))

  for _, testSrc := range gd.properties.Srcs {
    if matches, err := ctx.GlobWithDeps(testSrc, []string{}); err == nil {
      inputs = append(inputs, matches...)
    } else {
      ctx.PropertyErrorf("srcs", "Cannot resolve files that match pattern %s", testSrc)
      return
    }
  }

  ctx.Build(pctx, blueprint.BuildParams{
    Description: fmt.Sprintf("Generate documentation for %s", name),
    Rule:        goDoc,
    Outputs:     []string{outputPath},
    Implicits:   inputs,
    Args: map[string]string{
      "workDir":    ctx.ModuleDir(),
      "outputPath": outputPath,
      "pkg":        gd.properties.Pkg,
    },
  })
}

//DocFactory is a go_doc module factory function,
//which is invoked to instantiate a new Module object
//to handle the build action generation for the module.
func DocFactory() (blueprint.Module, []interface{}) {
  dm := &goDocModule{}
  return dm, []interface{}{&dm.SimpleName.Properties, &dm.properties}
}