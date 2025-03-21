package rbd

// #cgo LDFLAGS: -lrbd
// #include <errno.h>
// #include <stdlib.h>
// #include <rados/librados.h>
// #include <rbd/librbd.h>
import "C"

import (
	"fmt"
	"unsafe"

	"github.com/ceph/go-ceph/internal/cutil"
	"github.com/ceph/go-ceph/rados"
)

// MirrorGroupEnable will enable mirroring for a group using the specified mode.
//
// Implements:
//
//	int rbd_mirror_group_enable(rados_ioctx_t p, const char *name,
//	  							rbd_mirror_image_mode_t mirror_image_mode)
func MirrorGroupEnable(groupIoctx *rados.IOContext, groupName string, mode ImageMirrorMode) error {
	cGroupName := C.CString(groupName)
	defer C.free(unsafe.Pointer(cGroupName))

	ret := C.rbd_mirror_group_enable(
		cephIoctx(groupIoctx),
		cGroupName,
		C.rbd_mirror_image_mode_t(mode))

	return getError(ret)
}

// MirrorGroupDisable will disable mirroring for a group.
//
// Implements:
//
//	int rbd_mirror_group_disable(rados_ioctx_t p, const char *name,
//	  							bool force)
func MirrorGroupDisable(groupIoctx *rados.IOContext, groupName string, force bool) error {
	cGroupName := C.CString(groupName)
	defer C.free(unsafe.Pointer(cGroupName))

	ret := C.rbd_mirror_group_disable(
		cephIoctx(groupIoctx),
		cGroupName,
		C.bool(force))

	return getError(ret)
}

// MirrorGroupPromote will promote the mirrored group to primary state.
//
// Implements:
//
//	int rbd_mirror_group_promote(rados_ioctx_t p, const char *name,
//	  							 bool force)
func MirrorGroupPromote(groupIoctx *rados.IOContext, groupName string, force bool) error {
	cGroupName := C.CString(groupName)
	defer C.free(unsafe.Pointer(cGroupName))

	ret := C.rbd_mirror_group_promote(
		cephIoctx(groupIoctx),
		cGroupName,
		C.bool(force))

	return getError(ret)
}

// MirrorGroupDemote will demote the mirrored group to secondary state.
//
// Implements:
//
//	int rbd_mirror_group_demote(rados_ioctx_t p, const char *name)
func MirrorGroupDemote(groupIoctx *rados.IOContext, groupName string) error {
	cGroupName := C.CString(groupName)
	defer C.free(unsafe.Pointer(cGroupName))

	ret := C.rbd_mirror_group_demote(
		cephIoctx(groupIoctx),
		cGroupName)

	return getError(ret)
}

// MirrorGroupResync is used to manually resolve split-brain status by triggering
// resynchronization.
//
// Implements:
//
//	int rbd_mirror_group_resync(rados_ioctx_t p, const char *name)
func MirrorGroupResync(groupIoctx *rados.IOContext, groupName string) error {
	cGroupName := C.CString(groupName)
	defer C.free(unsafe.Pointer(cGroupName))

	ret := C.rbd_mirror_group_resync(
		cephIoctx(groupIoctx),
		cGroupName)

	return getError(ret)
}

// MirrorGroupState represents the current state of the mirrored group.
type MirrorGroupState C.rbd_mirror_group_state_t

// String representation of MirrorGroupState.
func (mgs MirrorGroupState) String() string {
	switch mgs {
	case MirrorGroupEnabled:
		return "enabled"
	case MirrorGroupDisabled:
		return "disabled"
	case MirrorGroupEnabling:
		return "enabling"
	case MirrorGroupDisabling:
		return "disabling"
	default:
		return "<unknown>"
	}
}

const (
	// MirrorGroupDisabling is the representation of
	// RBD_MIRROR_GROUP_DISABLING from librbd.
	MirrorGroupDisabling = MirrorGroupState(C.RBD_MIRROR_GROUP_DISABLING)
	// MirrorGroupEnabling is the representation of
	// RBD_MIRROR_GROUP_ENABLING from librbd
	MirrorGroupEnabling = MirrorGroupState(C.RBD_MIRROR_GROUP_ENABLING)
	// MirrorGroupEnabled is the representation of
	// RBD_MIRROR_IMAGE_ENABLED from librbd.
	MirrorGroupEnabled = MirrorGroupState(C.RBD_MIRROR_GROUP_ENABLED)
	// MirrorGroupDisabled is the representation of
	// RBD_MIRROR_GROUP_DISABLED from librbd.
	MirrorGroupDisabled = MirrorGroupState(C.RBD_MIRROR_GROUP_DISABLED)
)

