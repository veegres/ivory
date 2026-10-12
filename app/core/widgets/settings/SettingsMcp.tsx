import {Box, FormControl, InputLabel, MenuItem, Select, ToggleButton, Tooltip} from "@mui/material"
import {useState} from "react"

import {useRouterAuthToken} from "../../../features/auth/api/AuthHook"
import {useRouterInfo} from "../../../features/management/api/ManagementHook"
import {AlertAlign} from "../../../shared/component/box/AlertAlign"
import {CopyField} from "../../../shared/component/box/CopyField"
import {TitleBox} from "../../../shared/component/box/TitleBox"
import {SimpleButton} from "../../../shared/component/button/SimpleButton"
import {CodeField} from "../../../shared/component/input/CodeField"
import {SxPropsMap} from "../../../shared/helper/HelperType"
import {buildRouteUrl} from "../../../shared/helper/HelperUrl"
import {DateTimeFormatter} from "../../../shared/helper/HelperUtils"

const SX: SxPropsMap = {
    scroll: {display: "flex", flexDirection: "column", padding: "0px 15px", gap: 2},
    section: {display: "flex", flexDirection: "column", gap: 1, padding: "4px 6px"},
    row: {display: "flex", alignItems: "center", gap: 1},
    toggles: {display: "flex", alignItems: "center", gap: 0.5},
    toggle: {whiteSpace: "nowrap"},
    grow: {flexGrow: 1, minWidth: "150px"},
}

const HOUR = 60 * 60
const DAY = 24 * HOUR
const EXPIRATIONS = [
    {label: "1h", description: "Expires in 1 hour", value: HOUR},
    {label: "6h", description: "Expires in 6 hours", value: 6 * HOUR},
    {label: "1d", description: "Expires in 1 day", value: DAY},
    {label: "15d", description: "Expires in 15 days", value: 15 * DAY},
    {label: "30d", description: "Expires in 30 days", value: 30 * DAY},
]

// TOKEN_PLACEHOLDER stands in the config until a token has been generated. It
// is the same {{...}} shape a deployment command uses for a value nobody has
// answered yet, so CodeEditor highlights it as one - a reader sees the one part
// of the config that is not real.
const TOKEN_PLACEHOLDER = "{{ivory_token}}"

