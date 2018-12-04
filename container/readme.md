# Container 


A container for golang dependency injection.

	c := container.New()

	c.BindValue("conn_str", "root:root@/my_db?charset=utf8")
	c.Singleton(func(c *container.Container) (*UserRepo, error) {
		connStr, err := c.Get("conn_str")
		if err != nil {
			return nil, err
		}

		return &UserRepo{connStr: connStr.(string)}, nil
	})
	c.Prototype(func(userRepo *UserRepo) (*UserService, error) {
		return &UserService{repo: userRepo}, nil
	})

	if err := c.Resolve(func(userService *UserService) {
		if userService.GetUser() != expectedValue {
			t.Error("test failed")
		}
	}); err != nil {
		panic(err)
	}

	userService, err := c.Get((*UserService)(nil))
	if err != nil {
		panic(err)
	}

	if userService.(*UserService).GetUser() != expectedValue {
		panic(err)
	}