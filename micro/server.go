package micro

import (
	"strings"

	"github.com/angrypie/tie/parser"
	"github.com/angrypie/tie/template"
	. "github.com/dave/jennifer/jen"
)

const gomicro = "github.com/micro/go-micro/v2"
const gomicroClient = "github.com/micro/go-micro/v2/client"
const microModuleId = "GoMicro"
const microUtils = "github.com/angrypie/tie-modules/micro/microutils"

type PackageInfo = template.PackageInfo

func NewModule(p *parser.Parser, services []string) template.Module {
	if p.GetPackageName() == "main" {
		return NewUpgradedModule(p, services)
	}

	deps := []template.Module{
		NewClientModule(p),
		NewUpgradedModule(p, services),
	}
	return template.NewStandartModule("micromod", GenerateServer, p, deps)
}

func NewUpgradedModule(p *parser.Parser, services []string) template.Module {
	gen := func(p *parser.Parser) *template.Package {
		return GenerateUpgraded(p, services)
	}
	return template.NewStandartModule("upgraded", gen, p, nil)
}

func GenerateUpgraded(p *parser.Parser, services []string) (pkg *template.Package) {
	p.UpgradeApiImports(services, func(i string) (n string) {
		return i + "/tie_modules/micromod/client"
	})
	files := p.ToFiles()
	pkg = &template.Package{Name: "upgraded", Files: files}
	return
}

func GenerateServer(p *parser.Parser) *template.Package {
	info := template.NewPackageInfoFromParser(p)
	info.SetServicePath(info.Service.Name + "/tie_modules/micromod/upgraded")
	f := NewFile(strings.ToLower(microModuleId))

	f.Add(template.TemplateServer(info, func(g *Group, resource, instance string) {
		g.Id("service").Op(":=").Qual(gomicro, "NewService").Call(
			Qual(gomicro, "Name").Call(Lit(resource)),
		)
		g.Id("service").Dot("Init").Call()

		g.Qual(gomicro, "RegisterHandler").Call(Id("service").Dot("Server").Call(), Id(instance))
		g.Id("service").Dot("Run").Call()
	}))

	return &template.Package{
		Name:  "micromod",
		Files: [][]byte{[]byte(f.GoString())},
	}
}