// CLIENTS carries each client's whole config: withToken for an Ivory that
// authenticates, withoutToken for one running open. Every one is written out in
// full, including the two that bridge through mcp-remote identically - they are
// data, and the point is that a reader sees the config the user will paste,
// without resolving a helper first. The duplication is deliberate.
const CLIENTS = [
    {
        key: "claude-code", label: "Claude Code", where: "your project, as a command rather than a file",
        withToken: (url: string, token: string) =>
            `claude mcp add --transport http ivory ${url} \\\n  --header "Authorization: Bearer ${token}"`,
        withoutToken: (url: string) => `claude mcp add --transport http ivory ${url}`,
        envToken: false,
    },
    {
        key: "claude-code-file", label: "Claude Code (file)", where: ".mcp.json, committed with the project",
        withToken: (url: string, token: string) => `{
  "mcpServers": {
    "ivory": {
      "type": "http",
      "url": "${url}",
      "headers": {"Authorization": "Bearer ${token}"}
    }
  }
}`,
        withoutToken: (url: string) => `{
  "mcpServers": {
    "ivory": {
      "type": "http",
      "url": "${url}"
    }
  }
}`,
        envToken: false,
    },
    {
        key: "cursor", label: "Cursor", where: ".cursor/mcp.json, or ~/.cursor/mcp.json for every project",
        withToken: (url: string, token: string) => `{
  "mcpServers": {
    "ivory": {
      "url": "${url}",
      "headers": {"Authorization": "Bearer ${token}"}
    }
  }
}`,
        withoutToken: (url: string) => `{
  "mcpServers": {
    "ivory": {
      "url": "${url}"
    }
  }
}`,
        envToken: false,
    },
    {
        key: "vscode", label: "VS Code", where: ".vscode/mcp.json",
        withToken: (url: string, token: string) => `{
  "servers": {
    "ivory": {
      "type": "http",
      "url": "${url}",
      "headers": {"Authorization": "Bearer ${token}"}
    }
  }
}`,
        withoutToken: (url: string) => `{
  "servers": {
    "ivory": {
      "type": "http",
      "url": "${url}"
    }
  }
}`,
        envToken: false,
    },
    {
        key: "codex", label: "Codex", where: "~/.codex/config.toml",
        withToken: (url: string) =>
            `[mcp_servers.ivory]\nurl = "${url}"\nbearer_token_env_var = "IVORY_TOKEN"`,
        withoutToken: (url: string) => `[mcp_servers.ivory]\nurl = "${url}"`,
        envToken: true,
    },
    {
        key: "antigravity", label: "Antigravity",
        where: ".agents/mcp_config.json, or ~/.gemini/config/mcp_config.json globally",
        withToken: (url: string, token: string) => `{
  "mcpServers": {
    "ivory": {
      "serverUrl": "${url}",
      "headers": {"Authorization": "Bearer ${token}"}
    }
  }
}`,
        withoutToken: (url: string) => `{
  "mcpServers": {
    "ivory": {
      "serverUrl": "${url}"
    }
  }
}`,
        envToken: false,
    },
    {
        key: "intellij", label: "IntelliJ IDEA",
        where: "Settings | Tools | AI Assistant | Model Context Protocol (MCP)",
        withToken: (url: string, token: string) => `{
  "mcpServers": {
    "ivory": {
      "command": "npx",
      "args": ["mcp-remote", "${url}", "--header", "Authorization:Bearer \${IVORY_TOKEN}"],
      "env": {"IVORY_TOKEN": "${token}"}
    }
  }
}`,
        withoutToken: (url: string) => `{
  "mcpServers": {
    "ivory": {
      "url": "${url}"
    }
  }
}`,
        envToken: true,
    },
    {
        key: "claude-desktop", label: "Claude Desktop", where: "claude_desktop_config.json",
        withToken: (url: string, token: string) => `{
  "mcpServers": {
    "ivory": {
      "command": "npx",
      "args": ["mcp-remote", "${url}", "--header", "Authorization:Bearer \${IVORY_TOKEN}"],
      "env": {"IVORY_TOKEN": "${token}"}
    }
  }
}`,
        withoutToken: (url: string) => `{
  "mcpServers": {
    "ivory": {
      "command": "npx",
      "args": ["mcp-remote", "${url}"]
    }
  }
}`,
        envToken: true,
    },
    {
        key: "other", label: "Any other client", where: "whatever that client calls its config",
        withToken: (url: string, token: string) =>
            `transport  streamable HTTP\nurl        ${url}\nheader     Authorization: Bearer ${token}`,
        withoutToken: (url: string) => `transport  streamable HTTP\nurl        ${url}`,
        envToken: false,
    },
] as const

type ClientKey = typeof CLIENTS[number]["key"]

