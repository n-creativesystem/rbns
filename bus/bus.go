package bus

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/n-creativesystem/rbns/ncsfw/logger"
	"github.com/n-creativesystem/rbns/ncsfw/tracer"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type ImplName struct{}

var (
	ErrHandlerNotFound  = errors.New("handler not found")
	ContextImplementKey = ImplName{}
)

type HandlerFunc interface{}

type Msg interface{}

type Bus interface {
	Dispatch(ctx context.Context, msg Msg) error
	PublishCtx(ctx context.Context, msg Msg) error

	AddHandler(implName string, handler HandlerFunc)
	AddEventListenerCtx(handler HandlerFunc)
}

type InProcBus struct {
	log logger.Logger
	// handlers            map[string]HandlerFunc
	handlersWithContext map[string]HandlerFunc
	// listeners           map[string][]HandlerFunc
	listenersWithCtx map[string][]HandlerFunc
}

func newInProcBus() *InProcBus {
	return &InProcBus{
		log: logger.New("bus"),
		// handlers:            make(map[string]HandlerFunc),
		handlersWithContext: make(map[string]HandlerFunc),
		// listeners:           make(map[string][]HandlerFunc),
		listenersWithCtx: make(map[string][]HandlerFunc),
	}
}

func (b *InProcBus) Dispatch(ctx context.Context, msg Msg) error {
	var msgName = reflect.TypeOf(msg).Elem().Name()
	ctx, span := tracer.Start(ctx, fmt.Sprintf("bus - %s", msgName), trace.WithAttributes(attribute.String("msg", msgName)))
	defer span.End()

	var handler = b.handlersWithContext[msgName]
	if handler == nil {
		return ErrHandlerNotFound
	}

	var params = []reflect.Value{reflect.ValueOf(ctx), reflect.ValueOf(msg)}

	ret := reflect.ValueOf(handler).Call(params)
	err := ret[0].Interface()
	if err == nil {
		return nil
	}
	return err.(error)
}

func (b *InProcBus) PublishCtx(ctx context.Context, msg Msg) error {
	var msgName = reflect.TypeOf(msg).Elem().Name()
	ctx, span := tracer.Start(ctx, "bus - "+msgName, trace.WithAttributes(attribute.String("msg", msgName)))
	defer span.End()

	var params = []reflect.Value{}
	if listeners, exists := b.listenersWithCtx[msgName]; exists {
		params = append(params, reflect.ValueOf(ctx))
		params = append(params, reflect.ValueOf(msg))
		if err := callListeners(listeners, params); err != nil {
			return err
		}
	}
	return nil
}

func callListeners(listeners []HandlerFunc, params []reflect.Value) error {
	for _, listenerHandler := range listeners {
		ret := reflect.ValueOf(listenerHandler).Call(params)
		e := ret[0].Interface()
		if e != nil {
			err, ok := e.(error)
			if ok {
				return err
			}
			return fmt.Errorf("expected listener to return an error, got '%T'", e)
		}
	}
	return nil
}

func (b *InProcBus) AddHandler(implName string, handler HandlerFunc) {
	handlerType := reflect.TypeOf(handler)
	queryTypeName := handlerType.In(1).Elem().Name()
	b.handlersWithContext[queryTypeName] = handler
}

func (b *InProcBus) AddEventListenerCtx(handler HandlerFunc) {
	handlerType := reflect.TypeOf(handler)
	eventName := handlerType.In(1).Elem().Name()
	_, exists := b.listenersWithCtx[eventName]
	if !exists {
		b.listenersWithCtx[eventName] = make([]HandlerFunc, 0)
	}
	b.listenersWithCtx[eventName] = append(b.listenersWithCtx[eventName], handler)
}

func (b *InProcBus) GetHandlerCtx(name string) HandlerFunc {
	return b.handlersWithContext[name]
}

var globalBus = newInProcBus()

func AddHandler(implName string, handler HandlerFunc) {
	globalBus.AddHandler(implName, handler)
}

func AddEventListenerCtx(handler HandlerFunc) {
	globalBus.AddEventListenerCtx(handler)
}

func Dispatch(ctx context.Context, msg Msg) error {
	return globalBus.Dispatch(ctx, msg)
}

func PublishCtx(ctx context.Context, msg Msg) error {
	return globalBus.PublishCtx(ctx, msg)
}

func GetHandlerCtx(name string) HandlerFunc {
	return globalBus.GetHandlerCtx(name)
}

func ClearBusHandlers() {
	globalBus = newInProcBus()
}

func GetBus() Bus {
	return globalBus
}
