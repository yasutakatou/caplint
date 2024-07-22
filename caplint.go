package main

//set PATH=%PATH%;c:\Program Files\Tesseract-OCR

import (
	"bytes"
	"flag"
	"fmt"
	"image"
	"image/png"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"os/signal"
	"regexp"
	"strings"
	"syscall"
	"time"
	"unsafe"

	"github.com/kbinani/screenshot"
	"github.com/moutend/go-hook/pkg/keyboard"
	"github.com/moutend/go-hook/pkg/types"
	"github.com/nfnt/resize"
	"golang.design/x/clipboard"
	"gopkg.in/ini.v1"
)

var (
	debug, logging, nodel                           bool
	ocr, linter                                     string
	shortcutwindow, shortcutclipboard, shortcutexit int
	targetHwnd                                      uintptr
	rs1Letters                                      = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type (
	HANDLE uintptr
	HWND   HANDLE
)

var (
	user32                  = syscall.MustLoadDLL("user32.dll")
	procEnumWindows         = user32.MustFindProc("EnumWindows")
	procGetWindowTextW      = user32.MustFindProc("GetWindowTextW")
	procSetActiveWindow     = user32.MustFindProc("SetActiveWindow")
	procSetForegroundWindow = user32.MustFindProc("SetForegroundWindow")
	procGetForegroundWindow = user32.MustFindProc("GetForegroundWindow")
	procGetWindowRect       = user32.MustFindProc("GetWindowRect")
)

type _RECT struct {
	Left   int32
	Top    int32
	Right  int32
	Bottom int32
}

func main() {
	_Debug := flag.Bool("debug", false, "[-debug=debug mode (true is enable)]")
	_Logging := flag.Bool("log", false, "[-log=logging mode (true is enable)]")
	_Config := flag.String("config", "caplint.ini", "[-config=config file)]")
	_File := flag.String("file", "text.png", "[-file=existing png files)]")
	_Clipboard := flag.Bool("clipboard", false, "[-clipboard=Mmode for reading images from the clipboard (true is enable)]")
	_Shortcut := flag.Bool("shortcut", false, "[-shortcut=keep activated and use shortcut keys to perform readings in this mode (true is enable)]")
	_Shortcutwindow := flag.Int("shortcutwindow", 65, "[-shortcutwindow=shortcut key to read from a focused window (default 'a')]")
	_Shortcutclipboard := flag.Int("shortcutclipboard", 90, "[-shortcutclipboatrd=shortcut keys for reading from the clipboard (default 'z')]")
	_Shortcutexit := flag.Int("shortexit", 81, "[-shortcutexit=key to exit shortcut key mode (default 'q')]")
	_Resize := flag.Int("resize", 2, "[-resize=temporary resize and enlarge multiples (default 2)]")
	_NoDelete := flag.Bool("nodelete", false, "[-nodelete=leave temporary files undeleted mode. for debugging purposes (true is enable)]")

	flag.Parse()

	debug = bool(*_Debug)
	logging = bool(*_Logging)
	shortcutwindow = int(*_Shortcutwindow)
	shortcutclipboard = int(*_Shortcutclipboard)
	shortcutexit = int(*_Shortcutexit)
	nodel = bool(*_NoDelete)

	if Exists(*_Config) == true {
		loadConfig(*_Config)
	} else {
		fmt.Printf("Fail to read config file: %v\n", *_Config)
		os.Exit(1)
	}

	filename := ""

	if *_Shortcut == true {
		go func() {
			ShortCutDo(*_Resize)
		}()

		for {
			fmt.Println("Forground window short cut: " + string(rune(*_Shortcutwindow)))
			fmt.Println("Clipboard image short cut: " + string(rune(*_Shortcutclipboard)))
			fmt.Println("exit short cut mode. " + string(rune(*_Shortcutexit)))
			time.Sleep(time.Second * 360)
		}
		os.Exit(0)
	} else if *_Clipboard == true {
		filename = resizeImage(loadFromClipboard(), *_Resize)
		debugLog("filename: " + filename)
	} else {
		filename = resizeImage(loadFromFile(*_File), *_Resize)
		debugLog("filename: " + filename)
	}

	ocrDo(filename)
	lintDo(filename)
	if nodel == false {
		os.Remove(filename + ".png")
		os.Remove(filename + ".txt")
	}
	os.Exit(0)
}

func ocrDo(filename string) {
	command := strings.Replace(ocr, "{}", filename, -1)
	debugLog("ocr: " + command)

	cmd := exec.Command("cmd", "/c", command)
	output, err := cmd.CombinedOutput()
	if err != nil {
		debugLog(fmt.Sprint(err) + ": " + string(output))
	} else {
		debugLog(string(output))
	}
}

func lintDo(filename string) {
	command := strings.Replace(linter, "{}", filename, -1)
	debugLog("linter: " + command)

	cmd := exec.Command("cmd", "/c", command)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + string(output))
		debugLog(fmt.Sprint(err) + ": " + string(output))
	} else {
		fmt.Println(" -- No Lint! -- ")
		debugLog(string(output))
	}
}

