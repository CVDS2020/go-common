package config

import (
	"encoding/json"
	"gitee.com/sy_183/common/errors"
	"gitee.com/sy_183/common/uns"
	"gopkg.in/yaml.v3"
	"reflect"
	"strconv"
	"time"
)

type handler interface {
	canHandle(c any) bool
	handle(c any) (nc any, modified bool, err error)
	handleAfterChildren(c any) (nc any, modified bool, err error)
}

type PreModifyConfig interface {
	PreModify() (nc any, modified bool)
}

type PreHandlerConfig interface {
	PreHandle()
}

type PreModifyAfterChildrenConfig interface {
	PreModifyAfterChildren() (nc any, modified bool)
}

type PreHandleAfterChildrenConfig interface {
	PreHandleAfterChildren()
}

type PostModifyConfig interface {
	PostModify() (nc any, modified bool, err error)
}

type PostHandlerConfig interface {
	PostHandle() error
}

type PostModifyAfterChildrenConfig interface {
	PostModifyAfterChildren() (nc any, modified bool, err error)
}

type PostHandleAfterChildrenConfig interface {
	PostHandleAfterChildren() error
}

type preHandler struct{}

func (preHandler) canHandle(c any) bool {
	if _, ok := c.(PreModifyConfig); ok {
		return true
	}
	if _, ok := c.(PreHandlerConfig); ok {
		return true
	}
	return false
}

func (p preHandler) handle(c any) (nc any, modified bool, err error) {
	nc = c
	if pmc, ok := nc.(PreModifyConfig); ok {
		if nc, modified = pmc.PreModify(); !modified {
			nc = c
		}
	}
	if phc, ok := nc.(PreHandlerConfig); ok {
		phc.PreHandle()
	}
	return
}

func (p preHandler) handleAfterChildren(c any) (nc any, modified bool, err error) {
	nc = c
	if pmc, ok := nc.(PreModifyAfterChildrenConfig); ok {
		if nc, modified = pmc.PreModifyAfterChildren(); !modified {
			nc = c
		}
	}
	if phc, ok := nc.(PreHandleAfterChildrenConfig); ok {
		phc.PreHandleAfterChildren()
	}
	return
}

type postHandler struct{}

func (postHandler) canHandle(c any) bool {
	if _, ok := c.(PostModifyConfig); ok {
		return true
	}
	if _, ok := c.(PostHandlerConfig); ok {
		return true
	}
	return false
}

func (postHandler) handle(c any) (nc any, modified bool, err error) {
	nc = c
	if pmc, ok := nc.(PostModifyConfig); ok {
		if nc, modified, err = pmc.PostModify(); err != nil {
			return
		} else if !modified {
			nc = c
		} else {
			ct := reflect.TypeOf(c)
			nct := reflect.TypeOf(nc)
			if ct.Kind() == reflect.Ptr && ct.Elem() == nct {
				nec := reflect.New(nct)
				nec.Elem().Set(reflect.ValueOf(nc))
				nc = nec.Interface()
			}
		}
	}
	if phc, ok := nc.(PostHandlerConfig); ok {
		err = phc.PostHandle()
	}
	return
}

func (postHandler) handleAfterChildren(c any) (nc any, modified bool, err error) {
	nc = c
	if pmc, ok := nc.(PostModifyAfterChildrenConfig); ok {
		if nc, modified, err = pmc.PostModifyAfterChildren(); err != nil {
			return
		} else if !modified {
			nc = c
		} else {
			ct := reflect.TypeOf(c)
			nct := reflect.TypeOf(nc)
			if ct.Kind() == reflect.Ptr && ct.Elem() == nct {
				ncp := reflect.New(nct)
				ncp.Elem().Set(reflect.ValueOf(nc))
				nc = ncp.Interface()
			}
		}
	}

	if phc, ok := nc.(PostHandleAfterChildrenConfig); ok {
		err = phc.PostHandleAfterChildren()
	}
	return
}

