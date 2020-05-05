package processor

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"code.videolan.org/videolan/CrashDragon/internal/config"
	"code.videolan.org/videolan/CrashDragon/internal/database"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

var rchan = make(chan database.Report, 5000)

// QueueSize returns the number of reports in the queue
func QueueSize() int {
	return len(rchan)
}

// StartQueue runs the processor queue
func StartQueue() {
	// Spawn 4 processors
	for i := 0; i < 4; i++ {
		go processHandler()
	}
}

// AddToQueue adds new items to the queue
func AddToQueue(Report database.Report) {
	select {
	case rchan <- Report:
	default:
		log.Println("Channel full. Discarding report")
	}
}

// Reprocess is a direct way to spawn a single processor which reprocesses a single report
func Reprocess(Report database.Report) {
	processReport(Report, true)
}

// ProcessText adds the text version of the report to the database, which is only used when the text button is clicked
func ProcessText(Report *database.Report) {
	filepth := filepath.Join(config.C.ContentDirectory, "TXT", Report.ID.String()[0:2], Report.ID.String()[0:4])
	err := os.MkdirAll(filepth, 0755)
	if err != nil {
		return
	}
	f, err := os.Create(filepath.Join(filepth, Report.ID.String()+".txt"))
	if err != nil {
		return
	}
	defer f.Close()

	file := filepath.Join(config.C.ContentDirectory, "Reports", Report.ID.String()[0:2], Report.ID.String()[0:4], Report.ID.String()+".dmp")
	symbolsPath := filepath.Join(config.C.ContentDirectory, "Symfiles", Report.Product.Slug, Report.Version.Slug)

	dataTXT, err := runProcessor(file, symbolsPath, "txt")
	if err != nil {
		return
	}

	_, err = f.Write(dataTXT)
	if err != nil {
		return
	}
}

func processHandler() {
	for {
		r := <-rchan
		log.Printf("Unprocessed reports: %d", len(rchan))
		processReport(r, false)
	}
}

func runProcessor(minidumpFile string, symbolsPath string, format string) ([]byte, error) {
	cmd := exec.Command(config.C.SymbolicatorPath, "-f", format, minidumpFile, symbolsPath)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	if err = cmd.Start(); err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(stdout)
	if err != nil {
		return nil, err
	}

	if err := cmd.Wait(); err != nil {
		return nil, err
	}
	return data, nil
}

func processReport(Report database.Report, reprocess bool) {
	start := time.Now()

	file := filepath.Join(config.C.ContentDirectory, "Reports", Report.ID.String()[0:2], Report.ID.String()[0:4], Report.ID.String()+".dmp")
	symbolsPath := filepath.Join(config.C.ContentDirectory, "Symfiles", Report.Product.Slug, Report.Version.Slug)

	dataJSON, err := runProcessor(file, symbolsPath, "json")
	tx := database.Db.Begin()
	if err != nil {
		os.Remove(file)
		tx.Unscoped().Delete(&Report)
		tx.Commit()
		return
	}

	Report.Report = database.ReportContent{}
	err = json.Unmarshal(dataJSON, &Report.Report)
	if err != nil {
		os.Remove(file)
		tx.Unscoped().Delete(&Report)
		tx.Commit()
		return
	}

	if Report.Report.Status != "OK" {
		Report.Processed = false
	} else {
		Report.Processed = true
	}

	Report.Os = Report.Report.SystemInfo.Os
	Report.OsVersion = Report.Report.SystemInfo.OsVer
	Report.Arch = Report.Report.SystemInfo.CPUArch

	if reprocess {
		Report.Signature = ""
		Report.Module = ""
		Report.CrashLocation = ""
		Report.CrashPath = ""
		Report.CrashLine = 0
	}

	if len(Report.Report.Threads) > Report.Report.CrashInfo.CrashingThread {
		for _, Frame := range Report.Report.Threads[Report.Report.CrashInfo.CrashingThread].Frames {
			if Frame.File == "" && Report.Signature != "" {
				continue
			}
			if Report.Module == "" || (Report.Signature == "" && Frame.Function != "") {
				Report.Module = strings.TrimSuffix(Frame.Module, filepath.Ext(Frame.Module))
				Report.Signature = Frame.Function
			}
			if Frame.File == "" {
				continue
			}
			Report.CrashLocation = Frame.File + ":" + strconv.Itoa(Frame.Line)
			Report.CrashPath = Frame.File
			Report.CrashLine = Frame.Line
			break
		}
	} else {
		log.Printf("Crashing thread %d is out of index in Threads!", Report.Report.CrashInfo.CrashingThread)
	}

	if !reprocess {
		Report.CreatedAt = time.Now()
	}

	var Crash database.Crash
	processCrash(tx, Report, reprocess, &Crash)
	Report.CrashID = Crash.ID

	Report.ProcessingTime = time.Since(start).Seconds()

	if reprocess {
		tx.Save(&Report)
	} else {
		tx.Create(&Report)
	}

	tx.Save(&Crash)
	tx.Commit()
}

func processCrash(tx *gorm.DB, Report database.Report, reprocess bool, Crash *database.Crash) {
	if reprocess && Report.CrashID != uuid.Nil {
		database.Db.First(&Crash, "id = ?", Report.CrashID)
		Crash.Signature = Report.Signature
		Crash.Module = Report.Module
	} else {
		database.Db.FirstOrInit(&Crash, "signature = ? AND module = ?", Report.Signature, Report.Module)
	}

	if Crash.ID == uuid.Nil {
		Crash.ID = uuid.NewV4()

		Crash.FirstReported = Report.CreatedAt
		Crash.Signature = Report.Signature
		Crash.Module = Report.Module

		Crash.ProductID = Report.ProductID

		Crash.Fixed = nil

		tx.Create(&Crash)
		reprocess = false
	}
	if !reprocess || Report.CrashID == uuid.Nil {
		Crash.LastReported = Report.CreatedAt
	}

	tx.Model(&Crash).Related(&Crash.Versions, "Versions")
	for _, Version := range Crash.Versions {
		if Version.ID == Report.Version.ID {
			break
		}
		Crash.Fixed = nil
	}

	tx.Model(&Crash).Association("Versions").Append(&Report.Version)
}