func loadConfig(configFile string) {
	loadOptions := ini.LoadOptions{}
	loadOptions.UnparseableSections = []string{"tesseract", "textlint"}

	cfg, err := ini.LoadSources(loadOptions, configFile)
	if err != nil {
		fmt.Printf("Fail to read config file: %v", err)
		os.Exit(1)
	}

	ocr = ""
	linter = ""

	setStructs("tesseract", cfg.Section("tesseract").Body(), 1)
	setStructs("textlint", cfg.Section("textlint").Body(), 2)
}

func debugLog(message string) {
	var file *os.File
	var err error

	if debug == true {
		fmt.Println(message)
	}

	if logging == false {
		return
	}

	const layout = "2006-01-02_15"
	const layout2 = "2006/01/02 15:04:05"
	t := time.Now()
	filename := t.Format(layout) + ".log"
	logHead := "[" + t.Format(layout2) + "] "

	if Exists(filename) == true {
		file, err = os.OpenFile(filename, os.O_WRONLY|os.O_APPEND, 0666)
	} else {
		file, err = os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0666)
	}

	if err != nil {
		log.Fatal(err)
		return
	}
	defer file.Close()
	fmt.Fprintln(file, logHead+message)
}

func setStructs(configType, datas string, flag int) {
	debugLog(" -- " + configType + " --")

	for _, v := range regexp.MustCompile("\r\n|\n\r|\n|\r").Split(datas, -1) {
		if len(v) > 0 {
			if flag == 1 {
				ocr = v
				debugLog(v)
			} else if flag == 2 {
				linter = v
				debugLog(v)
			}
		}
	}
}

func ShortCutDo(resizer int) {
	keyboardChan := make(chan types.KeyboardEvent, 100)

	if err := keyboard.Install(nil, keyboardChan); err != nil {
		log.Fatal(err)
		return
	}
	defer keyboard.Uninstall()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	shiftFlag := false
	ctrlFlag := false

	//https://www.k-cube.co.jp/wakaba/server/ascii_code.html
	for {
		select {
		case <-signalChan:
		case k := <-keyboardChan:
			//fmt.Printf("Received %v %v\n", k.Message, k.VKCode)
			if k.Message == types.WM_KEYDOWN {
				if k.VKCode == types.VK_LSHIFT || k.VKCode == types.VK_RSHIFT {
					shiftFlag = true
				}
				if k.VKCode == types.VK_LCONTROL || k.VKCode == types.VK_RCONTROL {
					ctrlFlag = true
				}
				if shiftFlag == true && ctrlFlag == true {
					if int(k.VKCode) == shortcutwindow {
						debugLog("Shortcut: capture from forground window!")
						getScreenCapture(resizer)
					} else if int(k.VKCode) == shortcutclipboard {
						debugLog("Shortcut: capture from clipboard")
						shortcutclipboardDo(resizer)
					} else if int(k.VKCode) == shortcutexit {
						debugLog("Shortcut exit.")
						os.Exit(0)
					}
				}
			}
			if k.Message == types.WM_KEYUP {
				if k.VKCode == types.VK_LSHIFT || k.VKCode == types.VK_RSHIFT {
					shiftFlag = false
				}
				if k.VKCode == types.VK_LCONTROL || k.VKCode == types.VK_RCONTROL {
					ctrlFlag = false
				}
			}
			continue
		}
	}
}

