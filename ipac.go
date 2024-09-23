package ipac

import (
	"fmt"
	"image/color"
	"io"

	"github.com/bearsh/hid"
)

const (
	espressifVendorID  = 0xd209
	espressifProductID = 0x0412
	hidInterface       = 2
)

type IPacUltimateIO struct {
	hid *hid.Device
}

func NewIPacUltimateIO() *IPacUltimateIO {
	return &IPacUltimateIO{}
}

func (pd *IPacUltimateIO) InitDevices() error {
	devices := hid.Enumerate(espressifVendorID, espressifProductID)
	if len(devices) == 0 {
		return fmt.Errorf("no HID devices found")
	}

	fmt.Printf("Found %d HID devices\n", len(devices))

	// Open each device and initialize
	for _, devInfo := range devices {
		if devInfo.Interface != hidInterface {
			continue
		}

		device, err := devInfo.Open()
		if err != nil {
			return fmt.Errorf("failed to open device: %v", err)
		}

		pd.hid = device
	}

	return nil
}

func (pd *IPacUltimateIO) Close() error {
	return pd.hid.Close()
}

type Report struct {
	ReportID     byte
	ReportBuffer [4]byte
}

func (r *Report) Write(writer io.Writer) error {
	data := []byte{r.ReportID, r.ReportBuffer[0], r.ReportBuffer[1], r.ReportBuffer[2], r.ReportBuffer[3]}

	_, err := writer.Write(data)
	if err != nil {
		return err
	}

	return nil
}

func (d *IPacUltimateIO) SetLEDColor(group int, c color.Color) (bool, error) {
	r32, g32, b32, _ := c.RGBA()

	r := uint8(r32 >> 8)
	g := uint8(g32 >> 8)
	b := uint8(b32 >> 8)

	first := group * 3
	d.SetLEDIntensity(first, r)
	d.SetLEDIntensity(first+1, g)
	d.SetLEDIntensity(first+2, b)

	return true, nil
}

func (d *IPacUltimateIO) SetLEDIntensity(port int, intensity byte) (bool, error) {
	bitMask := byte(0x7F)

	// Create the output report
	var r Report
	r.ReportID = 3
	r.ReportBuffer[0] = 0x80

	if port != -1 {
		r.ReportBuffer[0] = byte(port) & bitMask
	}

	r.ReportBuffer[1] = intensity
	r.ReportBuffer[2] = 0
	r.ReportBuffer[3] = 0

	if err := r.Write(d.hid); err != nil {
		return false, fmt.Errorf("error setting lediD %d to %d", -1, intensity)
	}

	return true, nil
}

func (d *IPacUltimateIO) SetLEDFadeTime(fadeTime int) (bool, error) {
	var r Report
	r.ReportID = 3
	r.ReportBuffer[0] = 0xc0
	r.ReportBuffer[1] = byte(fadeTime)
	r.ReportBuffer[2] = 0
	r.ReportBuffer[3] = 0

	if err := r.Write(d.hid); err != nil {
		return false, fmt.Errorf("error setting le")
	}

	return true, nil
}
