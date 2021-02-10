package app

// GetSudoers returns a map of sudo users and their public ssh keys for a matching server groups
func (c *Config) GetSudoers(serverGroups []string) map[string]string {
	return c.getUniqueUsers(serverGroups, true)
}

// GetUsers returns a map of users and their public ssh keys for a matching server groups
func (c *Config) GetUsers(serverGroups []string) map[string]string {
	return c.getUniqueUsers(serverGroups, false)
}

func (c *Config) getUniqueUsers(serverGroups []string, sudoers bool) map[string]string {
	users := make([]string, 0)
	for _, group := range serverGroups {
		if sudoers {
			users = combineUnique(users, c.ServerGroups[group].Sudoers)
		} else {
			users = combineUnique(users, c.ServerGroups[group].Users)
		}
	}

	output := make(map[string]string)
	for _, user := range users {
		output[user] = c.Users[user]
	}

	return output
}

func combineUnique(a []string, b []string) []string {
	check := make(map[string]int)
	d := append(a, b...)
	res := make([]string, 0)

	for _, val := range d {
		check[val] = 1
	}

	for letter := range check {
		res = append(res, letter)
	}

	return res
}
