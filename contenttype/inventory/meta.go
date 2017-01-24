package inventory

import (
	"github.com/leeola/kala/contenttype"
	"github.com/leeola/kala/index"
)

type Meta struct {
	contenttype.Meta
	InventoryMeta
}

type InventoryMeta struct {
	// A human description of the item.
	//
	// Note that if you want to have it be matchable against a direct key,
	// use the Name field. Description complements Name.
	Description string `json:"description"`

	// The container's anchor of the this inventory item.
	//
	// The meta that resolves from this anchor will also be InventoryMeta.
	// These can be nested arbintrarily deep.
	Container string `json:"container"`

	// An optional image hash of the inventory item in question.
	//
	// This mainly serves to help humans identify the item if the name/description
	// are not identifying alone.
	Image string `json:"image,omitempty"`
}

// ToMetadata implements Indexable.
func (m Meta) ToMetadata(im index.Metadata) {
	m.Meta.ToMetadata(im)
	m.InventoryMeta.ToMetadata(im)
}

// ToMetadata implements Indexable.
func (m InventoryMeta) ToMetadata(im index.Metadata) {
	if m.Description != "" {
		im["description"] = m.Description
	}
	if m.Container != "" {
		im["container"] = m.Container
	}
	if m.Image != "" {
		im["image"] = m.Image
	}
}

func (m *Meta) FromChanges(c contenttype.Changes) {
	m.Meta.FromChanges(c)
	m.InventoryMeta.FromChanges(c)
}

func (m *InventoryMeta) FromChanges(c contenttype.Changes) {
	if v, ok := c.GetString("description"); ok {
		m.Description = v
	}
	if v, ok := c.GetString("container"); ok {
		m.Container = v
	}
	if v, ok := c.GetString("image"); ok {
		m.Image = v
	}
}
