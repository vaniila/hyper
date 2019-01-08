package dataloader

type server struct {
	id      string
	batches []Batch
	opts    Options
}

func (v *server) Start() error {
	return nil
}

func (v *server) Stop() error {
	return nil
}

func (v *server) Instance() DataLoaders {
	loaders := make(map[Batch]DataLoader, len(v.batches))
	for _, batch := range v.batches {
		loaders[batch] = newBatchedLoader(batch, v.opts)
	}
	return &Loaders{loaders}
}

func (v *server) String() string {
	return "Hyper::DataLoader"
}
