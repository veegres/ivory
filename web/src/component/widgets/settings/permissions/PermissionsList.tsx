import {Search} from "@mui/icons-material"
import {Box, InputBase} from "@mui/material"
import {useMemo, useState} from "react"

import {PermissionStatus} from "../../../../api/permission/type"
import {SxPropsMap} from "../../../../app/type"
import {NoBox} from "../../../view/box/NoBox"
import {PermissionsButtons} from "./PermissionsButtons"
import {PermissionsListItem} from "./PermissionsListItem"

const SX: SxPropsMap = {
    box: {display: "flex", flexDirection: "column", width: "100%"},
    header: {display: "flex", justifyContent: "space-between", alignItems: "center", padding: "5px 8px"},
    search: {display: "flex", alignItems: "center", gap: 1, flexGrow: 1},
    filter: {width: "100%"},
}

type Props = {
    username: string,
    permissions: { [permission: string]: PermissionStatus },
    view?: "admin" | "user",
}

export function PermissionsList(props: Props) {
    const {permissions, username, view = "user"} = props
    const [search, setSearch] = useState("")
    const filtered = useMemo(handleMemoSearch, [permissions, search])

    return (
        <Box sx={SX.box}>
            <Box sx={SX.header}>
                <Box sx={SX.search}>
                    <Search fontSize={"small"}/>
                    <InputBase
                        sx={SX.filter}
                        slotProps={{input: {sx: {padding: "0px"}}}}
                        size={"small"}
                        placeholder={"Filter by permission name"}
                        onChange={e => setSearch(e.target.value)}
                    />
                </Box>
                <PermissionsButtons
                    username={username}
                    permissions={filtered}
                    approve={view === "admin"}
                    reject={view === "admin"}
                    request={view === "user"}
                    count={true}
                />
            </Box>
            <Box>{renderList()}</Box>
        </Box>
    )

    function renderList() {
        if (filtered.length === 0) return <NoBox text={"no permissions"}/>
        return filtered.map(([name, status]) => (
            <PermissionsListItem key={name} username={username} name={name} status={status} view={view}/>
        ))
    }

    function handleMemoSearch() {
        return Object.entries(permissions).filter(([name]) => name.includes(search))
    }
}