package sync

type config struct {
	namespace     string
	aliases       []string
	name          string
	documentation string
	summary       string
	authorize     AuthorizeFunc
	catch         HandlerFunc
	handlers      []Handler
	middleware    HandlerFuncs
	channels      Channels
}

func (v *config) Namespace() string {
	return v.namespace
}

func (v *config) Aliases() []string {
	return v.aliases
}

func (v *config) Name() string {
	return v.name
}

func (v *config) Summary() string {
	return v.summary
}

func (v *config) Authorize() AuthorizeFunc {
	return v.authorize
}

func (v *config) Middlewares() []HandlerFunc {
	return v.middleware
}

func (v *config) Handlers() []Handler {
	return v.handlers
}

func (v *config) Catch() HandlerFunc {
	return v.catch
}

func (v *config) Doc() string {
	return v.documentation
}

func (v *config) Channels() Channels {
	return v.channels
}
