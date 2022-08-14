package redis

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"path"

	"github.com/gadumitrachioaiei/slotserver/slot"

	"github.com/go-redis/redis/v9"
)

//go:embed lua
var luaFS embed.FS

// Users represents the users service.
type Users struct {
	c            *redis.Client
	defaultChips int

	scripts map[string]string
}

// NewUsers returns an initialized users service.
func NewUsers(c *redis.Client, defaultChips int) (*Users, error) {
	u := &Users{c: c, defaultChips: defaultChips, scripts: map[string]string{}}
	if err := u.loadScripts(context.Background()); err != nil {
		return nil, err
	}
	return u, nil
}

// Update updates the user chips with the given amount.
func (u *Users) Update(ctx context.Context, id string, amount int) (slot.User, error) {
	script := u.scripts["update_chips.lua"]
	if script == "" {
		return slot.User{}, errors.New("script for updating user not found")
	}
	v, err := u.c.EvalSha(ctx, script, []string{id}, u.defaultChips, amount).Result()
	if err != nil {
		return slot.User{}, fmt.Errorf("cannot eval script: %v", err)
	}
	chips, ok := v.(int64)
	if !ok {
		return slot.User{}, fmt.Errorf("script does not return an int: %T", v)
	}
	if chips == -1 {
		return slot.User{}, fmt.Errorf("script does not allow updating the user: %v", err)
	}
	return slot.User{ID: id, Chips: int(chips)}, nil
}

func (u *Users) loadScripts(ctx context.Context) error {
	dirPath := "lua"
	entries, err := luaFS.ReadDir(dirPath)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		content, err := luaFS.ReadFile(path.Join(dirPath, entry.Name()))
		if err != nil {
			return fmt.Errorf("cannot read embedded file: %v", err)
		}
		sha, err := u.c.ScriptLoad(ctx, string(content)).Result()
		if err != nil {
			return err
		}
		u.scripts[entry.Name()] = sha
	}
	return nil
}
