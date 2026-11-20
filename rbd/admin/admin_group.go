//go:build ceph_preview && ceph_group_mirror
// +build ceph_preview,ceph_group_mirror

// +ceph:api:skip

package admin

import (
	"fmt"
)

// NewGroupLevelSpec is used to construct a LevelSpec given a pool and
// optional namespace and group names.
func NewGroupLevelSpec(pool, namespace, group string) LevelSpec {
	var s string
	if group != "" && namespace != "" {
		s = fmt.Sprintf("%s/%s/%s", pool, namespace, group)
	} else if group != "" {
		s = fmt.Sprintf("%s/%s", pool, group)
	} else if namespace != "" {
		s = fmt.Sprintf("%s/%s/", pool, namespace)
	} else {
		s = fmt.Sprintf("%s/", pool)
	}
	return LevelSpec{s}
}
