package sync

type namespace struct {
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
	config        NamespaceConfig
}

func (v *namespace) Alias(ns ...string) Namespace {
	for _, s := range ns {
		if len(s) > 0 {
			v.aliases = append(v.aliases, s)
		}
	}
	return v
}

func (v *namespace) Name(s string) Namespace {
	v.name = s
	return v
}

func (v *namespace) Summary(s string) Namespace {
	v.summary = s
	return v
}

func (v *namespace) Doc(s string) Namespace {
	v.documentation = s
	return v
}

func (v *namespace) Authorize(f AuthorizeFunc) Namespace {
	v.authorize = f
	return v
}

func (v *namespace) Middleware(s ...HandlerFunc) Namespace {
	for _, f := range s {
		if f != nil {
			v.middleware = append(v.middleware, f)
		}
	}
	return v
}

func (v *namespace) Handle(s string, f HandlerFunc) Namespace {
	if s != "" && f != nil {
		h := &handler{
			name: s,
			fn:   f,
		}
		v.handlers = append(v.handlers, h)
	}
	return v
}

func (v *namespace) Catch(f HandlerFunc) Namespace {
	v.catch = f
	return v
}

func (v *namespace) Channels() Channels {
	return v.channels
}

func (v *namespace) Config() NamespaceConfig {
	if v.config == nil {
		v.config = &config{
			namespace:     v.namespace,
			aliases:       v.aliases,
			name:          v.name,
			documentation: v.documentation,
			summary:       v.summary,
			authorize:     v.authorize,
			catch:         v.catch,
			handlers:      v.handlers,
			middleware:    v.middleware,
			channels:      v.channels,
		}
	}
	return v.config
}
