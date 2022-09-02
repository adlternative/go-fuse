// Copyright 2019 the Go-FUSE Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fs

import (
	"time"

	"github.com/hanwen/go-fuse/v2/fuse"
)

// Mount mounts the given NodeFS on the directory, and starts serving
// requests. This is a convenience wrapper around NewNodeFS and
// fuse.NewServer.  If nil is given as options, default settings are
// applied, which are 1 second entry and attribute timeout.
/* dir-> xxx-fs root->the path we want to look  */
func Mount(dir string, root InodeEmbedder, options *Options) (*fuse.Server, error) {
	if options == nil {
		oneSec := time.Second
		options = &Options{
			EntryTimeout: &oneSec,
			AttrTimeout:  &oneSec,
		}
	}
	/* 1. create fs (only root now) */
	rawFS := NewNodeFS(root, options)
	/* 2. create server (handle fuse request) */
	server, err := fuse.NewServer(rawFS, dir, &options.MountOptions)
	if err != nil {
		return nil, err
	}
	/* 3. 循环读取处理请求 */
	go server.Serve()
	/* 处理第一个请求 (为了避免竞争) */
	if err := server.WaitMount(); err != nil {
		// we don't shutdown the serve loop. If the mount does
		// not succeed, the loop won't work and exit.
		return nil, err
	}

	return server, nil
}
