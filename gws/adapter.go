package gws

import (
	"time"

	"github.com/vaniila/hyper/router"
)

type adaptor struct {
	s Service
}

func (v *adaptor) Emit(e router.GQLEvent) error {
	var idx int
	d := &Distribution{
		Field:   e.Field(),
		Payload: e.Payload(),
		Filters: make([]*Filter, len(e.Filters())),
		Strict:  e.Strict(),
	}
	for k, v := range e.Filters() {
		o := &Filter{
			Key: k,
		}
		switch i := v.(type) {
		case string:
			o.ValOneof = &Filter_StringValue{i}
		case int:
			o.ValOneof = &Filter_IntValue{int64(i)}
		case int8:
			o.ValOneof = &Filter_IntValue{int64(i)}
		case int16:
			o.ValOneof = &Filter_IntValue{int64(i)}
		case int32:
			o.ValOneof = &Filter_IntValue{int64(i)}
		case int64:
			o.ValOneof = &Filter_IntValue{i}
		case uint:
			o.ValOneof = &Filter_IntValue{int64(i)}
		case uint8:
			o.ValOneof = &Filter_IntValue{int64(i)}
		case uint16:
			o.ValOneof = &Filter_IntValue{int64(i)}
		case uint32:
			o.ValOneof = &Filter_IntValue{int64(i)}
		case uint64:
			o.ValOneof = &Filter_IntValue{int64(i)}
		case float32:
			o.ValOneof = &Filter_FloatValue{float64(i)}
		case float64:
			o.ValOneof = &Filter_FloatValue{i}
		case bool:
			o.ValOneof = &Filter_BoolValue{i}
		case []byte:
			o.ValOneof = &Filter_BytesValue{i}
		case time.Time:
			o.ValOneof = &Filter_TimeValue{i.UnixNano()}
		}
		d.Filters[idx] = o
		idx++
	}
	if len(e.EqIDs()) > 0 || len(e.NeIDs()) > 0 || len(e.EqKeys()) > 0 || len(e.NeKeys()) > 0 {
		d.Condition = &Condition{
			EqIDs:  e.EqIDs(),
			NeIDs:  e.NeIDs(),
			EqKeys: e.EqKeys(),
			NeKeys: e.NeKeys(),
		}
	}
	return v.s.Publish(d)
}
