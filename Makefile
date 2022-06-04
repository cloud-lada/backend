TAG?=latest

kustomize:
	cd manifests && kustomize edit set image ghcr.io/cloud-lada/persistor:$(TAG)
	cd manifests && kustomize edit set image ghcr.io/cloud-lada/ingestor:$(TAG)
	cd manifests && kustomize edit set image ghcr.io/cloud-lada/dumper:$(TAG)
	cd manifests && kustomize edit set image ghcr.io/cloud-lada/api:$(TAG)

	kustomize build manifests -o deploy.yaml

release:
	goreleaser release --rm-dist
