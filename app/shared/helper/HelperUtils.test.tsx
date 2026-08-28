import {createTheme} from "@mui/material"
import {describe, expect, it} from "@rstest/core"
import {AxiosError} from "axios"
import dayjs from "dayjs"
import utc from "dayjs/plugin/utc"

import {NodeConfig} from "../../features/cluster/api/ClusterType"
import {DeployVar, KeeperDeploySpecResponse} from "../../features/node/api/NodeType"
import {
    DateTimeFormatter,
    DeployPasswordMask,
    getDomain,
    getDomains,
    getErrorMessage,
    getKeeperDefaultPort,
    getNodeConfig,
    getNodeConfigs,
    getPlaceholders,
    getShortUuid,
    getUnknownPlaceholders,
    interpolateCommand,
    isConnectionEqual,
    NodeInputFormat,
    randomUnicodeAnimal,
    SizeFormatter,
    SxPropsFormatter
} from "./HelperUtils"

// Extend dayjs with UTC plugin for DateTimeFormatter tests
dayjs.extend(utc)

describe("shortUuid", () => {
  it("should return the first 8 characters of a UUID", () => {
    const uuid = "12345678-90ab-cdef-1234-567890abcdef"
    expect(getShortUuid(uuid)).toBe("12345678")
  })
})

describe("getDomain", () => {
  it("should return the domain string from a NodeConfig object", () => {
    const config: NodeConfig = {host: "localhost", keeperPort: 8008}
    expect(getDomain(config)).toBe("localhost:8008::")
  })

  it("should return the domain string without port if keeperPort is undefined", () => {
    const config: NodeConfig = {host: "localhost"}
    expect(getDomain(config)).toBe("localhost:::")
  })
})

describe("getDomains", () => {
    it("should return an array of domain strings from an array of NodeConfig objects", () => {
        const configs: NodeConfig[] = [
            {host: "localhost", keeperPort: 8008},
            {host: "127.0.0.1", keeperPort: 8008},
        ]
        expect(getDomains(configs)).toEqual(["localhost:8008::", "127.0.0.1:8008::"])
    })
})

describe("getNodeConnection", () => {
    it("should return a NodeConfig object from a domain string", () => {
        const domain = "localhost:8008"
        expect(getNodeConfig(domain)).toEqual({
            host: "localhost",
            keeperPort: 8008,
        })
    })

    it("should return a NodeConfig object with default port if port is not in domain string", () => {
        const domain = "localhost"
        expect(getNodeConfig(domain)).toEqual({
            host: "localhost",
        })
    })
})

describe("getNodeConnections", () => {
    it("should return an array of NodeConfig objects from an array of domain strings", () => {
        const domains = ["localhost:8008", "127.0.0.1"]
        expect(getNodeConfigs(domains)).toEqual([
            {host: "localhost", keeperPort: 8008},
            {host: "127.0.0.1"},
        ])
    })
})

describe("getNodeConfig with format", () => {
    const withKeeperPort: NodeInputFormat = {
        withKeeperPort: true,
        defaults: {keeperPort: 8008, dbPort: 5432, sshPort: 22},
    }
    const withoutKeeperPort: NodeInputFormat = {
        withKeeperPort: false,
        defaults: {dbPort: 2379, sshPort: 22},
    }

    it("should fill missing segments from defaults", () => {
        expect(getNodeConfig("node1", withKeeperPort)).toEqual({
            host: "node1", keeperPort: 8008, dbPort: 5432, sshPort: 22,
        })
    })

    it("should keep provided segments over defaults", () => {
        expect(getNodeConfig("node1:8009:5433:2222", withKeeperPort)).toEqual({
            host: "node1", keeperPort: 8009, dbPort: 5433, sshPort: 2222,
        })
    })

    it("should parse host:dbPort:sshPort and mirror keeperPort from dbPort", () => {
        expect(getNodeConfig("node1:2381:2222", withoutKeeperPort)).toEqual({
            host: "node1", keeperPort: 2381, dbPort: 2381, sshPort: 2222,
        })
    })

    it("should mirror the default dbPort into keeperPort when only host is given", () => {
        expect(getNodeConfig("node1", withoutKeeperPort)).toEqual({
            host: "node1", keeperPort: 2379, dbPort: 2379, sshPort: 22,
        })
    })
})

