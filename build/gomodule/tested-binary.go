package gomodule

import (
	"fmt"
	"path"

	"github.com/google/blueprint"
	"github.com/roman-mazur/bood"
)

var (
	pctx = blueprint.NewPackageContext("github.com/gogaeva/build-system/build/gomodule")

	goBuild = pctx.StaticRule("binaryBuild", blueprint.RuleParams{
		Command:     "cd $workDir && go build -gcflags=\"all=-N -l\" -o $outputPath $pkg",
		Description: "build go command $pkg",
	}, "workDir", "outputPath", "pkg")

	goVendor = pctx.StaticRule("vendor", blueprint.RuleParams{
		Command:     "cd $workDir && go mod vendor",
		Description: "vendor dependencies of $name",
	}, "workDir", "name")

	goTest = pctx.StaticRule("test", blueprint.RuleParams{
		Command:     "cd $workDir && go test -v $pkg > $outputPath",
		Description: "test $pkg",
	}, "workDir", "pkg", "outputPath")
)

type testedBinaryModule struct {
	blueprint.SimpleName

	properties struct {
		Pkg         string
		Srcs        []string
		SrcsExclude []string
		TestPkg     string
		TestSrcs    []string
		VendorFirst bool
		Deps        []string
	}
}

func (tb *testedBinaryModule) DynamicDependencies(blueprint.DynamicDependerModuleContext) []string {
	return tb.properties.Deps
}

func (tb *testedBinaryModule) GenerateBuildActions(ctx blueprint.ModuleContext) {
	name := ctx.ModuleName()
	config := bood.ExtractConfig(ctx)
	config.Debug.Printf("Adding build actions for go binary module %s", name)

	outputPath := path.Join(config.BaseOutputDir, "bin", name)
	outputTestsPath := path.Join(config.BaseOutputDir, "reports", name, "test.txt")

	var inputs []string
	var testInputs []string
	inputErrors := false
	for _, src := range tb.properties.Srcs {
		if matches, err := ctx.GlobWithDeps(src, tb.properties.SrcsExclude); err == nil {
			inputs = append(inputs, matches...)
		} else {
			ctx.PropertyErrorf("srcs", "Cannot resolve files that match pattern %s", src)
			inputErrors = true
		}
	}
	for _, testSrc := range tb.properties.TestSrcs {
		if matches, err := ctx.GlobWithDeps(testSrc, []string{}); err == nil {
			testInputs = append(testInputs, matches...)
		} else {
			ctx.PropertyErrorf("srcs", "Cannot resolve files that match pattern %s", testSrc)
			inputErrors = true
		}
	}
	if inputErrors {
		return
	}

	if tb.properties.VendorFirst {
		vendorDirPath := path.Join(ctx.ModuleDir(), "vendor")
		ctx.Build(pctx, blueprint.BuildParams{
			Description: fmt.Sprintf("Vendor dependencies of %s", name),
			Rule:        goVendor,
			Outputs:     []string{vendorDirPath},
			Implicits:   []string{path.Join(ctx.ModuleDir(), "go.mod")},
			Optional:    true,
			Args: map[string]string{
				"workDir": ctx.ModuleDir(),
				"name":    name,
			},
		})
		inputs = append(inputs, vendorDirPath)
	}

	ctx.Build(pctx, blueprint.BuildParams{
		Description: fmt.Sprintf("Build %s as Go binary", name),
		Rule:        goBuild,
		Outputs:     []string{outputPath},
		Implicits:   inputs,
		Args: map[string]string{
			"workDir":    ctx.ModuleDir(),
			"pkg":        tb.properties.Pkg,
			"outputPath": outputPath,
		},
	})

	if tb.properties.TestPkg != "" {
		ctx.Build(pctx, blueprint.BuildParams{
			Description: fmt.Sprintf("Test %s", name),
			Rule:        goTest,
			Outputs:     []string{outputTestsPath},
			Implicits:   testInputs,
			Args: map[string]string{
				"workDir":    ctx.ModuleDir(),
				"pkg":        tb.properties.TestPkg,
				"outputPath": outputTestsPath,
			},
		})
	}
}

func TestedBinFactory() (blueprint.Module, []interface{}) {
	tbm := &testedBinaryModule{}
	return tbm, []interface{}{&tbm.SimpleName.Properties, &tbm.properties}
}
