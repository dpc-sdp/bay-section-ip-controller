package util

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

type SectionAPI struct {
	Username               string
	Token                  string
	Host                   string
	AccountId              string
	BlockedIps             SectionIpRestrictionSchema
	ActionableEnvironments []string
	ActionableApplications []string
	Applications           []SectionApp
	Client                 *http.Client
}

type BlockedIP struct {
	IP   string
	Date time.Time
}

type SectionAccount struct {
	Id          int    `json:"id"`
	Name        string `json:"account_name"`
	Requires2FA bool   `json:"requires_2fa"`
	Owner       struct {
		Id       int    `json:"id"`
		Email    string `json:"email"`
		Verified bool   `json:"verified"`
	} `json:"owner"`
}

type SectionApp struct {
	Id              int    `json:"id"`
	Href            string `json:"href"`
	ApplicationName string `json:"application_name"`
	Environments    map[string]SectionEnvironment
}

type SectionEnvironment struct {
	Id          int    `json:"id"`
	Name        string `json:"environment_name"`
	IpBlocklist SectionIpRestrictionSchema
}

type SectionIpRestrictionSchema struct {
	IpBlacklist []string `json:"ip_blacklist"`
}

type Log struct {
	Status          string   `json:"status"`
	Message         string   `json:"message"`
	ApplicationName string   `json:"application,omitempty"`
	EnvironmentName string   `json:"environment,omitempty"`
	IPList          []string `json:"iplist,omitempty"`
}

func (section *SectionAPI) Init() (bool, error) {
	section.GetApplications()

	for i := range section.Applications {
		section.GetApplicationEnvironments(&section.Applications[i])
	}

	msg, _ := json.Marshal(Log{
		Status:  "success",
		Message: "Initialised application list.",
	})

	fmt.Println(string(msg))

	return true, nil
}

func (section *SectionAPI) GetAccount() SectionAccount {
	url := fmt.Sprintf("%s/account/%s", section.Host, section.AccountId)
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.SetBasicAuth(section.Username, section.Token)
	resp, err := section.Client.Do(req)

	var account SectionAccount

	if err != nil {
		return account
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return account
	}

	json.Unmarshal(body, &account)
	return account
}

func (section *SectionAPI) GetApplications() (bool, error) {
	url := fmt.Sprintf("%s/account/%s/application", section.Host, section.AccountId)
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.SetBasicAuth(section.Username, section.Token)
	resp, err := section.Client.Do(req)

	if err != nil {
		msg, _ := json.Marshal(Log{
			Status:  "fail",
			Message: string(err.Error()),
		})
		fmt.Println(string(msg))
		return false, err
	}

	if resp.StatusCode != http.StatusOK {
		msg, _ := json.Marshal(Log{
			Status:  "fail",
			Message: fmt.Sprintf("Invalid API status (%d)", resp.StatusCode),
		})
		fmt.Println(string(msg))
		return false, errors.New("Unexpected API Response")
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		msg, _ := json.Marshal(Log{
			Status:  "fail",
			Message: string(err.Error()),
		})
		fmt.Println(string(msg))
		return false, err
	}

	err = json.Unmarshal(body, &section.Applications)

	if err != nil {
		msg, _ := json.Marshal(Log{
			Status:  "fail",
			Message: string(err.Error()),
		})
		fmt.Println(string(msg))
		return false, err
	}

	return true, nil
}

func (section *SectionAPI) GetApplicationEnvironments(app *SectionApp) (bool, error) {
	url := fmt.Sprintf("%s/account/%s/application/%s/environment", section.Host, section.AccountId, strconv.Itoa(app.Id))
	envReq, _ := http.NewRequest(http.MethodGet, url, nil)
	envReq.SetBasicAuth(section.Username, section.Token)
	envResp, err := section.Client.Do(envReq)

	if err != nil {
		msg, _ := json.Marshal(Log{
			Status:  "fail",
			Message: fmt.Sprintf("Unable to fetch environments for %s", app.ApplicationName),
		})
		fmt.Println(string(msg))
		return false, errors.New(string(msg))
	}

	defer envResp.Body.Close()
	eb, err := ioutil.ReadAll(envResp.Body)

	if err != nil {
		msg, _ := json.Marshal(Log{
			Status:  "fail",
			Message: fmt.Sprintf("Unable to fetch environments for %s", app.ApplicationName),
		})
		fmt.Println(string(msg))
		return false, errors.New(string(msg))
	}

	var envs []SectionEnvironment
	json.Unmarshal(eb, &envs)

	app.Environments = make(map[string]SectionEnvironment)

	for _, e := range envs {
		e.IpBlocklist, _ = section.GetIPBlocklist(app, &e)
		app.Environments[e.Name] = e
	}

	return true, nil
}

