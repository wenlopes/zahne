# AI Development Agent Instructions

Este documento fornece diretrizes abrangentes para assistentes de IA ao trabalhar no projeto. Cobre padrões de desenvolvimento Go, práticas de teste, validação e fluxos de trabalho específicos da equipe para garantir contribuições de código consistentes e de alta qualidade.

- Você é um especialista em Go e familiarizado com a linguagem de programação Go, sua biblioteca padrão e bibliotecas comuns usadas no ecossistema Go.
- Mantenha as respostas concisas e focadas em soluções técnicas.

## Arquitetura de Código

### Estrutura de Arquivos

Abaixo está um exemplo de estrutura de arquivos recomendada para pacotes Go, alguns blocos são opcionais dependendo da complexidade do pacote, quantidade de interfaces e implementações, etc., mas considerar seguir este padrão para manter a consistência:

```go
// Validação de interface
var _ Interface = (*Implementation)(nil)

// Construtor
func NewImplementation(...) *Implementation { ... }

// Struct de implementação
type Implementation struct { ... }

// Interface definition
type SomeInterface interface { ... }

// Input structs
type SomeInput struct { ... }
```

### Interfaces e Contratos

- Implemente validação de interface com `var _ Interface = (*Implementation)(nil)` (interface guard).
- Defina interfaces após as implementações quando no mesmo arquivo.
- Métodos sempre recebem `context.Context` como primeiro parâmetro.
- Prefira input structs em vez de múltiplos parâmetros; para construtores com muitos parâmetros opcionais, considere o Functional Options Pattern.
- Retorne apenas `error` para operações de validação/execução.
- Aceite interfaces e retorne tipos concretos.
- Crie structs dedicadas para inputs de métodos (ex: `ValidationInput`).
- Mantenha inputs próximos à interface que os utiliza.

Interface Guards garantem que uma struct implementa uma interface em tempo de compilação:

```go
// Garante que ContractValidator implementa interface Validator
var _ Validator = (*ContractValidator)(nil)

// Múltiplas implementações de interface
var (
    _ http.Handler = (*ValidationHandler)(nil)
    _ Operation    = (*HTTP)(nil)
    _ Operation    = (*Async)(nil)
)
```

### Injeção de Dependências

- Construa dependências complexas dentro do construtor usando `New*` functions.
- Ao usar uma configuração global, injete uma struct de ambientes (`Environments`), e não valores individuais.

### Functional Options Pattern

Use o padrão Functional Options para construtores que necessitam de configuração flexível. Este padrão permite evolução da API sem quebrar retrocompatibilidade.

**Quando usar:**

- Construtores com muitos parâmetros opcionais.
- APIs que precisam manter retrocompatibilidade ao adicionar novos parâmetros.
- Configurações que se beneficiam de legibilidade no código consumidor.

**Estrutura padrão:**

```go
// ServerOption configura parâmetros do Server.
type ServerOption func(*Server)

// WithTimeout define o timeout do servidor.
func WithTimeout(timeout time.Duration) ServerOption {
	return func(s *Server) {
		s.timeout = timeout
	}
}

// NewServer cria um novo Server com as opções fornecidas.
func NewServer(opts ...ServerOption) *Server {
	s := &Server{
		timeout: 30 * time.Second, // valor padrão
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}
```

**Uso:**

```go
server := NewServer(
	WithTimeout(time.Minute),
	WithMaxConn(100),
)
```

**Convenções:**

- Prefixe funções de opção com `With`.
- Documente cada função de opção com comentário godoc.
- Defina valores padrão sensatos no construtor.
- Use `<Type>Option` para o tipo da função de opção.

## Padrões de Código

### Convenções Gerais

- Comentários devem terminar com ponto final.
- **SEMPRE** use `any` em vez de `interface{}`.
- Structs que implementam interfaces devem definir um Interface Guard.
- Use nomes de variáveis significativos e evite variáveis de uma letra, exceto para loops de curta duração.
  - **Bom**: `schemaProcessor`, `contractValidator`, `openAPIData`
  - **Ruim**: `sp`, `cv`, `d` (a menos que em contextos de curta duração)
  - **Aceitável**: `i` em `for i := 0; i < len(items); i++`
