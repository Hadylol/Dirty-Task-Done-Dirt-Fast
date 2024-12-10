package userstate

import (
	"context"
	"fmt"

	"github.com/looplab/fsm"
)

type UserFSM struct{
	StateMachine *fsm.FSM
	UserID int64
}

func NewUserFSM(userID int64) *UserFSM{
	transitions := fsm.Events{
		{Name: "start", Src:[]string{"idle"},Dst: "waiting_for_url" },
		{Name: "receive_url",Src: []string{"waiting_for_url"}, Dst: "url_shortened"},
		{Name: "rest",Src: []string{"url_shortened","waiting_for_url"},Dst: "idle"},
	}
	callbacks :=fsm.Callbacks{
		"enter_state": func(_ context.Context,e *fsm.Event){
			fmt.Printf("User %d transitioned to state : %s \n",userID,e.Dst)
		},
	}
	return &UserFSM{
		StateMachine: fsm.NewFSM("idle",transitions,callbacks),
		UserID: userID,
	}
}
var userFSMs = make(map[int64]*UserFSM)

func GetUserFSM(userID int64)*UserFSM{
	if UserFSM,exists := userFSMs[userID];exists{
		return UserFSM
	}
	userFSM := NewUserFSM(userID)
	userFSMs[userID]=userFSM
	return userFSM
}