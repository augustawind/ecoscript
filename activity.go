package main

type Activity struct {
	ticks       int
	ticksNeeded int
	exec        action
	active      bool
}

func NewActivity() *Activity {
	return new(Activity)
}

func (act *Activity) InProgress() bool {
	return act.active
}

func (act *Activity) Begin(delay int, exec action) (done bool) {
	act.ticks = 0
	act.ticksNeeded = delay
	act.exec = exec
	act.active = true
	return act.Continue()
}

func (act *Activity) Continue() (done bool) {
	act.ticks++
	if act.ticks >= act.ticksNeeded {
		act.exec()
		act.active = false
		done = true
	}
	return
}