func (section *SectionAPI) GetIPBlocklist(app *SectionApp, env *SectionEnvironment) (SectionIpRestrictionSchema, error) {
	var ips SectionIpRestrictionSchema
	ipRestrictionAddr := fmt.Sprintf("%s/account/%s/application/%s/environment/%s/ipRestrictions", section.Host, section.AccountId, strconv.Itoa(app.Id), env.Name)
	req, _ := http.NewRequest(http.MethodGet, ipRestrictionAddr, nil)
	req.SetBasicAuth(section.Username, section.Token)
	resp, err := section.Client.Do(req)

	if err != nil {
		return ips, err
	}

	defer resp.Body.Close()
	e, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(e, &ips)
	return ips, nil
}

func (section *SectionAPI) AddAllIPBlocklist(i SectionIpRestrictionSchema) (bool, error) {
	for _, app := range section.Applications {
		action := false
		if len(section.ActionableApplications) > 0 {
			for _, a := range section.ActionableApplications {
				if app.ApplicationName == a {
					action = true
					break
				}
			}
		} else {
			action = true
		}

		if !action {
			msg, _ := json.Marshal(Log{
				Status:  "info",
				Message: fmt.Sprintf("Skipping %s due to ActionableApplications exclusion", app.ApplicationName),
			})
			fmt.Println(string(msg))
			continue
		}

		for _, e := range section.ActionableEnvironments {
			env, exists := app.Environments[e]
			if !exists {
				continue
			}
			go section.AddIPBlocklist(app, env, i)
		}

	}
	return true, nil
}

func (section *SectionAPI) AddIPBlocklist(app SectionApp, env SectionEnvironment, i SectionIpRestrictionSchema) (bool, error) {
	payload, _ := json.Marshal(i)
	ipRestrictionAddr := fmt.Sprintf("%s/account/%s/application/%s/environment/%s/ipRestrictions", section.Host, section.AccountId, strconv.Itoa(app.Id), env.Name)
	req, _ := http.NewRequest(http.MethodPost, ipRestrictionAddr, bytes.NewBuffer(payload))
	req.SetBasicAuth(section.Username, section.Token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := section.Client.Do(req)
	if err != nil {
		msg, _ := json.Marshal(Log{
			Status:  "fail",
			Message: string(err.Error()),
		})
		fmt.Println(string(msg))
		return false, errors.New(string(msg))
	}
	if resp.StatusCode == http.StatusUnauthorized {
		msg, _ := json.Marshal(Log{
			Status:  "fail",
			Message: fmt.Sprintf("addIPBlocklist: unauthorised (%s:%s)", app.ApplicationName, env.Name),
		})
		fmt.Println(string(msg))
		return false, errors.New(string(msg))
	}
	if resp.StatusCode != http.StatusOK {
		msg, _ := json.Marshal(Log{
			Status:  "fail",
			Message: fmt.Sprintf("addIPBlocklist: unauthorised (%s:%s)", app.ApplicationName, env.Name),
		})
		fmt.Println(string(msg))
		return false, errors.New(string(msg))
	}
	msg, _ := json.Marshal(Log{
		Status:          "success",
		Message:         fmt.Sprintf("Updated IP blocklist"),
		ApplicationName: app.ApplicationName,
		EnvironmentName: env.Name,
		IPList:          i.IpBlacklist,
	})
	fmt.Println(string(msg))

	return true, nil
}