export function SettingsMcp() {
    const info = useRouterInfo(false)
    const token = useRouterAuthToken()
    const [client, setClient] = useState<ClientKey>("claude-code")
    const [expiresIn, setExpiresIn] = useState(HOUR)
    const authEnabled = (info.data?.auth.supported.length ?? 0) > 0
    const url = buildRouteUrl("/api/mcp")
    const entry = getClient()

    return (
        <Box sx={SX.scroll}>
            <AlertAlign align={"justify"} text={renderIntro()}/>
            <TitleBox label={"Endpoint"} island collapsible={false}>
                <Box sx={SX.section}>
                    <CopyField value={url} message={"Endpoint copied to clipboard!"}/>
                </Box>
            </TitleBox>
            <TitleBox label={"Token"} island collapsible={false}>
                <Box sx={SX.section}>{renderToken()}</Box>
            </TitleBox>
            <TitleBox label={"Connect a client"} island collapsible={false}>
                <Box sx={SX.section}>{renderClient()}</Box>
            </TitleBox>
            <AlertAlign align={"justify"} text={renderPermissions()}/>
        </Box>
    )

    function renderIntro() {
        return (<>
            Ivory speaks the <b>Model Context Protocol</b>, so an AI agent can read a cluster's health
            and operate it with the same permissions you have. Everything Ivory already knows - node
            roles, replication lag, pending restarts, its own drift warnings, VM metrics, container logs
            and your saved queries - becomes readable in one conversation instead of one screen at a time.
        </>)
    }

    function renderToken() {
        if (!authEnabled) return <AlertAlign align={"justify"} text={renderTokenNotNeeded()}/>
        if (!token.data) return renderTokenButton()
        return (<>
            <CopyField value={token.data.token} message={"Token copied to clipboard!"}/>
            <AlertAlign align={"justify"} severity={"warning"} text={renderTokenWarning(token.data.expiresAt)}/>
        </>)
    }

    function renderTokenNotNeeded() {
        return (<>
            Ivory is running with authentication switched off, so a client needs no token and every tool
            is available - the same as the REST API, which is already open in this mode.
        </>)
    }

    function renderTokenButton() {
        return (<>
            <Box sx={SX.row}>
                <Box sx={SX.toggles}>{EXPIRATIONS.map(renderExpiration)}</Box>
                <Box sx={SX.grow}>
                    <SimpleButton
                        fullWidth={true}
                        variant={"contained"}
                        color={"primary"}
                        loading={token.isPending}
                        onClick={handleGenerate}
                    >
                        Generate token
                    </SimpleButton>
                </Box>
            </Box>
            <AlertAlign align={"justify"} severity={"info"} text={renderExpirationHint()}/>
        </>)
    }

    function renderExpirationHint() {
        return (<>
            Ivory <b>cannot revoke</b> a token once it is issued - it stays valid until it expires, even
            if it leaks. Choose the <b>shortest expiration</b> your agent needs, and be cautious with the
            longer ones.
        </>)
    }

    function renderExpiration(option: typeof EXPIRATIONS[number]) {
        return (
            <Tooltip key={option.value} title={option.description} placement={"top"} arrow disableInteractive>
                <Box component={"span"}>
                    <ToggleButton
                        sx={SX.toggle}
                        value={option.value}
                        selected={expiresIn === option.value}
                        onClick={() => setExpiresIn(option.value)}
                    >
                        {option.label}
                    </ToggleButton>
                </Box>
            </Tooltip>
        )
    }

    function renderTokenWarning(expiresAt: string) {
        return (<>
            Anyone holding this token can do everything you may do in Ivory, so treat it as your
            password. It expires on <b>{DateTimeFormatter.utc(expiresAt)}</b> and cannot be
            refreshed - when an agent starts failing with <b>401</b>, generate another one. Ivory shows
            it once, so copy it now.
        </>)
    }

    function renderClient() {
        return (<>
            <FormControl fullWidth>
                <InputLabel id={"select-mcp-client"}>Client</InputLabel>
                <Select
                    labelId={"select-mcp-client"}
                    label={"Client"}
                    value={client}
                    onChange={(e) => setClient(e.target.value as ClientKey)}
                    fullWidth={true}
                    variant={"outlined"}
                >
                    {CLIENTS.map(({key, label}) => <MenuItem key={key} value={key}>{label}</MenuItem>)}
                </Select>
            </FormControl>
            <CodeField
                label={"Config"}
                value={getConfig()}
                editable={false}
                copy={true}
                hint={`where it goes - ${entry.where}`}
            />
            {renderEnvToken()}
        </>)
    }

    function renderEnvToken() {
        if (!authEnabled || !entry.envToken) return
        return (
            <CodeField
                label={"Environment"}
                value={getEnvToken()}
                editable={false}
                copy={true}
                hint={"some clients can work with, or require, an environment variable to hold the auth token - you can copy it here"}
            />
        )
    }

    function renderPermissions() {
        return (<>
            A client only ever sees the tools you hold permission for, and the list is recomputed on
            every call - so revoking a permission takes effect on the agent's next request without
            reconnecting. Tools that can break something are marked destructive, which is what your
            client shows you at its approval prompt.
        </>)
    }

    function handleGenerate() {
        token.mutate({expiresIn})
    }

    function getClient() {
        return CLIENTS.find((c) => c.key === client) ?? CLIENTS[0]
    }

    function getConfig() {
        return authEnabled ? entry.withToken(url, getToken()) : entry.withoutToken(url)
    }

    function getEnvToken() {
        return `export IVORY_TOKEN="${getToken()}"`
    }

    function getToken() {
        return token.data?.token ?? TOKEN_PLACEHOLDER
    }
}
