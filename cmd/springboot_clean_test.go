package cmd

import (
	"reflect"
	"testing"

	"github.com/lahat-m/serv-kun/internal/compose"
)

func TestParseComposeVolumes(t *testing.T) {
	got, err := parseComposeVolumes(`{
		"name": "splm-mis",
		"volumes": {
			"postgres-data": {"name": "splm-mis_postgres-data"},
			"certbot-certs": {"name": "splm-mis_certbot-certs"},
			"shared": {"name": "shared", "external": true}
		}
	}`)
	if err != nil {
		t.Fatal(err)
	}
	want := []composeVolume{
		{Key: "certbot-certs", Name: "splm-mis_certbot-certs"},
		{Key: "postgres-data", Name: "splm-mis_postgres-data"},
		{Key: "shared", Name: "shared", External: true},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %+v, want %+v", got, want)
	}
}

func TestParseComposeVolumesNone(t *testing.T) {
	got, err := parseComposeVolumes(`{"name": "x", "services": {}}`)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 0 {
		t.Errorf("got %+v, want no volumes", got)
	}
}

func TestCleanDefaultsKeepEveryVolume(t *testing.T) {
	volumes := []composeVolume{
		{Key: "postgres-data", Name: "p_postgres-data"},
		{Key: "shared", Name: "shared", External: true},
	}
	options, defaults := cleanOptions(volumes)
	if len(options) != 6 {
		t.Errorf("got %d options, want everything + 4 fixed + 1 non-external volume: %q", len(options), options)
	}
	plan := planFromSelection(defaults, volumes)
	want := cleanPlan{Containers: true, Images: true, BuilderCache: true, Dangling: true}
	if !reflect.DeepEqual(plan, want) {
		t.Errorf("default plan = %+v, want %+v", plan, want)
	}
}

func TestPlanFromSelectionPicksVolumes(t *testing.T) {
	data := composeVolume{Key: "postgres-data", Name: "p_postgres-data"}
	certs := composeVolume{Key: "certbot-certs", Name: "p_certbot-certs"}
	plan := planFromSelection([]string{cleanOptContainers, volumeOption(certs)},
		[]composeVolume{data, certs})
	want := cleanPlan{Containers: true, Volumes: []composeVolume{certs}}
	if !reflect.DeepEqual(plan, want) {
		t.Errorf("plan = %+v, want %+v", plan, want)
	}
}

func TestCleanPlanValidate(t *testing.T) {
	vol := []composeVolume{{Key: "postgres-data", Name: "p_postgres-data"}}
	cases := []struct {
		plan cleanPlan
		ok   bool
	}{
		{cleanPlan{BuilderCache: true, Dangling: true}, true},
		{cleanPlan{Containers: true, Images: true, Volumes: vol}, true},
		{cleanPlan{Images: true}, false},
		{cleanPlan{Volumes: vol}, false},
	}
	for _, c := range cases {
		if err := c.plan.validate(); (err == nil) != c.ok {
			t.Errorf("validate(%+v) = %v, want ok=%v", c.plan, err, c.ok)
		}
	}
}

func TestCleanNextStepsBringsUpTheWholeStack(t *testing.T) {
	proj := compose.Project{Dir: "/opt/splm-mis", File: "docker-compose.prod.yml", EnvFile: ".env.prod"}
	got := cleanNextSteps(cleanPlan{Containers: true, Images: true}, proj)
	want := []string{
		"cd /opt/splm-mis",
		"sudo docker compose -f docker-compose.prod.yml --env-file .env.prod up -d",
		"sudo serv-kun springboot status --dir /opt/splm-mis",
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestCleanNextStepsReissuesDeletedCerts(t *testing.T) {
	proj := compose.Project{Dir: "/opt/app", File: "compose.yml", EnvFile: ".env"}
	plan := cleanPlan{Containers: true, Volumes: []composeVolume{{Key: "certbot-certs", Name: "app_certbot-certs"}}}
	got := cleanNextSteps(plan, proj)
	if len(got) != 4 {
		t.Fatalf("got %q, want tls before the three usual steps", got)
	}
	if want := "sudo serv-kun springboot tls --dir /opt/app --compose-file compose.yml --env-file .env"; got[0] != want {
		t.Errorf("first step = %q, want %q", got[0], want)
	}
}

func TestEverythingSelectsEveryOwnedVolume(t *testing.T) {
	data := composeVolume{Key: "postgres-data", Name: "p_postgres-data"}
	certs := composeVolume{Key: "certbot-certs", Name: "p_certbot-certs"}
	shared := composeVolume{Key: "shared", Name: "shared", External: true}
	volumes := []composeVolume{data, certs, shared}

	plan := planFromSelection([]string{cleanOptEverything}, volumes)
	want := cleanPlan{Containers: true, Images: true, BuilderCache: true, Dangling: true,
		Volumes: []composeVolume{data, certs}}
	if !reflect.DeepEqual(plan, want) {
		t.Errorf("plan = %+v, want %+v", plan, want)
	}
	if !plan.isFull(volumes) {
		t.Error("everything should count as a full reset")
	}
	if err := plan.validate(); err != nil {
		t.Errorf("validate() = %v", err)
	}
}

func TestIsFullNeedsEveryVolume(t *testing.T) {
	data := composeVolume{Key: "postgres-data", Name: "p_postgres-data"}
	certs := composeVolume{Key: "certbot-certs", Name: "p_certbot-certs"}
	volumes := []composeVolume{data, certs}

	partial := cleanPlan{Containers: true, Images: true, BuilderCache: true, Dangling: true,
		Volumes: []composeVolume{data}}
	if partial.isFull(volumes) {
		t.Error("a plan that keeps certbot-certs is not a full reset")
	}
	if !fullCleanPlan(volumes).isFull(volumes) {
		t.Error("fullCleanPlan should be a full reset")
	}
}
