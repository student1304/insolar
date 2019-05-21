//
// Copyright 2019 Insolar Technologies GmbH
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package preprocessor

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/build"
	"go/format"
	"go/parser"
	"go/printer"
	"go/token"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"text/template"

	"github.com/insolar/insolar/platformpolicy"

	"github.com/insolar/insolar/insolar"

	"github.com/pkg/errors"
)

var foundationPath = "github.com/insolar/insolar/logicrunner/goplugin/foundation"
var proxyctxPath = "github.com/insolar/insolar/logicrunner/goplugin/proxyctx"
var corePath = "github.com/insolar/insolar/insolar"

var immutableFlag = "//ins:immutable"

const (
	TemplateDirectory = "templates"
)

// ParsedFile struct with prepared info we extract from source code
type ParsedFile struct {
	name        string
	code        []byte
	fileSet     *token.FileSet
	node        *ast.File
	machineType insolar.MachineType

	types        map[string]*ast.TypeSpec
	methods      map[string][]*ast.FuncDecl
	constructors map[string][]*ast.FuncDecl
	contract     string
}

// ParseFile parses a file as Go source code of a smart contract
// and returns it as `ParsedFile`
func ParseFile(fileName string, machineType insolar.MachineType) (*ParsedFile, error) {
	res := &ParsedFile{
		name:        fileName,
		machineType: machineType,
	}
	sourceCode, err := slurpFile(fileName)
	if err != nil {
		return nil, errors.Wrap(err, "Can't read file")
	}
	res.code = sourceCode

	res.fileSet = token.NewFileSet()
	node, err := parser.ParseFile(res.fileSet, res.name, res.code, parser.ParseComments)
	if err != nil {
		return nil, errors.Wrapf(err, "Can't parse %s", fileName)
	}
	res.node = node

	err = res.parseTypes()
	if err != nil {
		return nil, errors.Wrap(err, "")
	}

	err = res.parseFunctionsAndMethods()
	if err != nil {
		return nil, errors.Wrap(err, "")
	}
	if res.contract == "" {
		return nil, errors.New("Only one smart contract must exist")
	}

	return res, nil
}