describe("getKeeperDefaultPort", () => {
    it("should return the keeper port when the plugin has a separate one", () => {
        const spec: KeeperDeploySpecResponse = {dbPort: 5432, keeperPort: 8008, credentials: true, dbUser: "postgres"}
        expect(getKeeperDefaultPort(spec)).toBe(8008)
    })

    it("should fall back to the db port when there is no separate keeper endpoint", () => {
        const spec: KeeperDeploySpecResponse = {dbPort: 2379, credentials: true, dbUser: "root"}
        expect(getKeeperDefaultPort(spec)).toBe(2379)
    })
})

describe("getPlaceholders", () => {
    it("should list every variable in order, deduplicated", () => {
        const command = "docker run --name {{name}} -p {{dbPort}}:{{dbPort}} --hostname {{host}}"
        expect(getPlaceholders(command)).toEqual(["{{name}}", "{{dbPort}}", "{{host}}"])
    })

    it("should return nothing for a command with no variables", () => {
        expect(getPlaceholders("docker run -d redis:7")).toEqual([])
    })

    it("should not match docker's own template syntax", () => {
        expect(getPlaceholders("--format '{{json .}}'")).toEqual([])
    })
})

describe("interpolateCommand", () => {
    const values = {cluster: "main", name: "etcd-1", host: "10.0.0.1", sshPort: 22, keeperPort: 8008, dbPort: 2379}

    it("should fill every value a deploy supplies", () => {
        const command = "--name {{name}} --hostname {{host}} -p {{dbPort}} -e TOKEN={{cluster}}"
        expect(interpolateCommand(command, values)).toBe("--name etcd-1 --hostname 10.0.0.1 -p 2379 -e TOKEN=main")
    })

    it("should leave the credentials unresolved when the form has none", () => {
        expect(interpolateCommand("-e PASS={{dbPass}} -u {{dbUser}}", values)).toBe("-e PASS={{dbPass}} -u {{dbUser}}")
    })

    // NOTE: the username is the form's own, so a preview may show it; the
    // password never leaves the server, so what a preview gets is the mask
    it("should fill the username and mask the password", () => {
        const command = "-e PASS={{dbPass}} -u {{dbUser}}"
        expect(interpolateCommand(command, {...values, dbUser: "postgres", dbPass: DeployPasswordMask}))
            .toBe("-e PASS=***** -u postgres")
    })

    it("should leave a value that has not been filled in yet", () => {
        expect(interpolateCommand("--name {{name}} --hostname {{host}}", {name: "etcd-1"}))
            .toBe("--name etcd-1 --hostname {{host}}")
    })

    it("should not touch docker's own template syntax", () => {
        expect(interpolateCommand("--format '{{json .}}'", values)).toBe("--format '{{json .}}'")
    })
})

describe("getUnknownPlaceholders", () => {
    it("should accept every variable in the closed vocabulary", () => {
        const command = Object.values(DeployVar).join(" ")
        expect(getUnknownPlaceholders(command)).toEqual([])
    })

    // NOTE: this is what keeps the vocabulary closed - a typo has to be an
    // error, never a silently-introduced new variable
    it("should report a misspelled variable", () => {
        expect(getUnknownPlaceholders("-e USER={{dbUesr}}")).toEqual(["{{dbUesr}}"])
    })

    // NOTE: a peer port, a member list, a coordinator address and the leader's
    // host are written literally now - they are values only the operator
    // knows, not variables
    it("should report a variable retired from the vocabulary", () => {
        const command = "-p {{peerPort}} -e HOSTS={{clusterHosts}} -e DCS={{dcs}} -h {{leaderHost}}"
        expect(getUnknownPlaceholders(command))
            .toEqual(["{{peerPort}}", "{{clusterHosts}}", "{{dcs}}", "{{leaderHost}}"])
    })
})



describe("isConnectionEqual", () => {
    it("should return true if connections are equal", () => {
        const c1: NodeConfig = {host: "localhost", keeperPort: 8008}
        const c2: NodeConfig = {host: "localhost", keeperPort: 8008}
        expect(isConnectionEqual(c1, c2)).toBe(true)
    })

    it("should return false if connections are not equal", () => {
        const c1: NodeConfig = {host: "localhost", keeperPort: 8008}
        const c2: NodeConfig = {host: "localhost", keeperPort: 8009}
        expect(isConnectionEqual(c1, c2)).toBe(false)
    })
})

