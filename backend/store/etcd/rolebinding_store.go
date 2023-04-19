package etcd

import (
	"context"

	v2 "github.com/sensu/core/v2"
	"github.com/sensu/sensu-go/backend/store"
)

var (
	roleBindingsPathPrefix = "rbac/rolebindings"
	roleBindingKeyBuilder  = store.NewKeyBuilder(roleBindingsPathPrefix)
)

func getRoleBindingPath(roleBinding *v2.RoleBinding) string {
	return roleBindingKeyBuilder.WithResource(roleBinding).Build(roleBinding.Name)
}

// GetRoleBindingsPath gets the path of the role binding store.
func GetRoleBindingsPath(ctx context.Context, name string) string {
	return roleBindingKeyBuilder.WithContext(ctx).Build(name)
}

// CreateRoleBinding ...
func (s *Store) CreateRoleBinding(ctx context.Context, roleBinding *v2.RoleBinding) error {
	if err := roleBinding.Validate(); err != nil {
		return &store.ErrNotValid{Err: err}
	}
	return Create(ctx, s.client, getRoleBindingPath(roleBinding), roleBinding.Namespace, roleBinding)
}

// CreateOrUpdateRoleBinding ...
func (s *Store) CreateOrUpdateRoleBinding(ctx context.Context, roleBinding *v2.RoleBinding) error {
	if err := roleBinding.Validate(); err != nil {
		return &store.ErrNotValid{Err: err}
	}
	return CreateOrUpdate(ctx, s.client, getRoleBindingPath(roleBinding), roleBinding.Namespace, roleBinding)
}

// DeleteRoleBinding ...
func (s *Store) DeleteRoleBinding(ctx context.Context, name string) error {
	return Delete(ctx, s.client, GetRoleBindingsPath(ctx, name))
}

// GetRoleBinding ...
func (s *Store) GetRoleBinding(ctx context.Context, name string) (*v2.RoleBinding, error) {
	role := &v2.RoleBinding{}
	err := Get(ctx, s.client, GetRoleBindingsPath(ctx, name), role)
	return role, err
}

// ListRoleBindings ...
func (s *Store) ListRoleBindings(ctx context.Context, pred *store.SelectionPredicate) ([]*v2.RoleBinding, error) {
	roles := []*v2.RoleBinding{}
	err := List(ctx, s.client, GetRoleBindingsPath, &roles, pred)
	return roles, err
}

// UpdateRoleBinding ...
func (s *Store) UpdateRoleBinding(ctx context.Context, roleBinding *v2.RoleBinding) error {
	if err := roleBinding.Validate(); err != nil {
		return &store.ErrNotValid{Err: err}
	}
	return Update(ctx, s.client, getRoleBindingPath(roleBinding), roleBinding.Namespace, roleBinding)
}