func handleDefault(v reflect.Value, zerop, def *string, cs map[any]struct{}) error {
	vt := v.Type()
	switch vt.Kind() {
	case reflect.String:
		if v.CanSet() && def != nil {
			if v.IsZero() {
				v.SetString(*def)
			}
		}
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
		if v.CanSet() && def != nil {
			if vt == reflect.TypeOf(time.Duration(0)) {
				du, err := time.ParseDuration(*def)
				if err != nil {
					return err
				}
				if v.IsZero() {
					v.SetInt(int64(du))
				}
			} else {
				parsed, err := strconv.ParseInt(*def, 0, int(vt.Size())*8)
				if err != nil {
					return err
				}
				if v.IsZero() {
					v.SetInt(parsed)
				}
			}
		}
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
		if v.CanSet() && def != nil {
			parsed, err := strconv.ParseUint(*def, 0, int(vt.Size())*8)
			if err != nil {
				return err
			}
			if v.IsZero() {
				v.SetUint(parsed)
			}
		}
	case reflect.Float32, reflect.Float64:
		if v.CanSet() && def != nil {
			parsed, err := strconv.ParseFloat(*def, int(vt.Size())*8)
			if err != nil {
				return err
			}
			if v.IsZero() {
				v.SetFloat(parsed)
			}
		}
	case reflect.Complex64, reflect.Complex128:
		if v.CanSet() && def != nil {
			parsed, err := strconv.ParseComplex(*def, int(vt.Size())*8)
			if err != nil {
				return err
			}
			if v.IsZero() {
				v.SetComplex(parsed)
			}
		}
	case reflect.Bool:
		if v.CanSet() && def != nil {
			var zero bool
			parsed, err := strconv.ParseBool(*def)
			if err != nil {
				return err
			}
			if v.Bool() == zero {
				v.SetBool(parsed)
			}
		}
	case reflect.Array, reflect.Struct:
		if def != nil && v.IsZero() {
			if vt == reflect.TypeOf(time.Time{}) {
				if v.CanSet() {
					parsed, err := time.Parse(time.RFC3339Nano, *def)
					if err != nil {
						return err
					}
					v.Set(reflect.ValueOf(parsed))
				}
			} else if *def != "" {
				// noCopy flag represents the pointer using the current array or struct
				var noCopy bool
				var pv reflect.Value
				var pvi any
				if v.CanAddr() {
					// pointer using the current array or struct
					pv = v.Addr()
					noCopy = true
					if pv.CanInterface() {
						pvi = pv.Interface()
					}
				} else if v.CanSet() {
					// new array pointer
					pv = reflect.New(vt)
					pvi = pv.Interface()
				}
				if pvi != nil {
					// parse array use yaml or json
					err := yaml.Unmarshal(uns.StringToBytes(*def), pvi)
					if err != nil {
						err2 := json.Unmarshal(uns.StringToBytes(*def), pvi)
						if err2 != nil {
							return errors.Append(err, err2)
						}
					}
					if !noCopy {
						// use the pointer of the current array, and do not
						// need to set
						v.Set(pv.Elem())
					}
				}
			}
		}
		switch vt.Kind() {
		case reflect.Array:
			// handle array element
			l := v.Len()
			for i := 0; i < l; i++ {
				if err := handleDefault(v.Index(i), nil, nil, cs); err != nil {
					return err
				}
			}
		case reflect.Struct:
			nf := v.NumField()
			for i := 0; i < nf; i++ {
				fv := v.Field(i)
				if !vt.Field(i).IsExported() {
					continue
				}
				tag := vt.Field(i).Tag
				if fdef, has := tag.Lookup("default"); has {
					if fzero, has := tag.Lookup("zero"); has {
						if err := handleDefault(fv, &fzero, &fdef, cs); err != nil {
							return err
						}
					} else if err := handleDefault(fv, nil, &fdef, cs); err != nil {
						return err
					}
				} else if err := handleDefault(fv, nil, nil, cs); err != nil {
					return err
				}
			}
		}
	case reflect.Slice, reflect.Map:
		kind := vt.Kind()
		if def != nil && v.CanSet() {
			if *def != "" && v.Len() == 0 {
				var pv reflect.Value
				var pvi any
				if v.CanAddr() {
					if pv = v.Addr(); pv.CanInterface() {
						// use current pointer of slice or map
						if !v.IsNil() {
							// map or slice is nil, make it first
							switch kind {
							case reflect.Slice:
								v.Set(reflect.MakeSlice(vt, 0, 0))
							case reflect.Map:
								v.Set(reflect.MakeMap(vt))
							}
						}
						pvi = pv.Interface()
					}
				}
				if pvi == nil {
					// new pointer of slice or map
					pv = reflect.New(vt)
					switch kind {
					case reflect.Slice:
						pv.Elem().Set(reflect.MakeSlice(vt, 0, 0))
					case reflect.Map:
						pv.Elem().Set(reflect.MakeMap(vt))
					}
					pvi = pv.Interface()
				}
				// parse slice or map use yaml or json
				err := yaml.Unmarshal(uns.StringToBytes(*def), pvi)
				if err != nil {
					err2 := json.Unmarshal(uns.StringToBytes(*def), pvi)
					if err2 != nil {
						return errors.Append(err, err2)
					}
				}
				v.Set(pv.Elem())
			}
		}
		switch kind {
		case reflect.Slice:
			// handle slice element
			l := v.Len()
			for i := 0; i < l; i++ {
				if err := handleDefault(v.Index(i), nil, nil, cs); err != nil {
					return err
				}
			}
		case reflect.Map:
			// handle map value
			for iter := v.MapRange(); iter.Next(); {
				if err := handleDefault(iter.Value(), nil, nil, cs); err != nil {
					return err
				}
			}
		}
	case reflect.Ptr:
		if def != nil && v.CanSet() {
			nv := reflect.New(vt.Elem())
			if err := handleDefault(nv.Elem(), zerop, def, cs); err != nil {
				return err
			}
			if v.CanInterface() {
				// update recursive call set
				if !v.IsNil() {
					delete(cs, v.Interface())
				}
				v.Set(nv)
				cs[v.Interface()] = struct{}{}
			}
		} else if v.IsNil() {
			return nil
		} else {
			if v.CanInterface() {
				// check and update recursive call set
				if _, repeat := cs[v.Interface()]; repeat {
					return nil
				}
				cs[v.Interface()] = struct{}{}
			}
			if err := handleDefault(v.Elem(), zerop, def, cs); err != nil {
				return err
			}
		}
	}
	return nil
}