func (pf *ParsedFile) parseTypes() error {
	pf.types = make(map[string]*ast.TypeSpec)
	for _, decl := range pf.node.Decls {
		tDecl, ok := decl.(*ast.GenDecl)
		if !ok || tDecl.Tok != token.TYPE {
			continue
		}

		for _, e := range tDecl.Specs {
			typeNode := e.(*ast.TypeSpec)

			err := pf.parseTypeSpec(typeNode)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (pf *ParsedFile) parseTypeSpec(typeSpec *ast.TypeSpec) error {
	if isContractTypeSpec(typeSpec) {
		if pf.contract != "" {
			return errors.New("more than one contract in a file")
		}
		pf.contract = typeSpec.Name.Name
	} else {
		pf.types[typeSpec.Name.Name] = typeSpec
	}

	return nil
}

func (pf *ParsedFile) parseFunctionsAndMethods() error {
	pf.methods = make(map[string][]*ast.FuncDecl)
	pf.constructors = make(map[string][]*ast.FuncDecl)
	for _, decl := range pf.node.Decls {
		fd, ok := decl.(*ast.FuncDecl)
		if !ok || !fd.Name.IsExported() {
			continue
		}

		var err error
		if fd.Recv == nil || fd.Recv.NumFields() == 0 {
			err = pf.parseConstructor(fd)
		} else {
			err = pf.parseMethod(fd)
		}
		if err != nil {
			return err
		}
	}

	return nil
}

func (pf *ParsedFile) parseConstructor(fd *ast.FuncDecl) error {
	name := fd.Name.Name
	if !strings.HasPrefix(name, "New") {
		return nil // doesn't look like a constructor
	}

	res := fd.Type.Results

	if res.NumFields() != 2 {
		return errors.Errorf("Constructor %q should return exactly two values", name)
	}

	if pf.typeName(res.List[1].Type) != "error" {
		return errors.Errorf("Constructor %q should return 'error'", name)
	}

	typename := pf.typeName(res.List[0].Type)
	pf.constructors[typename] = append(pf.constructors[typename], fd)

	return nil
}

func (pf *ParsedFile) parseMethod(fd *ast.FuncDecl) error {
	name := fd.Name.Name

	res := fd.Type.Results
	if res.NumFields() < 1 {
		return errors.Errorf("Method %q should return at least one result (error)", name)
	}

	lastResType := pf.typeName(res.List[res.NumFields()-1].Type)
	if lastResType != "error" {
		return errors.Errorf(
			"Method %q should return 'error' as last value, but it's %q",
			name, lastResType,
		)
	}

	typename := pf.typeName(fd.Recv.List[0].Type)
	pf.methods[typename] = append(pf.methods[typename], fd)

	return nil
}

// ProxyPackageName guesses user friendly contract "name" from file name
// and/or package in the file
func (pf *ParsedFile) ProxyPackageName() (string, error) {
	match := regexp.MustCompile("([^/]+)/([^/]+).(go|insgoc)$").FindStringSubmatch(pf.name)
	if match == nil {
		return "", errors.New("couldn't match filename without extension and path")
	}

	packageName := pf.node.Name.Name

	proxyPackageName := packageName
	if proxyPackageName == "main" {
		proxyPackageName = match[2]
	}
	if proxyPackageName == "main" {
		proxyPackageName = match[1]
	}
	return proxyPackageName, nil
}

// ContractName returns name of the contract
func (pf *ParsedFile) ContractName() string {
	return pf.node.Name.Name
}

func checkMachineType(machineType insolar.MachineType) error {
	if machineType != insolar.MachineTypeGoPlugin &&
		machineType != insolar.MachineTypeBuiltin {
		return errors.New("Unsupported machine type")
	}
	return nil
}

func templatePathConstruct(tplType string) string {
	return path.Join(TemplateDirectory, tplType+".go.tpl")
}

func formatAndWrite(out io.Writer, templateName string, data map[string]interface{}) error {
	templatePath := templatePathConstruct(templateName)
	tmpl, err := openTemplate(templatePath)
	if err != nil {
		return errors.Wrap(err, "couldn't open template file for wrapper")
	}

	var buff bytes.Buffer

	err = tmpl.Execute(&buff, data)
	if err != nil {
		return errors.Wrap(err, "couldn't write code output handle")
	}

	fmtOut, err := format.Source(buff.Bytes())
	if err != nil {

		return errors.Wrap(err, "couldn't format code "+buff.String())
	}

	_, err = out.Write(fmtOut)
	if err != nil {
		return errors.Wrap(err, "couldn't write code to output")
	}

	return nil
}

// WriteWrapper generates and writes into `out` source code
// of wrapper for the contract
func (pf *ParsedFile) WriteWrapper(out io.Writer, packageName string) error {
	if err := checkMachineType(pf.machineType); err != nil {
		return err
	}

	data := map[string]interface{}{
		"Package":            packageName,
		"ContractType":       pf.contract,
		"Methods":            pf.functionInfoForWrapper(pf.methods[pf.contract]),
		"Functions":          pf.functionInfoForWrapper(pf.constructors[pf.contract]),
		"ParsedCode":         pf.code,
		"FoundationPath":     foundationPath,
		"Imports":            pf.generateImports(true),
		"GenerateInitialize": pf.machineType == insolar.MachineTypeBuiltin,
	}

	return formatAndWrite(out, "wrapper", data)
}

func (pf *ParsedFile) functionInfoForWrapper(list []*ast.FuncDecl) []map[string]interface{} {
	var res []map[string]interface{}
	for _, fun := range list {
		info := map[string]interface{}{
			"Name":                fun.Name.Name,
			"ArgumentsZeroList":   generateZeroListOfTypes(pf, "args", fun.Type.Params),
			"Arguments":           numberedVars(fun.Type.Params, "args"),
			"Results":             numberedVars(fun.Type.Results, "ret"),
			"ErrorInterfaceInRes": typeIndexes(pf, fun.Type.Results, "error"),
			"Immutable":           isImmutable(fun), // only for methods, not constructors
		}
		res = append(res, info)
	}
	return res
}

func generateTextReference(pulse insolar.PulseNumber, code []byte) *insolar.Reference {
	hasher := platformpolicy.NewPlatformCryptographyScheme().ReferenceHasher()
	codeHash := hasher.Hash(code)
	return insolar.NewReference(insolar.ID{}, *insolar.NewID(pulse, codeHash))
}

// WriteProxy generates and writes into `out` source code of contract's proxy
func (pf *ParsedFile) WriteProxy(classReference string, out io.Writer) error {
	proxyPackageName, err := pf.ProxyPackageName()
	if err != nil {
		return err
	}

	if classReference == "" {
		classReference = generateTextReference(0, pf.code).String()
	}

	_, err = insolar.NewReferenceFromBase58(classReference)
	if err != nil {
		return errors.Wrap(err, "can't write proxy: ")
	}

	if err := checkMachineType(pf.machineType); err != nil {
		return err
	}

	methodsProxies := pf.functionInfoForProxy(pf.methods[pf.contract])
	constructorProxies := pf.functionInfoForProxy(pf.constructors[pf.contract])

	data := map[string]interface{}{
		"PackageName":         proxyPackageName,
		"Types":               generateTypes(pf),
		"ContractType":        pf.contract,
		"MethodsProxies":      methodsProxies,
		"ConstructorsProxies": constructorProxies,
		"ClassReference":      classReference,
		"Imports":             pf.generateImports(false),
	}

	return formatAndWrite(out, "proxy", data)
}

func (pf *ParsedFile) functionInfoForProxy(list []*ast.FuncDecl) []map[string]interface{} {
	var res []map[string]interface{}

	for _, fun := range list {
		info := map[string]interface{}{
			"Name":            fun.Name.Name,
			"Arguments":       genFieldList(pf, fun.Type.Params, true),
			"InitArgs":        generateInitArguments(fun.Type.Params),
			"ResultZeroList":  generateZeroListOfTypes(pf, "ret", fun.Type.Results),
			"Results":         numberedVars(fun.Type.Results, "ret"),
			"ErrorVar":        fmt.Sprintf("ret%d", fun.Type.Results.NumFields()-1),
			"ResultsWithErr":  commaAppend(numberedVarsI(fun.Type.Results.NumFields()-1, "ret"), "err"),
			"ResultsNilError": commaAppend(numberedVarsI(fun.Type.Results.NumFields()-1, "ret"), "nil"),
			"ResultsTypes":    genFieldList(pf, fun.Type.Results, false),
			"Immutable":       isImmutable(fun),
		}
		res = append(res, info)
	}
	return res
}

// ChangePackageToMain changes package of the parsed code to "main"
func (pf *ParsedFile) ChangePackageToMain() {
	pf.node.Name.Name = "main"
}

// Write prints `out` contract's code, it could be changed with a few methods
func (pf *ParsedFile) Write(out io.Writer) error {
	return printer.Fprint(out, pf.fileSet, pf.node)
}

// codeOfNode returns source code of an AST node
func (pf *ParsedFile) codeOfNode(n ast.Node) string {
	return string(pf.code[n.Pos()-1 : n.End()-1])
}

func (pf *ParsedFile) typeName(t ast.Expr) string {
	if tmp, ok := t.(*ast.StarExpr); ok { // *type
		t = tmp.X
	}
	return pf.codeOfNode(t)
}

func (pf *ParsedFile) generateImports(wrapper bool) map[string]bool {
	imports := make(map[string]bool)
	imports[fmt.Sprintf(`"%s"`, proxyctxPath)] = true
	if !wrapper {
		imports[fmt.Sprintf(`"%s"`, corePath)] = true
	}
	for _, method := range pf.methods[pf.contract] {
		extendImportsMap(pf, method.Type.Params, imports)
		if !wrapper {
			extendImportsMap(pf, method.Type.Results, imports)
		}
	}
	for _, fun := range pf.constructors[pf.contract] {
		extendImportsMap(pf, fun.Type.Params, imports)
		if !wrapper {
			extendImportsMap(pf, fun.Type.Results, imports)
		}
	}

	return imports
}

func openTemplate(fileName string) (*template.Template, error) {
	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		return nil, errors.Wrap(nil, "couldn't find info about current file")
	}
	templateDir := filepath.Join(filepath.Dir(currentFile), fileName)
	tmpl, err := template.ParseFiles(templateDir)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't parse template for output")
	}
	return tmpl, nil
}

func numberedVars(list *ast.FieldList, name string) string {
	if list == nil || list.NumFields() == 0 {
		return ""
	}
	return numberedVarsI(list.NumFields(), name)
}

func commaAppend(l string, r string) string {
	if l == "" {
		return r
	}
	return l + ", " + r
}

func numberedVarsI(n int, name string) string {
	if n == 0 {
		return ""
	}

	res := ""
	for i := 0; i < n; i++ {
		res = commaAppend(res, name+strconv.Itoa(i))
	}
	return res
}

func typeIndexes(parsed *ParsedFile, list *ast.FieldList, t string) []int {
	if list == nil || list.NumFields() == 0 {
		return []int{}
	}

	rets := []int{}
	for i, e := range list.List {
		if parsed.codeOfNode(e.Type) == t {
			rets = append(rets, i)
		}
	}
	return rets
}

func isContractTypeSpec(typeNode *ast.TypeSpec) bool {
	baseContract := "foundation.BaseContract"
	st, ok := typeNode.Type.(*ast.StructType)
	if !ok {
		return false
	}
	if st.Fields == nil || st.Fields.NumFields() == 0 {
		return false
	}
	for _, fd := range st.Fields.List {
		if len(fd.Names) != 0 {
			continue // named struct field
		}
		selectField, ok := fd.Type.(*ast.SelectorExpr)
		if !ok {
			continue
		}
		pack := selectField.X.(*ast.Ident).Name
		class := selectField.Sel.Name
		if baseContract == (pack + "." + class) {
			return true
		}
	}

	return false
}

func generateTypes(parsed *ParsedFile) []string {
	var types []string
	for _, t := range parsed.types {
		types = append(types, "type "+parsed.codeOfNode(t))
	}

	return types
}

func extendImportsMap(parsed *ParsedFile, params *ast.FieldList, imports map[string]bool) {
	if params == nil || params.NumFields() == 0 {
		return
	}

	for _, e := range params.List {
		if parsed.codeOfNode(e.Type) == "error" {
			imports[fmt.Sprintf(`"%s"`, foundationPath)] = true
		}
	}

	for _, e := range params.List {
		tname := parsed.codeOfNode(e.Type)
		tname = strings.Trim(tname, "*")
		tnameFrom := strings.Split(tname, ".")

		if len(tnameFrom) < 2 {
			continue
		}

		for _, imp := range parsed.node.Imports {
			var importAlias string
			var impValue string

			if imp.Name != nil {
				importAlias = imp.Name.Name
				impValue = fmt.Sprintf(`%s %s`, importAlias, imp.Path.Value)
			} else {
				impValue = imp.Path.Value
				importString := strings.Trim(impValue, `"`)
				importAlias = filepath.Base(importString)
			}

			if importAlias == tnameFrom[0] {
				imports[impValue] = true
				break
			}
		}
	}
}

func generateZeroListOfTypes(parsed *ParsedFile, name string, list *ast.FieldList) string {
	if list == nil || list.NumFields() == 0 {
		return fmt.Sprintf("%s := []interface{}{}\n", name)
	}

	text := fmt.Sprintf("%s := [%d]interface{}{}\n", name, list.NumFields())

	for i, arg := range list.List {
		tname := parsed.codeOfNode(arg.Type)
		if tname == "error" {
			tname = "*foundation.Error"
		}

		text += fmt.Sprintf("\tvar %s%d %s\n", name, i, tname)
		text += fmt.Sprintf("\t%s[%d] = &%s%d\n", name, i, name, i)
	}

	return text
}

func genFieldList(parsed *ParsedFile, params *ast.FieldList, withNames bool) string {
	res := ""
	if params == nil {
		return res
	}
	for i, e := range params.List {
		if i > 0 {
			res += ", "
		}
		if withNames {
			res += e.Names[0].Name + " "
		}
		res += parsed.codeOfNode(e.Type)
	}
	return res
}

func generateInitArguments(list *ast.FieldList) string {
	initArgs := ""
	initArgs += fmt.Sprintf("var args [%d]interface{}\n", list.NumFields())
	for i, arg := range list.List {
		initArgs += fmt.Sprintf("\targs[%d] = %s\n", i, arg.Names[0].Name)
	}
	return initArgs
}

// GetRealApplicationDir return application dir path
func GetRealApplicationDir(dir string) (string, error) {
	gopath := build.Default.GOPATH
	if gopath == "" {
		return "", errors.Errorf("GOPATH is not set")
	}
	contractsPath := ""
	for _, p := range strings.Split(gopath, ":") {
		contractsPath = path.Join(p, "src/github.com/insolar/insolar/application/", dir)
		_, err := os.Stat(contractsPath)
		if err == nil {
			return contractsPath, nil
		}
	}
	return "", errors.Errorf("Not found github.com/insolar/insolar in GOPATH")
}

// GetRealContractsNames returns names of all real smart contracts
func GetRealContractsNames() ([]string, error) {
	pathWithContracts, err := GetRealApplicationDir("contract")
	if err != nil {
		return nil, errors.Wrap(err, "[ GetContractNames ]")
	}
	if len(pathWithContracts) == 0 {
		return nil, errors.New("[ GetContractNames ] There are contracts dir")
	}
	var result []string
	files, err := ioutil.ReadDir(pathWithContracts)
	if err != nil {
		return nil, err
	}
	for _, f := range files {
		if f.IsDir() {
			result = append(result, f.Name())
		}
	}

	return result, nil
}

func slurpFile(fileName string) ([]byte, error) {
	file, err := os.OpenFile(fileName, os.O_RDONLY, 0)
	if err != nil {
		return nil, errors.Wrap(err, "Can't open file '"+fileName+"'")
	}
	defer file.Close() //nolint: errcheck

	res, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, errors.Wrap(err, "Can't read file '"+fileName+"'")
	}
	return res, nil
}

