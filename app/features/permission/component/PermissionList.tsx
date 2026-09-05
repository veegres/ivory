import {Search} from "@mui/icons-material"
import {Box, InputBase} from "@mui/material"
import {useMemo, useState} from "react"

import {NoBox} from "../../../shared/component/box/NoBox"
import {SxPropsMap} from "../../../shared/helper/HelperType"
import {Feature} from "../../Feature"
import {PermissionMap, Status} from "../api/PermissionType"
import {PermissionButtons} from "./PermissionButtons"
import {PermissionListItem} from "./PermissionListItem"

const SX: SxPropsMap = {
    box: {display: "flex", flexDirection: "column", width: "100%"},
    header: {
        display: "flex", justifyContent: "space-between", alignItems: "center", padding: "0px 8px",
        height: "35px", position: "sticky", top: 0, backgroundImage: "var(--Paper-overlay)",
        backgroundColor: "background.default", zIndex: 1,
    },
    search: {display: "flex", alignItems: "center", gap: 1, flexGrow: 1, height: "27px"},
    filter: {width: "100%"},
}

type Props = {
    username: string,
    permissions?: PermissionMap,
    view?: "admin" | "user",
}

export function PermissionList(props: Props) {
    const {permissions = {}, username, view = "user"} = props
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
                <PermissionButtons
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
            <PermissionListItem key={name} username={username} name={name} status={status} view={view}/>
        ))
    }

    function handleMemoSearch() {
        return Object.entries(permissions).filter(([name]) => name.includes(search)) as [Feature, Status][]
    }
}