func HandleDefault(c any) error {
	return handleDefault(reflect.ValueOf(c), nil, nil, make(map[interface{}]struct{}))
}

func doHandlerPtr(pi any, pv reflect.Value, setFn func(v reflect.Value), handleFn func(c any) (nc any, modified bool, err error), cs map[any]struct{}) (err error, repeat bool) {
	pvt := pv.Type()
	nc, mod, err := handleFn(pi)
	if err != nil {
		return err, false
	}
	if mod {
		nct := reflect.TypeOf(nc)
		ncv := reflect.ValueOf(nc)
		if nct == pvt && ncv != pv {
			if pv.CanSet() {
				pv.Set(ncv)
			} else if setFn != nil {
				setFn(ncv)
			}
			delete(cs, pi)
			if _, repeat := cs[nc]; repeat {
				return nil, true
			}
			cs[nc] = struct{}{}
		} else if nct == pvt.Elem() {
			if ve := pv.Elem(); ve.CanSet() {
				ve.Set(ncv)
			}
		}
	}
	return nil, false
}

func doHandle(v, pv reflect.Value, copied bool, setFn func(v reflect.Value), handleFn func(c any) (nc any, modified bool, err error)) error {
	if pv.CanInterface() {
		nc, mod, err := handleFn(pv.Interface())
		if err != nil {
			return err
		}
		if mod {
			nct := reflect.TypeOf(nc)
			ncv := reflect.ValueOf(nc)
			if v.CanSet() {
				if nct == v.Type() {
					v.Set(ncv)
				} else if nct == pv.Type() && (copied || ncv != pv) {
					v.Set(ncv.Elem())
				}
			} else if setFn != nil {
				if nct == v.Type() {
					setFn(ncv)
				} else if nct == pv.Type() && (copied || ncv != pv) {
					setFn(ncv.Elem())
				}
			}
		}
	}
	return nil
}

func handleChildren(v reflect.Value, handler handler, cs map[any]struct{}) error {
	vt := v.Type()
	switch v.Kind() {
	case reflect.Slice, reflect.Array:
		ll := v.Len()
		for i := 0; i < ll; i++ {
			if err := handle(v.Index(i), nil, handler, cs); err != nil {
				return err
			}
		}
	case reflect.Map:
		for iter := v.MapRange(); iter.Next(); {
			if err := handle(iter.Value(), func(mv reflect.Value) {
				v.SetMapIndex(iter.Key(), mv)
			}, handler, cs); err != nil {
				return err
			}
		}
	case reflect.Struct:
		nf := v.NumField()
		for i := 0; i < nf; i++ {
			fv := v.Field(i)
			if !vt.Field(i).IsExported() {
				continue
			}
			if err := handle(fv, nil, handler, cs); err != nil {
				return err
			}
		}
	case reflect.Ptr:
		if err := handle(v, nil, handler, cs); err != nil {
			return err
		}
	}
	return nil
}

