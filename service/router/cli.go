package router

import (
	"bufio"
	"github.com/gin-gonic/gin"
	. "ivory/model"
	"net/http"
	"os/exec"
	"strconv"
)

func (r routes) CliGroup(group *gin.RouterGroup) {
	node := group.Group("/cli")
	node.POST("/pgcompacttable", runPgcompacttable)
}

func runPgcompacttable(context *gin.Context) {
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

	cmd := exec.Command("pgcompacttable", sb...)
	pipe, err := cmd.StdoutPipe()
	if err := cmd.Start(); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	reader := bufio.NewReader(pipe)
	line, _, err := reader.ReadLine()
	var output []string
	for err == nil {
		output = append(output, string(line))
		line, _, err = reader.ReadLine()
	}

	context.JSON(http.StatusOK, gin.H{"response": output})
}
