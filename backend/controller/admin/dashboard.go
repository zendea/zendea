package admin

import (
	"runtime"
	"time"
	"github.com/gin-gonic/gin"

	"zendea/config"
	"zendea/controller"
	"zendea/util"
)

// initTime is the time when the application was initialized.
var initTime = time.Now()

// DashboardController dashboard controller
type DashboardController struct {
	controller.BaseController
}

// GetSysteminfo get system info
func (c *DashboardController) Systeminfo(ctx *gin.Context) {
	c.Success(ctx, gin.H{
		"appName":     config.AppName,
		"appVersion":  config.AppVersion,
		"buildTime":   config.BuildTime,
		"buildCommit": config.BuildCommit,
		"upTime":      util.TimeSincePro(initTime),
		"os":          runtime.GOOS,
		"arch":        runtime.GOARCH,
		"numCpu":      runtime.NumCPU(),
		"goversion":   runtime.Version(),
	})
}