func handle(v reflect.Value, setFn func(v reflect.Value), handler handler, cs map[any]struct{}) error {
	vt := v.Type()
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return nil
		}
		if v.CanInterface() {
			c := v.Interface()
			if _, repeat := cs[c]; repeat {
				return nil
			}
			cs[c] = struct{}{}
			if err, repeat := doHandlerPtr(c, v, setFn, handler.handle, cs); err != nil {
				return err
			} else if repeat {
				return nil
			}
		}
		if err := handleChildren(v.Elem(), handler, cs); err != nil {
			return err
		}
		if v.CanInterface() {
			c := v.Interface()
			if err, repeat := doHandlerPtr(c, v, setFn, handler.handleAfterChildren, cs); err != nil {
				return err
			} else if repeat {
				return nil
			}
		}
	} else if handler.canHandle(reflect.NewAt(vt, nil).Interface()) {
		var pv reflect.Value
		var copied bool
		if v.CanAddr() {
			pv = v.Addr()
		} else {
			pv = reflect.New(vt)
			pv.Elem().Set(v)
			copied = true
		}
		if pv.CanInterface() {
			if err := doHandle(v, pv, copied, setFn, handler.handle); err != nil {
				return err
			}
		}
		if err := handleChildren(v, handler, cs); err != nil {
			return err
		}
		if pv.CanInterface() {
			if err := doHandle(v, pv, copied, setFn, handler.handleAfterChildren); err != nil {
				return err
			}
		}
	} else if err := handleChildren(v, handler, cs); err != nil {
		return err
	}

	return nil
}

func PreHandle(c any) error {
	return handle(reflect.ValueOf(c), nil, preHandler{}, make(map[interface{}]struct{}))
}

func PostHandle(c any) error {
	return handle(reflect.ValueOf(c), nil, postHandler{}, make(map[interface{}]struct{}))
}

