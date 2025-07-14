OPENCONFIG_VERSION = v5.2.0
OPENCONFIG_DIR = ../openconfig-public
YANG_DIR = yangModels
REPO_DIR = goModels
GEN_DIR = $(REPO_DIR)/$(PACKAGE_NAME)
PACKAGE_NAME = $(subst -,_,$(TARGET_MODEL))
TARGET_MODEL = openconfig-network-instance


install-ygot:
	go install github.com/openconfig/ygot/generator@latest

clone-openconfig:
	git clone https://github.com/openconfig/public.git $(OPENCONFIG_DIR)
	cd $(OPENCONFIG_DIR) && git checkout tags/$(OPENCONFIG_VERSION)

clone-ietf:
	# git clone https://github.com/YangModels/yang.git ../yang
	cp ../yang/standard/ietf/RFC/*.yang $(YANG_DIR)/

prepare-yang:
	mkdir -p yangModels
	find ../openconfig-public/release/models -name '*.yang' -exec cp {} yangModels/ \;

check-imports:
	@echo "Checking for missing YANG modules..."
	@for mod in $(shell grep -h "^import " $(YANG_DIR)/*.yang | awk '{print $$2}' | sort -u); do \
		if [ ! -f $(YANG_DIR)/$$mod.yang ]; then \
			echo "Missing: $$mod.yang"; \
		fi \
	done
 

validate-yang:
	find yangModels -name '*.yang' -exec goyang -path yangModels {} \;

generate-models:
	mkdir -p $(GEN_DIR)
	generator \
		-path=$(YANG_DIR) \
		-output_file=$(GEN_DIR)/$(PACKAGE_NAME).go \
		-path_structs_output_file=$(GEN_DIR)/$(PACKAGE_NAME)Paths.go \
		-package_name=$(PACKAGE_NAME) \
		-generate_fakeroot \
		-generate_path_structs \
		-compress_paths=true \
		--generate_simple_unions=true \
		-exclude_modules=ietf-interfaces,openconfig-interfaces \
		$(YANG_DIR)/$(TARGET_MODEL).yang