package gosimplicate

type RegistrationCollection struct {
	registrations []Registration
	projects      map[string][]Registration
}

func NewRegistrationCollection() RegistrationCollection {
	rc := RegistrationCollection{}
	rc.projects = make(map[string][]Registration)
	return rc
}

func (rc *RegistrationCollection) Add(r Registration) {
	rc.registrations = append(rc.registrations, r)

	_, present := rc.projects[r.Project.Name]
	if !present {
		rc.projects[r.Project.Name] = []Registration{}
	}

	rc.projects[r.Project.Name] = append(rc.projects[r.Project.Name], r)
}

func (rc RegistrationCollection) GetByProject() map[string][]Registration {
	return rc.projects
}