describe("SxPropsFormatter", () => {
    describe("merge", () => {
        it("should merge two sx props objects", () => {
            const sx1 = {color: "red"}
            const sx2 = {backgroundColor: "blue"}
            const result = SxPropsFormatter.merge(sx1, sx2)
            expect(result).toEqual([sx1, sx2])
        })

        it("should merge sx props arrays", () => {
            const sx1 = [{color: "red"}, {fontSize: 14}]
            const sx2 = [{backgroundColor: "blue"}]
            const result = SxPropsFormatter.merge(sx1, sx2)
            expect(result).toEqual([{color: "red"}, {fontSize: 14}, {backgroundColor: "blue"}])
        })

        it("should handle undefined sx props", () => {
            const sx1 = {color: "red"}
            const result = SxPropsFormatter.merge(sx1, undefined)
            expect(result).toEqual([sx1, undefined])
        })
    })

    describe("style", () => {
        it("should have paper style", () => {
            expect(SxPropsFormatter.style.paper).toEqual({
                backgroundImage: "linear-gradient(rgba(255, 255, 255, 0.09), rgba(255, 255, 255, 0.09))"
            })
        })

        it("should generate bgImageError style", () => {
            const theme = createTheme({palette: {mode: "dark"}})
            const bgImageError = SxPropsFormatter.style.bgImageError
            if (typeof bgImageError === "function") {
                const result = bgImageError(theme) as {backgroundImage?: string}
                expect(result.backgroundImage).toBeDefined()
                expect(result.backgroundImage).toContain(theme.palette.error.dark)
            }
        })

        it("should generate bgImageSelected style", () => {
            const theme = createTheme({palette: {mode: "dark"}})
            const bgImageSelected = SxPropsFormatter.style.bgImageSelected
            if (typeof bgImageSelected === "function") {
                const result = bgImageSelected(theme) as {backgroundImage?: string}
                expect(result.backgroundImage).toBeDefined()
                expect(result.backgroundImage).toContain(theme.palette.action.hover)
            }
        })
    })
})

describe("DateTimeFormatter", () => {
    describe("utc", () => {
        it("should convert UTC time to local time with timezone", () => {
            const utcTime = "2024-01-15 10:30:00"
            const result = DateTimeFormatter.utc(utcTime)
            // Result should be in format "YYYY-MM-DD HH:mm Z"
            expect(result).toMatch(/\d{4}-\d{2}-\d{2} \d{2}:\d{2} [+-]\d{2}:\d{2}/)
        })
    })
})

describe("SizeFormatter", () => {
    describe("pretty", () => {
        it("should format bytes to human readable format", () => {
            const result1 = SizeFormatter.pretty(1024)
            const result2 = SizeFormatter.pretty(1048576)
            const result3 = SizeFormatter.pretty(500)
            // Just check that it returns a string with the expected unit
            expect(result1).toContain("K")
            expect(result2).toContain("M")
            expect(result3).toContain("B")
        })

        it("should handle zero bytes", () => {
            const result = SizeFormatter.pretty(0)
            expect(result).toContain("B")
        })
    })
})

describe("randomUnicodeAnimal", () => {
    it("should return a unicode animal string", () => {
        const animal = randomUnicodeAnimal()
        expect(typeof animal).toBe("string")
        expect(animal.length).toBeGreaterThan(0)
    })
})

describe("getErrorMessage", () => {
    it("should extract message from axios error with error field", () => {
        const error = new AxiosError("Request failed")
        error.response = {
            data: {error: "Test error message"},
            status: 400,
            statusText: "Bad Request",
            headers: {},
            config: {} as any,
        }
        expect(getErrorMessage(error)).toBe("Test error message")
    })

    it("should use error message if no response", () => {
        const error = new AxiosError("Network error")
        expect(getErrorMessage(error)).toBe("Network error")
    })

    it("should use status text if response has no data", () => {
        const error = new AxiosError("Request failed")
        error.response = {
            data: null,
            status: 500,
            statusText: "Internal Server Error",
            headers: {},
            config: {} as any,
        }
        expect(getErrorMessage(error)).toBe("500 Internal Server Error")
    })

    it("should convert string errors to string", () => {
        const error = "Simple error string"
        expect(getErrorMessage(error)).toBe("Simple error string")
    })

    it("should return unknown for unhandled error types", () => {
        const error = {someField: "value"}
        expect(getErrorMessage(error)).toBe("unknown")
    })
})