//func handle3(handler handler, c any, cs map[any]struct{}) (nc any, modified bool, err error) {
//	var mod bool
//	nc = c
//	// ct: config type
//	ct := reflect.TypeOf(c)
//	// cv: config value
//	cv := reflect.ValueOf(nc)
//	v := cv
//
//	switch ct.Kind() {
//	case reflect.Ptr:
//		if cv.IsNil() {
//			return c, false, nil
//		}
//		v = v.Elem()
//		if _, repeat := cs[nc]; repeat {
//			return
//		}
//	}
//
//	old := nc
//	if nc, mod, err = handler.handle(nc); err != nil {
//		return
//	}
//	modified = modified || mod
//
//	if ct.Kind() == reflect.Ptr && mod {
//		// handler modify config
//		delete(cs, old)
//		cs[nc] = struct{}{}
//	}
//
//	switch v.Kind() {
//	case reflect.Struct:
//		nf := v.NumField()
//		// walk fields and handle filed
//		for i := 0; i < nf; i++ {
//			fv := v.Field(i)
//			if fv.Kind() == reflect.Ptr {
//				if fv.CanInterface() {
//					fi := fv.Interface()
//					if nfi, mod, err := handle3(handler, fi, cs); err != nil {
//						return nc, modified, err
//					} else if mod {
//						fv.Set(reflect.ValueOf(nfi))
//					}
//				}
//				continue
//			}
//			// field not pointer and field has address
//			if fv.CanAddr() {
//				// field has address
//				// fpv: field pointer value
//				fpv := fv.Addr()
//				if fpv.CanInterface() {
//					// fpi: field pointer interface
//					fpi := fpv.Interface()
//					// nfpi: new field pointer interface
//					if nfpi, mod, err := handle3(handler, fpi, cs); err != nil {
//						return nc, modified, err
//					} else if mod {
//						fpv.Set(reflect.ValueOf(nfpi))
//					}
//				}
//				continue
//			}
//			// field not have address
//			if fv.CanInterface() {
//				fpv := reflect.New(fv.Type())
//				fpv.Elem().Set(fv)
//				if fpv.CanInterface() {
//					// mpi: map pointer interface
//					fpi := fpv.Interface()
//					// npmi: new map pointer interface
//					if nfpi, mod, err := handle3(handler, fpi, cs); err != nil {
//						return nc, modified, err
//					} else if mod {
//						fpv = reflect.ValueOf(nfpi)
//					}
//					fv.Set(fpv.Elem())
//				}
//			}
//		}
//
//	case reflect.Map:
//		for iter := v.MapRange(); iter.Next(); {
//			// mv: map value
//			mv := iter.Value()
//			if mv.Kind() == reflect.Ptr {
//				if mv.CanInterface() {
//					// mvi: map value interface
//					mvi := mv.Interface()
//					// nmvi: new map value interface
//					// mod: map value modified
//					if nmvi, mod, err := handle3(handler, mvi, cs); err != nil {
//						return nc, modified, err
//					} else if mod {
//						v.SetMapIndex(iter.Key(), reflect.ValueOf(nmvi))
//					}
//				}
//				continue
//			}
//			if mv.CanAddr() {
//				// mpv: map pointer value
//				mpv := mv.Addr()
//				if mpv.CanInterface() {
//					// mpi: map pointer interface
//					mpi := mpv.Interface()
//					// npmi: new map pointer interface
//					if nmpi, mod, err := handle3(handler, mpi, cs); err != nil {
//						return nc, modified, err
//					} else if mod {
//						mpv.Set(reflect.ValueOf(nmpi))
//					}
//				}
//				continue
//			}
//			if mv.CanInterface() {
//				// mvi: map value interface
//				mpv := reflect.New(mv.Type())
//				mpv.Elem().Set(mv)
//				if mpv.CanInterface() {
//					// mpi: map pointer interface
//					mpi := mpv.Interface()
//					// npmi: new map pointer interface
//					if nmpi, mod, err := handle3(handler, mpi, cs); err != nil {
//						return nc, modified, err
//					} else if mod {
//						mpv = reflect.ValueOf(nmpi)
//					}
//					v.SetMapIndex(iter.Key(), mpv.Elem())
//				}
//			}
//		}
//
//	case reflect.Slice, reflect.Array:
//		ll := v.Len()
//		for i := 0; i < ll; i++ {
//			lv := v.Index(i)
//			if lv.Kind() == reflect.Ptr {
//				if lv.CanInterface() {
//					// lvi: list value interface
//					lvi := lv.Interface()
//					// nlvi: new list value interface
//					// mod: list value modified
//					if nlvi, mod, err := handle3(handler, lvi, cs); err != nil {
//						return nc, modified, err
//					} else if mod {
//						lv.Set(reflect.ValueOf(nlvi))
//					}
//				}
//				continue
//			}
//			if lv.CanAddr() {
//				// lpv: list pointer value
//				lpv := lv.Addr()
//				if lpv.CanInterface() {
//					// lpi: list pointer interface
//					lpi := lpv.Interface()
//					// nlmi: new list pointer interface
//					if nlpi, mod, err := handle3(handler, lpi, cs); err != nil {
//						return nc, modified, err
//					} else if mod {
//						lpv.Set(reflect.ValueOf(nlpi))
//					}
//				}
//				continue
//			}
//			if lv.CanInterface() {
//				lpv := reflect.New(lv.Type())
//				lpv.Elem().Set(lv)
//				if lpv.CanInterface() {
//					// mpi: map pointer interface
//					lpi := lpv.Interface()
//					// npmi: new map pointer interface
//					if nlpi, mod, err := handle3(handler, lpi, cs); err != nil {
//						return nc, modified, err
//					} else if mod {
//						lpv = reflect.ValueOf(nlpi)
//					}
//					lv.Set(lpv.Elem())
//				}
//			}
//		}
//	}
//
//	old = nc
//	if nc, mod, err = handler.handleAfterChildren(nc); err != nil {
//		return
//	}
//	modified = modified || mod
//
//	if ct.Kind() == reflect.Ptr && mod {
//		// handler modify config
//		delete(cs, old)
//		cs[nc] = struct{}{}
//	}
//
//	return
//}

func Handle(c any) error {
	if err := HandleDefault(c); err != nil {
		return err
	}
	if err := PreHandle(c); err != nil {
		return err
	}
	if err := PostHandle(c); err != nil {
		return err
	}
	return nil
}