// MirrorGroupInfo represents the mirroring status information of group.
type MirrorGroupInfo struct {
	GlobalID        string
	State           MirrorGroupState
	MirrorImageMode ImageMirrorMode
	Primary         bool
}

// GetMirrorGroupInfo returns the mirroring info of the group.
//
// Implements:
//
//	int rbd_mirror_group_get_info(rados_ioctx_t p, const char *name,
//								  rbd_mirror_group_info_t *mirror_group_info,
//								  size_t info_size)
func GetMirrorGroupInfo(groupIoctx *rados.IOContext, groupName string) (*MirrorGroupInfo, error) {
	var cgInfo C.rbd_mirror_group_info_t
	cGroupName := C.CString(groupName)
	defer C.free(unsafe.Pointer(cGroupName))

	ret := C.rbd_mirror_group_get_info(
		cephIoctx(groupIoctx),
		cGroupName,
		&cgInfo,
		C.sizeof_rbd_mirror_group_info_t)

	if err := getError(ret); err != nil {
		return &MirrorGroupInfo{}, err
	}

	info := convertMirrorGroupInfo(&cgInfo)

	// free C memory allocated by C.rbd_mirror_group_get_info call
	C.rbd_mirror_group_get_info_cleanup(&cgInfo)
	return &info, nil
}

func convertMirrorGroupInfo(cgInfo *C.rbd_mirror_group_info_t) MirrorGroupInfo {
	return MirrorGroupInfo{
		GlobalID:        C.GoString(cgInfo.global_id),
		MirrorImageMode: ImageMirrorMode(cgInfo.mirror_image_mode),
		State:           MirrorGroupState(cgInfo.state),
		Primary:         bool(cgInfo.primary),
	}
}

// MirrorGroupStatusState is used to indicate the state of a mirrored group
// within the site status info.
type MirrorGroupStatusState int64

const (
	// MirrorGroupStatusStateUnknown is equivalent to MIRROR_GROUP_STATUS_STATE_UNKNOWN
	MirrorGroupStatusStateUnknown = MirrorGroupStatusState(C.MIRROR_GROUP_STATUS_STATE_UNKNOWN)
	// MirrorGroupStatusStateError is equivalent to MIRROR_GROUP_STATUS_STATE_ERROR
	MirrorGroupStatusStateError = MirrorGroupStatusState(C.MIRROR_GROUP_STATUS_STATE_ERROR)
	// MirrorGroupStatusStateStartingReplay is equivalent to MIRROR_GROUP_STATUS_STATE_STARTING_REPLAY
	MirrorGroupStatusStateStartingReplay = MirrorGroupStatusState(C.MIRROR_GROUP_STATUS_STATE_STARTING_REPLAY)
	// MirrorGroupStatusStateReplaying is equivalent to MIRROR_GROUP_STATUS_STATE_REPLAYING
	MirrorGroupStatusStateReplaying = MirrorGroupStatusState(C.MIRROR_GROUP_STATUS_STATE_REPLAYING)
	// MirrorGroupStatusStateStoppingReplay is equivalent to MIRROR_GROUP_STATUS_STATE_STOPPING_REPLAY
	MirrorGroupStatusStateStoppingReplay = MirrorGroupStatusState(C.MIRROR_GROUP_STATUS_STATE_STOPPING_REPLAY)
	// MirrorGroupStatusStateStopped is equivalent to MIRROR_GROUP_STATUS_STATE_STOPPED
	MirrorGroupStatusStateStopped = MirrorGroupStatusState(C.MIRROR_GROUP_STATUS_STATE_STOPPED)
)

// String represents the MirrorImageStatusState as a short string.
func (state MirrorGroupStatusState) String() string {
	switch state {
	case MirrorGroupStatusStateUnknown:
		return "unknown"
	case MirrorGroupStatusStateError:
		return "error"
	case MirrorGroupStatusStateStartingReplay:
		return "starting_replay"
	case MirrorGroupStatusStateReplaying:
		return "replaying"
	case MirrorGroupStatusStateStoppingReplay:
		return "stopping_replay"
	case MirrorGroupStatusStateStopped:
		return "stopped"
	default:
		return fmt.Sprintf("unknown(%d)", state)
	}
}

// SiteMirrorGroupStatus contains information pertaining to the status of
// a mirrored group within a site.
type SiteMirrorGroupStatus struct {
	MirrorUUID           string
	State                MirrorGroupStatusState
	MirrorImageCount     int
	MirrorImagePoolIDs   int64
	MirrorImageGlobalIDs string
	MirrorImages         []SiteMirrorImageStatus
	Description          string
	LastUpdate           int64
	Up                   bool
}

