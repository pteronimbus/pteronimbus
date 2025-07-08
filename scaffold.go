package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var structure = []string{
	"pteronimbus/apps/backend/cmd/",
	"pteronimbus/apps/backend/internal/",
	"pteronimbus/apps/backend/go.mod",
	"pteronimbus/apps/backend/Dockerfile",

	"pteronimbus/apps/controller/api/",
	"pteronimbus/apps/controller/controllers/",
	"pteronimbus/apps/controller/config/",
	"pteronimbus/apps/controller/cmd/",
	"pteronimbus/apps/controller/go.mod",
	"pteronimbus/apps/controller/Dockerfile",

	"pteronimbus/apps/frontend/public/",
	"pteronimbus/apps/frontend/pages/",
	"pteronimbus/apps/frontend/components/",
	"pteronimbus/apps/frontend/tailwind.config.js",
	"pteronimbus/apps/frontend/package.json",

	"pteronimbus/charts/pteronimbus/",
	"pteronimbus/charts/controller/",

	"pteronimbus/config/crd/",
	"pteronimbus/config/rbac/",
	"pteronimbus/config/samples/",

	"pteronimbus/deployments/examples/single-node-k3s.yaml",

	"pteronimbus/docker/backend.Dockerfile",
	"pteronimbus/docker/controller.Dockerfile",
	"pteronimbus/docker/frontend.Dockerfile",

	"pteronimbus/internal/auth/",
	"pteronimbus/internal/config/",
	"pteronimbus/internal/types/",

	"pteronimbus/proto/handshake.proto",

	"pteronimbus/scripts/generate-token.sh",
	"pteronimbus/scripts/setup-local.sh",

	"pteronimbus/docs/architecture.md",
	"pteronimbus/docs/README.user.md",
	"pteronimbus/docs/README.cursor.md",

	"pteronimbus/Makefile",

	"pteronimbus/.github/workflows/ci.yml",
}

func main() {
	for _, path := range structure {
		if strings.HasSuffix(path, "/") {
			// It's a directory
			if err := os.MkdirAll(path, 0755); err != nil {
				fmt.Printf("Failed to create dir: %s\n", path)
			}
			f, err := os.Create(path + ".keep")
                        if err != nil {
                                fmt.Printf("Failed to create file: %s\n", path)
                        } else {
                                fmt.Printf("Created: %s\n", path)
                                f.Close()
                        }
		} else {
			// It's a file
			dir := filepath.Dir(path)
			if err := os.MkdirAll(dir, 0755); err != nil {
				fmt.Printf("Failed to create dir for file: %s\n", path)
				continue
			}
			f, err := os.Create(path)
			if err != nil {
				fmt.Printf("Failed to create file: %s\n", path)
			} else {
				fmt.Printf("Created: %s\n", path)
				f.Close()
			}
		}
	}
}

