package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClien(t *testing.T) {
	mockServer := NewSimpleMockServer("testdata/1.json", "application/json", 200)
	defer mockServer.Close()

	// setup client with mock server and tags for command
	scheme, host, port, _ := parseURL(mockServer.URL())
	client := NewClient(scheme, host, port)
	res := Response{}
	client.Get("foo", &res)

	expectedResponse := Response{
		Resolution: Resolution{
			Bin:    "create-react-app",
			Npm:    "create-react-app",
			Github: "facebook/create-react-app",
		},
		Keywords: []string{
			"react",
			"create",
			"boilerplate",
			"scaffold",
			"webpack",
			"project",
		},
		Inputs: Inputs{
			Args: []Arg{
				Arg{
					Description: "App name",
					Type:        "string",
					IsRequired:  true,
				},
			},
			Flags: []Flag{
				Flag{
					Description: "Template name",
					Short:       "",
					Long:        "template",
					Type:        "string",
					IsRequired:  false,
					Default:     "",
				},
			},
		},
		Recipes: []Recipe{
			Recipe{
				Description: "Create React app with no build configuration",
				Keywords:    []string{},
				Inputs: Inputs{
					Args:  []Arg(nil),
					Flags: []Flag(nil),
				},
			},
			Recipe{
				Description: "Create React app with no build configuration using TypeScript",
				Keywords:    []string{"typescript"},
				Inputs: Inputs{
					Args: []Arg{
						Arg{
							Description: "App name",
							Type:        "string",
							IsRequired:  true,
						},
					},
					Flags: []Flag{
						Flag{
							Description: "Template name",
							Short:       "",
							Long:        "template",
							Type:        "string",
							IsRequired:  false,
							Default:     "",
						},
					},
				},
			},
		},
	}
	assert.Equal(t, expectedResponse, res)
}
