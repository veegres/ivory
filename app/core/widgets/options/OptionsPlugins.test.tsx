import {fireEvent, render, screen} from "@testing-library/react"
import {describe, expect, it, vi} from "vitest"

import {Plugins} from "../../../features/cluster/api/type"
import {KeeperPlugin} from "../../../features/node/api/type"
import {DbPlugin} from "../../../features/query/api/type"
import {OptionsPlugins} from "./OptionsPlugins"

const PLUGINS: Plugins = {keeper: KeeperPlugin.PATRONI_POSTGRES, database: DbPlugin.POSTGRES}

function openSelect(label: string) {
    fireEvent.mouseDown(screen.getByRole("combobox", {name: label}))
}

describe("OptionsPlugins", () => {

    it("should render both selects with the current selection", () => {
        render(<OptionsPlugins plugins={PLUGINS} onUpdate={() => {}}/>)

        expect(screen.getByRole("combobox", {name: "Keeper Plugin"})).toHaveTextContent("Patroni Postgres")
        expect(screen.getByRole("combobox", {name: "Database Plugin"})).toHaveTextContent("Postgres")
    })

    it("should update the keeper plugin and keep the database plugin", () => {
        const onUpdate = vi.fn()
        render(<OptionsPlugins plugins={PLUGINS} onUpdate={onUpdate}/>)

        openSelect("Keeper Plugin")
        fireEvent.click(screen.getByRole("option", {name: "Native Etcd"}))

        expect(onUpdate).toHaveBeenCalledWith({keeper: KeeperPlugin.NATIVE_ETCD, database: DbPlugin.POSTGRES})
    })

    it("should update the database plugin and keep the keeper plugin", () => {
        const onUpdate = vi.fn()
        render(<OptionsPlugins plugins={{keeper: KeeperPlugin.NATIVE_ETCD, database: DbPlugin.POSTGRES}} onUpdate={onUpdate}/>)

        openSelect("Database Plugin")
        fireEvent.click(screen.getByRole("option", {name: "Etcd"}))

        expect(onUpdate).toHaveBeenCalledWith({keeper: KeeperPlugin.NATIVE_ETCD, database: DbPlugin.ETCD})
    })

    it("should offer all keeper plugins as options", () => {
        render(<OptionsPlugins plugins={PLUGINS} onUpdate={() => {}}/>)

        openSelect("Keeper Plugin")

        expect(screen.getByRole("option", {name: "Patroni Postgres"})).toBeInTheDocument()
        expect(screen.getByRole("option", {name: "Native Postgres"})).toBeInTheDocument()
        expect(screen.getByRole("option", {name: "Native Etcd"})).toBeInTheDocument()
    })
})
