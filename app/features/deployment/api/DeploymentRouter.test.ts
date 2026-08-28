import {describe, expect, it} from "@rstest/core"

import {KeeperPlugin, PlatformPlugin} from "../../node/api/NodeType"
import {DeploymentApi} from "./DeploymentRouter"

describe("DeploymentApi.template.list.key", () => {

    it("should separate cache by keeper and platform", () => {
        const etcd = DeploymentApi.template.list.key({keeper: KeeperPlugin.NATIVE_ETCD, platform: PlatformPlugin.DOCKER})
        const redis = DeploymentApi.template.list.key({keeper: KeeperPlugin.NATIVE_REDIS, platform: PlatformPlugin.DOCKER})
        expect(etcd).not.toEqual(redis)
    })

    // NOTE: refetch matches on a common prefix, so keyCommon has to be a
    // genuine prefix of every filtered key. key() is not one - it appends the
    // filter (or undefined), so invalidating with it leaves a created or
    // deleted template invisible until the dialog is reopened.
    it("should invalidate every filtered list from keyCommon", () => {
        const common = DeploymentApi.template.list.keyCommon()
        const filtered = [
            DeploymentApi.template.list.key({keeper: KeeperPlugin.NATIVE_ETCD, platform: PlatformPlugin.DOCKER}),
            DeploymentApi.template.list.key({keeper: KeeperPlugin.NATIVE_REDIS}),
            DeploymentApi.template.list.key(),
        ]
        for (const key of filtered) {
            expect(key.slice(0, common.length)).toEqual(common)
        }
    })

    it("should not be able to invalidate a filtered list from key()", () => {
        const filtered = DeploymentApi.template.list.key({keeper: KeeperPlugin.NATIVE_ETCD})
        expect(DeploymentApi.template.list.key()).not.toEqual(filtered.slice(0, 4))
    })
})
