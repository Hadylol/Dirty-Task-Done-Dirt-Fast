package states

var UsersStates = make(map[int64]string)

func UpdateUserState(id int64, state string) {

	UsersStates[id] = state

}
