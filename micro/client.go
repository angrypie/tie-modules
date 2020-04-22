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
	//TODO all modules noods to create upgraded subpackage to make ServicePath reusable,
	info.SetServicePath(info.Service.Name + "/tie_modules/micromod/upgraded")
	f := NewFile(strings.ToLower(microModuleId))

	f.Add(template.CreateReqRespTypes(info))
	f.Add(template.CreateTypeAliases(info))

	makeClientAPI(info, f)

	return &template.Package{
		Name:  "client",
		Files: [][]byte{[]byte(f.GoString())},
	}
}

func makeClientAPI(info *PackageInfo, f *File) {

	cb := func(receiverType string, constructor *parser.Function) {
		f.Type().Id(receiverType).Struct()
	}
	template.MakeForEachReceiver(info, cb)

	template.ForEachFunction(info, true, func(fn *parser.Function) {
		args := fn.Arguments

		body := func(g *Group) {
			rpcMethodName, requestType, responseType := template.GetMethodTypes(fn)
			request, response := template.ID("request"), template.ID("response")

			g.Id(response).Op(":=").New(Id(responseType))
			g.Id(request).Op(":=").New(Id(requestType))

			if len(args) != 0 {
				g.ListFunc(template.CreateArgsListFunc(args, request)).Op("=").
					ListFunc(template.CreateArgsListFunc(args))
			}

			resourceName := template.GetResourceName(info)
			g.Err().Op("=").Qual(microUtils, "NewClient").Call().Dot("Call").Call(
				Lit(resourceName),
				Lit(fmt.Sprintf("%s.%s", resourceName, rpcMethodName)),
				Id(request), Id(response),
			)
			template.AddIfErrorGuard(g, nil, nil)

			g.Return(ListFunc(template.CreateArgsListFunc(fn.Results, response)))
		}

		f.Func().ListFunc(func(g *Group) {
			if template.HasReceiver(fn) {
				g.Params(Id("resource").Id(fn.Receiver.GetLocalTypeName()))
				return
			}
		}).Id(fn.Name).
			ParamsFunc(template.CreateSignatureFromArgs(args, info)).
			ParamsFunc(template.CreateSignatureFromArgs(fn.Results, info)).BlockFunc(body)
	})
}