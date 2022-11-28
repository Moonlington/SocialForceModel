package main

// Behavior defines the behavior of a person.
type Behavior interface {
	GetGoal() *Goal
	Update(p *Person, dt float64)
}

// FollowerBehavior defines the behavior of a person that follows a person.
type FollowerBehavior struct {
	Target *Person
	Goal   *Goal
}

// NewFollowerBehavior creates a new follow behavior.
func NewFollowerBehavior(target *Person) *FollowerBehavior {
	return &FollowerBehavior{Target: target, Goal: newGoal(target.Position, 0, 0)}
}

// GetGoal returns the goal of the behavior.
func (b *FollowerBehavior) GetGoal() *Goal {
	return b.Goal
}

// Update updates the behavior.
func (b *FollowerBehavior) Update(p *Person, dt float64) {
	b.Goal = newGoal(b.Target.Position, 0, 0)
	if p.Position.To(b.Target.Position).Len() <= 1.5*(p.Radius+b.Target.Radius) {
		b.Goal = nil
	}
	return
}

// PathBehavior defines the behavior of a person that follows a path.
type PathBehavior struct {
	Path        *Path
	CurrentGoal *Goal
	LastGoal    *Goal
	Loitered    float64
}

// NewPathBehavior creates a new path behavior.
func NewPathBehavior(path *Path) *PathBehavior {
	return &PathBehavior{Path: path, CurrentGoal: path.NextGoal()}
}

// GetGoal returns the goal of the behavior.
func (b *PathBehavior) GetGoal() *Goal {
	return b.CurrentGoal
}

// Update updates the behavior.
func (b *PathBehavior) Update(p *Person, dt float64) {
	if b.CurrentGoal != nil {
		if p.Position.To(b.CurrentGoal.Target).Len() > b.CurrentGoal.Range {
			return
		}

		if b.CurrentGoal.LoiterAfter != 0 {
			b.Loitered += dt
			if b.Loitered <= b.CurrentGoal.LoiterAfter {
				return
			}
		}
	}

	if b.Path.IsEmpty() {
		if b.LastGoal == nil && b.CurrentGoal != nil {
			b.LastGoal = b.CurrentGoal
		}
		b.CurrentGoal = nil
		// If the person is far enough away from their last goal, set their current goal to their last goal
		if b.LastGoal != nil && p.Position.To(b.LastGoal.Target).Len() > b.LastGoal.Range*2 {
			b.CurrentGoal = b.LastGoal
		}
		return
	}
	b.CurrentGoal = b.Path.NextGoal()
}
