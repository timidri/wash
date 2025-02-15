package volume

import (
	"context"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/puppetlabs/wash/datastore"
	"github.com/puppetlabs/wash/plugin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockDirEntry struct {
	plugin.EntryBase
	mock.Mock
}

func (m *mockDirEntry) VolumeList(ctx context.Context, path string) (DirMap, error) {
	arger := m.Called(ctx, path)
	return arger.Get(0).(DirMap), arger.Error(1)
}

func (m *mockDirEntry) VolumeOpen(context.Context, string) (plugin.SizedReader, error) {
	return nil, nil
}

func (m *mockDirEntry) VolumeStream(context.Context, string) (io.ReadCloser, error) {
	return nil, nil
}

func (m *mockDirEntry) Schema() *plugin.EntrySchema {
	return nil
}

func TestVolumeDir(t *testing.T) {
	dmap, err := StatParseAll(strings.NewReader(fixture), mountpoint, mountpoint, mountDepth)
	assert.Nil(t, err)

	plugin.SetTestCache(datastore.NewMemCache())
	entry := mockDirEntry{EntryBase: plugin.NewEntry("mine")}
	entry.SetTestID("/mine")
	ctx := context.Background()

	assert.NotNil(t, dmap[RootPath]["path"])
	vd := newDir("path", dmap[RootPath]["path"], &entry, "/path")
	attr := plugin.Attributes(vd)
	assert.Equal(t, 0755|os.ModeDir, attr.Mode())

	assert.NotNil(t, dmap[RootPath]["path1"])
	vd = newDir("path", dmap[RootPath]["path1"], &entry, "/path1")
	entry.On("VolumeList", ctx, "/path1").Return(dmap, nil).Once()
	entries, err := vd.List(ctx)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(entries))
	assert.Equal(t, "a file", plugin.Name(entries[0]))
	if entry, ok := entries[0].(*file); assert.Equal(t, true, ok) {
		assert.Equal(t, "/path1/a file", entry.path)
	}

	assert.NotNil(t, dmap[RootPath]["path2"])
	vd = newDir("path", dmap[RootPath]["path2"], &entry, "/path2")
	vd.dirmap = dmap
	entries, err = vd.List(ctx)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(entries))
	assert.Equal(t, "dir", plugin.Name(entries[0]))
	if entry, ok := entries[0].(*dir); assert.Equal(t, true, ok) {
		assert.Equal(t, "/path2/dir", entry.path)
	}

	plugin.UnsetTestCache()
}
