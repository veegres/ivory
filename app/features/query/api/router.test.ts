import {describe, expect, it} from "vitest"

import {QueryApi} from "./router"
import {DbPlugin, Type} from "./type"

describe("QueryApi.list.key", () => {

    it("should separate cache by type and plugin", () => {
        expect(QueryApi.list.key(Type.BLOAT, DbPlugin.POSTGRES)).toEqual(["query", "list", Type.BLOAT, DbPlugin.POSTGRES])
        expect(QueryApi.list.key(Type.BLOAT, DbPlugin.ETCD)).toEqual(["query", "list", Type.BLOAT, DbPlugin.ETCD])
    })

    it("should share the [query, list] prefix used for broad invalidation", () => {
        const key = QueryApi.list.key(Type.ACTIVITY, DbPlugin.ETCD)
        expect(key.slice(0, 2)).toEqual(["query", "list"])
    })
})
