/*
 * Copyright (C) distroy
 */

package ldref

import (
	"reflect"

	"github.com/distroy/ldgo/v2/lderr"
)

type CopyConfig struct {
	Clone     bool   // is clone if same type
	TargetTag string // default: json
	SourceTag string // default: json
}

// Copy will copy the data from source to target by the tags
//   - Copy(*int, uint)
//   - Copy(*int, *uint)
//   - Copy(*structA, structB)
//   - Copy(*structA, *structB)
//   - Copy(*structA, map)
//   - Copy(*map, structA)
func Copy(target, source interface{}, cfg ...*CopyConfig) error {
	return copyV2(target, source, cfg...)
}

// DeepCopy will copy the data from source to target by the tags
//   - DeepCopy(*int, uint)
//   - DeepCopy(*int, *uint)
//   - DeepCopy(*structA, structB)
//   - DeepCopy(*structA, *structB)
//   - DeepCopy(*structA, map)
//   - DeepCopy(*map, structA)
func DeepCopy(target, source interface{}, cfg ...*CopyConfig) error {
	return deepCopyV2(target, source, cfg...)
}

func copyV1(target, source interface{}, cfg ...*CopyConfig) error {
	c := &copyContext{
		context:    getContext(),
		CopyConfig: &CopyConfig{},
	}
	defer putContext(c.context)

	if len(cfg) > 0 && cfg[0] != nil {
		c.CopyConfig = cfg[0]
	}

	return copyWithCheckTarget(c, target, source)
}
func deepCopyV1(target, source interface{}, cfg ...*CopyConfig) error {
	c := &CopyConfig{
		Clone: true,
	}
	if len(cfg) > 0 && cfg[0] != nil {
		c = cfg[0]
		c.Clone = true
	}
	return copyV1(target, source, c)
}

func copyV2(target, source interface{}, cfg ...*CopyConfig) error {
	c := &copyContext{
		context:    getContext(),
		CopyConfig: &CopyConfig{},
	}
	defer putContext(c.context)

	if len(cfg) > 0 && cfg[0] != nil {
		c.CopyConfig = cfg[0]
	}

	return copyWithCheckTargetV2(c, target, source)
}

func deepCopyV2(target, source interface{}, cfg ...*CopyConfig) error {
	c := &CopyConfig{
		Clone: true,
	}
	if len(cfg) > 0 && cfg[0] != nil {
		c = cfg[0]
		c.Clone = true
	}
	return copyV2(target, source, c)
}

type copyContext struct {
	*context
	*CopyConfig
}

func copyWithCheckTarget(c *copyContext, target, source interface{}) error {
	sVal := valueOf(source)
	// sVal, _ = valueElment(sVal)

	if tVal, ok := target.(reflect.Value); ok {
		if !tVal.CanAddr() {
			if tVal.Kind() != reflect.Ptr {
				return lderr.ErrReflectTargetNotPtr

			} else if tVal.IsNil() {
				return lderr.ErrReflectTargetNilPtr
			}
		}

		copyReflect(c, tVal, sVal)
		return c.Error()
	}

	tVal := reflect.ValueOf(target)
	if tVal.Kind() != reflect.Ptr {
		return lderr.ErrReflectTargetNotPtr
	}
	if tVal.IsNil() {
		return lderr.ErrReflectTargetNilPtr
	}

	// tVal = tVal.Elem()

	copyReflect(c, tVal, sVal)
	return c.Error()
}

func copyWithCheckTargetV2(c *copyContext, target, source interface{}) error {
	sVal := valueOf(source)
	// sVal, _ = valueElment(sVal)

	if tVal, ok := target.(reflect.Value); ok {
		if !tVal.CanAddr() {
			if tVal.Kind() != reflect.Ptr {
				return lderr.ErrReflectTargetNotPtr

			} else if tVal.IsNil() {
				return lderr.ErrReflectTargetNilPtr
			}
		}

		copyReflectV2(c, tVal, sVal)
		return c.Error()
	}

	tVal := reflect.ValueOf(target)
	if tVal.Kind() != reflect.Ptr {
		return lderr.ErrReflectTargetNotPtr
	}
	if tVal.IsNil() {
		return lderr.ErrReflectTargetNilPtr
	}

	// tVal = tVal.Elem()

	copyReflectV2(c, tVal, sVal)
	return c.Error()
}

func valueOf(v interface{}) reflect.Value {
	if vv, ok := v.(reflect.Value); ok {
		return vv
	}

	return reflect.ValueOf(v)
}

func indirectType(_type reflect.Type) (typ reflect.Type, lvl int) {
	typ = _type
	lvl = 0
	for typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
		lvl++
	}
	return
}

func indirectCopySource(_source reflect.Value) (source reflect.Value, lvl int) {
	source = _source
	lvl = 0
	for source.Kind() == reflect.Ptr && !source.IsNil() {
		source = source.Elem()
		lvl++
	}
	return
}

