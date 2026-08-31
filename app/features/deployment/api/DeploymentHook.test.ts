import {describe, expect, it} from "@rstest/core"
import {act, renderHook} from "@testing-library/react"

import {KeeperPlugin, PlatformPlugin} from "../../node/api/NodeType"
import {getUnknownCommandPlaceholders, useTemplateForm} from "./DeploymentHook"
import {TemplateCommand, TemplateRequest} from "./DeploymentType"

const etcdCommand: TemplateCommand = {
    command: `docker run -d --name {{name}} -p 2380:2380 -e ETCD_INITIAL_CLUSTER="etcd-1=http://etcd-1:2380" etcd`,
    postScripts: [`etcdctl user add {{dbUser}}:{{dbPass}}`],
}

function initialTemplate(): TemplateRequest {
    return {name: "etcd", keeper: KeeperPlugin.NATIVE_ETCD, platform: PlatformPlugin.DOCKER, commands: [etcdCommand]}
}

describe("getUnknownCommandPlaceholders", () => {

    it("should accept a command built only from the closed vocabulary", () => {
        expect(getUnknownCommandPlaceholders(etcdCommand)).toEqual([])
    })

    it("should report an unknown variable from either the command or the post script", () => {
        const command: TemplateCommand = {command: "docker run {{nope}}", postScripts: ["echo {{alsoNope}}"]}
        expect(getUnknownCommandPlaceholders(command)).toEqual(["{{nope}}", "{{alsoNope}}"])
    })

    // NOTE: everything only the operator knows - ports, member lists,
    // coordinators, which node leads - is literal text in the command now
    it("should reject the variables that became literal values", () => {
        const command: TemplateCommand = {command: "docker run -p {{peerPort}} -e HOSTS={{clusterHosts}}"}
        expect(getUnknownCommandPlaceholders(command)).toEqual(["{{peerPort}}", "{{clusterHosts}}"])
    })
})

describe("useTemplateForm", () => {

    it("should start valid for a named template of known variables", () => {
        const {result} = renderHook(() => useTemplateForm(initialTemplate()))
        expect(result.current.valid).toBe(true)
        expect(result.current.unknown).toEqual([[]])
    })

    // NOTE: kept per command, because a template has several and the warning
    // has to say which one it came from
    it("should report unknown variables against the command that used them", () => {
        const {result} = renderHook(() => useTemplateForm(initialTemplate()))

        act(() => result.current.addCommand({command: "docker run {{nope}}"}))

        expect(result.current.unknown).toEqual([[], ["{{nope}}"]])
        expect(result.current.valid).toBe(false)
    })

    it("should be invalid without a name", () => {
        const {result} = renderHook(() => useTemplateForm({...initialTemplate(), name: "   "}))
        expect(result.current.valid).toBe(false)
    })

    it("should be invalid with no commands left", () => {
        const {result} = renderHook(() => useTemplateForm(initialTemplate()))

        act(() => result.current.removeCommand(0))

        expect(result.current.template.commands).toEqual([])
        expect(result.current.valid).toBe(false)
    })

    // NOTE: a blank command is a node the template counts and cannot run, and
    // adding one is how a new node starts - so the form stays invalid until it
    // is written
    it("should be invalid while any command is blank", () => {
        const {result} = renderHook(() => useTemplateForm(initialTemplate()))

        act(() => result.current.addCommand())

        expect(result.current.valid).toBe(false)

        act(() => result.current.updateCommand(1, {command: "docker run -d --name {{name}} etcd"}))

        expect(result.current.valid).toBe(true)
    })

    it("should add an empty command by default and a copy of one when given", () => {
        const {result} = renderHook(() => useTemplateForm(initialTemplate()))

        act(() => result.current.addCommand())
        act(() => result.current.addCommand(etcdCommand))

        expect(result.current.template.commands).toEqual([etcdCommand, {command: ""}, etcdCommand])
    })

    it("should replace only the command at the given index", () => {
        const {result} = renderHook(() => useTemplateForm(initialTemplate()))
        const replacement: TemplateCommand = {command: "docker run -d --name {{name}} redis"}

        act(() => result.current.addCommand())
        act(() => result.current.updateCommand(1, replacement))

        expect(result.current.template.commands).toEqual([etcdCommand, replacement])
    })
})
