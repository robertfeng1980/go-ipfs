// Package iface defines IPFS Core API which is a set of interfaces used to
// interact with IPFS nodes.
package iface

import (
	"context"
	"errors"
	"io"

	options "github.com/ipfs/go-ipfs/core/coreapi/interface/options"

	ipld "gx/ipfs/QmNwUEK7QbwSqyKBu3mMtToo8SUc6wQJ7gdZq4gGGJqfnf/go-ipld-format"
	cid "gx/ipfs/QmeSrf6pzut73u6zLQkRFQ3ygt3k6XFT2kjdYP8Tnkwwyg/go-cid"
)

// Path is a generic wrapper for paths used in the API. A path can be resolved
// to a CID using one of Resolve functions in the API.
type Path interface {
	String() string
	Cid() *cid.Cid
	Root() *cid.Cid
	Resolved() bool
}

// TODO: should we really copy these?
//       if we didn't, godoc would generate nice links straight to go-ipld-format
type Node ipld.Node
type Link ipld.Link

type Reader interface {
	io.ReadSeeker
	io.Closer
}

// CoreAPI defines an unified interface to IPFS for Go programs.
type CoreAPI interface {
	// Unixfs returns an implementation of Unixfs API
	Unixfs() UnixfsAPI

	// ResolvePath resolves the path using Unixfs resolver
	ResolvePath(context.Context, Path) (Path, error)

	// ResolveNode resolves the path (if not resolved already) using Unixfs
	// resolver, gets and returns the resolved Node
	ResolveNode(context.Context, Path) (Node, error)
}

// UnixfsAPI is the basic interface to immutable files in IPFS
type UnixfsAPI interface {
	// Add imports the data from the reader into merkledag file
	Add(context.Context, io.Reader) (Path, error)

	// Cat returns a reader for the file
	Cat(context.Context, Path) (Reader, error)

	// Ls returns the list of links in a directory
	Ls(context.Context, Path) ([]*Link, error)
}

//TODO: Should this use paths instead of cids?
type ObjectAPI interface {
	New(context.Context, ...options.ObjectNewOption) (Node, error)
	WithType(string) options.ObjectNewOption

	Put(context.Context, Node) (Path, error)
	Get(context.Context, Path) (Node, error)
	Data(context.Context, Path) (io.Reader, error)
	Links(context.Context, Path) ([]*Link, error)
	Stat(context.Context, Path) (*ObjectStat, error)

	AddLink(ctx context.Context, base Path, name string, child Path, opts ...options.ObjectAddLinkOption) (Node, error)
	WithCreate(create bool) options.ObjectAddLinkOption

	RmLink(context.Context, Path, string) (Node, error)
	AppendData(context.Context, Path, io.Reader) (Node, error)
	SetData(context.Context, Path, io.Reader) (Node, error)
}

type ObjectStat struct {
	Cid            *cid.Cid
	NumLinks       int
	BlockSize      int
	LinksSize      int
	DataSize       int
	CumulativeSize int
}

var ErrIsDir = errors.New("object is a directory")
var ErrOffline = errors.New("can't resolve, ipfs node is offline")
