package control

import (
	"os"
	"testing"
)

func TestApp(t *testing.T) {
	token := os.Getenv("ABLY_ACCOUNT_TOKEN")
	client, _, err := NewClient(token)
	if err != nil {
		t.Error(err)
		return
	}
	apps, err := client.Apps()
	if err != nil {
		t.Error(err)
		return
	}
	for _, app := range apps {
		t.Logf("%#v", app)
	}

	apps, err = client.Apps()
	if err != nil {
		t.Error(err)
		return
	}
	var tapp *App
	for _, v := range apps {
		if v.Name == "test" && v.Status != "deleted" {
			tapp = &v
			break
		}
	}

	if tapp != nil {
		_ = client.DeleteApp(tapp.ID)
	}
	app, err := client.CreateApp(&App{Name: "test"})
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("created app: %#v", app)

	a, err := client.UpdateApp(app.ID, &App{Name: "test", TLSOnly: false})
	if err != nil {
		t.Error(err)
		return
	}
	if a.TLSOnly {
		t.Error("shouldnt be tls only")
		return
	}
	a, err = client.UpdateApp(app.ID, &App{Name: "test", TLSOnly: true})
	if err != nil {
		t.Error(err)
		return
	}
	if !a.TLSOnly {
		t.Error("should be tls only")
		return
	}

	err = client.DeleteApp(app.ID)
	if err != nil {
		t.Error(err)
		return
	}

	keys, err := client.Keys(apps[0].ID)
	if err != nil {
		t.Error(err)
		return
	}
	var key *Key
	for _, k := range keys {
		if k.Name == "test" && k.Status == 0 {
			key = &k
		}
		t.Logf("%#v", k)
	}
	if key != nil {
		_, err = client.UpdateKey(key.AppID, key.ID, NewKey{Capability: map[string][]string{"[*]*": []string{"subscribe"}}})
		if err != nil {
			t.Error(err)
			return
		}
	}

	namespaces, err := client.Namespaces(key.AppID)
	if err != nil {
		t.Error(err)
		return
	}
	for _, n := range namespaces {
		t.Logf("%#v", n)

	}
	rules, err := client.Rules(key.AppID)
	if err != nil {
		t.Error(err)
		return
	}
	for _, r := range rules {
		t.Logf("%#v", r)

	}

}
