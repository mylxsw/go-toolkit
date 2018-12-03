package container_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/mylxsw/go-toolkit/container"
)

type GetUserInterface interface {
	GetUser() string
}

type UserService struct {
	repo *UserRepo
}

func (u *UserService) GetUser() string {
	return fmt.Sprintf("get user from connection: %s", u.repo.connStr)
}

type UserRepo struct {
	connStr string
}

var expectedValue = "get user from connection: root:root@/my_db?charset=utf8"

func TestPrototype(t *testing.T) {
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

	if err := c.Invoke(func(userService *UserService) {
		if userService.GetUser() != expectedValue {
			t.Error("test failed")
		}
	}); err != nil {
		t.Errorf("test failed: %s", err)
	}

	userService, err := c.Get(reflect.TypeOf((*UserService)(nil)))
	if err != nil {
		t.Error(err)
	}

	if userService.(*UserService).GetUser() != expectedValue {
		t.Error("test failed")
	}
}

func TestInterfaceInjection(t *testing.T) {
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

	if err := c.Invoke(func(userService GetUserInterface) {
		if userService.GetUser() != expectedValue {
			t.Error("test failed")
		}
	}); err != nil {
		t.Errorf("test failed: %s", err)
	}
}
