package micro

import (
	"fmt"
	"strings"

	"github.com/angrypie/tie/parser"
	"github.com/angrypie/tie/template"
	. "github.com/dave/jennifer/jen"
)

func NewClientModule(p *parser.Parser) template.Module {
	return template.NewStandartModule("client", GenerateClient, p, nil)
}

func GenerateClient(p *parser.Parser) (pkg *template.Package) {
	info := template.NewPackageInfoFromParser(p)
	//TODO all modules needs to create upgraded subpackage to make ServicePath reusable,
	info.SetServicePath(info.Service.Name + "/tie_modules/micromod/upgraded")
	f := NewFile(strings.ToLower(microModuleId))

	code := template.TemplateClient(info, func(ids template.ClientMethodIds) *Statement {
		code := Comment("go-micro specific call").Line()
		code.Id(ids.Err).Op("=").Qual(microUtils, "NewClient").
			Call().Dot("Call").Call(
			Lit(ids.Resource),
			Lit(fmt.Sprintf("%s.%s", ids.Resource, ids.Method)),
			Id(ids.Request), Id(ids.Response),
		)
		return code
	})

	f.Add(code)

	return &template.Package{
		Name:  "client",
		Files: [][]byte{[]byte(f.GoString())},
	}
}
