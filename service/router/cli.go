package router

import (
	"bufio"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"io"
	. "ivory/model"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
)

func (r routes) CliGroup(group *gin.RouterGroup) {
	jobs := make(map[uuid.UUID]chan []byte)
	node := group.Group("/cli")
	node.POST("/pgcompacttable", postPgcompacttable(jobs))
	node.GET("/pgcompacttable/:uuid/stream", stream(jobs))
}

func postPgcompacttable(jobs map[uuid.UUID]chan []byte) func(context *gin.Context) {
	return func(context *gin.Context) {
		var cli PgCompactTable
		_ = context.ShouldBindJSON(&cli)

		sb := []string{
			"--host", cli.Connection.Host,
			"--port", strconv.Itoa(cli.Connection.Port),
			"--user", cli.Connection.Username,
			"--password", cli.Connection.Password,
		}
		if cli.Target != nil {
			if cli.Target.DbName != "" {
				sb = append(sb, "--dbname", cli.Target.DbName)
			}
			if cli.Target.Schema != "" {
				sb = append(sb, "--schema", cli.Target.Schema)
			}
			if cli.Target.Table != "" {
				sb = append(sb, "--table", cli.Target.Table)
			}
			if cli.Target.ExcludeSchema != "" {
				sb = append(sb, "--excludeSchema", cli.Target.ExcludeSchema)
			}
			if cli.Target.ExcludeTable != "" {
				sb = append(sb, "--excludeTable", cli.Target.ExcludeTable)
			}
		}
		if len(sb) == 8 {
			sb = append(sb, "--all")
		}
		if cli.Ratio != 0 {
			sb = append(sb, "--delay-ratio", strconv.Itoa(cli.Ratio))
		}
		sb = append(sb, "--verbose")

		jobUuid := uuid.New()
		go func() {
			cmd := exec.Command("pgcompacttable", sb...)
			pipe, err := cmd.StdoutPipe()
			if err := cmd.Start(); err != nil {
				context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			}
			jobs[jobUuid] = make(chan []byte)
			reader := bufio.NewReader(pipe)
			line, _, err := reader.ReadLine()
			var output []string
			for err == nil {
				output = append(output, string(line))
				jobs[jobUuid] <- line
				line, _, err = reader.ReadLine()
			}
			close(jobs[jobUuid])
		}()

		context.JSON(http.StatusOK, gin.H{"response": gin.H{
			"cmd":  "pgcompacttable " + strings.Join(sb, " "),
			"uuid": jobUuid,
		}})
	}
}

func stream(jobs map[uuid.UUID]chan []byte) func(context *gin.Context) {
	return func(context *gin.Context) {
		jobUuid, err := uuid.Parse(context.Param("uuid"))
		job := jobs[jobUuid]
		if err != nil || job == nil {
			context.AbortWithStatus(404)
			return
		}
		// notify proxy that it shouldn't enable any caching
		context.Writer.Header().Set("Cache-Control", "no-transform")
		// force using correct event-stream if there is no proxy
		context.Writer.Header().Set("Content-Type", "text/event-stream")
		context.SSEvent("status", "start")
		context.Writer.Flush()

		context.Stream(func(w io.Writer) bool {
			if line, ok := <-job; ok {
				context.SSEvent("log", string(line))
				return true
			}

			return false
		})

		delete(jobs, jobUuid)
		context.SSEvent("status", "finish")
	}
}
