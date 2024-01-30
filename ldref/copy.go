/*
 * Copyright (C) distroy
 */

package ldref

import (
	"reflect"

	"github.com/distroy/ldgo/v2/lderr"
)

type CopyConfig struct {
	Clone bool
}

func Copy(target, source interface{}, cfg ...*CopyConfig) lderr.Error {
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

func DeepCopy(target, source interface{}) lderr.Error {
	cfg := &CopyConfig{
		Clone: true,
	}
	return Copy(target, source, cfg)
}

type copyContext struct {
	*context
	*CopyConfig
}

func copyWithCheckTarget(c *copyContext, target, source interface{}) lderr.Error {
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

func copyReflectWithIndirect(c *copyContext, target, source reflect.Value) bool {
	for {
		pair := copyPair{To: target.Kind(), From: source.Kind()}
		fnCopy := copyFuncMap[pair]
		if fnCopy != nil {
			return fnCopy(c, target, source)
		}

		switch source.Kind() {
		case reflect.Interface:
			source = reflect.ValueOf(source.Interface())
			continue

		case reflect.Ptr:
			source, _ = indirectCopySource(source)
			continue
		}

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