func indirectCopyTarget(_target reflect.Value) (target reflect.Value, lvl int) {
	target = _target
	lvl = 0
	for target.Kind() == reflect.Ptr {
		if target.IsNil() {
			target.Set(reflect.New(target.Type().Elem()))
		}

		target = target.Elem()
		lvl++
	}

	return
}

func copyReflect(c *copyContext, target, source reflect.Value) bool {
	_target := target
	_source := source

	if !target.CanAddr() {
		target = target.Elem()
	}

	if end := copyReflectWithIndirect(c, target, source); end {
		return end
	}

	// clear target
	if !_target.CanAddr() {
		_target.Elem().Set(reflect.Zero(_target.Elem().Type()))
	} else {
		_target.Set(reflect.Zero(_target.Type()))
	}

	c.AddErrorf("%s can not copy to %s", typeNameOfReflect(_source), typeNameOfReflect(_target))
	return false
}

func copyReflectV2(c *copyContext, target, source reflect.Value) bool {
	pf, done := getCopyFunc(c, refTypeOfValue(target), refTypeOfValue(source))
	done()
	return (*pf)(c, target, source)
}

func _getCopyFuncReflect(c *copyContext, tTyp, sTyp reflect.Type) copyFuncType {
	pf, done := getCopyFuncIndirect(c, tTyp, sTyp)
	var pfe *copyFuncType
	var de func()
	if refKindOfType(tTyp) == reflect.Ptr {
		pfe, de = getCopyFuncIndirect(c, tTyp.Elem(), sTyp)
	}
	return func(c *copyContext, target, source reflect.Value) (end bool) {
		_target := target
		_source := source

		fn := *pf
		d := done
		if !target.CanAddr() {
			target = target.Elem()
			fn = *pfe
			d = de
		}

		d()
		if end := fn(c, target, source); end {
			return end
		}

		// clear target
		if !_target.CanAddr() {
			_target.Elem().Set(reflect.Zero(_target.Elem().Type()))
		} else {
			_target.Set(reflect.Zero(_target.Type()))
		}

		c.AddErrorf("%s can not copy to %s", typeNameOfReflect(_source), typeNameOfReflect(_target))
		return false
	}
}

func copyReflectWithIndirect(c *copyContext, target, source reflect.Value) bool {
	for {
		pair := copyPair{To: target.Kind(), From: source.Kind()}
		fnCopy := copyFuncMap[pair]
		if fnCopy != nil {
			return fnCopy(c, target, source)
		}

		switch source.Kind() {
		case reflect.Interface:
			source = source.Elem()
			continue

		case reflect.Ptr:
			source, _ = indirectCopySource(source)
			continue
		}

		return false
	}
}
func copyReflectWithIndirectV2(c *copyContext, target, source reflect.Value) bool {
	pf, done := getCopyFuncIndirect(c, target.Type(), source.Type())
	done()
	return (*pf)(c, target, source)
}
func _getCopyFuncReflectWithIndirect(c *copyContext, tTyp, sTyp reflect.Type) copyFuncType {
	pair := copyPair{To: refKindOfType(tTyp), From: refKindOfType(sTyp)}
	fnGet := getCopyFuncMap[pair]
	if fnGet != nil {
		return fnGet(c, tTyp, sTyp)
	}

	switch refKindOfType(sTyp) {
	case reflect.Interface:
		return func(c *copyContext, target, source reflect.Value) (end bool) {
			source = source.Elem()
			return copyReflectWithIndirectV2(c, target, source)
		}

	case reflect.Ptr:
		for sTyp.Kind() == reflect.Ptr {
			sTyp = sTyp.Elem()
		}

		pfe, done := getCopyFuncIndirect(c, tTyp, sTyp)
		tZero := reflect.Zero(tTyp)
		return func(c *copyContext, target, source reflect.Value) (end bool) {
			source, _ = indirectCopySource(source)
			if source.Kind() == reflect.Ptr {
				target.Set(tZero)
				return true
			}
			done()
			return (*pfe)(c, target, source)
		}
	}

	return func(c *copyContext, target, source reflect.Value) (end bool) {
		return false
	}
}

func isCopyTypeConvertible(toType, fromType reflect.Type) bool {
	toType, _ = indirectType(toType)
	fromType, _ = indirectType(fromType)

	pair := copyPair{To: toType.Kind(), From: fromType.Kind()}
	_, ok := copyFuncMap[pair]
	return ok
}

func isCopyTypeConvertibleV2(toType, fromType reflect.Type) bool {
	toType, _ = indirectType(toType)
	fromType, _ = indirectType(fromType)

	pair := copyPair{To: toType.Kind(), From: fromType.Kind()}
	_, ok := getCopyFuncMap[pair]
	return ok
}