func shortcutclipboardDo(resizer int) {
	filename := resizeImage(loadFromClipboard(), resizer)
	debugLog("filename: " + filename)

	ocrDo(filename)
	lintDo(filename)
	if nodel == false {
		os.Remove(filename + ".png")
		os.Remove(filename + ".txt")
	}
}

func loadFromFile(filename string) image.Image {
	fileData, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	img, _, err := image.Decode(fileData)
	if err != nil {
		log.Fatal(err)
	}
	fileData.Close()
	return img
}

func GetWindowRect(hwnd HWND, rect *_RECT) (err error) {
	r1, _, e1 := syscall.Syscall(procGetWindowRect.Addr(), 7, uintptr(hwnd), uintptr(unsafe.Pointer(rect)), 0)
	if r1 == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func getWindow(funcName string) uintptr {
	hwnd, _, _ := syscall.Syscall(procGetForegroundWindow.Addr(), 6, 0, 0, 0)
	if debug == true {
		fmt.Printf("currentWindow: handle=0x%x\n", hwnd)
	}
	return hwnd
}

func getScreenCapture(resizer int) {
	filenametmp := RandStr(8)

	var rect _RECT
	GetWindowRect(HWND(getWindow("GetForegroundWindow")), &rect)
	if debug == true {
		fmt.Printf("window rect: ")
		fmt.Println(rect)
	}

	img, err := screenshot.Capture(int(rect.Left)+10, int(rect.Top)+10, int(rect.Right-rect.Left)-20, int(rect.Bottom-rect.Top)-20)
	if err != nil {
		panic(err)
	}
	save(img, filenametmp+".png")

	filename := resizeImage(loadFromFile(filenametmp+".png"), resizer)
	debugLog("filename: " + filename)

	ocrDo(filename)
	lintDo(filename)
	if nodel == false {
		os.Remove(filenametmp + ".png")
		os.Remove(filename + ".png")
		os.Remove(filename + ".txt")
	}
}

func save(img *image.RGBA, filePath string) {
	file, err := os.Create(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	png.Encode(file, img)
}

func loadFromClipboard() image.Image {
	var b []byte

	b = clipboard.Read(clipboard.FmtImage)
	if b == nil {
		fmt.Println("clipboard is empty!")
		os.Exit(1)
	}

	buf := bytes.NewReader(b)
	img, _, err := image.Decode(buf)
	if err != nil {
		log.Fatal(err)
	}
	return img
}

func RandStr(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = rs1Letters[rand.Intn(len(rs1Letters))]
	}
	return string(b)
}

func resizeImage(img image.Image, resizecount int) string {
	filename := RandStr(8)

	bounds := img.Bounds()
	w := bounds.Dx()
	h := bounds.Dy()
	//fmt.Printf("(image size) w: %d h:%d\n", w, h)

	resizedImg := resize.Resize(uint(w*resizecount), uint(h*resizecount), img, resize.NearestNeighbor)

	output, err := os.Create(filename + ".png")
	if err != nil {
		log.Fatal(err)
	}
	defer output.Close()

	if err := png.Encode(output, resizedImg); err != nil {
		log.Fatal(err)
	}
	return filename
}

func Exists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}