func isImmutable(decl *ast.FuncDecl) bool {
	var isImmutable = false
	if decl.Doc != nil && decl.Doc.List != nil {
		for _, comment := range decl.Doc.List {
			if comment.Text == immutableFlag {
				isImmutable = true
			}
		}
	}
	return isImmutable
}

type ContractListEntry struct {
	Name       string
	Path       string
	Parsed     *ParsedFile
	ImportPath string
	Version    int
}

const (
	CodeType      = "code"
	PrototypeType = "prototype"
)

func (e *ContractListEntry) GenerateReference(tp string) *insolar.Reference {
	contractID := fmt.Sprintf("%s::%s::v%02d", tp, e.Name, e.Version)
	return generateTextReference(insolar.BuiltinPulseNumber, []byte(contractID))
}

type ContractList []ContractListEntry

func generateContractList(contracts ContractList) interface{} {
	importList := make([]interface{}, 0)
	for _, contract := range contracts {
		data := map[string]interface{}{
			"Name":               contract.Name,
			"ImportName":         contract.Name,
			"ImportPath":         contract.ImportPath,
			"CodeReference":      contract.GenerateReference(CodeType).String(),
			"PrototypeReference": contract.GenerateReference(PrototypeType).String(),
		}
		importList = append(importList, data)
	}
	return importList
}

func GenerateInitializationList(out io.Writer, contracts ContractList) error {
	data := map[string]interface{}{
		"Contracts": generateContractList(contracts),
		"Package":   "builtin",
	}

	return formatAndWrite(out, "initialize", data)
}