// GlobalMirrorGroupStatus contains information pertaining to the global
// status of a mirrored group. It contains general information as well
// as per-site information stored in the SiteStatuses slice.
type GlobalMirrorGroupStatus struct {
	Name              string
	Info              MirrorGroupInfo
	SiteStatusesCount int
	SiteStatuses      []SiteMirrorGroupStatus
}

// LocalStatus returns one SiteMirrorGroupStatus item from the SiteStatuses
// slice that corresponds to the local site's status. If the local status
// is not found than the error ErrNotExist will be returned.
func (gmis GlobalMirrorGroupStatus) LocalStatus() (SiteMirrorGroupStatus, error) {
	for i := range gmis.SiteStatuses {
		// I couldn't find it explicitly documented, but a site mirror uuid
		// of an empty string indicates that this is the local site.
		// This pattern occurs in both the pybind code and ceph c++.
		if gmis.SiteStatuses[i].MirrorUUID == "" {
			return gmis.SiteStatuses[i], nil
		}
	}
	return SiteMirrorGroupStatus{}, ErrNotExist
}

type groupSiteArray [cutil.MaxIdx]C.rbd_mirror_group_site_status_t

// GetGlobalMirrorGroupStatus returns status information pertaining to the state
// of a groups's mirroring.
//
// Implements:
//
//	int rbd_mirror_group_get_global_status(
//		rados_ioctx_t p,
//		const char *group_name
//		rbd_mirror_group_global_status_t *mirror_group_status,
//		size_t status_size);
func GetGlobalMirrorGroupStatus(ioctx *rados.IOContext, groupName string) (GlobalMirrorGroupStatus, error) {
	s := C.rbd_mirror_group_global_status_t{}
	cGroupName := C.CString(groupName)
	defer C.free(unsafe.Pointer(cGroupName))

	ret := C.rbd_mirror_group_get_global_status(
		cephIoctx(ioctx),
		(*C.char)(cGroupName),
		&s,
		C.sizeof_rbd_mirror_group_global_status_t)
	if err := getError(ret); err != nil {
		return GlobalMirrorGroupStatus{}, err
	}

	status := newGlobalMirrorGroupStatus(&s)
	return status, nil
}

func newGlobalMirrorGroupStatus(s *C.rbd_mirror_group_global_status_t) GlobalMirrorGroupStatus {
	// Initializing the status object
	status := GlobalMirrorGroupStatus{
		Name:              C.GoString(s.name),
		Info:              convertMirrorGroupInfo(&s.info),
		SiteStatusesCount: int(s.site_statuses_count),
		SiteStatuses:      make([]SiteMirrorGroupStatus, s.site_statuses_count),
	}

	// Check if site statuses are not null before using them
	if s.site_statuses != nil && s.site_statuses_count > 0 {
		gsscs := (*groupSiteArray)(unsafe.Pointer(s.site_statuses))[:s.site_statuses_count:s.site_statuses_count]
		for i := C.uint32_t(0); i < s.site_statuses_count; i++ {
			gss := gsscs[i]
			// Ensure that fields are valid before using them
			if gss.mirror_uuid != nil && gss.mirror_image_global_ids != nil {
				status.SiteStatuses[i] = SiteMirrorGroupStatus{
					MirrorUUID:           C.GoString(gss.mirror_uuid),
					MirrorImageGlobalIDs: C.GoString(*gss.mirror_image_global_ids),
					MirrorImagePoolIDs:   int64(*gss.mirror_image_pool_ids),
					State:                MirrorGroupStatusState(gss.state),
					Description:          C.GoString(gss.description),
					MirrorImageCount:     int(gss.mirror_image_count),
					LastUpdate:           int64(gss.last_update),
					MirrorImages:         make([]SiteMirrorImageStatus, gss.mirror_image_count),
					Up:                   bool(gss.up),
				}

				// Check if the mirror_images pointer is valid
				if gss.mirror_images != nil && gss.mirror_image_count > 0 {

					sscs := (*siteArray)(unsafe.Pointer(gss.mirror_images))[:gss.mirror_image_count:gss.mirror_image_count]
					for j := C.uint32_t(0); j < gss.mirror_image_count; j++ {
						ss := sscs[j]

						// Ensure that fields are valid before using them
						if ss.mirror_uuid != nil {
							status.SiteStatuses[i].MirrorImages[j] = SiteMirrorImageStatus{
								MirrorUUID:  C.GoString(ss.mirror_uuid),
								State:       MirrorImageStatusState(ss.state),
								Description: C.GoString(ss.description),
								LastUpdate:  int64(ss.last_update),
								Up:          bool(ss.up),
							}
						}
					}
				}
			}
		}
	}
	return status
}
