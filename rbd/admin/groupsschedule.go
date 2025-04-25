//go:build ceph_preview
// +build ceph_preview

package admin

import (
	ccom "github.com/ceph/go-ceph/common/commands"
	"github.com/ceph/go-ceph/internal/commands"
)

// GroupSnapshotScheduleAdmin encapsulates management functions for
// ceph rbd mirror group snapshot schedules.
type GroupSnapshotScheduleAdmin struct {
	conn ccom.MgrCommander
}

// GroupSnapshotSchedule returns a GroupSnapshotScheduleAdmin type for
// managing ceph rbd mirror group snapshot schedules.
func (ra *RBDAdmin) GroupSnapshotSchedule() *GroupSnapshotScheduleAdmin {
	return &GroupSnapshotScheduleAdmin{conn: ra.conn}
}

// Add a new group snapshot schedule to the given pool/group based on the supplied
// level spec.
//
// Similar To:
//
//	rbd mirror group snapshot schedule add <level_spec> <interval> <start_time>
func (mss *GroupSnapshotScheduleAdmin) Add(l LevelSpec, i Interval, s StartTime) error {
	m := map[string]string{
		"prefix":     "rbd mirror group snapshot schedule add",
		"level_spec": l.spec,
		"format":     "json",
	}
	if i != NoInterval {
		m["interval"] = string(i)
	}
	if s != NoStartTime {
		m["start_time"] = string(s)
	}
	return commands.MarshalMgrCommand(mss.conn, m).NoData().End()
}
