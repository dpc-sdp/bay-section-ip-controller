package util_test

import (
	"testing"

	"github.com/dpc-sdp/bay-section-ip-controller/internal/util"
	"github.com/stretchr/testify/assert"
)

func TestIsActionAbleEnvironment(t *testing.T) {
	s := util.Section{ActionableEnvironments: []string{"Develop"}}

	assert.True(t, s.IsActionableEnvironment("Develop"))
	assert.False(t, s.IsActionableEnvironment("develop"))
	assert.False(t, s.IsActionableEnvironment("Master"))
	assert.False(t, s.IsActionableEnvironment("Production"))

	s = util.Section{ActionableEnvironments: []string{"Develop", "Master"}}

	assert.True(t, s.IsActionableEnvironment("Develop"))
	assert.True(t, s.IsActionableEnvironment("Master"))
	assert.False(t, s.IsActionableEnvironment("Production"))

	s = util.Section{ActionableEnvironments: []string{}}

	assert.True(t, s.IsActionableEnvironment("Develop"))
	assert.True(t, s.IsActionableEnvironment("Master"))
	assert.True(t, s.IsActionableEnvironment("Production"))

	s = util.Section{ActionableEnvironments: []string{""}}

	assert.True(t, s.IsActionableEnvironment("Develop"))
	assert.True(t, s.IsActionableEnvironment("Master"))
	assert.True(t, s.IsActionableEnvironment("Production"))
}

func TestIsActionableApplication(t *testing.T) {
	s := util.Section{ActionableApplications: []string{"www.test.com"}}

	assert.True(t, s.IsActionableApplication("www.test.com"))
	assert.False(t, s.IsActionableApplication("www.test2.com"))
	assert.False(t, s.IsActionableApplication("www.test.io"))

	s = util.Section{ActionableApplications: []string{}}

	assert.True(t, s.IsActionableApplication("www.test.com"))
	assert.True(t, s.IsActionableApplication("www.test2.com"))
	assert.True(t, s.IsActionableApplication("www.test.io"))

	s = util.Section{ActionableApplications: []string{""}}

	assert.True(t, s.IsActionableApplication("www.test.com"))
	assert.True(t, s.IsActionableApplication("www.test2.com"))
	assert.True(t, s.IsActionableApplication("www.test.io"))
}

func TestIsActionableAccount(t *testing.T) {
	s := util.Section{ActionableAccounts: []string{"1000"}}

	assert.True(t, s.IsActionableAccount(int32(1000)))
	assert.False(t, s.IsActionableAccount(int32(1001)))

	s = util.Section{ActionableAccounts: []string{}}

	assert.True(t, s.IsActionableAccount(int32(1000)))
	assert.True(t, s.IsActionableAccount(int32(1001)))

	s = util.Section{ActionableAccounts: []string{""}}

	assert.True(t, s.IsActionableAccount(int32(1000)))
	assert.True(t, s.IsActionableAccount(int32(1001)))
}
