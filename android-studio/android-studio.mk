ANDROID_STUDIO_IMAGE_DOCKER_BUILDX_LOCAL_CACHE=$(CURDIR)/.containers/docker-buildx-cache/android-studio
$(ANDROID_STUDIO_IMAGE_DOCKER_BUILDX_LOCAL_CACHE):;	@mkdir -p "$@"
ifneq ($(wildcard $(ANDROID_STUDIO_IMAGE_DOCKER_BUILDX_LOCAL_CACHE)/index.json),)
ANDROID_STUDIO_IMAGE_DOCKER_BUILDX_LOCAL_CACHE_READ = type=local,src=$(ANDROID_STUDIO_IMAGE_DOCKER_BUILDX_LOCAL_CACHE)
endif
ANDROID_STUDIO_IMAGE_DOCKER_BUILDX_LOCAL_CACHE_WRITE= type=local,dest=$(ANDROID_STUDIO_IMAGE_DOCKER_BUILDX_LOCAL_CACHE)
ANDROID_STUDIO_IMAGE = $(PUBLIC_IMAGES_BASE):android-studio-latest

local-android-studio-docker-compose-up: local-android-studio-docker-compose-rm
	VERSION=latest $(LOCAL_DIR)/android-studio/docker-compose-up.sh &
	@$(LOCAL_DOCKER_COMPOSE) logs android-studio
	@echo
	@echo Started Androïd studio display: 
	@echo Run "'make $(@:%-up=%-logs)'" to retrieve more logs
.PHONY: local-android-studio-docker-compose-up

LOCAL_ANDROID_STUDIO_DOCKER_COMPOSE_EXEC = $(LOCAL_DOCKER_COMPOSE) exec android-studio

LOCAL_ANDROID_STUDIO_ANDROID_AVD_NAME ?= Pixel_6_Playstore
LOCAL_ANDROID_STUDIO_ANDROID_GPU_MODE ?= swiftshader_indirect

local-android-studio-appium-emulator-avd-wipe-data:
	@$(LOCAL_ANDROID_STUDIO_DOCKER_COMPOSE_EXEC) bash -c ' \
		emulator -avd $(LOCAL_ANDROID_STUDIO_ANDROID_AVD_NAME) -no-snapshot-save -wipe-data \
.PHONY: local-android-studio-appium-emulator-avd-wipe-data

local-android-studio-appium-emulator-avd-restart:
	@$(LOCAL_ANDROID_STUDIO_DOCKER_COMPOSE_EXEC) bash -c ' \
	  echo "[*] Killing emulator & adb..." ; \
	  pkill -9 emulator ; \
	  pkill -9 qemu-system ; \
	  adb kill-server ; \
	  echo "[*] Cleaning up locks..." ; \
	  rm -rf ~/.android/avd/*/*.lock ; \
	  rm -f ~/.android/*.lock ; \
	  rm -f ~/.android/adb*.ini.lock ; \
	  rm -f /tmp/.X1-lock ; \
	  echo "[*] Restarting ADB..." ; \
	  adb start-server ; \
	  echo "[*] Launching emulator..." ; \
	  echo "Starting android-studio emulator..." ; \
	  ANDROID_AVD_NAME=$(LOCAL_ANDROID_STUDIO_ANDROID_AVD_NAME) \
	  LOCAL_ANDROID_GPU_MODE=$(LOCAL_ANDROID_STUDIO_ANDROID_GPU_MODE) \
	  	appium-emulator-avd.sh ; \
	  sleep infinity ; \
	'
.PHONY: local-android-studio-appium-emulator-avd-restart

local-android-studio-emulator-logs:
	@$(LOCAL_ANDROID_STUDIO_DOCKER_COMPOSE_EXEC) adb logcat -T 10
.PHONY: local-android-studio-emulator-logs

local-android-studio-emulator-kernel:
	$(LOCAL_ANDROID_STUDIO_DOCKER_COMPOSE_EXEC) bash -c ' \
	  echo "[*] Showing emulator kernel..." ; \
	  emulator -avd $(LOCAL_ANDROID_STUDIO_ANDROID_AVD_NAME) -no-snapshot-save -wipe-data -show-kernel ; \
	  echo "[*] Emulator kernel shown." ; \
	'
.PHONY: local-android-studio-emulator-kernel

local-android-studio-appium-emulator-avd:
	@$(LOCAL_ANDROID_STUDIO_DOCKER_COMPOSE_EXEC) appium-emulator-avd.sh
.PHONY: local-android-studio-appium-emulator-avd

local-android-studio-emulator-dump: 
	@$(LOCAL_ANDROID_STUDIO_DOCKER_COMPOSE_EXEC) bash -c ' \
	  set -e ; \
	  output_dir=$(LOCAL_DIR)/android-studio/_emulator_dump ; \
	  mkdir -p $$output_dir ; \
	  cd $$output_dir ; \
	  echo Capture android-studio mobile, outputs folder : $$(pwd) ; \
	  adb shell uiautomator dump --compressed ; \
	  adb pull /sdcard/window_dump.xml ; \
	  adb shell rm /sdcard/window_dump.xml ; \
	  adb shell screencap -p /sdcard/popup.png ; \
	  adb pull /sdcard/popup.png ; \
	  adb shell rm /sdcard/popup.png ; \
	  adb shell uiautomator dump /sdcard/dump.xml ; \
	  adb pull /sdcard/dump.xml ./dump.xml ; \
	  adb shell rm /sdcard/dump.xml ; \
	  sudo chmod o+rw -R $$(pwd) ; \
	  echo "Capture android-studio mobile done, outputs folder : $$(pwd)" ; \
	'
.PHONY: local-android-studio-emulator-dump

local-android-studio-emulator-data-load:
	@$(LOCAL_ANDROID_STUDIO_DOCKER_COMPOSE_EXEC) \
	make -C ../.. local-android-studio-emulator-data-load-mobile-images
	@echo "Data loaded to android-studio emulator"
.PHONY: local-android-studio-emulator-data-load

local-android-studio-emulator-data-load-mobile-images:
	@bash -c ' \
	set -e ; \
	echo "Load android-studio emulator data, inputs folder : $$(pwd)" ; \
	$(LOCAL_DIR)/android-studio/emulator-data-load.sh \
		$(APPLICATION_DIR)/tests/mobile_images ; \
	'
.PHONY: local-android-studio-emulator-data-load-mobile-images

local-android-studio-emulator-app-sha1-fingerprint:
	@echo "Android Studio Emulator SHA1 fingerprint:" 
	$(LOCAL_ANDROID_STUDIO_DOCKER_COMPOSE_EXEC) \
	  keytool -list -v -keystore ~/.android/debug.keystore -alias androiddebugkey -storepass android -keypass android
.PHONY: local-android-studio-emulator-app-sha1-fingerprint