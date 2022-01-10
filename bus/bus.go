package bus

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/n-creativesystem/rbns/config"
	"github.com/n-creativesystem/rbns/logger"
	"github.com/n-creativesystem/rbns/tracer"
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

type TransactionManager interface {
	InTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}

type Bus interface {
	Dispatch(ctx context.Context, msg Msg) error
	PublishCtx(ctx context.Context, msg Msg) error

	InTransaction(ctx context.Context, fn func(ctx context.Context) error) error
	AddHandler(implName string, handler HandlerFunc)
	AddEventListenerCtx(handler HandlerFunc)

	SetTransaction(txMng TransactionManager)
}

type InProcBus struct {
	log logger.Logger
	// handlers            map[string]map[string]HandlerFunc
	handlersWithContext map[string]map[string]HandlerFunc
	// listeners           map[string][]HandlerFunc
	listenersWithCtx map[string][]HandlerFunc
	txMng            TransactionManager
}

func newInProcBus() *InProcBus {
	return &InProcBus{
		log: logger.New("bus"),
		// handlers:            make(map[string]map[string]HandlerFunc),
		handlersWithContext: make(map[string]map[string]HandlerFunc),
		// listeners:           make(map[string][]HandlerFunc),
		listenersWithCtx: make(map[string][]HandlerFunc),
		txMng:            &noopTransactionManager{},
	}
}

func (b *InProcBus) Dispatch(ctx context.Context, msg Msg) error {
	implName, ok := ctx.Value(ContextImplementKey).(string)
	if !ok {
		implName = config.ImplName
	}
	var msgName = reflect.TypeOf(msg).Elem().Name()
	ctx, span := tracer.Start(ctx, fmt.Sprintf("bus - %s - %s", implName, msgName),
		trace.WithAttributes(attribute.String("implement name", implName), attribute.String("msg", msgName)))
	defer span.End()

	// withCtx := true
	var handler = b.handlersWithContext[implName][msgName]
	if handler == nil {
		return ErrHandlerNotFound
		// withCtx = false
		// handler = b.handlers[msgName]
		// if handler == nil {
		// }
	}

	var params = []reflect.Value{reflect.ValueOf(ctx), reflect.ValueOf(msg)}
	// params = append(params, reflect.ValueOf(msg))

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

func (b *InProcBus) InTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return b.txMng.InTransaction(ctx, fn)
}

func (b *InProcBus) AddHandler(implName string, handler HandlerFunc) {
	handlerType := reflect.TypeOf(handler)
	queryTypeName := handlerType.In(1).Elem().Name()
	if b.handlersWithContext[implName] == nil {
		b.handlersWithContext[implName] = make(map[string]HandlerFunc)
	}
	b.handlersWithContext[implName][queryTypeName] = handler
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

func (b *InProcBus) SetTransaction(txMng TransactionManager) {
	b.txMng = txMng
}

type noopTransactionManager struct{}

func (*noopTransactionManager) InTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return fn(ctx)
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