- Nomes de pacotes devem ser minúsculos, palavras únicas quando possível.
- Coloque os métodos públicos antes dos privados na definição da struct.
- Documente todos os métodos públicos com comentários godoc.
- Use [`gofmt`](https://golang.org/cmd/gofmt/) e [`goimports`](https://golang.org/x/tools/cmd/goimports) para formatar código consistentemente.
- Use `bytes.NewReader` (imutável) a menos que você realmente precise modificar `jsonData` depois, via `bytes.NewBuffer`.

### Tratamento de Erros

- **SEMPRE** trate erros explicitamente; **NUNCA** os ignore.
- Use `pperr.Wrap` e `pperr.New` para criar e envolver erros.
- Se necessário, use `fmt.Errorf` com verbo `%w` para envolver erros ao adicionar contexto.
- Crie tipos de erro personalizados para erros específicos do domínio.
- Use `errors.Is()` e `errors.As()` para verificação de erros.

### Logging

- Use logging estruturado com níveis de log apropriados.
- Escreva mensagens de log limpas e inclua contexto relevante como atributo:
```go
		m.instrument.Log().Error(ctx,
			"something bad happened",
			err,
			slog.Int("http_status", resp.StatusCode),
			slog.String("this_should_be_useful", usefulVariable),
		)
```
- Use `Log().With()` para criar um logger com atributos pré-configurados quando múltiplos logs compartilham o mesmo contexto:
```go
		loggerCtx := n.instrument.Log().With(
			slog.String("fallback_channel", n.fallbackChannel),
			slog.String("original_channel", origChannel),
		)

		// Reutilize o logger em múltiplas chamadas.
		loggerCtx.Info(ctx, "attempting fallback")
		loggerCtx.Error(ctx, "fallback failed", err)
```
- Log no nível apropriado: DEBUG, INFO, WARN, ERROR.

## Testes

### Padrões Gerais

- Use o comando `make test` para executar todos os testes (unit + integration + coverage).
- **SEMPRE** use `t.Context()` em vez de `context.Background()` em testes para propagação adequada de contexto.
- Para nomes de subteste (`t.Run()`): seja legível, descritivo, conciso, evite caracteres especiais; espaços substituídos por underscores.
- Use testes orientados por tabela para testar múltiplos cenários.
- Teste tanto o caminho feliz quanto casos de erro.
- Use `t.Parallel()` para testes que podem executar em paralelo.
- **SEMPRE** atualize expectativas de teste ao refatorar código que muda formatos de saída, valores de retorno ou assinaturas de função.
- Teste as interfaces públicas, não implementações privadas. Teste métodos privados indiretamente através de interfaces públicas.
- Prefira o uso de `testify/assert` para aumentar a legibilidade dos testes.
- Evite usar `gomock.Any()` em expectativas de mock, a menos que absolutamente necessário. Seja específico sobre os argumentos esperados.

### Nomenclatura e Estrutura

| Aspecto | Melhor Prática |
| --- | --- |
| Nomes de teste de nível superior | Test + nome da função/método, PascalCase para casos específicos (ex: `TestFromOpenAPIOneOfWithInlineObject`). |
| Nomes de teste de método | Formato TestType_Method (ex: `TestHTTP_OfKind`). |
| Nomes de subteste | Legível, descritivo, conciso, evite caracteres especiais; espaços substituídos por underscores. |
| Estrutura de teste | Table-driven com `t.Run()`; AAA explícito como fallback para testes não table-driven. |
| Testes complexos | Extraia funções auxiliares para reduzir complexidade cognitiva ≤10. |

### Mocks

- Use `uber-go/mock` para gerar mocks.
- Adicione a diretiva `//go:generate mockgen ...` na interface que deseja mockar.
- Execute `make mocks` para gerar/atualizar todos os mocks.
- Mocks gerados devem ser commitados.

### Documentação de API

A documentação completa da API está disponível em `docs/swagger.yaml` e `docs/asyncapi.yaml`. Ao implementar ou modificar endpoints, eventos e mensagens, **SEMPRE** revise e atualize a documentação seguindo as boas práticas:

- **Documentação GoDoc gerada**: Ao adicionar ou modificar pacotes, tipos ou funções públicas, verificar se o `Makefile` possui make target de geração de docs (ex: `docs-api`, `docs-build`) e executá-lo. Arquivos gerados devem ser commitados junto com o código.
  - Em toda e qualquer nova implementação que adiciona ou modifica pacotes, tipos ou funções públicas exportadas deve-se gerar os GoDocs e estes devem ser validados em codereviews.

- **Erros internos (500)**: Retornar mensagens genéricas (`pperr.New("failed to...", pperr.EINTERNAL)`) para evitar vazamento de informações sensíveis
- **Logging detalhado**: Sempre logar erros originais via `instrument.NoticeError(ctx, err)` para debugging
- **Input sanitization**: Aplicar sanitização em todos os campos de texto para prevenir XSS (usar `bluemonday`)
- **Validação rigorosa**: Validar todos os campos obrigatórios e retornar erros descritivos (400)
- **Documentação Swagger**: Incluir exemplos, descrições claras e constraints de validação
- **Status HTTP corretos**: Usar códigos apropriados (200, 201, 400, 500) conforme a operação

## Dependências e Fluxo de Trabalho

### Gerenciamento de Dependências

- Use Go modules para gerenciamento de dependências.
- Mantenha dependências mínimas e atualizadas.
- Prefira biblioteca padrão em vez de dependências externas quando possível.
- Use `make tidy` para limpar dependências.
- Prefira usar as abstrações fornecidas pelas shared-libs, não crie novas sem uma real necessidade.
- Considere compartilhar como módulo implementações de clientes diretos para APIs externas.

### Executar pós mudanças

Execute na ordem: `gofmt -l .`, `make tidy`, `make test`, `make lint`, `make docs-build`. Todos devem passar sem erros. Commite arquivos gerados de documentação.