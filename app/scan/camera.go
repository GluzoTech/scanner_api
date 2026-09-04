package scan

// cameraDevices binds each camera to the scanner device that feeds it. Devices
// are provisioned by hand today, so this is a hardcoded deployment fact with no
// table behind it — add a line here when a camera is installed.
var cameraDevices = map[string]string{
	"c1": "SCANNER-248A881D-5B66C9EE",
}

// deviceCameras is the reverse index used at lookup time, built once at start
// up. One camera per device: mapping two cameras to the same device id would
// keep only one of them.
var deviceCameras = buildDeviceCameras()

func buildDeviceCameras() map[string]string {
	cameras := make(map[string]string, len(cameraDevices))

	for camera, deviceID := range cameraDevices {
		cameras[deviceID] = camera
	}

	return cameras
}

// CameraForDevice returns the camera bound to a device id, or an empty string
// when the device has not been mapped to one.
func CameraForDevice(deviceID string) string {
	return deviceCameras[deviceID]
}
