package cluster

import (
	"fmt"
	"time"
)

func (rs *RaftSettings) handleVoteRequest(req struct {
	CandidateIdentity string `json:"candidate_identity"`
	CandidateEpoch    int64  `json:"candidate_epoch"`
	CandidateLastLog  int64  `json:"candidate_last_log"`
}) int {
	rs.last_heartbeat = time.Now()

	var log_forward_v1 bool = (req.CandidateEpoch > rs.max_epoch)
	var log_forward_v2 bool = (req.CandidateEpoch == rs.max_epoch && req.CandidateLastLog >= rs.self_last_log)
	var log_forward_v3 bool = (log_forward_v1 || log_forward_v2)

	var can_vote_v1 bool = (rs.voted_for == "" || rs.voted_for == req.CandidateIdentity)
	var can_vote_v2 bool = (req.CandidateEpoch > rs.max_epoch)
	var can_vote_v3 bool = (can_vote_v1 || can_vote_v2)

	if log_forward_v3 && can_vote_v3 {
		rs.voted_for = req.CandidateIdentity
		rs.max_epoch = req.CandidateEpoch
		rs.leader_identity = ""
		return RaftVoteGranted
	}
	return RaftVoteDenied
}

func (rs *RaftSettings) handleLogSyncRequest(req struct {
	LeaderIdentity string `json:"leader_identity"`
	LeaderEpoch    int64  `json:"leader_epoch"`
	NewEvent       string `json:"leader_last_event"`
	LeaderLastLog  int64  `json:"leader_last_log"`
}) int {
	rs.last_heartbeat = time.Now()

	if req.LeaderEpoch < rs.max_epoch {
		return RaftLeaderLaggingBehind
	}

	if req.LeaderEpoch > rs.max_epoch {
		rs.max_epoch = req.LeaderEpoch
		rs.voted_for = ""
	}

	rs.leader_identity = req.LeaderIdentity
	var log_diff int64 = req.LeaderLastLog - rs.self_last_log

	if log_diff == 1 || log_diff == 0 {
		if log_diff != 0 {
			rs.save_event_callback(req.NewEvent)
		}
		rs.self_last_log = req.LeaderLastLog
		fmt.Print(rs.GetLeaderDescriptor())
		return RaftSyncAccepted
	}

	if log_diff < 0 {
		return RaftLeaderLaggingBehind
	}

	return RaftFollowerLaggingBehind
}

func (rs *RaftSettings) handleNewEventRequest(req struct {
	NewEvent string `json:"new_event"`
}) int {
	if rs.leader_identity == rs.self_identity {
		rs.save_event_callback(req.NewEvent)
		rs.self_last_log += 1
		return RaftEventWritten
	}
	return RaftEventIgnored
}
