package toggl

type TogglClients []*TogglClientElement

type GetClientOpts struct {
	Name string
}

type TogglClientElement struct {
	ID       int64  `json:"id"`
	Wid      int64  `json:"wid"`
	Archived bool   `json:"archived"`
	Name     string `json:"name"`
	At       string `json:"at"`
}

// ByName implements sort.Interface based on the Name field of TogglClientElement.
type ClientsByName TogglClients

func (a ClientsByName) Len() int           { return len(a) }
func (a ClientsByName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ClientsByName) Less(i, j int) bool { return a[i].Name < a[j].Name }
