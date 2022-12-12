package postgres

import (
	"reflect"
	"testing"
	"time"

	corev2 "github.com/sensu/core/v2"
	corev3 "github.com/sensu/core/v3"
)

func TestEntityWrapperSQLParams(t *testing.T) {
	want := len(new(EntityWrapper).SQLParams())
	got := reflect.ValueOf(EntityWrapper{}).NumField()
	if got > want {
		t.Errorf("field added to EntityWrapper but not SQLParams: got %d, want %d", got, want)
	}
	if got < want {
		t.Errorf("field removed from EntityWrapper, but not SQLParams: got %d, want %d", got, want)
	}
}

func TestEntityWrapperUnwrap(t *testing.T) {
	wrapper := EntityWrapper{
		Namespace:         "default",
		Name:              "name",
		LastSeen:          time.Now().Unix(),
		Selectors:         []byte(`{"labels.foo":"bar"}`),
		Annotations:       []byte(`{"anno":"t8n"}`),
		Hostname:          "localhost",
		OS:                "skyos",
		Platform:          "skynet",
		PlatformFamily:    "the machine AI collective",
		Arch:              "arm64",
		LibCType:          "gcc",
		VMSystem:          "kvm",
		VMRole:            "host",
		FloatType:         "hardware",
		SensuAgentVersion: "10.9.8",
		NetworkNames:      []string{"one"},
		NetworkMACs:       []string{"two"},
		NetworkAddresses:  []string{`["foobar"]`},
	}
	resource, err := wrapper.Unwrap()
	if err != nil {
		t.Fatal(err)
	}
	entity := resource.(*corev3.EntityState)
	if err := entity.Validate(); err != nil {
		t.Fatal(err)
	}
	wantMeta := &corev2.ObjectMeta{
		Namespace:   "default",
		Name:        "name",
		Labels:      map[string]string{"foo": "bar"},
		Annotations: map[string]string{"anno": "t8n"},
	}
	if got, want := entity.Metadata, wantMeta; !reflect.DeepEqual(got, want) {
		t.Errorf("bad Metadata: got %v, want %v", got, want)
	}
	if got, want := entity.LastSeen, wrapper.LastSeen; got != want {
		t.Errorf("bad LastSeen: got %v, want %v", got, want)
	}
	if got, want := entity.SensuAgentVersion, "10.9.8"; got != want {
		t.Errorf("bad SensuAgentVersion: got %v, want %v", got, want)
	}
	wantSystem := corev2.System{
		Hostname:       "localhost",
		OS:             "skyos",
		Platform:       "skynet",
		PlatformFamily: "the machine AI collective",
		Arch:           "arm64",
		LibCType:       "gcc",
		VMSystem:       "kvm",
		VMRole:         "host",
		FloatType:      "hardware",
		Network: corev2.Network{
			Interfaces: []corev2.NetworkInterface{
				{
					Name:      "one",
					MAC:       "two",
					Addresses: []string{"foobar"},
				},
			},
		},
	}
	if got, want := entity.System, wantSystem; !reflect.DeepEqual(got, want) {
		t.Errorf("bad System: got %v, want %v", got, want)
	}
}

func TestEntityWrapperWrapUnwrap(t *testing.T) {
	entity := corev3.FixtureEntityState("testent")
	// processes not supported
	entity.System.Processes = nil
	got := WrapEntityState(entity)
	entityCopy, err := got.Unwrap()
	if err != nil {
		t.Fatal(err)
	}
	if err := entityCopy.Validate(); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(entity, entityCopy) {
		t.Errorf("wrap/unwrap cycle yielded different results: got %#v, want %#v", entityCopy, entity)
	}
}
