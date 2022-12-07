// Code generated by scripts/gengraphql.go. DO NOT EDIT.

package schema

import (
	errors "errors"
	graphql1 "github.com/graphql-go/graphql"
	graphql "github.com/sensu/sensu-go/graphql"
)

// UserFieldResolvers represents a collection of methods whose products represent the
// response values of the 'User' type.
type UserFieldResolvers interface {
	// Username implements response to request for 'username' field.
	Username(p graphql.ResolveParams) (string, error)

	// Groups implements response to request for 'groups' field.
	Groups(p graphql.ResolveParams) ([]string, error)

	// Disabled implements response to request for 'disabled' field.
	Disabled(p graphql.ResolveParams) (bool, error)

	// HasPassword implements response to request for 'hasPassword' field.
	HasPassword(p graphql.ResolveParams) (bool, error)
}

// UserAliases implements all methods on UserFieldResolvers interface by using reflection to
// match name of field to a field on the given value. Intent is reduce friction
// of writing new resolvers by removing all the instances where you would simply
// have the resolvers method return a field.
type UserAliases struct{}

// Username implements response to request for 'username' field.
func (_ UserAliases) Username(p graphql.ResolveParams) (string, error) {
	val, err := graphql.DefaultResolver(p.Source, p.Info.FieldName)
	ret, ok := val.(string)
	if err != nil {
		return ret, err
	}
	if !ok {
		return ret, errors.New("unable to coerce value for field 'username'")
	}
	return ret, err
}

// Groups implements response to request for 'groups' field.
func (_ UserAliases) Groups(p graphql.ResolveParams) ([]string, error) {
	val, err := graphql.DefaultResolver(p.Source, p.Info.FieldName)
	ret, ok := val.([]string)
	if err != nil {
		return ret, err
	}
	if !ok {
		return ret, errors.New("unable to coerce value for field 'groups'")
	}
	return ret, err
}

// Disabled implements response to request for 'disabled' field.
func (_ UserAliases) Disabled(p graphql.ResolveParams) (bool, error) {
	val, err := graphql.DefaultResolver(p.Source, p.Info.FieldName)
	ret, ok := val.(bool)
	if err != nil {
		return ret, err
	}
	if !ok {
		return ret, errors.New("unable to coerce value for field 'disabled'")
	}
	return ret, err
}

// HasPassword implements response to request for 'hasPassword' field.
func (_ UserAliases) HasPassword(p graphql.ResolveParams) (bool, error) {
	val, err := graphql.DefaultResolver(p.Source, p.Info.FieldName)
	ret, ok := val.(bool)
	if err != nil {
		return ret, err
	}
	if !ok {
		return ret, errors.New("unable to coerce value for field 'hasPassword'")
	}
	return ret, err
}

// UserType User describes an operator in the system
var UserType = graphql.NewType("User", graphql.ObjectKind)

// RegisterUser registers User object type with given service.
func RegisterUser(svc *graphql.Service, impl UserFieldResolvers) {
	svc.RegisterObject(_ObjectTypeUserDesc, impl)
}
func _ObjTypeUserUsernameHandler(impl interface{}) graphql1.FieldResolveFn {
	resolver := impl.(interface {
		Username(p graphql.ResolveParams) (string, error)
	})
	return func(frp graphql1.ResolveParams) (interface{}, error) {
		return resolver.Username(frp)
	}
}

func _ObjTypeUserGroupsHandler(impl interface{}) graphql1.FieldResolveFn {
	resolver := impl.(interface {
		Groups(p graphql.ResolveParams) ([]string, error)
	})
	return func(frp graphql1.ResolveParams) (interface{}, error) {
		return resolver.Groups(frp)
	}
}

func _ObjTypeUserDisabledHandler(impl interface{}) graphql1.FieldResolveFn {
	resolver := impl.(interface {
		Disabled(p graphql.ResolveParams) (bool, error)
	})
	return func(frp graphql1.ResolveParams) (interface{}, error) {
		return resolver.Disabled(frp)
	}
}

func _ObjTypeUserHasPasswordHandler(impl interface{}) graphql1.FieldResolveFn {
	resolver := impl.(interface {
		HasPassword(p graphql.ResolveParams) (bool, error)
	})
	return func(frp graphql1.ResolveParams) (interface{}, error) {
		return resolver.HasPassword(frp)
	}
}

func _ObjectTypeUserConfigFn() graphql1.ObjectConfig {
	return graphql1.ObjectConfig{
		Description: "User describes an operator in the system",
		Fields: graphql1.Fields{
			"disabled": &graphql1.Field{
				Args:              graphql1.FieldConfigArgument{},
				DeprecationReason: "",
				Description:       "self descriptive",
				Name:              "disabled",
				Type:              graphql1.NewNonNull(graphql1.Boolean),
			},
			"groups": &graphql1.Field{
				Args:              graphql1.FieldConfigArgument{},
				DeprecationReason: "",
				Description:       "self descriptive",
				Name:              "groups",
				Type:              graphql1.NewNonNull(graphql1.NewList(graphql1.NewNonNull(graphql1.String))),
			},
			"hasPassword": &graphql1.Field{
				Args:              graphql1.FieldConfigArgument{},
				DeprecationReason: "",
				Description:       "self descriptive",
				Name:              "hasPassword",
				Type:              graphql1.NewNonNull(graphql1.Boolean),
			},
			"username": &graphql1.Field{
				Args:              graphql1.FieldConfigArgument{},
				DeprecationReason: "",
				Description:       "self descriptive",
				Name:              "username",
				Type:              graphql1.NewNonNull(graphql1.String),
			},
		},
		Interfaces: []*graphql1.Interface{},
		IsTypeOf: func(_ graphql1.IsTypeOfParams) bool {
			// NOTE:
			// Panic by default. Intent is that when Service is invoked, values of
			// these fields are updated with instantiated resolvers. If these
			// defaults are called it is most certainly programmer err.
			// If you're see this comment then: 'Whoops! Sorry, my bad.'
			panic("Unimplemented; see UserFieldResolvers.")
		},
		Name: "User",
	}
}

// describe User's configuration; kept private to avoid unintentional tampering of configuration at runtime.
var _ObjectTypeUserDesc = graphql.ObjectDesc{
	Config: _ObjectTypeUserConfigFn,
	FieldHandlers: map[string]graphql.FieldHandler{
		"disabled":    _ObjTypeUserDisabledHandler,
		"groups":      _ObjTypeUserGroupsHandler,
		"hasPassword": _ObjTypeUserHasPasswordHandler,
		"username":    _ObjTypeUserUsernameHandler,
	},
}
