package gosimplicate

type RegistrationCollection struct {
	registrations []Hours
	projects      map[string][]Hours
}

func NewRegistrationCollection() RegistrationCollection {
	rc := RegistrationCollection{}
	rc.projects = make(map[string][]Hours)
	return rc
}

func (rc *RegistrationCollection) Add(r Hours) {
	rc.registrations = append(rc.registrations, r)

	_, present := rc.projects[r.Project.Name]
	if !present {
		rc.projects[r.Project.Name] = []Hours{}
	}

	rc.projects[r.Project.Name] = append(rc.projects[r.Project.Name], r)
}

func (rc RegistrationCollection) GetByProject() map[string][]Hours {
	return rc.projects
}
