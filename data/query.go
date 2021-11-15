package data

func QueryLeaveByName(name string) (*[]Leave, error) {
	var leaves []Leave
	if errs := db.Find(leaves, "name = ?", name).GetErrors(); len(errs) != 0 {
		return nil, errs2err(errs)
	}
	return &leaves, nil
}
